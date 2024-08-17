CREATE TYPE StatusUser AS ENUM('active', 'inactive');
CREATE TYPE BuyOrSell AS ENUM('buy', 'sell');
CREATE TYPE MessageStatus AS ENUM('user', 'admin');
CREATE TYPE MessageReadStatus AS ENUM('false', 'true');

CREATE TABLE IF NOT EXISTS "super_admin"(
    "id" UUID NOT NULL PRIMARY KEY,
    "login" VARCHAR(255) UNIQUE NOT NULL,
    "password" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "images"(
    "id" UUID NOT NULL PRIMARY KEY,
    "image_link" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "coins"(
    "id" UUID NOT NULL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "coin_icon" UUID NOT NULL REFERENCES "images"("id"),
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
    "auth_date" VARCHAR,
    "hash" TEXT,
    "status" StatusUser DEFAULT 'active',
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
    "user_confirmation_img" UUID NOT NULL REFERENCES "images"("id"),
    "coin_price" VARCHAR NOT NULL,
    "coin_amount" VARCHAR NOT NULL,
    "all_price" VARCHAR NOT NULL,
    "status" BuyOrSell NOT NULL,
    "user_address" VARCHAR,
    "payment_card" VARCHAR,
    "message" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "messages"(
    "id" UUID NOT NULL PRIMARY KEY,
    "status" MessageStatus NOT NULL,
    "message" TEXT NOT NULL,
    "read" MessageReadStatus NOT NULL,
    "admin_id" UUID NOT NULL REFERENCES "admin"("id"),
    "user_id" UUID NOT NULL REFERENCES "users"("id"), 
    "file" UUID  REFERENCES "images"("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

-- INSERT INTO "admin"("id", "login", "password") VALUES('dbecf401-64b3-4b9b-829a-c8b061431286', 'bahodir2809', '123456789');

-- INSERT INTO "admin_address"("admin_id", "coin_id", "address") VALUES('dbecf401-64b3-4b9b-829a-c8b061431286', 'ecd98c25-4cd3-41f7-8526-5efe021533f7', 'addres$$TON');

-- CREATE TABLE IF NOT EXISTS "sell_coin"(
--     "user_id" UUID NOT NULL REFERENCES "users"("id"),
--     "coin_id" UUID NOT NULL REFERENCES "coins"("id"),
--     "address" VARCHAR NOT NULL,
--     "coin_amount" VARCHAR NOT NULL,
--     "number_of_card" VARCHAR NOT NULL,
--     "check_img" TEXT,
--     "price" VARCHAR NOT NULL,
--     "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     "updated_at" TIMESTAMP
-- );