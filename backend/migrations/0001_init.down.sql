-- Down migration for 0001_init. Drops in reverse dependency order.
-- Destructive: intended only for local development resets.

DROP TABLE IF EXISTS stock_transactions;
DROP FUNCTION IF EXISTS prevent_stock_tx_mutation();
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
