-- name: ListFeeds :many
SELECT
  feeds.name as name,
  feeds.url as url,
  users.name as username
FROM
  feeds
INNER JOIN users
  ON feeds.user_id = users.id;
