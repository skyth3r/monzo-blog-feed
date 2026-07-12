package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/skyth3r/monzo-blog-feed/internal/config"
	"github.com/skyth3r/monzo-blog-feed/internal/feed"
	"github.com/skyth3r/monzo-blog-feed/internal/models"
	"github.com/skyth3r/monzo-blog-feed/internal/scraper"
)

func main() {
	var wg sync.WaitGroup
	errChan := make(chan error, 10)

	wg.Add(1)
	go processFeed(config.Root, "blog", []string{"Technology"}, &wg, errChan)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func processFeed(baseURL, feedName string, tags []string, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	blogItems := &[]models.BlogItem{}
	blogTags := map[string]bool{}
	var pageWg sync.WaitGroup
	var tagsMutex sync.Mutex

	lastPage, err := scraper.GetLastPage(baseURL)
	if err != nil {
		errChan <- fmt.Errorf("error getting last page: %v", err)
		return
	}

	for i := 1; i <= lastPage; i++ {
		url := fmt.Sprintf("%s/page/%d", baseURL, i)
		pageWg.Add(1)
		go scraper.VisitPage(url, blogItems, blogTags, &pageWg, &tagsMutex)
	}
	pageWg.Wait()

	err = feed.GenerateFeeds(blogItems, feedName)
	if err != nil {
		errChan <- fmt.Errorf("error generating feeds: %v", err)
		return
	}

	for _, tag := range tags {
		subBlogItems := &[]models.BlogItem{}
		for _, item := range *blogItems {
			for _, itemTag := range item.Tags {
				if strings.Contains(itemTag, tag) {
					*subBlogItems = append(*subBlogItems, item)
					break
				}
			}
		}

		err = feed.GenerateFeeds(subBlogItems, fmt.Sprintf("%s_%s", feedName, strings.ToLower(tag)))
		if err != nil {
			errChan <- fmt.Errorf("error generating sub feeds: %v", err)
			return
		}
	}
}
