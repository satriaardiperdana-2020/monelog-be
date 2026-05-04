# Monelog
Monelog is a personal finance tracking app designed to give you complete control over your cash flow. Focusing on simplicity and transparency, Monelog helps you record every income and expense instantly, without confusing features.

Complete Step-by-Step Guide for monelog Backend with Golang, Echo, sqlc, JWT, and PostgreSQL:

📌 Ringkasan Teknologi:
    1. Go 1.21+ – bahasa pemrograman
    2. Echo – web framework
    3. sqlc – type-safe SQL code generator
    4. PostgreSQL – database
    5. JWT – autentikasi (golang-jwt)
    6. oapi-codegen – generate handler dari OpenAPI spec
    7. bcrypt – hashing password
    8. YAML config – konfigurasi aplikasi
    9. Personal coding asistant : deepseek, sometimes claude code.

#Features: user registration/login, category CRUD, transaction CRUD + soft delete, user soft delete, JWT authentication, reports (total summary & top categories).

🗂️ Project Structure
```
monelog/
├── cmd/
│   └── main.go
├── internal/
│   ├── api/                 # generated openapi code
│   ├── handlers/            # business logic implementation
│   ├── middleware/          # JWT auth middleware
│   ├── repository/
│   │   └── postgresql/      # generated sqlc code + connect.go
│   └── config/              # load YAML config
├── script/
│   └── sqlc/
│       ├── schema/
│       │   └── 001_create_tables.sql
│       ├── queries/
│       │   └── queries.sql
│       └── sqlc.yaml
├── api/
│   ├── openapi.yaml
│   └── generate.yaml
├── config-development.yml
├── go.mod
└── go.sum
```

# Database Structure

```sql
-- -- Tabel user (dari OAuth)
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    email       VARCHAR(255) UNIQUE NOT NULL,
    name        VARCHAR(255),
    picture     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabel kategori (per user)
CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(10) NOT NULL CHECK (type IN ('income','expense')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name, type)
);

-- Tabel transaksi (pemasukan/pengeluaran)
CREATE TABLE transactions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    amount      BIGINT NOT NULL, -- dalam satuan rupiah (integer)
    note        TEXT,
    date        DATE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

```


```text
monelog/
├── cmd/
│   └── api/
│       └── main.go               # Entry point
├── internal/
│   ├── api/                      # Handler (gin/echo)
│   │   ├── handler/
│   │   │   ├── auth.go
│   │   │   ├── transaction.go
│   │   │   ├── category.go
│   │   │   └── report.go
│   │   └── middleware/
│   │       └── auth.go           # Verifikasi token OAuth
│   ├── db/                       # sqlc output
│   │   ├── models.go
│   │   ├── queries.sql.go
│   │   └── db.go
│   ├── repository/               # (opsional wrapper, bisa pakai langsung sqlc)
│   ├── service/                  # business logic
│   └── config/
│       └── config.go             # env vars
├── sql/
│   ├── schema/
│   │   └── 001_create_tables.sql
│   ├── queries/
│   │   └── queries.sql           # SQL queries untuk sqlc
│   └── sqlc.yaml                 # konfigurasi sqlc
├── api/
│   └── openapi.yaml              # Spesifikasi OpenAPI 3.0
├── go.mod
├── go.sum
└── .env
```
**🚀 Step 1: Initialize the Project**
bash
```
mkdir monelog && cd monelog
go mod init github.com/username/monelog   # replace username
```
# Create folder structure
```
mkdir -p cmd internal/{handlers,middleware,repository/postgresql,config} script/sqlc/{schema,queries} api
```
**
📝 Step 2: Database Schema (script/sqlc/schema/001_create_tables.sql)**

```
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name          VARCHAR(255) NOT NULL,
    picture       TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted    BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(10) NOT NULL CHECK (type IN ('income','expense')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name, type)
);

CREATE TABLE transactions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    amount      BIGINT NOT NULL,
    note        TEXT,
    date        DATE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_transactions_user_deleted_date ON transactions (user_id, is_deleted, date DESC);
CREATE INDEX idx_users_email_deleted ON users (email, is_deleted);
```

Run migrations:
```
bash

createdb monelog_db  # or create via psql
psql -U postgres -d monelog_db -f script/sqlc/schema/001_create_tables.sql
```

