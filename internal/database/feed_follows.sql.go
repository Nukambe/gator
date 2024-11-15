// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    insert into feed_follows (created_at, updated_at, user_id, feed_id)
    values ($1, $2, $3, $4)
    returning id, created_at, updated_at, user_id, feed_id
)
select inserted_feed_follow.id, inserted_feed_follow.created_at, inserted_feed_follow.updated_at, inserted_feed_follow.user_id, inserted_feed_follow.feed_id, u.name as user_name, f.name as feed_name
from inserted_feed_follow
inner join users u on inserted_feed_follow.user_id = u.id
inner join feeds f on inserted_feed_follow.feed_id = f.id
`

type CreateFeedFollowParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    int32
}

type CreateFeedFollowRow struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    int32
	UserName  string
	FeedName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.UserName,
		&i.FeedName,
	)
	return i, err
}

const deleteFeedFollowByUserIdAndURL = `-- name: DeleteFeedFollowByUserIdAndURL :exec
delete
from feed_follows ff
using feeds f
where ff.user_id = $1 and f.url = $2
`

type DeleteFeedFollowByUserIdAndURLParams struct {
	UserID uuid.UUID
	Url    string
}

func (q *Queries) DeleteFeedFollowByUserIdAndURL(ctx context.Context, arg DeleteFeedFollowByUserIdAndURLParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollowByUserIdAndURL, arg.UserID, arg.Url)
	return err
}

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
select ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id, u.name as user_name, f.name as feed_name
from feed_follows ff
join users u on ff.user_id = u.id
join feeds f on ff.feed_id = f.id
where u.name = $1
`

type GetFeedFollowsForUserRow struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    int32
	UserName  string
	FeedName  string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, name string) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.UserName,
			&i.FeedName,
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
