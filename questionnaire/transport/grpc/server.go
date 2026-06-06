package grpc

import (
	"bytes"
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/evidra/evidra/api/gen/evidra/v1"
	"github.com/evidra/evidra/questionnaire/domain"
	"github.com/evidra/evidra/questionnaire/service"
)

type QuestionnaireServer struct {
	evidrav1.UnimplementedQuestionnaireServiceServer
	svc *service.QuestionnaireService
}

func NewServer(svc *service.QuestionnaireService) *QuestionnaireServer {
	return &QuestionnaireServer{svc: svc}
}

func (s *QuestionnaireServer) Upload(ctx context.Context, req *evidrav1.UploadRequest) (*evidrav1.Questionnaire, error) {
	tenantID, err := uuid.Parse(req.TenantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid tenant_id")
	}
	q, err := s.svc.UploadBytes(ctx, tenantID, req.Title, req.FileName, bytes.NewReader(req.Content), int64(len(req.Content)), req.FileType)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return questionnaireToProto(q), nil
}

func (s *QuestionnaireServer) GetQuestionnaire(ctx context.Context, req *evidrav1.GetQuestionnaireRequest) (*evidrav1.Questionnaire, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	q, err := s.svc.GetQuestionnaire(ctx, id)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return questionnaireToProto(q), nil
}

func (s *QuestionnaireServer) ListQuestionnaires(ctx context.Context, req *evidrav1.ListQuestionnairesRequest) (*evidrav1.ListQuestionnairesResponse, error) {
	tenantID, err := uuid.Parse(req.TenantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid tenant_id")
	}
	qs, err := s.svc.ListQuestionnaires(ctx, tenantID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	proto := make([]*evidrav1.Questionnaire, len(qs))
	for i, q := range qs {
		proto[i] = questionnaireToProto(q)
	}
	return &evidrav1.ListQuestionnairesResponse{Questionnaires: proto, Total: int32(len(qs))}, nil
}

func (s *QuestionnaireServer) ListQuestions(ctx context.Context, req *evidrav1.ListQuestionsRequest) (*evidrav1.ListQuestionsResponse, error) {
	qID, err := uuid.Parse(req.QuestionnaireId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid questionnaire_id")
	}
	questions, err := s.svc.GetQuestions(ctx, qID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	proto := make([]*evidrav1.Question, len(questions))
	for i, q := range questions {
		proto[i] = &evidrav1.Question{
			Id:              q.ID.String(),
			QuestionnaireId: q.QuestionnaireID.String(),
			Text:            q.Text,
			Type:            string(q.Type),
			Order:           int32(q.Order),
			Options:         q.Options,
			CreatedAt:       timestamppb.New(q.CreatedAt),
			UpdatedAt:       timestamppb.New(q.UpdatedAt),
		}
	}
	return &evidrav1.ListQuestionsResponse{Questions: proto}, nil
}

func questionnaireToProto(q *domain.Questionnaire) *evidrav1.Questionnaire {
	return &evidrav1.Questionnaire{
		Id:        q.ID.String(),
		TenantId:  q.TenantID.String(),
		Title:     q.Title,
		FileName:  q.FileName,
		FileUrl:   q.FileURL,
		FileType:  q.FileType,
		FileSize:  q.FileSize,
		Status:    string(q.Status),
		Version:   int32(q.Version),
		CreatedAt: timestamppb.New(q.CreatedAt),
		UpdatedAt: timestamppb.New(q.UpdatedAt),
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
