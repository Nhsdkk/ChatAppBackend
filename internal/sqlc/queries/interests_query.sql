-- name: GetManyInterestsById :many
SELECT * FROM interests
WHERE interests.id = ANY(@ids::uuid[]);

-- name: AssignInterestsToUser :exec
INSERT INTO user_interests
(user_id, interest_id)
VALUES
(@user_id, unnest(@interest_ids::uuid[]));

-- name: RemoveUserInterest :exec
DELETE FROM user_interests
WHERE
    user_id = @user_id AND
    interest_id = ANY (@interest_ids::uuid[]);

-- name: CreateInterest :one
INSERT INTO interests
(title, icon_file_name)
VALUES
(@title, @icon_file_name)
RETURNING *;

-- name: DeleteInterest :exec
DELETE FROM interests
WHERE interests.id = @id;

-- name: UpdateInterestIcon :one
UPDATE interests
SET
    icon_file_name = @icon_file_name
WHERE
    interests.id = @id
RETURNING *;
