package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/models/dto"
)

type PostClient struct {
	baseURL string
}

type Post = dto.PostDTO

func NewPostClient(baseURL string) *PostClient {
	return &PostClient{baseURL: baseURL}
}

func (c *PostClient) BatchGet(ctx context.Context, strIDs []string) ([]Post, error) {
	var ids []uint
	for _, s := range strIDs {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, uint(n))
	}

	body, err := json.Marshal(map[string][]uint{"ids": ids})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/posts/batch", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

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
