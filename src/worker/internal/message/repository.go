package message

import (
	"context"
	"errors"

	"github.com/tmlsergen/messaging-service-worker/internal/app"
	"gorm.io/gorm"
)

var ErrNotFound = app.Errorf("record not found")

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetPendingMessages(c *context.Context) ([]Message, error) {
	var messages []Message

	err := r.db.Model(&Message{}).Where("status = ?", pending).Limit(2).Find(&messages).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, ErrNotFound
		}

		return nil, app.ErrorWithCaller(err)
	}

	return messages, nil
}

func (r *repository) Update(c *context.Context, message Message) error {
	err := r.db.Model(&Message{}).Where("id = ?", message.ID).Updates(message).Error
	if err != nil {
		return app.ErrorWithCaller(err)
	}

	return nil
}
