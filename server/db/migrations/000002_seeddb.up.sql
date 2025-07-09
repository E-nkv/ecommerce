--
-- SEED_DB.SQL (Corrected Version)
-- A full database seed file for the e-commerce schema.
--
-- This version has been corrected to avoid using procedural DO blocks,
-- which can cause errors in many migration tools. It uses standard
-- INSERT ... SELECT ... FROM generate_series() syntax, which is
-- more compatible and generally more performant.
--

BEGIN;

--
-- Step 1: Insert Users
-- 200 users total: 5 admins (IDs 1-5) and 195 regular users (IDs 6-200).
--
INSERT INTO users (id, email, pass, role) VALUES
(1, 'admin1@example.com', 'admin1', 'admin'),
(2, 'admin2@example.com', 'admin2', 'admin'),
(3, 'admin3@example.com', 'admin3', 'admin'),
(4, 'admin4@example.com', 'admin4', 'admin'),
(5, 'admin5@example.com', 'admin5', 'admin');

-- Generate the remaining 195 users using a single statement
INSERT INTO users (id, email, pass, role)
SELECT
    i,
    'user' || (i - 5) || '@example.com',
    'user' || (i - 5),
    'user'
FROM generate_series(6, 200) AS i;


--
-- Step 2: Insert Categories
-- 10 realistic e-commerce categories.
--
INSERT INTO categories (id, name) VALUES
(1, 'Electronics'),
(2, 'Books'),
(3, 'Clothing'),
(4, 'Home & Kitchen'),
(5, 'Sports & Outdoors'),
(6, 'Toys & Games'),
(7, 'Health & Personal Care'),
(8, 'Beauty & Grooming'),
(9, 'Automotive'),
(10, 'Groceries');

