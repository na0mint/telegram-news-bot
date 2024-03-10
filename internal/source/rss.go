package source

import (
	"context"
	"tg-bot/internal/model"

	"github.com/SlyMarbo/rss"
)

type RSSSource struct {
	URL        string
	SourceID   int64
	SourceName string
}

func (s RSSSource) ID() int64 {
	return s.SourceID
}

func (s RSSSource) Name() string {
	return s.SourceName
}

func NewRSSSource(m model.Source) RSSSource {
	return RSSSource{
		URL:        m.FeedURL,
		SourceID:   m.ID,
		SourceName: m.Name,
	}
}

func (s RSSSource) Fetch(ctx context.Context) ([]model.RSSArticle, error) {
	feed, err := s.loadFeed(ctx, s.URL)
	if err != nil {
		return nil, err
	}

	var result []model.RSSArticle

	for _, item := range feed.Items {
		result = append(result, model.RSSArticle{
			Title:      item.Title,
			Categories: item.Categories,
			Link:       item.Link,
			Date:       item.Date,
			Summary:    item.Summary,
			SourceName: s.SourceName,
		})
	}

	return result, nil

}

func (s RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedCh  = make(chan *rss.Feed)
		errorCh = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(url)
		if err != nil {
			errorCh <- err
			return
		}

		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errorCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil
	}
}
