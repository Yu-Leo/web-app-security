-- name: CreateMLModel :one
INSERT INTO ml_models (name, model_data)
VALUES ($1, $2)
RETURNING id, name, model_data;

-- name: GetMLModel :one
SELECT id, name, model_data
FROM ml_models
WHERE id = $1;

-- name: ListMLModels :many
SELECT id, name, model_data
FROM ml_models
ORDER BY id;

-- name: UpdateMLModel :one
UPDATE ml_models
SET name = $2,
    model_data = $3
WHERE id = $1
RETURNING id, name, model_data;

-- name: DeleteMLModel :one
DELETE FROM ml_models
WHERE id = $1
RETURNING id;
