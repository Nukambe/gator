-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllFeeds :many
select f.name as feed_name, f.url as feed_url, u.name as user_name
from feeds f
join users u on f.user_id = u.id;

-- name: GetFeedByUrl :one
select * from feeds
where url = $1;