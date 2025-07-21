package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserService interface {
	GetFollowingUserIDs(userID string) ([]string, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (u *userService) GetFollowingUserIDs(userID string) ([]string, error) {
	url := fmt.Sprintf("http://localhost:8080/users/%s/following-ids", userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		IDs []string `json:"ids"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response body:", string(body))

	return result.IDs, nil
}
