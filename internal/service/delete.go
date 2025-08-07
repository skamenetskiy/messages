package service

import (
	"context"

	"github.com/skamenetskiy/messages/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Delete(ctx context.Context, req *api.Delete_Request) (*api.Delete_Response, error) {
	if len(req.GetIds()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing ids")
	}
	if err := s.r.Delete(ctx, req.GetIds()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.Delete_Response{}, nil
}
