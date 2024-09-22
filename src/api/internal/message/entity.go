package message

import "time"

type Status int

const (
	pending   Status = 0
	processed Status = 1
)

type Message struct {
	ID        uint64    `gorm:"primary_key;auto_increment"`
	MessageID string    `gorm:"type:UUID;unique;null"`
	Content   string    `gorm:"type:varchar(50);not null"`
	To        string    `gorm:"type:varchar(25);not null"`
	Status    Status    `gorm:"default:0;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (m *Message) TableName() string {
	return "messages"
}
