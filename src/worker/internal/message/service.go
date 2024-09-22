package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/tmlsergen/messaging-service-worker/internal/app"
	msgSendler "github.com/tmlsergen/messaging-service-worker/pkg/msg-sendler"
	"github.com/tmlsergen/messaging-service-worker/pkg/redis"
)

var ErrJSONMarshal = app.Errorf("error marshalling json")

const (
	sentMessageKey      = "sent:%s"
	processedMessageKey = "processed:%s"
)

type messageRepository interface {
	GetPendingMessages(*context.Context) ([]Message, error)
	Update(*context.Context, Message) error
}

type service struct {
	repo       messageRepository
	rds        *redis.RedisClient
	msgSendler *msgSendler.Client
	logger     *slog.Logger
}

func NewService(repo messageRepository, rds *redis.RedisClient, msgSendler *msgSendler.Client, logger *slog.Logger) *service {
	return &service{repo: repo, rds: rds, msgSendler: msgSendler, logger: logger}
}

func (s *service) SendPendingMessages(c context.Context) {
	messages, err := s.repo.GetPendingMessages(&c)
	if err != nil {
		s.logger.Error("failed to get pending messages", "error", err)
		return
	}

	ch := make(chan error, len(messages))

	for _, message := range messages {
		go s.sendMessage(&c, message, ch)
	}

	for range messages {
		err := <-ch
		if err != nil {
			s.logger.Error("failed to send message", "error", err)
		}
	}
}

type MessageResponse struct {
	MessageID string    `json:"message_id"`
	Message   string    `json:"message"`
	SentTime  time.Time `json:"sent_time"`
}

func (r *MessageResponse) ToJSON() (string, error) {
	json, err := json.Marshal(r)
	if err != nil {
		return "", ErrJSONMarshal
	}
	return string(json), nil
}

func (s *service) sendMessage(c *context.Context, message Message, ch chan error) {
	now := time.Now()
	resp, err := s.msgSendler.SendMessage(*c, message.Content, message.To)
	if err != nil {
		ch <- app.ErrorWithCaller(err)
		return
	}

	s.logger.Info("message sent", "message_id", resp.MessageID)

	err = s.repo.Update(c, Message{
		ID:        message.ID,
		MessageID: resp.MessageID,
		Status:    1,
	})
	if err != nil {
		ch <- app.ErrorWithCaller(err)
		return
	}

	messageKey := fmt.Sprintf(sentMessageKey, resp.MessageID)

	messageResp := MessageResponse{
		MessageID: resp.MessageID,
		Message:   resp.Message,
		SentTime:  now,
	}

	json, err := messageResp.ToJSON()
	if err != nil {
		ch <- app.ErrorWithCaller(err)
	}

	err = s.rds.Set(*c, messageKey, json)
	if err != nil {
		ch <- app.ErrorWithCaller(err)
		return
	}

	ch <- nil
}
