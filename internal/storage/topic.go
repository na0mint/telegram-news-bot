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
	selectAll = "SELECT * FROM topics"
)

type TopicPostgresStorage struct {
	db *sqlx.DB
}

func NewTopicStorage(db *sqlx.DB) *TopicPostgresStorage {
	return &TopicPostgresStorage{db: db}
}

func (t *TopicPostgresStorage) Topics(ctx context.Context) ([]model.Topic, error) {
	conn, err := t.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer utils.HandleCloseDbConnection(conn)

	var topics []dbTopic
	if err := conn.SelectContext(ctx, &topics, selectAll); err != nil {
		return nil, err
	}

	return lo.Map(topics, func(topic dbTopic, _ int) model.Topic { return model.Topic(topic) }), nil
}

func (t *TopicPostgresStorage) getConnection(ctx context.Context) (*sqlx.Conn, error) {
	conn, err := t.db.Connx(ctx)
	if err != nil {
		log.Printf("[ERROR] Failed to get connection to database: %v", err)
		return nil, err
	}

	return conn, nil
}

type dbTopic struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}
