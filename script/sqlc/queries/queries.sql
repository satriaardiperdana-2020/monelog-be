-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, picture, is_deleted)
VALUES ($1, $2, $3, $4, false)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND is_deleted = false;

-- name: SoftDeleteUser :exec
UPDATE users
SET is_deleted = true
WHERE id = $1;

-- name: CreateCategory :one
INSERT INTO categories (user_id, name, type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCategoriesByUser :many
SELECT * FROM categories WHERE user_id = $1;

-- name: GetCategoryByIdAndUser :one
SELECT * FROM categories WHERE id = $1 AND user_id = $2;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1 AND user_id = $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = COALESCE(sqlc.narg(name), name),
    type = COALESCE(sqlc.narg(type), type)
WHERE id = $1 AND user_id = $2
    RETURNING *;

-- name: CreateTransaction :one
INSERT INTO transactions (user_id, category_id, amount, note, date, is_deleted)
VALUES ($1, $2, $3, $4, $5, false)
RETURNING *;

-- name: GetTransactionsByUserAndDateRange :many
SELECT * FROM transactions
WHERE user_id = $1 AND is_delete = false AND date BETWEEN $2 AND $3
ORDER BY date DESC;

-- name: SoftDeleteTransaction :exec
UPDATE transactions
SET is_deleted = TRUE
WHERE id = $1 AND user_id = $2;

-- name: GetTotalByTypeAndDateRange :one
SELECT COALESCE(SUM(t.amount), 0)::bigint as total
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1 AND c.type = $2 AND t.date BETWEEN $3 AND $4 AND t.is_deleted = false;;

-- name: GetTopCategoryExpense :many
SELECT c.name, SUM(t.amount)::bigint as total
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1 AND c.type = 'expense' AND t.date BETWEEN $2 AND $3 AND t.is_deleted = false
GROUP BY c.name
ORDER BY total DESC
LIMIT 5;

-- name: UpdateTransaction :one
UPDATE transactions
SET category_id = $3, amount = $4, note = $5, date = $6
WHERE id = $1 AND user_id = $2
    RETURNING *;