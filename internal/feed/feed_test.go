package feed_test

import (
	"testing"
	"time"

	"github.com/skyth3r/monzo-blog-feed/internal/feed"
	"github.com/skyth3r/monzo-blog-feed/internal/models"
)

func GenerateFeeds_Test(t *testing.T) {
	mockBlogItems := []models.BlogItem{
		{
			Title:       "Test Title",
			Description: "Test Description",
			Link:        "http://example.com/test",
			PubDate:     time.Now(),
			Tags:        []string{"Test"},
		},
	}

	err := feed.GenerateFeeds(&mockBlogItems, "test_feed")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

}
