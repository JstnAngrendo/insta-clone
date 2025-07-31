package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FollowerServiceClient struct {
	baseURL string
}

func NewFollowerService(baseURL string) *FollowerServiceClient {
	return &FollowerServiceClient{baseURL: baseURL}
}

func (c *FollowerServiceClient) GetFollowers(ctx context.Context, userID int64) ([]int64, error) {
	url := fmt.Sprintf("%s/users/%d/following-ids", c.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		FollowingIDs []int64 `json:"following_ids"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.FollowingIDs, nil
}
