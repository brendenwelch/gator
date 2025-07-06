-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	)
	RETURNING *;

-- name: GetPostsByUser :many
SELECT posts.* FROM posts
	JOIN feeds ON posts.feed_id = feeds.id
	JOIN users ON feeds.user_id = users.id
	WHERE users.name = $1
	ORDER BY posts.published_at ASC
	LIMIT $2;
