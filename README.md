# Monelog
Monelog is a personal finance tracking app designed to give you complete control over your cash flow. Focusing on simplicity and transparency, Monelog helps you record every income and expense instantly, without confusing features.

#Database Structure
-- Tabel user (dari OAuth)
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    name          VARCHAR(255),
    picture       TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabel kategori (per user)
CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(10) NOT NULL CHECK (type IN ('income','expense')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name, type)
);

-- Tabel transaksi (pemasukan/pengeluaran)
CREATE TABLE transactions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id),
    amount      BIGINT NOT NULL, -- dalam satuan rupiah (integer)
    note        TEXT,
    date        DATE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
