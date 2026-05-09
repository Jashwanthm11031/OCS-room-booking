package models

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringArray handles PostgreSQL TEXT[] type
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "{}", nil
	}
	return "{" + strings.Join(s, ",") + "}", nil
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan type %T into StringArray", value)
	}
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")
	if str == "" {
		*s = StringArray{}
		return nil
	}
	*s = strings.Split(str, ",")
	return nil
}

type Room struct {
	ID               uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BlockID          uuid.UUID   `gorm:"type:uuid;not null" json:"block_id"`
	Block            *Block      `gorm:"foreignKey:BlockID" json:"block,omitempty"`
	RoomName         string      `gorm:"type:varchar(100);not null" json:"room_name"`
	Capacity         int         `gorm:"not null" json:"capacity"`
	IsAvailable      bool        `gorm:"default:true" json:"is_available"`
	AllowedPurposes  StringArray `gorm:"type:text[]" json:"allowed_purposes"`
	Notes            string      `gorm:"type:text" json:"notes"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.AllowedPurposes == nil {
		r.AllowedPurposes = StringArray{"OA", "Interview", "PPT"}
	}
	return nil
}
