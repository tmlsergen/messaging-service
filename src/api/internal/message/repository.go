package message

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tmlsergen/messaging-service-api/internal/app"
	"gorm.io/gorm"
)

var ErrNotFound = app.Errorf("record not found")

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetSendingMessages(c *fiber.Ctx, page, limit int) ([]Message, error) {
	offset := (page - 1) * limit

	var messages []Message

	err := r.db.Model(&Message{}).Where("status = ?", processed).Order("id asc").Limit(limit).Offset(offset).Find(&messages).Error
	if err != nil {
		return nil, app.ErrorWithCaller(err)
	}

	return messages, nil
}

func (r *repository) GetTotalCountOfSendingMessages(c *fiber.Ctx) (int, error) {
	var count int64

	err := r.db.Model(&Message{}).Where("status = ?", processed).Count(&count).Error
	if err != nil {
		return 0, app.ErrorWithCaller(err)
	}

	return int(count), nil
}

func (r *repository) Update(c *fiber.Ctx, message Message) error {
	err := r.db.Model(&Message{}).Where("id = ?", message.ID).Updates(message).Error
	if err != nil {
		return app.ErrorWithCaller(err)
	}

	return nil
}
