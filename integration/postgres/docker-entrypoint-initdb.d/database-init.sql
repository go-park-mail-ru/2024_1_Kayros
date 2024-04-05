-- Нужно будет найти самый короткий адрес в Москве
-- Узнать, как корректно оценить минимальную длину описания, почты
-- Спросить преподавателя насчет
CREATE TABLE IF NOT EXISTS "User"
(
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(20) CHECK (LENGTH(name) > 1)                          NOT NULL, -- Н-р, Ян
    email    VARCHAR(50) CHECK (LENGTH(email) >= 3)                        NOT NULL, --  x@x.xx
    phone    VARCHAR(18) CHECK (LENGTH(phone) = 18)                        NULL,     -- формат |+7 (989) 232 12 12|
    password VARCHAR(30) CHECK (LENGTH(password) > 8)                      NOT NULL,
    address  VARCHAR(100) CHECK (LENGTH(address) >= 14)                     NULL,     -- ул. Мира, д. 4 (самое короткое название улицы 4 символа)
    img_url  VARCHAR(60) DEFAULT 'http://localhost:9000/users/default.jpg' NULL
    );

CREATE TABLE IF NOT EXISTS "Category"
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
    );

CREATE TABLE IF NOT EXISTS "Restaurant"
(
    id                SERIAL PRIMARY KEY,
    name              VARCHAR(30) CHECK (LENGTH(name) > 1)                                NOT NULL,
    short_description VARCHAR(100) CHECK (LENGTH(short_description) > 10)                 NULL,
    long_description  VARCHAR(250) CHECK (LENGTH(long_description) > 20)                  NULL,
    address           VARCHAR(100) CHECK (LENGTH(address) >= 14)                           NOT NULL,
    img_url           VARCHAR(80) DEFAULT 'http://localhost:9000/restaurants/default.jpg' NOT NULL
    );

CREATE TABLE IF NOT EXISTS "Order"
(
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER CHECK (user_id > 0) REFERENCES "User" (id),
    date_order     TIMESTAMP WITHOUT TIME ZONE              NOT NULL,
    date_receiving TIMESTAMP WITHOUT TIME ZONE              NULL,
    status         VARCHAR(20) DEFAULT 'draft'              NOT NULL,
    address        VARCHAR(60) CHECK (LENGTH(address) >= 14) NULL,
    extra_address  VARCHAR(30) CHECK (LENGTH(address) > 1)  NULL,
    sum            INTEGER                                  NULL
    );

CREATE TABLE IF NOT EXISTS "Food"
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(60) CHECK (LENGTH(name) > 1)                         NOT NULL,
    description   VARCHAR(60)                                                  NULL,
    restaurant_id INTEGER CHECK (restaurant_id > 0) REFERENCES "Restaurant" (id),
    category_id   INTEGER CHECK (category_id > 0) REFERENCES "Category" (id),
    weight        INTEGER CHECK (weight > 0)                                   NULL,
    price         INTEGER CHECK (price > 0)                                    NOT NULL,
    proteins      NUMERIC CHECK (proteins >= 0)                                NULL,
    fats          NUMERIC CHECK (fats >= 0)                                    NULL,
    carbohydrates NUMERIC CHECK (carbohydrates >= 0)                           NULL,
    img_url       VARCHAR(80) DEFAULT 'http://localhost:9000/food/default.jpg' NULL
    );

CREATE TABLE IF NOT EXISTS "FoodOrder"
(
    food_id  INTEGER CHECK (food_id > 0) REFERENCES "Food" (id),
    order_id INTEGER CHECK (order_id > 0) REFERENCES "Order" (id),
    count    INTEGER CHECK (count > 0) NOT NULL,
    PRIMARY KEY (food_id, order_id)
);