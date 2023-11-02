package orm

import "github.com/google/uuid"

type ArenaRecord struct {
	ID            uuid.UUID    `gorm:"primaryKey;not null;column:id;type:uuid;default:uuid_generate_v4()"`
	LeaderboardID uuid.UUID    `gorm:"not null;column:leaderboard_id;type:uuid"`
	Leaderboard   *Leaderboard `gorm:"foreignKey:LeaderboardID"`
	CharacterID   int64        `gorm:"not null;column:character_id"`
	Character     *Character   `gorm:"foreignKey:CharacterID"`
	Rank          int          `gorm:"not null;column:rank"`
	Rating        int          `gorm:"not null;column:rating"`
	Won           int          `gorm:"not null;column:won"`
	Lost          int          `gorm:"not null;column:lost"`
}

func (ar *ArenaRecord) TableName() string {
	return "arena_records"
}
