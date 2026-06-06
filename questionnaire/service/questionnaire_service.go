package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/evidra/evidra/pkg/queue"
	"github.com/evidra/evidra/pkg/storage"
	"github.com/evidra/evidra/questionnaire/domain"
	"github.com/evidra/evidra/questionnaire/events"
	"github.com/evidra/evidra/questionnaire/parser"
	"github.com/evidra/evidra/questionnaire/repository"
)

type QuestionnaireService struct {
	questionnaires repository.QuestionnaireRepository
	questions      repository.QuestionRepository
	store          storage.FileStorage
	bus            queue.EventBus
	v              *validator.Validate
	bucket         string
	maxFileSize    int64
}

func New(
	questionnaires repository.QuestionnaireRepository,
	questions repository.QuestionRepository,
	store storage.FileStorage,
	bus queue.EventBus,
	bucket string,
	maxFileSize int64,
) *QuestionnaireService {
	return &QuestionnaireService{
		questionnaires: questionnaires,
		questions:      questions,
		store:          store,
		bus:            bus,
		v:              validator.New(),
		bucket:         bucket,
		maxFileSize:    maxFileSize,
	}
}

func (s *QuestionnaireService) Upload(ctx context.Context, tenantID uuid.UUID, title string, file multipart.File, header *multipart.FileHeader) (*domain.Questionnaire, error) {
	return s.UploadBytes(ctx, tenantID, title, header.Filename, file, header.Size, header.Header.Get("Content-Type"))
}

func (s *QuestionnaireService) UploadBytes(ctx context.Context, tenantID uuid.UUID, title, filename string, reader io.Reader, size int64, contentType string) (*domain.Questionnaire, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if _, err := parser.ExtractorFor(ext); err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrUnsupportedFile, ext)
	}

	if size > s.maxFileSize {
		return nil, fmt.Errorf("%w: max size is %d bytes", domain.ErrFileTooLarge, s.maxFileSize)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	objectKey := fmt.Sprintf("questionnaires/%s/%s%s", tenantID, uuid.New().String(), ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if err := s.store.Upload(ctx, s.bucket, objectKey, bytes.NewReader(data), int64(len(data)), contentType); err != nil {
		return nil, fmt.Errorf("upload to storage: %w", err)
	}

	fileURL := fmt.Sprintf("s3://%s/%s", s.bucket, objectKey)
	q := domain.NewQuestionnaire(tenantID, title, filename, fileURL, ext, size)

	if err := s.questionnaires.Create(ctx, q); err != nil {
		return nil, err
	}

	event := events.QuestionnaireUploaded{
		ID:       q.ID,
		TenantID: tenantID,
		FileURL:  fileURL,
		FileType: ext,
	}

	if err := s.bus.Publish(ctx, event); err != nil {
		slog.Warn("failed to publish upload event", "questionnaire_id", q.ID, "error", err)
	}

	q.Status = domain.StatusQueued
	q.UpdatedAt = time.Now()
	if err := s.questionnaires.Update(ctx, q); err != nil {
		slog.Warn("failed to update status to queued", "questionnaire_id", q.ID, "error", err)
	}

	return q, nil
}

func (s *QuestionnaireService) ProcessDocument(ctx context.Context, event events.QuestionnaireUploaded) error {
	q, err := s.questionnaires.GetByID(ctx, event.ID)
	if err != nil {
		return fmt.Errorf("get questionnaire: %w", err)
	}

	if err := q.TransitionStatus(domain.StatusParsing); err != nil {
		return err
	}
	if err := s.questionnaires.Update(ctx, q); err != nil {
		return fmt.Errorf("update status to parsing: %w", err)
	}

	reader, err := s.store.Download(ctx, s.bucket, strings.TrimPrefix(event.FileURL, fmt.Sprintf("s3://%s/", s.bucket)))
	if err != nil {
		s.failQuestionnaire(ctx, q, fmt.Errorf("download: %w", err))
		return err
	}
	defer reader.Close()

	extractor, err := parser.ExtractorFor(event.FileType)
	if err != nil {
		s.failQuestionnaire(ctx, q, err)
		return err
	}

	text, err := extractor.Extract(ctx, reader)
	if err != nil {
		s.failQuestionnaire(ctx, q, fmt.Errorf("extract: %w", err))
		return err
	}

	detector := parser.NewDetector()
	detected := detector.Find(text)

	var questions []*domain.Question
	for _, dq := range detected {
		questions = append(questions, domain.NewQuestion(q.ID, dq.Text, dq.Type, dq.Order, dq.Options))
	}

	if err := s.questions.CreateBatch(ctx, questions); err != nil {
		s.failQuestionnaire(ctx, q, fmt.Errorf("save questions: %w", err))
		return err
	}

	if err := q.TransitionStatus(domain.StatusParsed); err != nil {
		return err
	}
	if err := s.questionnaires.Update(ctx, q); err != nil {
		return fmt.Errorf("update status to parsed: %w", err)
	}

	parsedEvent := events.QuestionnaireParsed{
		ID:              q.ID,
		TenantID:        q.TenantID,
		QuestionCount:   len(questions),
		QuestionnaireID: q.ID,
	}
	if err := s.bus.Publish(ctx, parsedEvent); err != nil {
		slog.Warn("failed to publish parsed event", "questionnaire_id", q.ID, "error", err)
	}

	slog.Info("document processed",
		"questionnaire_id", q.ID,
		"questions", len(questions),
	)
	return nil
}

func (s *QuestionnaireService) failQuestionnaire(ctx context.Context, q *domain.Questionnaire, err error) {
	q.Status = domain.StatusFailed
	q.UpdatedAt = time.Now()
	if updateErr := s.questionnaires.Update(ctx, q); updateErr != nil {
		slog.Error("failed to update failed status", "error", updateErr)
	}

	failEvent := events.QuestionnaireFailed{
		ID:       q.ID,
		TenantID: q.TenantID,
		Error:    err.Error(),
	}
	if pubErr := s.bus.Publish(ctx, failEvent); pubErr != nil {
		slog.Warn("failed to publish fail event", "error", pubErr)
	}
}

func (s *QuestionnaireService) GetQuestionnaire(ctx context.Context, id uuid.UUID) (*domain.Questionnaire, error) {
	return s.questionnaires.GetByID(ctx, id)
}

func (s *QuestionnaireService) ListQuestionnaires(ctx context.Context, tenantID uuid.UUID) ([]*domain.Questionnaire, error) {
	return s.questionnaires.ListByTenant(ctx, tenantID)
}

func (s *QuestionnaireService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.questionnaires.Delete(ctx, id)
}

func (s *QuestionnaireService) GetQuestions(ctx context.Context, questionnaireID uuid.UUID) ([]*domain.Question, error) {
	return s.questions.ListByQuestionnaire(ctx, questionnaireID)
}
