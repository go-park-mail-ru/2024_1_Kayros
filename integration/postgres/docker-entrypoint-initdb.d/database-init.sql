-- Нужно будет найти самый короткий адрес в Москве
-- Узнать, как корректно оценить минимальную длину описания, почты
-- Спросить преподавателя насчет main_db
CREATE TABLE IF NOT EXISTS "user"
(
    id          INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        TEXT
        CONSTRAINT user_name_length CHECK (LENGTH(name) BETWEEN 2 AND 20)                              NOT NULL,
    email       TEXT
        CONSTRAINT user_email_domain CHECK (LENGTH(email) BETWEEN 6 AND 50) UNIQUE                     NOT NULL,
    phone       TEXT
        CONSTRAINT user_phone_domain CHECK (LENGTH(phone) = 18 OR phone IS NULL)                       NULL, -- формат |+7 (989) 232 12 12|
    password    BYTEA
        CONSTRAINT user_password_length CHECK (LENGTH(password) = 64),
    address     TEXT
        CONSTRAINT user_address_length CHECK ((LENGTH(address) BETWEEN 14 AND 100) OR address IS NULL) NULL, -- |ул. Мира, д. 4| (самое короткое название улицы в Москве 4 символа)
    img_url     TEXT
        CONSTRAINT user_img_url CHECK (LENGTH(img_url) <= 150) DEFAULT '/minio-api/users/default.jpg'   NOT NULL,
    card_number BYTEA
        CONSTRAINT user_card_number CHECK (LENGTH(password) = 64)                                      NULL,
    created_at  TIMESTAMPTZ
        CONSTRAINT user_time_create NOT NULL,
    updated_at  TIMESTAMPTZ
        CONSTRAINT user_time_last_updated NOT NULL
);

CREATE TABLE IF NOT EXISTS unauth_address
(
    unauth_id TEXT
        CONSTRAINT unauth_id_domain CHECK (length(unauth_id) = 36)                  NOT NULL,
    address   TEXT
        CONSTRAINT unauth_address_domain CHECK (LENGTH(address) BETWEEN 14 AND 100) NULL
);

CREATE TABLE IF NOT EXISTS category
(
    id   INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT
        CONSTRAINT category_name_length CHECK (LENGTH(name) BETWEEN 2 AND 30) UNIQUE NOT NULL,
    type TEXT
);

CREATE TABLE IF NOT EXISTS restaurant
(
    id                INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name              TEXT
        CONSTRAINT rest_name_length CHECK (LENGTH(name) BETWEEN 2 AND 30)                                        NOT NULL,
    short_description TEXT
        CONSTRAINT rest_short_description CHECK (LENGTH(short_description) BETWEEN 10 AND 100)                   NULL,
    long_description  TEXT
        CONSTRAINT rest_long_description CHECK (LENGTH(long_description) BETWEEN 20 AND 250)                     NULL,
    address           TEXT
        CONSTRAINT rest_address CHECK (LENGTH(address) BETWEEN 14 AND 100) UNIQUE                                NOT NULL,
    img_url           TEXT
        CONSTRAINT restaurant_img_url CHECK (LENGTH(img_url) <= 60) DEFAULT '/minio-api/restaurants/default.jpg' NOT NULL,
    rating            NUMERIC(2, 1)
        CONSTRAINT rest_rating CHECK (rating >= 0)                  DEFAULT 0                                    NOT NULL,
    comment_count     INTEGER
        CONSTRAINT rest_comment_count CHECK (comment_count >= 0)    DEFAULT 0,
    CONSTRAINT rest_unique UNIQUE (name, address)
);

-- Нужно будет удалить из отношение restaurant поле address
-- CREATE TABLE IF NOT EXISTS restaurant_address (
--     id   INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY REFERENCES restaurant(id),
--     address           TEXT CONSTRAINT rest_address CHECK (LENGTH(address) BETWEEN 14 AND 100) UNIQUE NOT NULL
-- );


CREATE TABLE IF NOT EXISTS rest_categories
(
    restaurant_id INTEGER
        CONSTRAINT foreign_key_rest CHECK (restaurant_id > 0) REFERENCES restaurant (id) ON DELETE CASCADE,
    category_id   INTEGER
        CONSTRAINT foreign_key_cat CHECK (category_id > 0) REFERENCES category (id) ON DELETE CASCADE,
    PRIMARY KEY (restaurant_id, category_id)
);

CREATE TABLE IF NOT EXISTS promocode(
    id                INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code              TEXT
        CONSTRAINT code_uinque                                    NOT NULL UNIQUE,
    date TIMESTAMPTZ NOT NULL,
    sale  INTEGER
        CONSTRAINT sale_range CHECK (sale > 0 AND sale < 100)                  NOT   NULL,
    type          TEXT
        CONSTRAINT type_length CHECK (LENGTH(type) > 0 ) NOT NULL,
    restaurant_id           INTEGER
        CONSTRAINT foreign_key_rest CHECK (restaurant_id > 0)      NULL REFERENCES restaurant (id) ON DELETE CASCADE,
    sum            INTEGER
        CONSTRAINT non_negatve_sum CHECK (sum > 0)  NULL
);

