package service

import (
	"context"
	"strings"

	"github.com/skamenetskiy/messages/api"
	entity2 "github.com/skamenetskiy/messages/internal/entity"
	"github.com/skamenetskiy/messages/internal/pkg/mentions"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Create(ctx context.Context, req *api.Create_Request) (*api.Create_Response, error) {
	content := strings.TrimSpace(req.GetContent())
	if content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is empty")
	}
	var m entity2.Mentions
	if len(req.Mentions) > 0 {
		m = mentions.Find(content, req.GetMentions())
	}
	e, err := s.r.Insert(ctx, entity2.Message{
		ThreadID:  req.GetThreadId(),
		AccountID: req.GetAccountId(),
		Mentions:  m,
		Content:   content,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.Create_Response{
		Message: messageToProto(e),
	}, nil
}
