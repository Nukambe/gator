-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    insert into feed_follows (created_at, updated_at, user_id, feed_id)
    values ($1, $2, $3, $4)
    returning *
)
select inserted_feed_follow.*, u.name as user_name, f.name as feed_name
from inserted_feed_follow
inner join users u on inserted_feed_follow.user_id = u.id
inner join feeds f on inserted_feed_follow.feed_id = f.id;

-- name: GetFeedFollowsForUser :many
select ff.*, u.name as user_name, f.name as feed_name
from feed_follows ff
join users u on ff.user_id = u.id
join feeds f on ff.feed_id = f.id
where u.name = $1;

-- name: DeleteFeedFollowByUserIdAndURL :exec
delete
from feed_follows ff
using feeds f
where ff.user_id = $1 and f.url = $2;