package model

import "time"

type RSSArticle struct {
	Title      string
	Categories []string
	Link       string
	Date       time.Time
	Summary    string
	SourceName string
}

type Source struct {
	ID        int64
	Name      string
	FeedURL   string
	TopicID   int64
	CreatedAt time.Time
}

type Article struct {
	ID          int64
	SourceID    int64
	Title       string
	Link        string
	Summary     string
	PublishedAt time.Time
	CreatedAt   time.Time
}

type Topic struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}
