package orm

import (
	"time"

	"github.com/google/uuid"
)

type TokenPriceLog struct {
	ID              uuid.UUID `gorm:"primaryKey;not null;column:id;type:uuid;default:uuid_generate_v4()"`
	LastUpdatedTime time.Time `gorm:"column:last_updated_time"`
	Price           int       `gorm:"column:price"`
}

func (TokenPriceLog) TableName() string {
	return "token_price_logs"
}
