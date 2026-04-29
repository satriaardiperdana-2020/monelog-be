# Monelog
Monelog is a personal finance tracking app designed to give you complete control over your cash flow. Focusing on simplicity and transparency, Monelog helps you record every income and expense instantly, without confusing features.

#Database Structure

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

# Database Structure

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
