-- name: CreateUser :one
INSERT INTO users (
    username, password
) VALUES (
    ?, ?
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = ?;

-- name: GetTodos :many
SELECT * FROM todos
WHERE username = ?;

-- name: SearchTodos :many
SELECT * FROM todos
WHERE username = ? AND id IN(
    SELECT rowid FROM todo_search
    WHERE todo_search = '' || sqlc.arg('query')
);

-- name: CreateTodo :one
INSERT INTO todos (
    username, category, title, description, is_done
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateTodo :one
UPDATE todos
SET
 username = sqlc.narg('username'),
 category = coalesce(sqlc.narg('category'), category),
 title = coalesce(sqlc.narg('title'), title),
 description = coalesce(sqlc.narg('description'), description),
 is_done = coalesce(sqlc.narg('is_done'), is_done)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE username = ? AND id = ?;
