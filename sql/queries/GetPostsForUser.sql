-- name: GetPostsForUser :many
SELECT
  posts.*
FROM
  feed_follows
INNER JOIN
  posts
ON
  posts.feed_id = feed_follows.feed_id
WHERE
  feed_follows.user_id = $1
ORDER BY
  published_at DESC
LIMIT $2;
