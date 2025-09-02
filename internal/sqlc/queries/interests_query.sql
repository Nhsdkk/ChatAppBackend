-- name: GetManyInterestsByFilters :many
SELECT id, title, icon_file_name, created_at, updated_at, description
FROM interests
WHERE
    (sqlc.narg('ids')::uuid[] IS NULL OR interests.id = ANY(sqlc.narg('ids')::uuid[]))
  AND
    (sqlc.narg('name')::text IS NULL OR interests.title ILIKE sqlc.narg('name')::text || '%');

-- name: GetInterestById :one
SELECT *
FROM interests
WHERE id = @id;

-- name: ExistenceCheck :one
SELECT COUNT(id)
FROM interests
WHERE id = ANY(@ids::uuid[]);

-- name: GetUserInterests :many
SELECT 
    interests.id,
    interests.title,
    interests.icon_file_name,
    interests.created_at,
    interests.updated_at,
    interests.description
FROM interests
JOIN user_interests on interests.id = user_interests.interest_id
WHERE user_interests.user_id = @id;

-- name: AssignInterestsToUser :exec
INSERT INTO user_interests
(user_id, interest_id)
VALUES
(@user_id, unnest(@interest_ids::uuid[]));

-- name: RemoveUserInterests :exec
DELETE FROM user_interests
WHERE
    user_id = @user_id;

-- name: CreateInterest :one
INSERT INTO interests
(title, icon_file_name, description)
VALUES
(@title, @icon_file_name, @description)
RETURNING *;

-- name: DeleteInterest :exec
DELETE FROM interests
WHERE interests.id = @id;

-- name: UpdateInterest :one
UPDATE interests
SET
    description = sqlc.narg('description')::text,
    icon_file_name = sqlc.narg('icon_file_name')::varchar(255)
WHERE
    id = @id
RETURNING *;
