package models

import "github.com/google/uuid"

type Block struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
}
