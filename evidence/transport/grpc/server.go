package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/evidra/evidra/api/gen/evidra/v1"
	"github.com/evidra/evidra/evidence/domain"
	"github.com/evidra/evidra/evidence/service"
)

type EvidenceServer struct {
	evidrav1.UnimplementedEvidenceServiceServer
	svc *service.EvidenceService
}

func NewServer(svc *service.EvidenceService) *EvidenceServer {
	return &EvidenceServer{svc: svc}
}

func (s *EvidenceServer) Create(ctx context.Context, req *evidrav1.CreateEvidenceRequest) (*evidrav1.EvidenceItem, error) {
	tenantID, err := uuid.Parse(req.TenantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid tenant_id")
	}
	ownerID, err := uuid.Parse(req.OwnerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid owner_id")
	}
	cat, ok := domain.ParseCategory(req.Category)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid category")
	}
	item, err := s.svc.Create(ctx, domain.CreateEvidenceInput{
		TenantID:  tenantID,
		Title:     req.Title,
		Content:   req.Content,
		Category:  cat,
		OwnerID:   ownerID,
		SourceURL: req.SourceUrl,
		Tags:      req.Tags,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return evidenceToProto(item), nil
}

func (s *EvidenceServer) Get(ctx context.Context, req *evidrav1.GetEvidenceRequest) (*evidrav1.EvidenceItem, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	item, err := s.svc.Get(ctx, id)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return evidenceToProto(item), nil
}

func (s *EvidenceServer) List(ctx context.Context, req *evidrav1.ListEvidenceRequest) (*evidrav1.ListEvidenceResponse, error) {
	ownerID, _ := uuid.Parse(req.OwnerId)
	filter := service.ListFilter{
		TenantID: uuid.MustParse(req.TenantId),
		Category: req.Category,
		Status:   req.Status,
		OwnerID:  ownerID,
		Limit:    int(req.PageSize),
		Offset:   int((req.Page - 1) * req.PageSize),
	}
	items, total, err := s.svc.List(ctx, filter)
	if err != nil {
		return nil, toGRPCError(err)
	}
	proto := make([]*evidrav1.EvidenceItem, len(items))
	for i, item := range items {
		proto[i] = evidenceToProto(item)
	}
	return &evidrav1.ListEvidenceResponse{Items: proto, Total: int32(total)}, nil
}

func (s *EvidenceServer) Update(ctx context.Context, req *evidrav1.UpdateEvidenceRequest) (*evidrav1.EvidenceItem, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	item, err := s.svc.Update(ctx, id, domain.CreateEvidenceInput{
		Title:     req.Title,
		Content:   req.Content,
		SourceURL: req.SourceUrl,
		Tags:      req.Tags,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return evidenceToProto(item), nil
}

func (s *EvidenceServer) Delete(ctx context.Context, req *evidrav1.DeleteEvidenceRequest) (*evidrav1.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	if err := s.svc.Delete(ctx, id); err != nil {
		return nil, toGRPCError(err)
	}
	return &evidrav1.Empty{}, nil
}

func (s *EvidenceServer) Submit(ctx context.Context, req *evidrav1.SubmitEvidenceRequest) (*evidrav1.EvidenceItem, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	item, err := s.svc.Submit(ctx, id)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return evidenceToProto(item), nil
}

func (s *EvidenceServer) Approve(ctx context.Context, req *evidrav1.ApproveEvidenceRequest) (*evidrav1.EvidenceItem, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	reviewerID, err := uuid.Parse(req.ReviewerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reviewer_id")
	}
	item, err := s.svc.Approve(ctx, id, reviewerID, req.Comment)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return evidenceToProto(item), nil
}

func (s *EvidenceServer) Reject(ctx context.Context, req *evidrav1.RejectEvidenceRequest) (*evidrav1.EvidenceItem, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	reviewerID, err := uuid.Parse(req.ReviewerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reviewer_id")
	}
	item, err := s.svc.Reject(ctx, id, reviewerID, req.Comment)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return evidenceToProto(item), nil
}

func (s *EvidenceServer) Export(ctx context.Context, req *evidrav1.ExportEvidenceRequest) (*evidrav1.EvidenceItem, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	item, err := s.svc.Export(ctx, id)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return evidenceToProto(item), nil
}

func (s *EvidenceServer) GetApprovalHistory(ctx context.Context, req *evidrav1.GetApprovalHistoryRequest) (*evidrav1.GetApprovalHistoryResponse, error) {
	eID, err := uuid.Parse(req.EvidenceId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid evidence_id")
	}
	approvals, err := s.svc.GetApprovalHistory(ctx, eID)
	if err != nil {
		return nil, toGRPCError(err)
	}
	proto := make([]*evidrav1.Approval, len(approvals))
	for i, a := range approvals {
		proto[i] = &evidrav1.Approval{
			Id:         a.ID.String(),
			EvidenceId: a.EvidenceID.String(),
			ReviewerId: a.ReviewerID.String(),
			Status:     string(a.Status),
			Comment:    a.Comment,
			CreatedAt:  timestamppb.New(a.CreatedAt),
		}
	}
	return &evidrav1.GetApprovalHistoryResponse{Approvals: proto}, nil
}

func evidenceToProto(item *domain.EvidenceItem) *evidrav1.EvidenceItem {
	return &evidrav1.EvidenceItem{
		Id:        item.ID.String(),
		TenantId:  item.TenantID.String(),
		Title:     item.Title,
		Content:   item.Content,
		Category:  string(item.Category),
		Status:    string(item.Status),
		OwnerId:   item.OwnerID.String(),
		SourceUrl: item.SourceURL,
		Tags:      item.Tags,
		Version:   int32(item.Version),
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidInput), errors.Is(err, domain.ErrInvalidCategory):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrForbidden), errors.Is(err, domain.ErrNotOwner):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