CREATE TABLE IF NOT EXISTS "order"
(
    id               INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id          INTEGER
        CONSTRAINT foreign_key CHECK (user_id > 0) REFERENCES "user" (id) ON DELETE CASCADE  NULL,
    sum              INTEGER
        CONSTRAINT positive_sum CHECK (sum >= 0)                                             NULL,
    status           TEXT
        CONSTRAINT order_status_domain NOT NULL DEFAULT 'draft',
    unauth_token     TEXT
        CONSTRAINT unauth_token_domain CHECK (length(unauth_token) = 36)                     NULL,
    address          TEXT
        CONSTRAINT order_address_length CHECK (LENGTH(address) BETWEEN 14 AND 100)           NULL,
    extra_address    TEXT
        CONSTRAINT order_extra_address_length CHECK (LENGTH(extra_address) BETWEEN 2 AND 30) NULL,
    created_at       TIMESTAMPTZ
        CONSTRAINT user_time_create NOT NULL,
    updated_at       TIMESTAMPTZ
        CONSTRAINT user_time_last_updated NOT NULL,
    received_at      TIMESTAMPTZ
        CONSTRAINT user_time_received NULL,
    order_created_at TIMESTAMPTZ
        CONSTRAINT order_time_payed NULL,
    delivered_at     TIMESTAMPTZ
        CONSTRAINT order_time_delivered NULL,
    commented        BOOLEAN                    DEFAULT FALSE,
    promocode_id          INTEGER
        CONSTRAINT foreign_promocode CHECK (promocode_id > 0) REFERENCES promocode (id) ON DELETE CASCADE  NULL
);

-- БЖУ хранятся в МГ
CREATE TABLE IF NOT EXISTS food
(
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name          TEXT
        CONSTRAINT food_name_length CHECK (LENGTH(name) BETWEEN 2 AND 60)                                  NOT NULL,
    description   TEXT
        CONSTRAINT food_description_length CHECK (LENGTH(description) BETWEEN 10 AND 100)                  NULL,
    restaurant_id INTEGER
        CONSTRAINT foreign_key_rest CHECK (restaurant_id > 0) REFERENCES restaurant (id) ON DELETE CASCADE,
    category_id   INTEGER
        CONSTRAINT foreign_key_cat CHECK (category_id > 0) REFERENCES category (id) ON DELETE CASCADE,
    weight        INTEGER
        CONSTRAINT positive_weight CHECK (weight > 0)                                                      NULL,
    price         INTEGER
        CONSTRAINT positive_price CHECK (price > 0)                                                        NOT NULL,
    proteins      INTEGER
        CONSTRAINT non_negative_prot CHECK (proteins >= 0)                                                 NULL,
    fats          INTEGER
        CONSTRAINT non_negative_fats CHECK (fats >= 0)                                                     NULL,
    carbohydrates INTEGER
        CONSTRAINT non_negative_carb CHECK (carbohydrates >= 0)                                            NULL,
    img_url       TEXT
        CONSTRAINT restaurant_img_url CHECK (LENGTH(img_url) <= 60) DEFAULT '/minio-api/foods/default.jpg' NOT NULL
);

CREATE TABLE IF NOT EXISTS food_order
(
    food_id    INTEGER
        CONSTRAINT foreign_key_food CHECK (food_id > 0) REFERENCES food (id) ON DELETE CASCADE,
    order_id   INTEGER
        CONSTRAINT foreign_key_order CHECK (order_id > 0) REFERENCES "order" (id) ON DELETE CASCADE,
    count      INTEGER
        CONSTRAINT food_count_in_order CHECK (count > 0) NOT NULL,
    created_at TIMESTAMPTZ
        CONSTRAINT time_create NOT NULL,
    updated_at TIMESTAMPTZ
        CONSTRAINT last_updated NOT NULL,
    PRIMARY KEY (food_id, order_id)
);


-- Хакатон
CREATE TABLE IF NOT EXISTS question
(
    id         INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       TEXT
        CONSTRAINT question_text CHECK (LENGTH(name) < 150),
    url        TEXT,
    focus_id   TEXT,
    param_type TEXT
);

CREATE TABLE IF NOT EXISTS quiz
(
    id          INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    question_id INTEGER
        CONSTRAINT foreign_key_question CHECK (question_id > 0) REFERENCES question (id) ON DELETE CASCADE,
    user_id     TEXT
        CONSTRAINT foreign_key_user CHECK (LENGTH(user_id) BETWEEN 1 AND 36),
    rating      INTEGER
        CONSTRAINT rating_domain CHECK (rating BETWEEN 0 AND 10),
    created_at  TIMESTAMPTZ
        CONSTRAINT quiz_time_create NOT NULL
);

CREATE TABLE IF NOT EXISTS "comment"
(
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id       INTEGER
        CONSTRAINT foreign_key CHECK (user_id > 0)                 NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    restaurant_id INTEGER
        CONSTRAINT foreign_key_rest CHECK (restaurant_id > 0)      NOT NULL REFERENCES restaurant (id) ON DELETE CASCADE,
    text          TEXT
        CONSTRAINT comment_text_length CHECK (LENGTH(text) <= 250) NULL,
    rating        INTEGER
        CONSTRAINT non_negative_rating CHECK (rating >= 0)         NOT NULL
);