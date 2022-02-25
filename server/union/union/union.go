package union

import (
	"context"
	auth "gService/share/auth"
	unionpb "gService/union/api/gen/v1"

	"go.uber.org/zap"
)

// Service is a union service.
type Service struct {
	Logger *zap.Logger
}

// GetUnionInfo is a gRPC method.
func (s *Service) GetUnionInfo(ctx context.Context, req *unionpb.UnionRequest) (*unionpb.UnionResponse, error) {
	s.Logger.Info("GetUnionInfo", zap.String("union_id", req.Id))
	userID, err := auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &unionpb.UnionResponse{
		Id:   userID,
		Name: "union",
	}, nil
}

// Ping is a gRPC method.
func (s *Service) Ping(ctx context.Context, req *unionpb.UnionRequest) (*unionpb.UnionResponse, error) {
	s.Logger.Info("Ping", zap.String("union_id", req.Id))
	return &unionpb.UnionResponse{
		Id:   "1111",
		Name: "union",
	}, nil
}
