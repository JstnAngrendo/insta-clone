package usecases

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jstnangrendo/instagram-clone/user-service/config"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/entities"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/repositories"
	"github.com/jstnangrendo/instagram-clone/user-service/utils"
)

type UserUsecase interface {
	Register(u *entities.User) error
	Login(email, pwd string) (string, error)
	GetProfile(userID uint) (*entities.User, error)
}

type userUsecase struct {
	repo repositories.UserRepository
}

func NewUserUsecase(r repositories.UserRepository) UserUsecase {
	return &userUsecase{r}
}

func (u *userUsecase) Register(user *entities.User) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)
	return u.repo.Create(user)
}

func (u *userUsecase) Login(email, pwd string) (string, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd)) != nil {
		return "", errors.New("invalid credentials")
	}

	tid := uuid.NewString()
	exp := time.Now().Add(72 * time.Hour)

	token, err := utils.GenerateJWT(tid, user.ID, exp)
	if err != nil {
		return "", err
	}

	at := &entities.AccessToken{
		ID:        tid,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: exp,
		Revoked:   false,
	}
	if err := config.DB.Create(at).Error; err != nil {
		return "", err
	}

	buf, _ := json.Marshal(at)
	config.RedisClient.Set(config.Ctx, "access_token:"+tid, buf, time.Until(exp))

	return token, nil
}

func (u *userUsecase) GetProfile(userID uint) (*entities.User, error) {
	return u.repo.FindByID(userID)
}
