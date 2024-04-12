-- Определение схемы
SET schema 'public';

-- Конфигурация
SET BYTEA_OUTPUT TO 'hex';

CREATE TABLE IF NOT EXISTS "user"
(
    id          INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        TEXT CONSTRAINT name_domain CHECK (LENGTH(name) BETWEEN 2 AND 20) NOT NULL,
    email       TEXT CONSTRAINT email_domain CHECK (LENGTH(email) BETWEEN 6 AND 50) UNIQUE NOT NULL,
    phone       TEXT CONSTRAINT phone_domain CHECK (LENGTH(phone) = 18), -- формат |+7 (989) 232 12 12|
    password    BYTEA CONSTRAINT password_domain CHECK (LENGTH(password) = 64) NOT NULL,
    address     TEXT CONSTRAINT address_domain CHECK (LENGTH(address) BETWEEN 14 AND 100),     -- |ул. Мира, д. 4| (самое короткое название улицы в Москве 4 символа)
    img_url     TEXT CONSTRAINT img_url_domain CHECK(LENGTH(img_url) <= 60) NOT NULL DEFAULT '/minio-api/users/default.jpg',  -- для генерации названия картинки используется uuid.V4 (36 символов в длину)
    card_number BYTEA CONSTRAINT card_number_domain CHECK (LENGTH(password) = 64),
    created_at  TIMESTAMPTZ CONSTRAINT time_create_domain NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ CONSTRAINT time_last_updated_domain NOT NULL DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE IF NOT EXISTS category
(
    id   INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT CONSTRAINT category_domain CHECK (LENGTH(name) BETWEEN 2 AND 30) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS restaurant
(
    id                INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name              TEXT CONSTRAINT name_domain CHECK (LENGTH(name) BETWEEN 2 AND 30) UNIQUE NOT NULL,
    short_description TEXT CONSTRAINT short_description_domain CHECK (LENGTH(short_description) BETWEEN 10 AND 100) NOT NULL,
    long_description  TEXT CONSTRAINT long_description_domain CHECK (LENGTH(long_description) BETWEEN 20 AND 250) NOT NULL,
    img_url           TEXT CONSTRAINT img_url_domain CHECK(LENGTH(img_url) <= 60) NOT NULL DEFAULT '/minio-api/restaurants/default.jpg'
);

CREATE TABLE IF NOT EXISTS restaurant_address (
    restaurant_id     INT  GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    address           TEXT CONSTRAINT rest_address CHECK (LENGTH(address) BETWEEN 14 AND 100) NOT NULL,
    CONSTRAINT rest_unique UNIQUE (restaurant_id, address)
);

CREATE TABLE IF NOT EXISTS "order"
(
    id             INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id        INTEGER CONSTRAINT foreign_key CHECK (user_id > 0) NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    sum            INTEGER CONSTRAINT positive_sum CHECK (sum >= 0) NOT NULL DEFAULT 0,
    status         TEXT CONSTRAINT order_status_domain NOT NULL DEFAULT 'draft',
    address        TEXT CONSTRAINT order_address_length CHECK (LENGTH(address) BETWEEN 14 AND 100),
    extra_address  TEXT CONSTRAINT order_extra_address_length CHECK (LENGTH(address) BETWEEN 2 AND 30),
    created_at     TIMESTAMPTZ CONSTRAINT user_time_create NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ CONSTRAINT user_time_last_updated NOT NULL DEFAULT CURRENT_TIMESTAMP,
    received_at    TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS food(
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_id INTEGER CONSTRAINT foreign_key_rest CHECK (restaurant_id > 0) NOT NULL REFERENCES restaurant (id) ON DELETE CASCADE,
    category_id   INTEGER CONSTRAINT foreign_key_category CHECK (category_id > 0) NOT NULL REFERENCES category (id) ON DELETE CASCADE,
    name          TEXT CONSTRAINT name_domain CHECK (LENGTH(name) BETWEEN 2 AND 60) NOT NULL,
    weight        INTEGER CONSTRAINT positive_weight CHECK (weight > 0) NOT NULL,
    price         INTEGER CONSTRAINT positive_price CHECK (price > 0) NOT NULL,
    proteins      NUMERIC CONSTRAINT non_negative_proteins CHECK (proteins >= 0) NOT NULL,
    fats          NUMERIC CONSTRAINT non_negative_fats CHECK (fats >= 0) NOT NULL,
    carbohydrates NUMERIC CONSTRAINT non_negative_carb CHECK (carbohydrates >= 0) NOT NULL,
    img_url       TEXT CONSTRAINT img_url_domain CHECK(LENGTH(img_url) <= 60) DEFAULT '/minio-api/foods/default.jpg' NOT NULL,
    CONSTRAINT unique_food_in_rests UNIQUE (name, restaurant_id)
);

CREATE TABLE IF NOT EXISTS food_order
(
    food_id     INTEGER CONSTRAINT foreign_key_food CHECK (food_id > 0) NOT NULL REFERENCES food (id) ON DELETE CASCADE,
    order_id    INTEGER CONSTRAINT foreign_key_order CHECK (order_id > 0) NOT NULL REFERENCES "order" (id) ON DELETE CASCADE,
    count       INTEGER CONSTRAINT food_count_in_order CHECK (count > 0) NOT NULL,
    created_at  TIMESTAMPTZ CONSTRAINT time_create NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ CONSTRAINT last_updated NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (food_id, order_id)
);