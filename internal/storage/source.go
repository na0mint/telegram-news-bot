package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"log"
	"tg-bot/internal/model"
	"tg-bot/internal/utils"
	"time"
)

const (
	selectAllSources string = "SELECT * from sources"
	findSourceById   string = "SELECT * from sources where id = $1"
	saveSource       string = "INSERT INTO sources (name, feed_url, topic_id, type) VALUES ($1, $2, $3, $4) RETURNING id"
	deleteSource     string = "DELETE FROM sources WHERE id = $1"
	sourcesByTopicId string = "SELECT * FROM sources where topic_id = $1"
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

func NewSourceStorage(db *sqlx.DB) *SourcePostgresStorage {
	return &SourcePostgresStorage{db: db}
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.getConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer utils.HandleCloseDbConnection(conn)

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, selectAllSources); err != nil {
		return nil, err
	}

	return lo.Map(sources, func(source dbSource, _ int) model.Source {
		return model.Source(source)
	}), nil
}

func (s *SourcePostgresStorage) SourceById(ctx context.Context, id int64) (*model.Source, error) {
	conn, err := s.getConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer utils.HandleCloseDbConnection(conn)

	var source dbSource
	if err := conn.GetContext(ctx, &source, findSourceById, id); err != nil {
		return nil, err
	}

	return (*model.Source)(&source), nil
}

func (s *SourcePostgresStorage) SourcesByTopicId(ctx context.Context, topicId int64) ([]model.Source, error) {
	conn, err := s.getConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer utils.HandleCloseDbConnection(conn)

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, sourcesByTopicId, topicId); err != nil {
		return nil, err
	}

	return lo.Map(sources, func(source dbSource, _ int) model.Source {
		return model.Source(source)
	}), nil
}

func (s *SourcePostgresStorage) Save(ctx context.Context, source model.Source) (int64, error) {
	conn, err := s.getConnection(ctx)
	if err != nil {
		return 0, err
	}
	defer utils.HandleCloseDbConnection(conn)

	var id int64

	row := conn.QueryRowxContext(ctx, saveSource,
		source.Name, source.FeedURL, source.TopicID, source.Type)

	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	conn, err := s.getConnection(ctx)
	if err != nil {
		return err
	}
	defer utils.HandleCloseDbConnection(conn)

	if _, err = conn.ExecContext(ctx, deleteSource, id); err != nil {
		return err
	}

	return nil
}

func (s *SourcePostgresStorage) getConnection(ctx context.Context) (*sqlx.Conn, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Printf("[ERROR] Failed to get connection to database: %v", err)
		return nil, err
	}

	return conn, nil
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	TopicID   int64     `db:"topic_id"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
}