**⚙️ Step 3: sqlc Configuration (script/sqlc/sqlc.yaml)**
```
yaml

version: "2"
sql:
  - engine: "postgresql"
    queries: "script/sqlc/queries/queries.sql"
    schema: "script/sqlc/schema/001_create_tables.sql"
    gen:
      go:
        package: "postgresql"
        out: "../../internal/repository/postgresql"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
        overrides:
          - db_type: "timestamptz"
            go_type: "time.Time"
            nullable: false
```
**📄 Step 4: SQL Queries (script/sqlc/queries/queries.sql)**
```
sql

-- Users
-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, picture, is_deleted)
VALUES ($1, $2, $3, $4, false)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND is_deleted = false;

-- name: SoftDeleteUser :execresult
UPDATE users SET is_deleted = true WHERE id = $1 AND is_deleted = false;

-- Categories
-- name: CreateCategory :one
INSERT INTO categories (user_id, name, type) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCategoriesByUser :many
SELECT * FROM categories WHERE user_id = $1 ORDER BY id;

-- name: GetCategoryByIdAndUser :one
SELECT * FROM categories WHERE id = $1 AND user_id = $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = COALESCE(sqlc.narg(name), name),
    type = COALESCE(sqlc.narg(type), type)
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1 AND user_id = $2;

-- Transactions
-- name: CreateTransaction :one
INSERT INTO transactions (user_id, category_id, amount, note, date, is_deleted)
VALUES ($1, $2, $3, $4, $5, false)
RETURNING *;

-- name: GetTransactionsByUserAndDateRange :many
SELECT * FROM transactions
WHERE user_id = $1 AND date BETWEEN $2 AND $3 AND is_deleted = false
ORDER BY date DESC;

-- name: UpdateTransaction :one
UPDATE transactions
SET category_id = $3, amount = $4, note = $5, date = $6
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: SoftDeleteTransaction :execresult
UPDATE transactions SET is_deleted = true WHERE id = $1 AND user_id = $2;

-- name: GetTotalByTypeAndDateRange :one
SELECT COALESCE(SUM(t.amount), 0)::bigint as total
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1 AND c.type = $2 AND t.date BETWEEN $3 AND $4 AND t.is_deleted = false;

-- name: GetTopCategoryExpense :many
SELECT c.name, SUM(t.amount)::bigint as total
FROM transactions t
JOIN categories c ON t.category_id = c.id
WHERE t.user_id = $1 AND c.type = 'expense' AND t.date BETWEEN $2 AND $3 AND t.is_deleted = false
GROUP BY c.name
ORDER BY total DESC
LIMIT 5;
```

Generate sqlc:
```
bash

cd script/sqlc
sqlc generate
cd ../..
```

**🔌 Step 5: Database Connection (internal/repository/postgresql/connect.go)**
```
package postgresql

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/username/monelog/internal/config"
)

func NewConnection(cfg *config.Config) (*pgxpool.Pool, error) {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
    return pgxpool.New(context.Background(), connStr)
}
```

**⚙️ Step 6: Application Configuration**
```
**config-development.yml (root)**
server:
  port: 8080
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: monelog_db
jwt:
  secret: rahasia12345
```

**internal/config/config.go**
```
package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Server   struct { Port string `yaml:"port"` } `yaml:"server"`
    Database struct {
        Host, Port, User, Password, Name string
    } `yaml:"database"`
    JWT struct { Secret string `yaml:"secret"` } `yaml:"jwt"`
}

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    if cfg.Server.Port == "" {
        cfg.Server.Port = "8080"
    }
    return &cfg, nil
}
```

**🔐 Step 7: JWT Middleware (internal/middleware/auth.go)**
```
package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

type contextKey string
const UserIDKey contextKey = "user_id"

func JWTAuth(secret []byte) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            auth := c.Request().Header.Get("Authorization")
            if auth == "" {
                return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
            }
            parts := strings.Split(auth, " ")
            if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
            }
            tokenString := parts[1]
            claims := jwt.MapClaims{}
            _, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
                return secret, nil
            })
            if err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
            }
            userID := int64(claims["user_id"].(float64))
            c.Set("user_id", userID)
            ctx := context.WithValue(c.Request().Context(), UserIDKey, userID)
            c.SetRequest(c.Request().WithContext(ctx))
            return next(c)
        }
    }
}
```
**🧩 Step 8: Auth Handler (internal/handlers/auth.go)**
```
package handlers

import (
    "net/http"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
    "golang.org/x/crypto/bcrypt"
    "github.com/username/monelog/internal/repository/postgresql"
)

type AuthHandler struct {
    Queries   *postgresql.Queries
    JWTSecret []byte
}

type RegisterRequest struct {
    Email, Password, Name string
}
type LoginRequest struct {
    Email, Password string
}
type AuthResponse struct {
    Token string `json:"token"`
    User  struct {
        ID        int64     `json:"id"`
        Email     string    `json:"email"`
        Name      string    `json:"name"`
        CreatedAt time.Time `json:"created_at"`
    } `json:"user"`
}

func (h *AuthHandler) Register(c echo.Context) error {
    var req RegisterRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
    }
    hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    user, err := h.Queries.CreateUser(c.Request().Context(), postgresql.CreateUserParams{
        Email: req.Email, PasswordHash: string(hashed), Name: req.Name, Picture: "",
    })
    if err != nil {
        return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID, "exp": time.Now().Add(24 * time.Hour).Unix(),
    })
    tokenString, _ := token.SignedString(h.JWTSecret)
    return c.JSON(http.StatusCreated, AuthResponse{
        Token: tokenString,
        User: struct {
            ID int64; Email, Name string; CreatedAt time.Time
        }{user.ID, user.Email, user.Name, user.CreatedAt},
    })
}

func (h *AuthHandler) Login(c echo.Context) error {
    var req LoginRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
    }
    user, err := h.Queries.GetUserByEmail(c.Request().Context(), req.Email)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID, "exp": time.Now().Add(24 * time.Hour).Unix(),
    })
    tokenString, _ := token.SignedString(h.JWTSecret)
    return c.JSON(http.StatusOK, AuthResponse{
        Token: tokenString,
        User: struct {
            ID int64; Email, Name string; CreatedAt time.Time
        }{user.ID, user.Email, user.Name, user.CreatedAt},
    })
}
```
**📋 Step 9: Category Handler (internal/handlers/category.go)**
```
[Full code provided previously – shortened for brevity but includes all necessary methods: CreateCategory, GetCategories, UpdateCategory, DeleteCategory, and conversion helper.]
```

