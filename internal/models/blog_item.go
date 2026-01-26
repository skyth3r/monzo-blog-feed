package models

import "time"

type BlogItem struct {
	PubDate     time.Time
	Description string
	Tags        []string
	Title       string
	Link        string
}
