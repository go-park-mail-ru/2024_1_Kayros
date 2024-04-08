-- Заполнение БД данными
INSERT INTO "user" (name, email, password, address, created_at, updated_at)
VALUES ('Иван', 'ivan@example.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8',
        'ул. Ленина, д. 1', current_timestamp, current_timestamp),
       ('Анна', 'anna@example.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8',
        'пр. Победы, д. 10', current_timestamp, current_timestamp),
       ('Петр', 'petr@example.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8',
        'ул. Мира, д. 5', current_timestamp, current_timestamp),
       ('Мария', 'maria@example.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8',
        'пр. Ленинградский, д. 15', current_timestamp, current_timestamp),
       ('Алексей', 'alex@example.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8',
        'ул. Советская, д. 25', current_timestamp, current_timestamp),
       ('Елена', 'elena@example.com', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8',
        'пр. Гагарина, д. 30', current_timestamp, current_timestamp);

-- Вставка данных в таблицу "Category"
INSERT INTO category (name)
VALUES ('Завтраки'),
       ('Закуски'),
       ('Салаты'),
       ('Пицца'),
       ('Паста'),
       ('Супы');


-- Вставка данных в таблицу "Restaurant"
INSERT INTO restaurant (name, short_description, long_description, address)
VALUES ('Ресторан "Вкусняшка"', 'Лучший ресторан в городе',
        'Здесь подают самую вкусную еду, приходите и наслаждайтесь!', 'пр. Победы, д. 20'),
       ('Кафе "Уют"', 'Уютное кафе с домашней атмосферой', 'У нас вы найдете самые вкусные десерты и ароматный кофе',
        'ул. Садовая, д. 15'),
       ('Пиццерия "Веселая пицца"', 'Широкий выбор пиццы',
        'Мы предлагаем разнообразные виды пиццы по самым вкусным ценам', 'ул. Московская, д. 5'),
       ('Кафе "Зеленый уголок"', 'Уютное кафе в центре города',
        'У нас вы найдете вкусные домашние блюда и дружелюбный персонал', 'пл. Пушкина, д. 10');

-- Вставка данных в таблицу "Order"
INSERT INTO "order" (user_id, created_at, updated_at, status, address, extra_address, sum)
VALUES (1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'processing', 'ул. Ленина, д. 1', NULL, 1500),
       (2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'draft', 'пр. Победы, д. 10', NULL, 2000),
       (3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'completed', 'ул. Мира, д. 5', NULL, 1000),
       (4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'processing', 'пр. Ленинградский, д. 15', NULL, 1200),
       (5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'draft', 'ул. Советская, д. 25', NULL, 1800),
       (6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'completed', 'пр. Гагарина, д. 30', NULL, 900);

-- Вставка данных в таблицу "Food"
INSERT INTO food (name, description, restaurant_id, category_id, weight, price, proteins, fats, carbohydrates)
VALUES ('Омлет', 'С яйцами и овощами', 1, 1, 200, 300, 15, 10, 5),
       ('Паста карбонара', 'С беконом и сливками', 1, 5, 300, 500, 20, 15, 10),
       ('Салат Цезарь', 'С курицей и сыром', 2, 3, 250, 400, 25, 20, 10),
       ('Пицца Маргарита', 'Соус, сыр, помидоры', 1, 4, 400, 600, 30, 20, 15),
       ('Борщ', 'Классический украинский борщ', 2, 6, 350, 450, 15, 10, 20),
       ('Картофель фри', 'Жареный картофель в виде палочек', 1, 2, 200, 300, 5, 10, 25);

-- Вставка данных в таблицу "FoodOrder"
INSERT INTO food_order (food_id, order_id, count, created_at, updated_at)
VALUES (1, 1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (2, 2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (3, 3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (4, 1, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

