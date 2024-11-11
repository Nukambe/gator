// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
insert into posts (title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5)
returning id, created_at, updated_at, title, url, description, published_at, feed_id
`

type CreatePostParams struct {
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      int32
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsForUser = `-- name: GetPostsForUser :many
select p.id, p.created_at, p.updated_at, title, url, description, published_at, p.feed_id, ff.id, ff.created_at, ff.updated_at, user_id, ff.feed_id, u.id, u.created_at, u.updated_at, name
from posts p
join feed_follows ff on p.feed_id = ff.feed_id
join users u on ff.user_id = u.id
where u.id = $1
order by published_at desc
limit $2
`

type GetPostsForUserParams struct {
	ID    uuid.UUID
	Limit int32
}

type GetPostsForUserRow struct {
	ID          int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      int32
	ID_2        int32
	CreatedAt_2 time.Time
	UpdatedAt_2 time.Time
	UserID      uuid.UUID
	FeedID_2    int32
	ID_3        uuid.UUID
	CreatedAt_3 time.Time
	UpdatedAt_3 time.Time
	Name        string
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]GetPostsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsForUserRow
	for rows.Next() {
		var i GetPostsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.UserID,
			&i.FeedID_2,
			&i.ID_3,
			&i.CreatedAt_3,
			&i.UpdatedAt_3,
			&i.Name,
		); err != nil {
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