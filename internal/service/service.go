package service

import (
	"context"

	"github.com/skamenetskiy/messages/api"
	"github.com/skamenetskiy/messages/internal/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func New(r repo) *Service {
	return &Service{
		r: r,
	}
}

type Service struct {
	api.UnimplementedMessagesAPIServer
	r repo
}

type repo interface {
	Insert(context.Context, entity.Message) (entity.Message, error)
	Get(context.Context, []uint64) ([]entity.Message, error)
	Delete(ctx context.Context, ids []uint64) error
}

func messageToProto(m entity.Message) *api.Message {
	return &api.Message{
		Id:        m.ID,
		ThreadId:  uint64p(m.ThreadID),
		AccountId: uint64p(m.AccountID),
		CreatedAt: timestamppb.New(m.CreatedAt),
		Mentions:  mentionsToProto(m.Mentions),
		Content:   m.Content,
	}
}

func messagesToProto(messages []entity.Message) []*api.Message {
	r := make([]*api.Message, len(messages))
	for i, m := range messages {
		r[i] = messageToProto(m)
	}
	return r
}

func mentionToProto(m entity.Mention) *api.Mention {
	return &api.Mention{
		Id:    m.ID,
		Type:  api.Mention_MentionType(m.Type),
		Start: m.Pos[0],
		End:   m.Pos[1],
	}
}

func mentionsToProto(m entity.Mentions) []*api.Mention {
	r := make([]*api.Mention, len(m))
	for i, v := range m {
		r[i] = mentionToProto(v)
	}
	return r
}

func uint64p(u uint64) *uint64 {
	if u == 0 {
		return nil
	}
	return &u
}
