// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, name, url, user_id
`

type CreateFeedParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const getAllFeeds = `-- name: GetAllFeeds :many
select f.name as feed_name, f.url as feed_url, u.name as user_name
from feeds f
join users u on f.user_id = u.id
`

type GetAllFeedsRow struct {
	FeedName string
	FeedUrl  string
	UserName string
}

func (q *Queries) GetAllFeeds(ctx context.Context) ([]GetAllFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllFeedsRow
	for rows.Next() {
		var i GetAllFeedsRow
		if err := rows.Scan(&i.FeedName, &i.FeedUrl, &i.UserName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedByUrl = `-- name: GetFeedByUrl :one
select id, created_at, updated_at, name, url, user_id from feeds
where url = $1
`

func (q *Queries) GetFeedByUrl(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUrl, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}