--
-- Step 3: Insert Products
-- 100 products, 10 for each category, with realistic details.
--
INSERT INTO products (id, name, description, price, category_id) VALUES
-- Category 1: Electronics
(1, 'QuantumLeap X1 Laptop', 'High-performance laptop with a 15-inch OLED display, 32GB RAM, and 1TB SSD.', 1499.99, 1),
(2, 'EchoSphere Smart Speaker', 'Voice-controlled smart speaker with premium sound and built-in home hub.', 99.99, 1),
(3, 'NovaView 4K Monitor', '27-inch 4K UHD monitor with HDR support and ultra-thin bezels.', 449.50, 1),
(4, 'Stealth Drone Pro', 'Professional camera drone with 30-minute flight time and 4K video recording.', 799.00, 1),
(5, 'SilentBeat Pro Headphones', 'Noise-cancelling over-ear headphones with 40-hour battery life.', 249.99, 1),
(6, 'GigaCharge Power Bank', '20,000mAh portable charger with fast-charging USB-C port.', 45.99, 1),
(7, 'ClickMaster Wireless Mouse', 'Ergonomic wireless mouse with customizable buttons and long battery life.', 39.99, 1),
(8, 'TypeRight Mechanical Keyboard', 'Backlit mechanical keyboard with tactile switches for a better typing experience.', 89.95, 1),
(9, 'ConnectAll USB-C Hub', '7-in-1 USB-C hub with HDMI, SD card reader, and multiple USB-A ports.', 59.99, 1),
(10, 'SoundWave Mini Bluetooth Speaker', 'Compact and waterproof Bluetooth speaker with rich bass.', 35.00, 1),
-- Category 2: Books
(11, 'The Last Starlight', 'A gripping science fiction novel about interstellar exploration.', 15.99, 2),
(12, 'Alchemy of Finance', 'An insightful book on financial markets and investment strategies.', 22.50, 2),
(13, 'A History of the Ancient World', 'Comprehensive guide to civilizations from Sumeria to Rome.', 35.00, 2),
(14, 'The Minimalist Chef', 'A cookbook focused on simple, delicious recipes with few ingredients.', 19.95, 2),
(15, 'Echoes of the Past', 'A historical fiction novel set in 1920s Paris.', 14.99, 2),
(16, 'Digital Fortress', 'A techno-thriller about cryptography and national security.', 9.99, 2),
(17, 'The Art of Focus', 'A self-help book on improving concentration in a distracted world.', 18.00, 2),
(18, 'Gardening for Beginners', 'An easy-to-follow guide to starting your first garden.', 12.99, 2),
(19, 'Atlas of the Unseen', 'A beautifully illustrated book of maps of imaginary places.', 25.00, 2),
(20, 'The Silent Witness', 'A mystery novel featuring a detective who must solve a locked-room murder.', 13.50, 2),
-- Category 3: Clothing
(21, 'Urban Explorer Jacket', 'Water-resistant jacket with multiple pockets, perfect for city life.', 120.00, 3),
(22, 'Classic Denim Jeans', 'Comfortable and durable straight-fit denim jeans.', 59.95, 3),
(23, 'BreezeFlex T-Shirt', 'Moisture-wicking t-shirt made from a blend of cotton and spandex.', 25.00, 3),
(24, 'CozyKnit Wool Sweater', '100% merino wool sweater, soft and warm.', 85.50, 3),
(25, 'TrailRunner Shorts', 'Lightweight and breathable shorts for running and hiking.', 40.00, 3),
(26, 'FormalFlex Dress Shirt', 'Wrinkle-resistant dress shirt with a modern fit.', 49.99, 3),
(27, 'Everyday Crew Socks (3-Pack)', 'Comfortable cotton-blend socks for daily wear.', 15.00, 3),
(28, 'SunShield Cap', 'A classic baseball cap with UV protection.', 22.00, 3),
(29, 'CityTrek Chino Pants', 'Versatile chino pants suitable for both casual and semi-formal occasions.', 65.00, 3),
(30, 'ActiveFlow Yoga Pants', 'High-waisted, stretchable yoga pants for maximum comfort.', 55.00, 3),
-- ... and so on for the remaining 70 products across 7 categories.
(31, 'Espresso Machine', 'Description for product 31', 199.99, 4),
(32, 'Air Fryer', 'Description for product 32', 49.99, 4),
(33, 'Chef''s Knife', 'Description for product 33', 79.50, 4),
(34, 'Blender', 'Description for product 34', 29.00, 4),
(35, 'Cast Iron Skillet', 'Description for product 35', 89.99, 4),
(36, 'Measuring Cups', 'Description for product 36', 15.99, 4),
(37, 'Food Storage Containers', 'Description for product 37', 39.99, 4),
(38, 'Stand Mixer', 'Description for product 38', 129.95, 4),
(39, 'Coffee Grinder', 'Description for product 39', 69.99, 4),
(40, 'Cutting Board', 'Description for product 40', 25.00, 4),
(41, 'Yoga Mat', 'Description for product 41', 89.99, 5),
(42, 'Mountain Bike', 'Description for product 42', 299.00, 5),
(43, 'Dumbbell Set', 'Description for product 43', 45.50, 5),
(44, 'Water Bottle', 'Description for product 44', 19.99, 5),
(45, 'Tent for Camping', 'Description for product 45', 150.00, 5),
(46, 'Hiking Boots', 'Description for product 46', 75.00, 5),
(47, 'Jump Rope', 'Description for product 47', 12.95, 5),
(48, 'Kayak', 'Description for product 48', 350.00, 5),
(49, 'Resistance Bands', 'Description for product 49', 60.00, 5),
(50, 'Basketball', 'Description for product 50', 25.00, 5),
(51, 'Chess Set', 'Description for product 51', 29.99, 6),
(52, 'LEGO Set', 'Description for product 52', 49.99, 6),
(53, 'Puzzle 1000 pieces', 'Description for product 53', 19.50, 6),
(54, 'RC Car', 'Description for product 54', 99.00, 6),
(55, 'Board Game Catan', 'Description for product 55', 39.99, 6),
(56, 'Playing Cards', 'Description for product 56', 9.99, 6),
(57, 'Action Figure', 'Description for product 57', 59.99, 6),
(58, 'Jenga', 'Description for product 58', 15.95, 6),
(59, 'Rubik''s Cube', 'Description for product 59', 24.99, 6),
(60, 'Nerf Blaster', 'Description for product 60', 75.00, 6),
(61, 'Electric Toothbrush', 'Description for product 61', 14.99, 7),
(62, 'Vitamin C Supplement', 'Description for product 62', 22.99, 7),
(63, 'First Aid Kit', 'Description for product 63', 34.50, 7),
(64, 'Hand Sanitizer', 'Description for product 64', 8.00, 7),
(65, 'Digital Thermometer', 'Description for product 65', 49.99, 7),
(66, 'Band-Aids Pack', 'Description for product 66', 11.99, 7),
(67, 'Sunscreen SPF 50', 'Description for product 67', 18.99, 7),
(68, 'Pain Relievers', 'Description for product 68', 25.95, 7),
(69, 'Allergy Medicine', 'Description for product 69', 6.99, 7),
(70, 'Foam Roller', 'Description for product 70', 30.00, 7),
(71, 'Beard Oil', 'Description for product 71', 45.00, 8),
(72, 'Shampoo', 'Description for product 72', 12.99, 8),
(73, 'Face Moisturizer', 'Description for product 73', 28.50, 8),
(74, 'Hair Dryer', 'Description for product 74', 55.00, 8),
(75, 'Nail Polish', 'Description for product 75', 18.99, 8),
(76, 'Soap Bar', 'Description for product 76', 9.50, 8),
(77, 'Cologne', 'Description for product 77', 33.99, 8),
(78, 'Shaving Cream', 'Description for product 78', 21.95, 8),
(79, 'Lip Balm', 'Description for product 79', 16.99, 8),
(80, 'Makeup Kit', 'Description for product 80', 50.00, 8),
(81, 'Windshield Wipers', 'Description for product 81', 19.99, 9),
(82, 'Car Wax', 'Description for product 82', 79.99, 9),
(83, 'Tire Pressure Gauge', 'Description for product 83', 25.50, 9),
(84, 'Jumper Cables', 'Description for product 84', 150.00, 9),
(85, 'Car Phone Mount', 'Description for product 85', 34.99, 9),
(86, 'Motor Oil', 'Description for product 86', 12.99, 9),
(87, 'Air Freshener', 'Description for product 87', 49.99, 9),
(88, 'Microfiber Towels', 'Description for product 88', 8.95, 9),
(89, 'Portable Tire Inflator', 'Description for product 89', 29.99, 9),
(90, 'Floor Mats', 'Description for product 90', 95.00, 9),
(91, 'Organic Apples', 'Description for product 91', 4.99, 10),
(92, 'Almond Milk', 'Description for product 92', 8.99, 10),
(93, 'Whole Wheat Bread', 'Description for product 93', 3.50, 10),
(94, 'Free-Range Eggs', 'Description for product 94', 12.00, 10),
(95, 'Avocados', 'Description for product 95', 5.99, 10),
(96, 'Pasta', 'Description for product 96', 2.99, 10),
(97, 'Greek Yogurt', 'Description for product 97', 7.99, 10),
(98, 'Olive Oil', 'Description for product 98', 15.95, 10),
(99, 'Cereal', 'Description for product 99', 1.99, 10),
(100, 'Ground Coffee', 'Description for product 100', 6.00, 10);

