-- Migration 0001: initial schema for the Inventory and Stock Management System.
-- Requirements coverage: FR-1..FR-5, AC-04..AC-09, AC-11..AC-14, NFR-09
-- (docs/requirements.md); data model per docs/architecture.md §3.
--
-- Composite foreign keys ((category_id, user_id) → categories, (product_id, user_id)
-- → products) make cross-user references impossible at the database level,
-- enforcing the per-user isolation of AC-02/AC-07 as defense in depth behind
-- the API's ownership checks.

CREATE TABLE users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email         TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Email is unique case-insensitively across the whole installation (Q4).
CREATE UNIQUE INDEX users_email_lower_unique ON users (lower(email));

CREATE TABLE categories (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users (id),
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- Required as the target of the products composite foreign key below.
    UNIQUE (id, user_id)
);

-- Category names are unique within each user's inventory (FR-2.2,
-- case-insensitive + trimmed per Q4).
CREATE UNIQUE INDEX categories_user_name_unique ON categories (user_id, lower(trim(name)));
CREATE INDEX categories_user_id_idx ON categories (user_id);

CREATE TABLE products (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users (id),
    category_id   BIGINT NOT NULL REFERENCES categories (id),
    name          TEXT NOT NULL,
    sku           TEXT NOT NULL,
    unit_price    NUMERIC(10,2) NOT NULL CHECK (unit_price >= 0),        -- AC-04
    opening_stock INTEGER NOT NULL CHECK (opening_stock >= 0),           -- AC-05
    current_stock INTEGER NOT NULL CHECK (current_stock >= 0),           -- AC-05, AC-11 backstop
    min_stock     INTEGER NOT NULL CHECK (min_stock >= 0),               -- AC-05
    status        TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')), -- AC-14
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- The composite foreign key guarantees a product's category belongs to the
    -- same user (AC-07 at the database level).
    FOREIGN KEY (category_id, user_id) REFERENCES categories (id, user_id),
    -- Required as the target of the stock_transactions composite foreign key below.
    UNIQUE (id, user_id)
);

-- SKU is unique within each user's inventory (AC-06, case-insensitive per Q4).
CREATE UNIQUE INDEX products_user_sku_unique ON products (user_id, lower(trim(sku)));
CREATE INDEX products_user_id_idx ON products (user_id);
CREATE INDEX products_category_id_idx ON products (category_id);

-- Serves the dashboard low-stock list (FR-5.2): active products at or below
-- their minimum stock level, per user.
CREATE INDEX products_low_stock_idx ON products (user_id)
    WHERE status = 'active' AND current_stock <= min_stock;

CREATE TABLE stock_transactions (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users (id),
    product_id BIGINT NOT NULL REFERENCES products (id),
    type       TEXT NOT NULL CHECK (type IN ('IN', 'OUT')),               -- FR-4.1
    quantity   INTEGER NOT NULL CHECK (quantity > 0),                     -- AC-09
    remarks    TEXT,                                                      -- optional, FR-4.2
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),                        -- FR-4.2
    -- A movement must reference the same user's product (defense in depth for AC-02).
    FOREIGN KEY (product_id, user_id) REFERENCES products (id, user_id)
);

-- AC-13: existing stock transactions must not be edited or deleted; corrections
-- happen through new transactions. Enforced at the database level as well as by
-- the absence of any update/delete API.
CREATE FUNCTION prevent_stock_tx_mutation() RETURNS trigger AS $$
BEGIN
    RAISE EXCEPTION 'stock_transactions are immutable: record a new transaction instead';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER stock_transactions_no_update
    BEFORE UPDATE ON stock_transactions
    FOR EACH ROW EXECUTE FUNCTION prevent_stock_tx_mutation();
CREATE TRIGGER stock_transactions_no_delete
    BEFORE DELETE ON stock_transactions
    FOR EACH ROW EXECUTE FUNCTION prevent_stock_tx_mutation();

CREATE INDEX stock_transactions_user_time_idx ON stock_transactions (user_id, created_at DESC);
CREATE INDEX stock_transactions_product_idx ON stock_transactions (product_id);
