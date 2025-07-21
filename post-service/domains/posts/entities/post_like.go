package entities

type PostLike struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	PostID uint
}