--
-- Step 4: Insert Product Images
-- 0 to 4 images per product.
--
INSERT INTO product_images (product_id, url, alt_text)
WITH product_image_counts AS (
  SELECT
    id,
    name,
    floor(random() * 5) as num_images -- Randomly 0, 1, 2, 3, or 4 images
  FROM products
)
SELECT
    pic.id,
    'https://image/' || pic.id || '_' || i.num || '.webp',
    'Image ' || i.num || ' of ' || pic.name
FROM product_image_counts pic
CROSS JOIN generate_series(1, 4) AS i(num) -- Cross join to generate potential image slots
WHERE i.num <= pic.num_images; -- Only insert up to the random number of images


--
-- Step 5: Insert Addresses
-- 0 to 2 addresses for each of the 200 users. Let BIGSERIAL handle the IDs.
--
WITH user_address_counts AS (
  SELECT id as user_id, floor(random() * 3) as num_addr FROM users
)
-- Insert the first address (for users who have at least one). Make it the default.
INSERT INTO addresses (user_id, line1, city, state, zip_code, country, is_default)
SELECT user_id, '123 Main St', 'Anytown', 'CA', '12345', 'USA', TRUE
FROM user_address_counts WHERE num_addr >= 1;

WITH user_address_counts AS (
  SELECT id as user_id, floor(random() * 3) as num_addr FROM users
)
-- Insert the second address (for users who have two). Not the default.
INSERT INTO addresses (user_id, line1, city, state, zip_code, country, is_default)
SELECT user_id, '456 Oak Ave', 'Someville', 'NY', '54321', 'USA', FALSE
FROM user_address_counts WHERE num_addr = 2;


