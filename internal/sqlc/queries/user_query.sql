-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE users.email = @email
LIMIT 1;

-- name: GetUserById :one
SELECT *
FROM users
WHERE users.id = @id;

-- name: CreateUser :one
INSERT INTO users
(full_name, birthday, gender, email, password, avatar_file_name, online)
VALUES
(@full_name, @birthday, @gender::gender, @email, @password, @avatar_file_name, true)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    online = coalesce(sqlc.narg('online'), online),
    full_name = coalesce(sqlc.narg('full_name'), full_name),
    birthday = coalesce(sqlc.narg('birthday'), birthday),
    gender = coalesce(sqlc.narg('gender'), gender),
    email = coalesce(sqlc.narg('email'), email),
    password = coalesce(sqlc.narg('password'), password),
    avatar_file_name = coalesce(sqlc.narg('avatar_file_name'), avatar_file_name),
    email_verified = case
        when sqlc.narg('email') is null then email_verified
        else false
    end,
    updated_at = now()
WHERE users.id = @id
RETURNING *;

-- name: RemoveUser :exec
DELETE FROM users
WHERE users.id = @id;

