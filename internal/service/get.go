package service

import (
	"context"

	"github.com/skamenetskiy/messages/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Get(ctx context.Context, req *api.Get_Request) (*api.Get_Response, error) {
	if len(req.GetId()) == 0 {
		return &api.Get_Response{Messages: make([]*api.Message, 0)}, nil
	}
	messages, err := s.r.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.Get_Response{
		Messages: messagesToProto(messages),
	}, nil
}