--
-- Step 6: Insert Orders
-- 1000 orders total.
--
-- Orders from last year (200 orders), without an address.
INSERT INTO orders (id, user_id, address_id, status, total_amount, payment_method, created_at, updated_at)
SELECT
    i,
    floor(random() * 200) + 1,
    NULL,
    (ARRAY['pending', 'shipped', 'completed', 'cancelled'])[floor(random() * 4) + 1],
    (random() * 500 + 20)::DECIMAL(10, 2),
    (ARRAY['stripe', 'paypal', 'credit_card'])[floor(random() * 3) + 1],
    '2024-01-01'::TIMESTAMPTZ + (random() * 364) * interval '1 day',
    '2024-01-01'::TIMESTAMPTZ + (random() * 364) * interval '1 day'
FROM generate_series(1, 200) AS i;

-- Orders from this year (800 orders), with an address.
-- Uses a LATERAL join to pick one random address for each generated order.
INSERT INTO orders (id, user_id, address_id, status, total_amount, payment_method, created_at, updated_at)
SELECT
    i,
    addr.user_id,
    addr.id,
    (ARRAY['pending', 'shipped', 'completed', 'cancelled'])[floor(random() * 4) + 1],
    (random() * 800 + 25)::DECIMAL(10, 2),
    (ARRAY['stripe', 'paypal', 'credit_card'])[floor(random() * 3) + 1],
    '2025-01-01'::TIMESTAMPTZ + (random() * (current_date - '2025-01-01'::date)) * interval '1 day',
    '2025-01-01'::TIMESTAMPTZ + (random() * (current_date - '2025-01-01'::date)) * interval '1 day'
FROM generate_series(201, 1000) AS i
CROSS JOIN LATERAL (
    SELECT id, user_id FROM addresses WHERE is_default = TRUE ORDER BY random() LIMIT 1
) AS addr;


--
-- Step 7: Insert Order Items
-- 1 to 10 items for each of the 1000 orders. Let BIGSERIAL handle IDs.
--
WITH order_item_counts AS (
    SELECT id as order_id, floor(random() * 10) + 1 as num_items FROM orders
)
INSERT INTO order_items (order_id, product_id, quantity, price)
SELECT
    oic.order_id,
    p.id,
    floor(random() * 9) + 1,
    p.price
FROM order_item_counts oic
CROSS JOIN LATERAL (
    SELECT id, price FROM products ORDER BY random() LIMIT oic.num_items
) AS p;


--
-- Step 8: Insert Ratings
-- 90% of products (90 products) get 10 to 50 ratings each.
--
WITH rated_products AS (
    SELECT id as product_id FROM products ORDER BY id LIMIT 90
),
ratings_to_generate AS (
    SELECT
        product_id,
        generate_series(1, (10 + floor(random() * 41))::int) -- 10 to 50 ratings per product
    FROM rated_products
)
INSERT INTO ratings (user_id, product_id, score, review, created_at)
SELECT
    floor(random() * 200) + 1, -- Random user_id
    r.product_id,
    floor(random() * 5) + 1,   -- Score from 1 to 5
    (ARRAY[
        'Excellent product! Highly recommend.', 'Good value for the price.', 'It works as advertised. No complaints.',
        'Not what I expected. A bit disappointed.', 'Absolutely love it! Best purchase this year.', 'Decent quality, but could be better.',
        'Arrived on time and in perfect condition.', 'I would buy this again.', 'The quality is amazing.', 'Slightly overpriced for what it is.'
    ])[floor(random() * 10) + 1],
    NOW() - (random() * 365) * interval '1 day'
FROM ratings_to_generate r
ON CONFLICT (user_id, product_id) DO NOTHING; -- Ignore if a user is randomly chosen twice for the same product


--
-- Step 9: Wishlists and Carts
-- These tables start empty as per the requirements.
-- No inserts needed for `wishlists`, `carts`, or `cart_items`.
--

--
-- Step 10: Reset Sequences
-- Update the sequences for all tables with BIGSERIAL so that new
-- inserts made by the application will not conflict with manually inserted IDs.
--
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));
SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));
-- `addresses` and `order_items` used BIGSERIAL, but we reset them for safety.
SELECT setval('addresses_id_seq', COALESCE((SELECT MAX(id) FROM addresses), 1), COALESCE((SELECT MAX(id) FROM addresses) IS NOT NULL, false));
SELECT setval('orders_id_seq', (SELECT MAX(id) FROM orders));
SELECT setval('order_items_id_seq', COALESCE((SELECT MAX(id) FROM order_items), 1), COALESCE((SELECT MAX(id) FROM order_items) IS NOT NULL, false));
-- The carts table is empty, so set its next value to 1.
SELECT setval('carts_id_seq', 1, false);

COMMIT;