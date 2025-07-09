-- GET PRODUCT (ID)
SELECT
	p.id,
	p.name,
	p.description,
	p.price,
	c.id AS category_id,
	c.name AS category_name,
	COALESCE(
		JSON_AGG(
			JSON_BUILD_OBJECT('url', pi.url, 'alt_text', pi.alt_text)
		) FILTER (WHERE pi.url IS NOT NULL),
		'[]'
		) AS images
FROM products p
LEFT JOIN product_images pi ON p.id = pi.product_id
LEFT JOIN categories c ON c.id = p.category_id
WHERE p.id = $1
GROUP BY p.id, p.name, c.id, c.name;
-- GET LIST OF PRODUCTS (FILT, PAGINATION)

:CATEGORYID = 3
:MIN_PRICE = 0
:MAX_PRICE = 500
:SEARCH = 'laptop'
:MIN_RATING = 3
:LIM = 10
:OFFSET = 0
SELECT 
    p.*, 
    c.name AS category_name,
    COALESCE(AVG(r.score), 0) AS avg_rating
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN ratings r ON r.product_id = p.id
WHERE
    ( p.category_id = 3)
    AND (:min_price IS NULL OR p.price >= :min_price)
    AND (:max_price IS NULL OR p.price <= :max_price)
    AND (:search IS NULL OR p.name ILIKE '%' || :search || '%' OR p.description ILIKE '%' || :search || '%')
GROUP BY p.id, c.name
HAVING (:min_rating IS NULL OR COALESCE(AVG(r.score), 0) >= :min_rating)
ORDER BY p.id
LIMIT :limit OFFSET :offset;
