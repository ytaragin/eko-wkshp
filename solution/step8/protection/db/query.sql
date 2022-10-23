-- name: AddVpg :one
INSERT INTO vpg (
    vpg_id,
    task_id,
    status
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: UpdateStatus :exec
UPDATE vpg SET status = $1 WHERE vpg_id = $2;

-- name: GetNonReadyVPGs :many
SELECT * FROM vpg WHERE status != 2;