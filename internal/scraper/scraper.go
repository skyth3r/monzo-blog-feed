package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

	"github.com/skyth3r/monzo-blog-feed/internal/models"
	"github.com/skyth3r/monzo-blog-feed/internal/utils"
)

func GetLastPage(url string) (int, error) {
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

func VisitPage(url string, items *[]models.BlogItem, tags map[string]bool, wg *sync.WaitGroup, mu *sync.Mutex) {
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.AllowedDomains("monzo.com"),
	)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		blogItems := e.DOM.Find("div[class*='CardList_CardList']")
		blogItems.Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(j int, a *goquery.Selection) {
				div := a.Find("div[class*='Card_card']")

				title := utils.FormatText(div.Find("div[class*='Card_titleContainer']").Text())

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

				description := utils.FormatText(div.Find("div[class*='Card_descriptionContainer']").Text())

				date := div.Find("div[class*='Card_dateContainer']").Text()
				date = strings.TrimSpace(date)
				parsedDate, err := time.Parse("2 January 2006", date)
				if err != nil {
					fmt.Println("Error parsing date:", err)
					return
				}

				href, _ := a.Attr("href")
				url := fmt.Sprintf("https://monzo.com%s", href)

				item := models.BlogItem{
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
