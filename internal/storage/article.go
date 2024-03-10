package storage

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"tg-bot/internal/model"
	"time"
)

const (
	saveArticle      string = "INSERT INTO articles (source_id, title, link, summary, published_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING"
	findAllNotPosted string = "SELECT * FROM articles where posted_at IS NULL AND published_at >= $1::timestamp ORDER BY published_at DESC LIMIT $2"
	markPosted              = "UPDATE articles SET posted_at = now() WHERE id = $1"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Save(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx,
		saveArticle,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *ArticlePostgresStorage) FindAllNotPosted(ctx context.Context, since time.Time, limit int64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(ctx, &articles, findAllNotPosted, since.UTC().Format(time.RFC3339), limit); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			ID:          article.ID,
			SourceID:    article.SourceID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary,
			PublishedAt: article.PublishedAt,
			CreatedAt:   article.CreatedAt,
		}
	}), nil
}

func (s *ArticlePostgresStorage) MarkPostedById(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, markPosted, id); err != nil {
		return err
	}

	return nil
}

type dbArticle struct {
	ID          int64        `db:"id"`
	SourceID    int64        `db:"source_id"`
	Title       string       `db:"title"`
	Link        string       `db:"link"`
	Summary     string       `db:"summary"`
	PublishedAt time.Time    `db:"published_at"`
	CreatedAt   time.Time    `db:"created_at"`
	PostedAt    sql.NullTime `db:"posted_at"`
}
