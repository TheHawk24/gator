-- name: CreateFeedFollows :one
WITH feed_follows_result AS (INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5) RETURNING *)
SELECT f.id, f.created_at, f.updated_at, f.user_id, f.feed_id, feeds.name AS feed_name, users.name AS user_name 
FROM feed_follows_result f 
INNER JOIN users ON f.user_id = users.id
INNER JOIN feeds ON f.feed_id = feeds.id;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows WHERE user_id = $1;

-- name: GetFeedFollowsForUser :many
SELECT  f.id, f.created_at, f.updated_at, f.user_id, f.feed_id, feeds.name AS feed_name
FROM feed_follows f
INNER JOIN feeds ON feeds.id = feed_id
WHERE f.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;
