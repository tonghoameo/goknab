-- Add down migration script here
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users ;


