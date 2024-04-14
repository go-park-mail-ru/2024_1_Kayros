-- Заполнение БД данными
INSERT INTO "user" (name, email, password, address, created_at, updated_at)
VALUES ('Иван', 'ivan@example.com',
        E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'ул. Ленина, д. 1', current_timestamp, current_timestamp),
       ('Анна', 'anna@example.com',
        E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'пр. Победы, д. 10', current_timestamp, current_timestamp),
       ('Петр', 'petr@example.com',
        E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'ул. Мира, д. 5', current_timestamp, current_timestamp),
       ('Мария', 'maria@example.com',
        E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'пр. Ленинградский, д. 15', current_timestamp, current_timestamp),
       ('Алексей', 'alex@example.com',
        E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'ул. Советская, д. 25', current_timestamp, current_timestamp),
       ('Елена', 'elena@example.com',
        E'\\x6AAA04147C1662A5D5608B40030CD4163F0ABAB80BF1B4F37863278237FB3F429DFA3527143F96DBB01971553A70EAF79D592B5BAFB4D229DF16EA67833D69F7',
        'пр. Гагарина, д. 30', current_timestamp, current_timestamp);

-- Вставка данных в таблицу "Category"
INSERT INTO category (name)
VALUES ('Завтраки'),
       ('Закуски'),
       ('Салаты'),
       ('Пицца'),
       ('Паста'),
       ('Супы'),
       ('Горячие блюда');


-- Вставка данных в таблицу "Restaurant"
INSERT INTO restaurant (name, short_description, long_description, address, img_url)
VALUES ('Горыныч', 'Лучший ресторан в городе',
        'Здесь подают самую вкусную еду, приходите и наслаждайтесь!', 'пр. Победы, д. 20',
        '/minio-api/restaurants/1.jpg'),

       ('Loona', 'Ресторан с итальянской атмосферой',
        'У нас вы найдете самую вкусную пиццу, пасту, закуски и  вкусные десерты',
        'ул. Садовая, д. 15', '/minio-api/restaurants/2.jpg'),

       ('Rustic Kitchen', 'Уютное семейное бистро с домашней кухней и гостеприимной атмосферой',
        'место, где вы чувствуете себя как дома. У нас готовят по-настоящему вкусную еду, используя только свежие продукты высокого качества',
        'ул. Невская, д. 1', '/minio-api/restaurants/3.jpg'),

       ('Густо', 'Ресторан с авторской европейской кухней',
        'Место, где каждый член семьи найдет что-то по душе. У нас вы попробуете неповторимые блюда европейской кухни, приготовленные с любовью к деталям',
        'ул. Садовая, д. 5', '/minio-api/restaurants/4.jpg'),

       ('Remy', 'Достижения французского кухни',
        'Шеф-повар создает удивительные шедевры из самых редких и изысканных ингредиентов, чтобы удовлетворить самые изысканные вкусы гостей. Ощутите волшебство атмосферы и насладитесь бесподобными блюдами, приготовленными с любовью и талантом',
        'ул. Московская, д. 5', '/minio-api/foods/5.jpg'),

       ('The Green Table', 'Экологичный ресторан с органической едой и вегетарианским меню',
        'Уникальное место, где ценятся здоровье и экология. У нас представлено широкое вегетарианское и веганское меню, а также блюда из свежих органических продуктов',
        'пр-кт. Вернадского, д. 10', '/minio-api/restaurants/6.jpg'),
       ('Cafe Rouge', 'Французское кафе с аутентичными блюдами и винной картой',
        'Настоящий кусочек Франции в самом центре города. У нас вы насладитесь изысканными блюдами французской кухни, а также отличным выбором местных и импортных вин',
        'пл. Пушкина, д. 10', '/minio-api/restaurants/7.jpg'),
        ('Bella Napoli', 'Уютное кафе в центре города',
        'Пиццерия Bella Napoli - это уютное заведение, где вас ждут ароматные итальянские пиццы, приготовленные по традиционным рецептам. Наслаждайтесь каждым укусом нашего лакомства, окунитесь в атмосферу итальянского города Неаполь',
        'пл. Победы, д. 7', '/minio-api/restaurants/8.jpg');