**💰 Step 10: Transaction Handler (internal/handlers/transaction.go)**
```
[Full code provided previously – includes CreateTransaction, GetTransactions, UpdateTransaction, SoftDeleteTransaction, and conversion helper.]
```

**👤 Step 11: User Handler (Soft Delete) (internal/handlers/user.go)**
```
package handlers

import (
    "context"
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/username/monelog/internal/api"
    "github.com/username/monelog/internal/middleware"
    "github.com/username/monelog/internal/repository/postgresql"
)

type UserHandler struct{ Queries *postgresql.Queries }

func (h *UserHandler) SoftDeleteUser(ctx context.Context, req api.SoftDeleteUserRequestObject) (api.SoftDeleteUserResponseObject, error) {
    currentUserID := ctx.Value(middleware.UserIDKey).(int64)
    if req.Id != currentUserID {
        return nil, echo.NewHTTPError(http.StatusForbidden, "You can only delete your own account")
    }
    cmdTag, err := h.Queries.SoftDeleteUser(ctx, req.Id)
    if err != nil {
        return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    if cmdTag.RowsAffected() == 0 {
        return nil, echo.NewHTTPError(http.StatusNotFound, "User not found or already deleted")
    }
    return api.SoftDeleteUser204Response{}, nil
}
```
**🔗 Step 12: Server Aggregator (internal/handlers/server.go)**
```
package: api
output: internal/api/openapi.gen.go
generate:
  echo-server: true
  models: true
  embedded-spec: true
  strict-server: true
```
**Generate:**
```
oapi-codegen -config api/generate.yaml api/openapi.yaml
```

**🚦 Step 14: Entry Point (cmd/main.go)**
```
code in main.go
```
**🔧 Step 15: Install Dependencies & Run**
```
go get github.com/labstack/echo/v4
go get github.com/jackc/pgx/v5
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get gopkg.in/yaml.v3
go get github.com/oapi-codegen/runtime
go mod tidy

# Generate sqlc & openapi
cd script/sqlc && sqlc generate && cd ../..
oapi-codegen -config api/generate.yaml api/openapi.yaml

# Run
go run cmd/main.go
```
Server ready at http://localhost:8080.

**🧪 Testing with curl or Postman******
# Register
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d '{"email":"user@test.com","password":"123456","name":"Test"}'

# Login (get token)
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{"email":"user@test.com","password":"123456"}'

# Create category
curl -X POST http://localhost:8080/api/v1/categories -H "Authorization: Bearer <token>" -H "Content-Type: application/json" -d '{"name":"Food","type":"expense"}'

# Create transaction
curl -X POST http://localhost:8080/api/v1/transactions -H "Authorization: Bearer <token>" -H "Content-Type: application/json" -d '{"category_id":1,"amount":50000,"date":"2026-05-04","note":"lunch"}'

# List transactions
curl "http://localhost:8080/api/v1/transactions?from=2026-05-01&to=2026-05-31" -H "Authorization: Bearer <token>"

# Soft delete transaction
curl -X DELETE http://localhost:8080/api/v1/transactions/1 -H "Authorization: Bearer <token>"

# Soft delete own user
curl -X PATCH http://localhost:8080/api/v1/users/1 -H "Authorization: Bearer <token>"










