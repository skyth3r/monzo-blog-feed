package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gorilla/feeds"
)

const (
	root   = "https://monzo.com/blog"
	usRoot = "https://monzo.com/us/blog"
)

type blogItem struct {
	PubDate     time.Time
	Description string
	Tags        []string
	Title       string
	Link        string
}

func main() {
	var wg sync.WaitGroup
	errChan := make(chan error, 10)

	wg.Add(1)
	go processFeed(root, "blog", []string{"Technology"}, &wg, errChan)
	wg.Add(1)
	go processFeed(usRoot, "us_blog", nil, &wg, errChan)

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
	blogItems := &[]blogItem{}
	blogTags := map[string]bool{}
	var pageWg sync.WaitGroup
	var tagsMutex sync.Mutex

	lastPage, err := getLastPage(baseURL)
	if err != nil {
		errChan <- fmt.Errorf("error getting last page: %v", err)
		return
	}

	for i := 1; i <= lastPage; i++ {
		url := fmt.Sprintf("%s/page/%d", baseURL, i)
		pageWg.Add(1)
		go vistPage(url, blogItems, blogTags, &pageWg, &tagsMutex)
	}
	pageWg.Wait()

	err = generateFeeds(blogItems, feedName)
	if err != nil {
		errChan <- fmt.Errorf("error generating feeds: %v", err)
		return
	}

	for _, tag := range tags {
		subBlogItems := &[]blogItem{}
		for _, item := range *blogItems {
			for _, itemTag := range item.Tags {
				if strings.Contains(itemTag, tag) {
					*subBlogItems = append(*subBlogItems, item)
					break
				}
			}
		}

		err = generateFeeds(subBlogItems, fmt.Sprintf("%s_%s", feedName, strings.ToLower(tag)))
		if err != nil {
			errChan <- fmt.Errorf("error generating sub feeds: %v", err)
			return
		}
	}
}

func getLastPage(url string) (int, error) {
	last := 1

	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.AllowedDomains("monzo.com"),
	)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		lastPageButton := e.DOM.Find("a[class*='Pagination_LastPageLinkDesktop']")
		lastPageLink := lastPageButton.AttrOr("href", "")
		if len(lastPageLink) > 0 {
			parts := strings.Split(lastPageLink, "/")
			if len(parts) > 0 {
				lastPart := parts[len(parts)-1]
				if lastPartInt, err := strconv.Atoi(lastPart); err == nil {
					last = lastPartInt
				}
			}
		}
	})

	err := c.Visit(url)
	if err != nil {
		return 0, fmt.Errorf("failed to start scraping: %v", err)
	}

	c.Wait()

	return last, nil
}

func vistPage(url string, items *[]blogItem, tags map[string]bool, wg *sync.WaitGroup, mu *sync.Mutex) {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.AllowedDomains("monzo.com"),
	)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		blogItems := e.DOM.Find("div[class*='CardList_CardList']")
		blogItems.Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(j int, a *goquery.Selection) {
				div := a.Find("div[class*='Card_card']")

				title := formatText(div.Find("div[class*='Card_titleContainer']").Text())

				tagsDivContainer := div.Find("div[class*='Card_tagContainer']")
				tagsList := tagsDivContainer.Find("div[class*='TagList_tagList']")
				tagsListWrapper := tagsList.Find("div[class*='TagList_tagWrapper']")
				var tag []string
				tagsListWrapper.Each(func(k int, t *goquery.Selection) {
					tagText := t.Text()
					tag = append(tag, tagText)

					mu.Lock()
					_, ok := tags[tagText]
					if !ok {
						tags[tagText] = true
					}
					mu.Unlock()
				})

				description := formatText(div.Find("div[class*='Card_descriptionContainer']").Text())

				date := div.Find("div[class*='Card_dateContainer']").Text()
				date = strings.TrimSpace(date)
				parsedDate, err := time.Parse("2 January 2006", date)
				if err != nil {
					fmt.Println("Error parsing date:", err)
					return
				}

				href, _ := a.Attr("href")
				url := fmt.Sprintf("https://monzo.com%s", href)

				item := blogItem{
					PubDate:     parsedDate,
					Description: description,
					Tags:        tag,
					Title:       title,
					Link:        url,
				}
				*items = append(*items, item)
			})
		})
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error visiting page:", err)
		return
	}

	c.Wait()
	wg.Done()
}

