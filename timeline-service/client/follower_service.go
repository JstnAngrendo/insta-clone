package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type followerServiceClient struct {
	baseURL string
}

func NewFollowerService(baseURL string) *followerServiceClient {
	return &followerServiceClient{baseURL: baseURL}
}

func (f *followerServiceClient) GetFollowers(ctx context.Context, userID int64) ([]int64, error) {
	url := fmt.Sprintf("%s/followers/%d", f.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var followers []int64
	if err := json.NewDecoder(resp.Body).Decode(&followers); err != nil {
		return nil, err
	}

	return followers, nil
}
