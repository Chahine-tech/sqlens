-- SQL Parser Go - Set Operations Examples
-- UNION, UNION ALL, INTERSECT, EXCEPT

-- ============================================================
-- UNION - Combines results and removes duplicates
-- ============================================================
-- Combining customer and prospect lists
SELECT customer_id, name, email FROM customers
UNION
SELECT prospect_id, name, email FROM prospects;

-- Combining data from multiple regions
SELECT product_id, product_name FROM us_products
UNION
SELECT product_id, product_name FROM eu_products
UNION
SELECT product_id, product_name FROM asia_products;

-- ============================================================
-- UNION ALL - Combines results keeping duplicates
-- ============================================================
-- Combining all transactions (keeping duplicates for accurate counts)
SELECT transaction_id, amount, transaction_date FROM online_sales
UNION ALL
SELECT transaction_id, amount, transaction_date FROM store_sales;

-- Historical data consolidation
SELECT user_id, action, timestamp FROM events_2023
UNION ALL
SELECT user_id, action, timestamp FROM events_2024;

-- ============================================================
-- INTERSECT - Returns only common records
-- ============================================================
-- Finding customers who are also employees
SELECT email FROM customers
INTERSECT
SELECT email FROM employees;

-- Products sold in all regions
SELECT product_id FROM us_sales
INTERSECT
SELECT product_id FROM eu_sales
INTERSECT
SELECT product_id FROM asia_sales;

-- ============================================================
-- EXCEPT - Returns records in first set but not in second
-- ============================================================
-- Customers who haven't made purchases
SELECT customer_id FROM customers
EXCEPT
SELECT customer_id FROM orders;

-- Active users who haven't completed onboarding
SELECT user_id FROM active_users
EXCEPT
SELECT user_id FROM completed_onboarding;

-- ============================================================
-- COMPLEX SET OPERATIONS
-- ============================================================
-- Chained UNION operations
SELECT user_id, 'Free' as tier FROM free_users
UNION
SELECT user_id, 'Premium' as tier FROM premium_users
UNION
SELECT user_id, 'Enterprise' as tier FROM enterprise_users
ORDER BY tier, user_id;

-- Mixing UNION and INTERSECT
SELECT product_id FROM featured_products
UNION
SELECT product_id FROM bestsellers
INTERSECT
SELECT product_id FROM in_stock_products;

-- ============================================================
-- PRACTICAL USE CASES
-- ============================================================
-- Creating a unified contact list
SELECT
    id,
    name,
    email,
    'Customer' as contact_type
FROM customers
WHERE active = 1
UNION ALL
SELECT
    id,
    name,
    email,
    'Lead' as contact_type
FROM leads
WHERE status = 'qualified'
UNION ALL
SELECT
    id,
    name,
    email,
    'Partner' as contact_type
FROM partners
WHERE relationship = 'active';

-- Finding exclusive products per region
-- Products only in US
SELECT product_id, product_name FROM us_inventory
EXCEPT
SELECT product_id, product_name FROM global_inventory;

-- Data quality check - finding duplicates across tables
SELECT email FROM users
INTERSECT
SELECT email FROM deleted_users;

-- ============================================================
-- SET OPERATIONS WITH AGGREGATIONS
-- ============================================================
-- Combining aggregated results
SELECT department, COUNT(*) as emp_count FROM full_time_employees GROUP BY department
UNION ALL
SELECT department, COUNT(*) as emp_count FROM part_time_employees GROUP BY department;

-- Year-over-year comparison
SELECT
    product_id,
    SUM(quantity) as total_sales,
    '2023' as year
FROM sales_2023
GROUP BY product_id
UNION ALL
SELECT
    product_id,
    SUM(quantity) as total_sales,
    '2024' as year
FROM sales_2024
GROUP BY product_id
ORDER BY product_id, year;

-- ============================================================
-- DIALECT-SPECIFIC EXAMPLES
-- ============================================================
-- MySQL with backticks
SELECT `user_id`, `name` FROM `active_users`
UNION
SELECT `user_id`, `name` FROM `pending_users`;

-- PostgreSQL with double quotes
SELECT "product_id", "name" FROM "products"
INTERSECT
SELECT "product_id", "name" FROM "featured_products";

-- SQL Server with brackets
SELECT [customer_id], [email] FROM [customers]
EXCEPT
SELECT [customer_id], [email] FROM [unsubscribed];

-- ============================================================
-- PERFORMANCE CONSIDERATIONS
-- ============================================================
-- Use UNION ALL when duplicates don't matter (faster)
SELECT log_id, message FROM error_logs_2024_01
UNION ALL
SELECT log_id, message FROM error_logs_2024_02
UNION ALL
SELECT log_id, message FROM error_logs_2024_03;

-- Use UNION when you need distinct results
SELECT tag_name FROM article_tags
UNION
SELECT tag_name FROM video_tags
UNION
SELECT tag_name FROM podcast_tags;