func generateFeeds(blogItems *[]blogItem, feedName string) error {
	sort.Slice(*blogItems, func(i, j int) bool {
		return (*blogItems)[i].PubDate.After((*blogItems)[j].PubDate)
	})

	err := generateItemFeed(blogItems, fmt.Sprintf("%s_feed_items", feedName))
	if err != nil {
		return err
	}

	feed, err := convertToFeed(blogItems, feedName)
	if err != nil {
		return err
	}

	err = generateRssFeed(feed, feedName)
	if err != nil {
		return err
	}
	err = generateJsonFeed(feed, feedName)
	if err != nil {
		return err
	}

	return nil
}

func generateItemFeed(items *[]blogItem, name string) error {
	jsonFile, err := os.Create(fmt.Sprintf("%s.json", name))
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(items)
	if err != nil {
		return err
	}

	moveFile(fmt.Sprintf("%s.json", name), fmt.Sprintf(".json/%s.json", name))

	return nil
}

func convertToFeed(items *[]blogItem, name string) (*feeds.Feed, error) {
	title := "Monzo"
	tagName := ""

	parts := strings.Split(name, "_")
	if len(parts) > 1 {
		tagName = parts[len(parts)-1]
		tagName = strings.ToUpper(tagName[:1]) + tagName[1:]
	}

	if tagName != "" {
		title = fmt.Sprintf("Monzo - %s", tagName)
	}

	feed := &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: root},
		Description: "An unofficial Monzo blog feed generated by Akash Goswami (https://akashgoswami.dev)",
		Created:     time.Now(),
	}

	if strings.Contains((*items)[0].Link, "https://monzo.com/us/blog") {
		feed.Title = "Monzo US"
		feed.Link = &feeds.Link{Href: usRoot}
		feed.Description = "An unofficial Monzo US blog feed generated by Akash Goswami (https://akashgoswami.dev)"
	}

	for _, item := range *items {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: item.Link},
			Description: item.Description,
			Created:     item.PubDate,
		})
	}

	sort.Slice(feed.Items, func(i, j int) bool {
		return feed.Items[i].Created.After(feed.Items[j].Created)
	})

	return feed, nil
}

func generateRssFeed(feed *feeds.Feed, name string) error {
	rss, err := feed.ToRss()
	if err != nil {
		return err
	}

	rssFile, err := os.Create(fmt.Sprintf("%s.rss", name))
	if err != nil {
		return err
	}
	defer rssFile.Close()

	_, err = rssFile.WriteString(rss)
	if err != nil {
		return err
	}

	moveFile(fmt.Sprintf("%s.rss", name), fmt.Sprintf("feeds/%s.rss", name))

	return nil
}

func generateJsonFeed(feed *feeds.Feed, name string) error {
	jsonFeed, err := feed.ToJSON()
	if err != nil {
		return err
	}

	jsonFile, err := os.Create(fmt.Sprintf("%s.json", name))
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	_, err = jsonFile.WriteString(jsonFeed)
	if err != nil {
		return err
	}

	moveFile(fmt.Sprintf("%s.json", name), fmt.Sprintf("feeds/%s.json", name))

	return nil
}

func formatText(text string) string {
	text = strings.NewReplacer(
		"’", "'",
	).Replace(text)

	text = strings.NewReplacer(
		"–", "-",
	).Replace(text)
	return text
}

func moveFile(fileName, filePath string) {
	if err := os.Rename(fileName, filePath); err != nil {
		log.Fatalf("unable to move %s to '%s'. Error: %v", fileName, filePath, err)
	}
}
