-- name: CreateFeedFollow :one
WITH feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
	VALUES (
	    $1,
	    $2,
	    $3,
	    $4,
	    $5
	)
	RETURNING *
) SELECT
    feed_follow.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follow
    INNER JOIN users ON feed_follow.user_id = users.id
    INNER JOIN feeds ON feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows WHERE user_id = $1;