-- Вставка данных в таблицу "Order"
INSERT INTO "order" (user_id, created_at, updated_at, status, address, extra_address, sum)
VALUES (1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'processing', 'ул. Ленина, д. 1', NULL, 1500),
       (2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'draft', 'пр. Победы, д. 10', NULL, 2000),
       (3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'completed', 'ул. Мира, д. 5', NULL, 1000),
       (4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'processing', 'пр. Ленинградский, д. 15', NULL, 1200),
       (5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'draft', 'ул. Советская, д. 25', NULL, 1800),
       (6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'completed', 'пр. Гагарина, д. 30', NULL, 900);

-- Вставка данных в таблицу "Food"
INSERT INTO food (name, restaurant_id, category_id, weight, price, proteins, fats, carbohydrates, img_url)
VALUES ('Салат с миксом листьев и ростками сои', 1, 3, 220, 450, NULL, NULL, NULL, '/minio-api/foods/1.jpg'),
       ('Салат с сыром Рокфор, грушей и грецкими орехами', 1, 3, 200, 500, NULL, NULL, NULL, '/minio-api/foods/2.jpg'),
       ('Суп-пюре из тыквы', 1, 6, 300, 420, NULL, NULL, NULL, '/minio-api/foods/3.jpg'),
       ('Террин из кролика с домашним хлебом', 1, 2, 180, 480, NULL, NULL, NULL, '/minio-api/foods/4.jpg'),
       ('Пицца с тремя видами сыра', 1, 4, 350, 520, NULL, NULL, NULL, '/minio-api/foods/5.jpg'),
       ('Паста с соусом песто и кедровыми орехами', 1, 5, 300, 450, NULL, NULL, NULL, '/minio-api/foods/6.jpg'),

       ('Салат с помидорами и моцареллой', 2, 3, 180, 400, NULL, NULL, NULL, '/minio-api/foods/7.jpg'),
       ('Салат "Антипасти"', 2, 3, 200, 550, NULL, NULL, NULL, '/minio-api/foods/8.jpg'),
       ('Минестроне', 2, 2, 300, 430, NULL, NULL, NULL, '/minio-api/foods/9.jpg'),
       ('Карпаччо из телятины', 2, 2, 150, 590, NULL, NULL, NULL, '/minio-api/foods/10.jpg'),
       ('Пицца "Диавола"', 2, 4, 400, 560, NULL, NULL, NULL, '/minio-api/foods/11.jpg'),
       ('Паста "Болоньезе"', 2, 5, 320, 460, NULL, NULL, NULL, '/minio-api/foods/12.jpg'),

       ('Салат с домашней курятиной и свежими овощами', 3, 3, 220, 420, NULL, NULL, NULL, '/minio-api/foods/13.jpg'),
       ('Зеленый салат с авокадо, шпинатом и семенами чиа', 3, 3, 190, 430, NULL, NULL, NULL,
        '/minio-api/foods/14.jpg'),
       ('Чечевичный суп', 3, 6, 250, 410, NULL, NULL, NULL, '/minio-api/foods/15.jpg'),
       ('Тосты с козьим сыром', 3, 2, 160, 450, NULL, NULL, NULL, '/minio-api/foods/16.jpg'),
       ('Пицца с грибами и курятиной на тонком тесте', 3, 4, 380, 540, NULL, NULL, NULL, '/minio-api/foods/17.jpg'),
       ('Карбонара', 3, 5, 350, 470, NULL, NULL, NULL, '/minio-api/foods/18.jpg'),

       ('Салат с бифштексом и рукколой', 4, 3, 210, 510, NULL, NULL, NULL, '/minio-api/foods/19.jpg'),
       ('Салат с морепродуктами', 4, 3, 200, 530, NULL, NULL, NULL, '/minio-api/foods/20.jpg'),
       ('Гаспачо с креветками', 4, 6, 280, 460, NULL, NULL, NULL, '/minio-api/foods/21.jpg'),
       ('Жюльен с курицей и грибами', 4, 2, 170, 420, NULL, NULL, NULL, '/minio-api/foods/22.jpg'),
       ('Пицца "Четыре сыра"', 4, 4, 360, 550, NULL, NULL, NULL, '/minio-api/foods/23.jpg'),
       ('Паста с курицей и грибами в сливочном соусе', 4, 5, 330, 490, NULL, NULL, NULL, '/minio-api/foods/24.jpg'),

       ('Салат с копченой уткой и ягодами', 5, 3, 210, 600, NULL, NULL, NULL, '/minio-api/foods/25.jpg'),
       ('Салат с морским окунем и свежими овощами', 5, 3, 200, 550, NULL, NULL, NULL, '/minio-api/foods/26.jpg'),
       ('Равиоли с трюфелями и сыром Бри', 5, 5, 6, 700, NULL, NULL, NULL, '/minio-api/foods/27.jpg'),
       ('Бульон с фуа-гра', 5, 6, 250, 650, NULL, NULL, NULL, '/minio-api/foods/28.jpg'),
       ('Пицца с тунцом и маслинами', 5, 4, 380, 570, NULL, NULL, NULL, '/minio-api/foods/29.jpg'),
       ('Паста "Рататуй"', 5, 5, 320, 520, NULL, NULL, NULL, '/minio-api/foods/30.jpg'),

       ('Салат с киноа и цитрусовым дрессингом', 6, 3, 220, 450, NULL, NULL, NULL, '/minio-api/foods/31.jpg'),
       ('Салат с белыми бобами и сушеными томатами', 6, 3, 210, 440, NULL, NULL, NULL, '/minio-api/foods/32.jpg'),
       ('Кремовый томатный суп с базиликом и моцареллой', 6, 6, 250, 420, NULL, NULL, NULL, '/minio-api/foods/33.jpg'),
       ('Хумус с овощами гриль и лавашом', 6, 2, 180, 430, NULL, NULL, NULL, '/minio-api/foods/34.jpg'),
       ('Веганская пицца с веганским сыром и овощами', 6, 4, 400, 560, NULL, NULL, NULL, '/minio-api/foods/35.jpg'),
       ('Паста с свежими помидорами и базиликом', 6, 5, 320, 450, NULL, NULL, NULL, '/minio-api/foods/36.jpg'),

       ('Салат с пармской ветчиной, моцареллой и руколой', 7, 3, 200, 520, NULL, NULL, NULL, '/minio-api/foods/37.jpg'),
       ('Салат с лососем и авокадо', 7, 3, 220, 540, NULL, NULL, NULL, '/minio-api/foods/38.jpg'),
       ('Крем-суп из моркови с имбирем', 7, 6, 280, 410, NULL, NULL, NULL, '/minio-api/foods/39.jpg'),
       ('Фуа-гра на тосте с голубикой', 7, 2, 150, 690, NULL, NULL, NULL, '/minio-api/foods/40.jpg'),
       ('Пицца с голубым сыром и грушей', 7, 4, 350, 580, NULL, NULL, NULL, '/minio-api/foods/41.jpg'),
       ('Паста с креветками', 7, 5, 300, 620, NULL, NULL, NULL, '/minio-api/foods/42.jpg'),

       ('Салат с томатами, моцареллой и прованскими травами', 8, 3, 180, 430, NULL, NULL, NULL,
        '/minio-api/foods/43.jpg'),
       ('Салат с морепродуктами и зеленью', 8, 3, 210, 560, NULL, NULL, NULL, '/minio-api/foods/44.jpg'),
       ('Суп с сезонными овощами и белой фасолью', 8, 6, 300, 420, NULL, NULL, NULL, '/minio-api/foods/45.jpg'),
       ('Брускетта с прошутто и рикоттой', 8, 2, 160, 450, NULL, NULL, NULL, '/minio-api/foods/46.jpg'),
       ('Пицца с ветчиной и грибами', 8, 4, 380, 550, NULL, NULL, NULL, '/minio-api/foods/47.jpg'),
       ('Паста с морепродуктами', 8, 5, 350, 580, NULL, NULL, NULL, '/minio-api/foods/48.jpg'),

       ('Тальята из курицы с томатами черри и кинзой', 1, 7, 230, 550, NULL, NULL, NULL, '/minio-api/foods/49.jpg'),
       ('Филе лосося с овощами', 1, 7, 240, 1250, NULL, NULL, NULL, '/minio-api/foods/50.jpg'),
       ('Тушеная говядина Бразато с трюфельным пюре', 1, 7, 380, 1290, NULL, NULL, NULL, '/minio-api/foods/51.jpg'),
       ('Индейка по-милански', 1, 7, 335, 790, NULL, NULL, NULL, '/minio-api/foods/52.jpg'),
       ('Ригатони с баклажанами и кремом из рикотты', 1, 5, 370, 690, NULL, NULL, NULL, '/minio-api/foods/53.jpg'),
       ('Тальятелле с домашним рагу', 1, 5, 320, 850, NULL, NULL, NULL, '/minio-api/foods/54.jpg'),
       ('Ригатони с томатами, оливками и каперсами', 1, 5, 270, 690, NULL, NULL, NULL, '/minio-api/foods/55.jpg'),
       ('Спагетти с веганским рагу', 1, 5, 340, 590, NULL, NULL, NULL, '/minio-api/foods/56.jpg'),
       ('Пицца Моцарелла с ветчиной', 1, 4, 460, 790, NULL, NULL, NULL, '/minio-api/foods/57.jpg'),
       ('Пицца Страчателла с пармой', 1, 4, 600, 1490, NULL, NULL, NULL, '/minio-api/foods/58.jpg'),
       ('Салат Цезарь с фермерским цыпленком и яйцом', 1, 3, 270, 590, NULL, NULL, NULL, '/minio-api/foods/59.jpg'),
       ('Салат с тунцом и яйцами перепелок', 1, 3, 285, 850, NULL, NULL, NULL, '/minio-api/foods/60.jpg'),
       ('Буррата с помидорами и базиликом', 1, 3, 410, 900, NULL, NULL, NULL, '/minio-api/foods/61.jpg'),
       ('Запеченные баклажаны со страчателлой', 1, 3, 260, 690, NULL, NULL, NULL, '/minio-api/foods/62.jpg'),
       ('Салат Цезарь с креветками и яйцом', 1, 3, 335, 990, NULL, NULL, NULL, '/minio-api/foods/63.jpg'),
       ('Салат с киноа, авокадо, брокколи', 1, 3, 270, 690, NULL, NULL, NULL, '/minio-api/foods/64.jpg'),
       ('Томатный суп со страчателлой', 1, 6, 230, 490, NULL, NULL, NULL, '/minio-api/foods/65.jpg'),
       ('Куриный бульон с цыпленком и перепелиными яйцами', 1, 6, 360, 550, NULL, NULL, NULL, '/minio-api/foods/66.jpg'),

       ('Медальон из индейки с горгонзолой', 2, 7, 310, 1150, NULL, NULL, NULL, '/minio-api/foods/67.jpg'),
       ('Грудка цыпленка с фисташкой и полентой', 2, 7, 275, 1350, NULL, NULL, NULL, '/minio-api/foods/68.jpg'),
       ('Стейк из говяжьей вырезки', 2, 7, 380, 1890, NULL, NULL, NULL, '/minio-api/foods/69.jpg'),
       ('Телячья печень с шалфеем и картофельным пюре', 2, 7, 270, 1750, NULL, NULL, NULL, '/minio-api/foods/70.jpg'),
       ('Спагетти Карбонара', 2, 5, 380, 690, NULL, NULL, NULL, '/minio-api/foods/71.jpg'),
       ('Спагетти с морепродуктами', 2, 5, 380, 850, NULL, NULL, NULL, '/minio-api/foods/72.jpg'),
       ('Пенне аррабиата', 2, 5, 350, 650, NULL, NULL, NULL, '/minio-api/foods/73.jpg'),
       ('Пенне с лососем', 2, 5, 370, 1390, NULL, NULL, NULL, '/minio-api/foods/74.jpg'),
       ('Зеленый салат со свежими овощами и пармезаном', 2, 3, 185, 750, NULL, NULL, NULL, '/minio-api/foods/75.jpg'),
       ('Салат с куриной печенью', 2, 3, 210, 910, NULL, NULL, NULL, '/minio-api/foods/76.jpg'),
       ('Руккола с креветками, черри помидорами и авокадо', 2, 3, 400, 1090, NULL, NULL, NULL, '/minio-api/foods/77.jpg'),
       ('Цезарь с креветками', 2, 3, 260, 970, NULL, NULL, NULL, '/minio-api/foods/78.jpg'),


       ('Салат с тыквой, соусом наршараб, фетой, киноа', 3, 3, 450, 200, NULL, NULL, NULL, '/minio-api/foods/79.jpg'),
       ('Перец рамиро, эстрагон, сливочный крем', 3, 3, 150, 470, NULL, NULL, NULL, '/minio-api/foods/80.jpg'),
       ('Лосось, понзу, трюфельное масло, каперсы', 3, 3, 120, 830, NULL, NULL, NULL, '/minio-api/foods/81.jpg'),
       ('Тартар из тунца, авокадо, шрирача', 3, 3, 200, 850, NULL, NULL, NULL, '/minio-api/foods/82.jpg'),
       ('Тартар из говядины, пармезан, трюфельное масло', 3, 3, 170, 790, NULL, NULL, NULL, '/minio-api/foods/83.jpg'),
       ('Микс-салат, копченый лосось, зеленый горошек, ялтинский лук', 3, 5, 240, 770, NULL, NULL, NULL, '/minio-api/foods/84.jpg'),
       ('Цезарь, куриное филе, пармезан, анчоусы', 3, 3, 230, 650, NULL, NULL, NULL, '/minio-api/foods/85.jpg'),
       ('Паста Болоньезе', 3, 5, 200, 690, NULL, NULL, NULL, '/minio-api/foods/86.jpg'),
       ('Паста Арабьята со страчателлой', 3, 5, 325, 790, NULL, NULL, NULL, '/minio-api/foods/87.jpg'),
       ('Спагетти, тефтели, томатный соус', 3, 5, 300, 6700, NULL, NULL, NULL, '/minio-api/foods/88.jpg'),
       ('Свекольное ризотто, утиное филе, чимичурри', 3, 5, 280, 950, NULL, NULL, NULL, '/minio-api/foods/89.jpg'),
       ('Пицца Салями Спьяната', 3, 3, 480, 850, NULL, NULL, NULL, '/minio-api/foods/90.jpg'),
       ('Пицца Четыре сыра', 3, 3, 410, 790, NULL, NULL, NULL, '/minio-api/foods/91.jpg'),
       ('Прошутто Котто, грибы', 3, 3, 520, 690, NULL, NULL, NULL, '/minio-api/foods/92.jpg'),
       ('Тунец, моцарелла, лук, чили', 3, 3, 480, 950, NULL, NULL, NULL, '/minio-api/foods/93.jpg'),
       ('Пицца Маргарита', 3, 3, 420, 750, NULL, NULL, NULL, '/minio-api/foods/94.jpg'),
       ('Пицца Груша-горгондзола', 3, 6, 490, 890, NULL, NULL, NULL, '/minio-api/foods/95.jpg');

-- Вставка данных в таблицу "FoodOrder"
INSERT INTO food_order (food_id, order_id, count, created_at, updated_at)
VALUES (1, 1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (2, 2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (3, 3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       (4, 1, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

