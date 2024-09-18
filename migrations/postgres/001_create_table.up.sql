CREATE TYPE StatusUser AS ENUM('active', 'inactive');
CREATE TYPE BuyOrSell AS ENUM('buy', 'sell');
CREATE TYPE MessageStatus AS ENUM('user', 'admin');
CREATE TYPE MessageReadStatus AS ENUM('false', 'true');
CREATE TYPE TransactionStatus AS ENUM('pending', 'success', 'error');

CREATE TABLE IF NOT EXISTS "coins"(
    "id" UUID NOT NULL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "coin_icon" VARCHAR NOT NULL,
    "coin_buy_price" VARCHAR NOT NULL,
    "coin_sell_price" VARCHAR NOT NULL,
    "address" VARCHAR,
    "card_number" VARCHAR,
    "status" BOOLEAN DEFAULT false,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "half_coins_price"(
    "coin_id" UUID NOT NULL REFERENCES "coins"("id"),
    "halfCoinAmount" VARCHAR NOT NULL,
    "halfCoinPrice" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS "users"(
    "id" UUID NOT NULL PRIMARY KEY,
    "first_name" VARCHAR NOT NULL,
    "last_name" VARCHAR,
    "username" VARCHAR,
    "status" StatusUser DEFAULT 'active',
    "telegram_id" VARCHAR UNIQUE NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "admin"(
    "id" UUID NOT NULL PRIMARY KEY,
    "login" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "admin_address"(
    "admin_id" UUID NOT NULL REFERENCES "admin"("id"),
    "coin_id" UUID NOT NULL REFERENCES "coins"("id"),
    "address" VARCHAR NOT NULL,
    "card_number" VARCHAR,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "user_transaction"(
    "id" UUID NOT NULL PRIMARY KEY,
    "coin_id" UUID NOT NULL REFERENCES "coins"("id"),
    "user_id" UUID NOT NULL REFERENCES "users"("id"),
    "user_confirmation_img" VARCHAR NOT NULL,
    "coin_price" VARCHAR NOT NULL,
    "coin_amount" VARCHAR NOT NULL,
    "all_price" VARCHAR NOT NULL,
    "status" BuyOrSell NOT NULL,
    "user_address" VARCHAR,
    "card_name" VARCHAR,
    "payment_card" VARCHAR,
    "message" TEXT,
    "transaction_status" TransactionStatus DEFAULT 'pending',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);


-- ALTER TABLE user_transaction
-- ADD transaction_status TransactionStatus DEFAULT 'pending';
CREATE TABLE IF NOT EXISTS "messages"(
    "id" UUID NOT NULL PRIMARY KEY,
    "status" MessageStatus NOT NULL,
    "message" TEXT NOT NULL,
    "read" MessageReadStatus NOT NULL,
    "admin_id" UUID NOT NULL REFERENCES "admin"("id"),
    "user_id" UUID NOT NULL REFERENCES "users"("id"), 
    "file" VARCHAR NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "pay_message"(
    "id" UUID NOT NULL PRIMARY KEY,
    "message" TEXT NOT NULL,
    "file" VARCHAR NOT NULL,
    "user_transaction_id" UUID,
    "premium_transaction_id" UUID,
    "user_id" UUID NOT NULL REFERENCES "users"("id"), 
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "premium"(
    "id" UUID NOT NULL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "card_number" VARCHAR NOT NULL,
    "img" VARCHAR NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "premium_price_month"(
    "id" UUID NOT NULL PRIMARY KEY,
    "month" VARCHAR NOT NULL,
    "price" VARCHAR NOT NULL,
    "premium_id" UUID NOT NULL REFERENCES "premium"("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "premium_transaction"(
    "id" UUID NOT NULL PRIMARY KEY,
    "phone_number" VARCHAR NOT NULL,
    "telegram_username" VARCHAR NOT NULL,
    "premium_id" UUID NOT NULL REFERENCES "premium"("id"),
    "user_id" UUID NOT NULL REFERENCES "users"("id"),
    "payment_img" VARCHAR NOT NULL,
    "transaction_status" TransactionStatus DEFAULT 'pending',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "coin_nft"(
    "id" UUID NOT NULL PRIMARY KEY,
    "nft_img" VARCHAR NOT NULL,
    "nft_price" VARCHAR NOT NULL,
    "nft_address" VARCHAR NOT NULL,
    "nft_name" VARCHAR NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "nft"(
    "id" UUID NOT NULL PRIMARY KEY,
    "nft_img" VARCHAR NOT NULL,
    "comment" VARCHAR NOT NULL,
    "user_id" UUID NOT NULL REFERENCES "users"("id"),
    "coin_nft_id" UUID NOT NULL REFERENCES "coin_nft"("id"),
    "telegram_id" VARCHAR NOT NULL,
    "card_number" VARCHAR NOT NULL,
    "card_number_name" VARCHAR NOT NULL,
    "status" TransactionStatus DEFAULT 'pending',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);



-- INSERT INTO "admin"("id", "login", "password") VALUES('dbecf401-64b3-4b9b-829a-c8b061431286', 'Sayusupov1972', 'sayusupov1972');
-- INSERT INTO "super_admin"("id","login","password") VALUES('690d15b1-b3bf-416f-83e1-02b183ccb2f2', 'azam1222', '938791222');
-- INSERT INTO "admin_address"("admin_id", "coin_id", "address") VALUES('dbecf401-64b3-4b9b-829a-c8b061431286', 'ecd98c25-4cd3-41f7-8526-5efe021533f7', 'addres$$TON');
-- [
--       {"HalfCoinAmount": "0.5", "HalfCoinPrice": "650000"},
--       {"HalfCoinAmount": "0.8", "HalfCoinPrice": "80000"}
-- ]





