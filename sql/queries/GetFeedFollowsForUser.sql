-- name: GetFeedFollowsForUser :many
SELECT
  feed_follows.*,
  feeds.name as feed_name,
  users.name as user_name
FROM feed_follows
INNER JOIN feeds
  ON feed_follows.feed_id = feeds.id
INNER JOIN users
  ON feed_follows.user_id = users.id
WHERE users.name = $1;
