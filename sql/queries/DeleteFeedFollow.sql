-- name: DeleteFeedFollow :one
DELETE FROM feed_follows
  WHERE feed_id = $1
RETURNING *;
