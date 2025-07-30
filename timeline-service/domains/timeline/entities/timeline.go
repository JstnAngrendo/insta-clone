package entities

type Timeline struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	PostID uint
}
