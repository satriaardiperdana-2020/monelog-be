-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, picture, is_deleted)
VALUES ($1, $2, $3, $4, false)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND is_deleted = false;

-- name: SoftDeleteUser :exec
UPDATE users
SET is_deleted = true
WHERE id = $1  AND is_deleted = false;

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

-- name: GetMainPageSummaryTransactions :many
SELECT
    t.date,
    COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount END), 0)::BIGINT AS total_income,
    COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount END), 0)::BIGINT AS total_expense,
    COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount ELSE -t.amount END), 0)::BIGINT AS balance
FROM transactions t
         JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1  AND t.is_deleted = false
GROUP BY t.date
ORDER BY t.date DESC
    LIMIT 10;

-- name: GetTransactionDetailsByDateLimit10 :many
SELECT
    t.date,
    c.name AS category_name,
    c.type AS category_type,
    t.note,
    COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount END), 0)::BIGINT AS total_income,
        COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount END), 0)::BIGINT AS total_expense,
        COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount ELSE -t.amount END), 0)::BIGINT AS balance
FROM transactions t
         JOIN categories c ON c.id = t.category_id
WHERE t.user_id = $1
  AND t.is_deleted = false
GROUP BY t.date, c.name, c.type, t.note
ORDER BY t.date DESC
    LIMIT 10;
-- name: GetLast7DaysDetail :many
SELECT
    t.date,
    t.note,
    COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount END), 0)::BIGINT AS total_income,
        COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount END), 0)::BIGINT AS total_expense,
        COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount ELSE -t.amount END), 0)::BIGINT AS balance
FROM transactions t
         JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1
  AND t.is_deleted = false
  AND t.date >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY t.date, t.note
ORDER BY t.date DESC;


-- name: GetLast30DaysDetail :many
SELECT
    t.date,
    t.note,
    COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount END), 0)::BIGINT AS total_income,
        COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount END), 0)::BIGINT AS total_expense,
        COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount ELSE -t.amount END), 0)::BIGINT AS balance
FROM transactions t
         JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1
  AND t.is_deleted = false
  AND t.date >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY t.date, t.note
ORDER BY t.date DESC;

-- name: GetTransactionsBetweenDatesWithNote :many
SELECT
    t.date,
    t.note,
    COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount ELSE 0 END), 0)::BIGINT AS total_income,
        COALESCE(SUM(CASE WHEN c.type = 'expense' THEN t.amount ELSE 0 END), 0)::BIGINT AS total_expense,
        COALESCE(SUM(CASE WHEN c.type = 'income' THEN t.amount ELSE -t.amount END), 0)::BIGINT AS balance
FROM transactions t
         JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1
  AND t.is_deleted = false
  AND t.date BETWEEN $2 AND $3
GROUP BY t.date, t.note
ORDER BY t.date DESC;

-- name: SoftDeleteTransaction :exec
UPDATE transactions
SET is_deleted = TRUE
WHERE id = $1 AND user_id = $2;

-- name: UpdateTransaction :one
UPDATE transactions
SET category_id = $3, amount = $4, note = $5, date = $6
WHERE id = $1 AND user_id = $2
    RETURNING *;

-- name: AddTokenToBlacklist :exec
INSERT INTO blacklisted_tokens (jti, expires_at)
VALUES ($1, $2);

-- name: IsTokenBlacklisted :one
SELECT EXISTS(SELECT 1 FROM blacklisted_tokens WHERE jti = $1) AS blacklisted;