package models

import "github.com/google/uuid"

type Permission struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Scope  string    `gorm:"type:varchar(100);not null" json:"scope"`
}
