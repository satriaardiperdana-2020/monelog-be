-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, picture)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateCategory :one
INSERT INTO categories (user_id, name, type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCategoriesByUser :many
SELECT * FROM categories WHERE user_id = $1;

-- name: CreateTransaction :one
INSERT INTO transactions (user_id, category_id, amount, note, date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTransactionsByUserAndDateRange :many
SELECT * FROM transactions
WHERE user_id = $1 AND date BETWEEN $2 AND $3
ORDER BY date DESC;

-- name: GetTotalByTypeAndDateRange :one
SELECT COALESCE(SUM(t.amount), 0)::bigint as total
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1 AND c.type = $2 AND t.date BETWEEN $3 AND $4;

-- name: GetTopCategoryExpense :many
SELECT c.name, SUM(t.amount)::bigint as total
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1 AND c.type = 'expense' AND t.date BETWEEN $2 AND $3
GROUP BY c.name
ORDER BY total DESC
LIMIT 5;