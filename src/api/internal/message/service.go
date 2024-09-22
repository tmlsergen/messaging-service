package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/tmlsergen/messaging-service-api/internal/app"
	"github.com/tmlsergen/messaging-service-api/pkg/redis"
)

const (
	sentMessageKey      = "sent:%s"
	processedMessageKey = "processed:%s"
	cronConfigKey       = "cron"
	messagesKey         = "messages_%d_%d"
)

var ErrJSONEncode = app.Errorf("failed to encode json")
var ErrJSONDecode = app.Errorf("failed to decode json")

type messageRepository interface {
	GetSendingMessages(*fiber.Ctx, int, int) ([]Message, error)
	Update(*fiber.Ctx, Message) error
	GetTotalCountOfSendingMessages(*fiber.Ctx) (int, error)
}

type service struct {
	repo messageRepository
	rds  *redis.RedisClient
}

func NewService(repo messageRepository, rds *redis.RedisClient) *service {
	return &service{repo: repo, rds: rds}
}

type GetMessages struct {
	ID        uint64    `json:"id"`
	MessageID string    `json:"message_id"`
	Content   string    `json:"content"`
	To        string    `json:"to"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetMessagesResponse struct {
	Messages  []GetMessages `json:"messages"`
	Total     int           `json:"total"`
	Page      int           `json:"page"`
	Limit     int           `json:"limit"`
	PageCount int           `json:"page_count"`
}

func (s *service) GetMessages(c *fiber.Ctx, page, limit int) (GetMessagesResponse, error) {
	cacheKey := fmt.Sprintf(messagesKey, page, limit)
	cache, err := s.rds.Get(c.Context(), cacheKey)
	if err != nil && !errors.Is(err, redis.RedisNilError) {
		return GetMessagesResponse{}, app.ErrorWithCaller(err)
	}

	if cache != "" {

		var resp GetMessagesResponse
		err = json.Unmarshal([]byte(cache), &resp)
		if err != nil {
			return GetMessagesResponse{}, ErrJSONDecode
		}

		if resp.Total > 0 {
			return resp, nil
		}
	}

	messages, err := s.repo.GetSendingMessages(c, page, limit)
	if err != nil {
		return GetMessagesResponse{}, err
	}

	getMessages := []GetMessages{}
	lo.ForEach(messages, func(message Message, _ int) {
		getMessages = append(getMessages, GetMessages(message))
	})

	total, err := s.repo.GetTotalCountOfSendingMessages(c)
	if err != nil {
		return GetMessagesResponse{}, err
	}

	resp := GetMessagesResponse{
		Messages:  getMessages,
		Total:     total,
		Page:      page,
		Limit:     limit,
		PageCount: len(getMessages),
	}

	json, err := json.Marshal(resp)
	if err != nil {
		return GetMessagesResponse{}, ErrJSONEncode
	}

	err = s.rds.Set(c.Context(), cacheKey, string(json))
	if err != nil {
		return GetMessagesResponse{}, app.ErrorWithCaller(err)
	}

	return resp, nil
}

func (s *service) HandleCronAction(c *fiber.Ctx, action string) error {
	err := s.rds.Set(c.Context(), cronConfigKey, action)
	if err != nil {
		return err
	}

	return nil
}
