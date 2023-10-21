package orm

import (
	"time"

	"github.com/google/uuid"
)

type Leaderboard struct {
	ID        uuid.UUID `gorm:"primaryKey;not null;column:id;type:uuid;default:uuid_generate_v4()"`
	SeasonID  int       `gorm:"not null;column:season_id"`
	Season    Season    `gorm:"foreignKey:SeasonID"`
	CreatedAt time.Time
	Bracket   string `gorm:"not null;column:bracket"`
}

func (l *Leaderboard) TableName() string {
	return "leaderboards"
}
