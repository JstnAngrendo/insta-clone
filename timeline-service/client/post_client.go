package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/models/dto"
)

type PostClient struct {
	baseURL string
}

type Post = dto.PostDTO

func NewPostClient(baseURL string) *PostClient {
	return &PostClient{baseURL: baseURL}
}

func (c *PostClient) BatchGet(ctx context.Context, ids []string) ([]Post, error) {
	url := fmt.Sprintf("%s/posts/batch", c.baseURL)
	body, _ := json.Marshal(map[string][]string{"ids": ids})
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var posts []Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}
	return posts, nil
}
