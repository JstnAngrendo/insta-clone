package entities

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	FollowerID uint `gorm:"not null;index;uniqueIndex:idx_follower_followee"`
	FolloweeID uint `gorm:"not null;index;uniqueIndex:idx_follower_followee"`
}
