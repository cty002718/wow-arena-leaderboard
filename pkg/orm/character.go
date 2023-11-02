package orm

type Character struct {
	ID       int64   `gorm:"primary_key;column:id"`
	ServerID int     `gorm:"not null;column:server_id"`
	Server   *Server `gorm:"foreignkey:ServerID"`
	Name     string  `gorm:"not null;column:name"`
	Faction  string  `gorm:"not null;column:faction"`
}

func (c *Character) TableName() string {
	return "characters"
}
