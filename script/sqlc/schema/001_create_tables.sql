-- Tabel users dengan password_hash untuk login lokal
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name          VARCHAR(255) NOT NULL,
    picture       TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted    BOOLEAN NOT NULL DEFAULT FALSE
);

-- Tabel kategori
CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(10) NOT NULL CHECK (type IN ('income','expense')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name, type)
);

-- Tabel transaksi
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

CREATE TABLE blacklisted_tokens (
                                    jti  TEXT PRIMARY KEY,
                                    expires_at TIMESTAMPTZ NOT NULL,
                                    created_at TIMESTAMPTZ DEFAULT NOW()
);