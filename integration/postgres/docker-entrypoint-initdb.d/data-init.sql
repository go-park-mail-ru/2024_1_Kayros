-- Заполнение БД данными
INSERT INTO "user" (name, email, password, address, created_at, updated_at)
VALUES ('Иван', 'ivan@example.com', E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'ул. Ленина, д. 1', current_timestamp, current_timestamp),
       ('Анна', 'anna@example.com', E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'пр. Победы, д. 10', current_timestamp, current_timestamp),
       ('Петр', 'petr@example.com', E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'ул. Мира, д. 5', current_timestamp, current_timestamp),
       ('Мария', 'maria@example.com', E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'пр. Ленинградский, д. 15', current_timestamp, current_timestamp),
       ('Алексей', 'alex@example.com', E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'ул. Советская, д. 25', current_timestamp, current_timestamp),
       ('Елена', 'elena@example.com', E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
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
INSERT INTO restaurant (name, short_description, long_description, address, img_url)
VALUES ('Горыныч"', 'Лучший ресторан в городе',
        'Здесь подают самую вкусную еду, приходите и наслаждайтесь!', 'пр. Победы, д. 20', '/minio-api/restaurants/1.jpg'),
       ('Loona"', 'Ресторан с итальянской атмосферой', 'У нас вы найдете самую вкусную пиццу, пасту, закуски и  вкусные десерты',
        'ул. Садовая, д. 15', '/minio-api/restaurants/2.jpg'),
       ('Rustic Kitchen', 'Уютное семейное бистро с домашней кухней и гостеприимной атмосферой',
        'место, где вы чувствуете себя как дома. У нас готовят по-настоящему вкусную еду, используя только свежие продукты высокого качества', 'ул. Невская, д. 1', '/minio-api/restaurants/3.jpg'),
       ('Густо', 'Ресторан с авторской европейской кухней',
        'Место, где каждый член семьи найдет что-то по душе. У нас вы попробуете неповторимые блюда европейской кухни, приготовленные с любовью к деталям', 'ул. Садовая, д. 5', '/minio-api/restaurants/4.jpg'),
       ('Remy', 'Достижения французского кухни',
        'Шеф-повар создает удивительные шедевры из самых редких и изысканных ингредиентов, чтобы удовлетворить самые изысканные вкусы гостей. Ощутите волшебство атмосферы и насладитесь бесподобными блюдами, приготовленными с любовью и талантом', 'ул. Московская, д. 5', '/minio-api/foods/5.jpg'),
       ('The Green Table', 'Экологичный ресторан с органической едой и вегетарианским меню',
        'Уникальное место, где ценятся здоровье и экология. У нас представлено широкое вегетарианское и веганское меню, а также блюда из свежих органических продуктов', 'пр-кт. Вернадского, д. 10', '/minio-api/restaurants/6.jpg'),
        ('Cafe Rouge', 'Французское кафе с аутентичными блюдами и винной картой',
        'Настоящий кусочек Франции в самом центре города. У нас вы насладитесь изысканными блюдами французской кухни, а также отличным выбором местных и импортных вин', 'пл. Пушкина, д. 10', '/minio-api/restaurants/7.jpg');
        ('Bella Napoli', 'Уютное кафе в центре города',
        'Пиццерия Bella Napoli - это уютное заведение, где вас ждут ароматные итальянские пиццы, приготовленные по традиционным рецептам. Наслаждайтесь каждым укусом нашего лакомства, окунитесь в атмосферу итальянского города Неаполь', 'пл. Победы, д. 7', '/minio-api/restaurants/8.jpg'),


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

