package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserClient struct {
	baseURL string
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{baseURL: baseURL}
}

func (c *UserClient) GetFollowers(ctx context.Context, userID int64) ([]int64, error) {
	url := fmt.Sprintf("%s/users/%d/followers", c.baseURL, userID)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Followers []int64 `json:"followers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Followers, nil
}
