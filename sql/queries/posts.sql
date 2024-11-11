-- name: CreatePost :one
insert into posts (title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetPostsForUser :many
select *
from posts p
join feed_follows ff on p.feed_id = ff.feed_id
join users u on ff.user_id = u.id
where u.id = $1
order by published_at desc
limit $2;