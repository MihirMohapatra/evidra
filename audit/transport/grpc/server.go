package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/evidra/evidra/api/gen/evidra/v1"
	"github.com/evidra/evidra/audit/domain"
	"github.com/evidra/evidra/audit/repository"
	"github.com/evidra/evidra/audit/service"
)

type AuditServer struct {
	evidrav1.UnimplementedAuditServiceServer
	svc *service.AuditService
}

func NewServer(svc *service.AuditService) *AuditServer {
	return &AuditServer{svc: svc}
}

func (s *AuditServer) ListEvents(ctx context.Context, req *evidrav1.ListAuditEventsRequest) (*evidrav1.ListAuditEventsResponse, error) {
	filter := repository.AuditFilter{}
	if req.TenantId != "" {
		filter.TenantID = uuid.MustParse(req.TenantId)
	}
	if req.ActorId != "" {
		filter.ActorID = uuid.MustParse(req.ActorId)
	}
	if req.Action != "" {
		filter.Action = domain.Action(req.Action)
	}
	filter.TargetID = req.TargetId
	if req.StartTime != nil {
		filter.Since = req.StartTime.AsTime()
	}
	if req.EndTime != nil {
		filter.Until = req.EndTime.AsTime()
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 50
	}
	filter.Limit = pageSize
	filter.Offset = int((req.Page - 1) * req.PageSize)
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	events, total, err := s.svc.List(ctx, filter)
	if err != nil {
		return nil, toGRPCError(err)
	}
	proto := make([]*evidrav1.AuditEvent, len(events))
	for i, e := range events {
		proto[i] = auditEventToProto(e)
	}
	return &evidrav1.ListAuditEventsResponse{Events: proto, Total: int32(total)}, nil
}

func (s *AuditServer) GetEvent(ctx context.Context, req *evidrav1.GetAuditEventRequest) (*evidrav1.AuditEvent, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	events, _, err := s.svc.List(ctx, repository.AuditFilter{Limit: 1})
	if err != nil {
		return nil, toGRPCError(err)
	}
	for _, e := range events {
		if e.ID == id {
			return auditEventToProto(e), nil
		}
	}
	return nil, status.Error(codes.NotFound, "not found")
}

func auditEventToProto(e *domain.AuditEvent) *evidrav1.AuditEvent {
	meta, _ := structpb.NewStruct(e.Metadata)
	return &evidrav1.AuditEvent{
		Id:        e.ID.String(),
		TenantId:  e.TenantID.String(),
		ActorId:   e.ActorID.String(),
		Action:    string(e.Action),
		TargetId:  e.TargetID,
		Timestamp: timestamppb.New(e.Timestamp),
		Metadata:  meta,
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
