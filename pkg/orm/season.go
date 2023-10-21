package orm

import "time"

type Season struct {
	ID        int       `gorm:"primary_key;column:id"`
	StartedAt time.Time `gorm:"not null;column:started_at"`
	EndedAt   time.Time `gorm:"column:ended_at"`
}

func (s *Season) TableName() string {
	return "seasons"
}
