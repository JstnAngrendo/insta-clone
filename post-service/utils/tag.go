package utils

import (
	"strings"

	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
)

func CleanTags(rawTags []string) []entities.Tag {
	tags := make([]entities.Tag, 0, len(rawTags))
	for _, t := range rawTags {
		clean := strings.TrimPrefix(t, "#")
		clean = strings.TrimPrefix(clean, "@")

		tags = append(tags, entities.Tag{Name: clean})
	}
	return tags
}
