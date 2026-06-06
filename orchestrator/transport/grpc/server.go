package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/evidra/evidra/api/gen/evidra/v1"
	"github.com/evidra/evidra/orchestrator/domain"
	"github.com/evidra/evidra/orchestrator/service"
)

type OrchestratorServer struct {
	evidrav1.UnimplementedOrchestratorServiceServer
	svc *service.OrchestratorService
}

func NewServer(svc *service.OrchestratorService) *OrchestratorServer {
	return &OrchestratorServer{svc: svc}
}

func (s *OrchestratorServer) Answer(ctx context.Context, req *evidrav1.AnswerRequest) (*evidrav1.AnswerResponse, error) {
	tenantID, err := uuid.Parse(req.TenantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid tenant_id")
	}
	result, err := s.svc.Answer(ctx, service.AnswerRequest{
		Question: req.Question,
		TenantID: tenantID,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}
	evIDs := make([]string, len(result.Evidence))
	for i, e := range result.Evidence {
		evIDs[i] = e.ID.String()
	}
	return &evidrav1.AnswerResponse{
		DraftId:     result.Draft.ID.String(),
		Answer:      result.Draft.Answer,
		Confidence:  result.Draft.Confidence,
		ModelUsed:   result.Draft.ModelUsed,
		EvidenceIds: evIDs,
		Reasoning:   result.Draft.Reasoning,
	}, nil
}

func (s *OrchestratorServer) GetDraft(ctx context.Context, req *evidrav1.GetDraftRequest) (*evidrav1.Draft, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	d, err := s.svc.GetDraft(ctx, id)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return draftToProto(d), nil
}

func (s *OrchestratorServer) ListDrafts(ctx context.Context, req *evidrav1.ListDraftsRequest) (*evidrav1.ListDraftsResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 {
		limit = 20
	}
	offset := int((req.Page - 1) * req.PageSize)
	if offset < 0 {
		offset = 0
	}
	drafts, err := s.svc.ListDrafts(ctx, limit, offset)
	if err != nil {
		return nil, toGRPCError(err)
	}
	proto := make([]*evidrav1.Draft, len(drafts))
	for i, d := range drafts {
		proto[i] = draftToProto(d)
	}
	return &evidrav1.ListDraftsResponse{Drafts: proto, Total: int32(len(drafts))}, nil
}

func draftToProto(d *domain.Draft) *evidrav1.Draft {
	evIDs := make([]string, len(d.EvidenceIDs))
	for i, eID := range d.EvidenceIDs {
		evIDs[i] = eID.String()
	}
	return &evidrav1.Draft{
		Id:           d.ID.String(),
		QuestionId:   d.QuestionID.String(),
		QuestionText: d.QuestionText,
		Answer:       d.Answer,
		Confidence:   d.Confidence,
		ModelUsed:    d.ModelUsed,
		EvidenceIds:  evIDs,
		Reasoning:    d.Reasoning,
		Status:       string(d.Status),
		Feedback:     d.Feedback,
		CreatedAt:    timestamppb.New(d.CreatedAt),
		UpdatedAt:    timestamppb.New(d.UpdatedAt),
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
