package orm

type Server struct {
	ID   int    `gorm:"primary_key;column:id"`
	Name string `gorm:"not null;column:name"`
	Slug string `gorm:"not null;column:slug"`
	Type string `gorm:"not null;column:type"`
}

func (s *Server) TableName() string {
	return "servers"
}
