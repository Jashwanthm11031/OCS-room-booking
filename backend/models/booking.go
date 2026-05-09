package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RoomID           uuid.UUID `gorm:"type:uuid;not null" json:"room_id"`
	Room             *Room     `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User             *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Date             string    `gorm:"type:date;not null" json:"date"`
	StartTime        string    `gorm:"type:time;not null" json:"start_time"`
	EndTime          string    `gorm:"type:time;not null" json:"end_time"`
	Purpose          string    `gorm:"type:varchar(50)" json:"purpose"`
	ParticipantCount int       `gorm:"not null" json:"participant_count"`
	Status           string    `gorm:"type:varchar(50);default:confirmed" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if b.Status == "" {
		b.Status = "confirmed"
	}
	return nil
}
