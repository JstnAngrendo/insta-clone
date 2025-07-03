package repositories

import (
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id uint) (*entities.User, error)
}

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(u *entities.User) error {
	return r.db.Create(u).Error
}

func (r *userRepo) FindByEmail(email string) (*entities.User, error) {
	var u entities.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByID(id uint) (*entities.User, error) {
	var u entities.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
