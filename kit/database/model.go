package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}
