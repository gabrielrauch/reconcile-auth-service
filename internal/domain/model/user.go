package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName string    `gorm:"not null" json:"first_name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"default:'user'" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
