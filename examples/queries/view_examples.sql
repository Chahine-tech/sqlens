-- SQL Parser Go - VIEW Examples
-- Demonstrates CREATE VIEW and DROP VIEW parsing across different dialects

-- ============================================================================
-- Simple CREATE VIEW
-- ============================================================================

-- Basic view
CREATE VIEW active_users AS
SELECT id, name, email, created_at
FROM users
WHERE active = 1;

-- View with aggregate functions
CREATE VIEW user_stats AS
SELECT user_id, COUNT(*) as order_count, SUM(total) as total_spent
FROM orders
GROUP BY user_id;

-- ============================================================================
-- CREATE VIEW with schema
-- ============================================================================

CREATE VIEW myschema.customer_orders AS
SELECT u.name, o.order_id, o.total, o.created_at
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE o.status = 'completed';

-- ============================================================================
-- CREATE OR REPLACE VIEW
-- ============================================================================

CREATE OR REPLACE VIEW high_value_orders AS
SELECT * FROM orders
WHERE total > 1000
ORDER BY total DESC;

-- ============================================================================
-- CREATE VIEW IF NOT EXISTS
-- ============================================================================

CREATE VIEW IF NOT EXISTS recent_orders AS
SELECT * FROM orders
WHERE created_at > DATE_SUB(NOW(), INTERVAL 30 DAY);

-- ============================================================================
-- CREATE VIEW with column list
-- ============================================================================

CREATE VIEW user_summary (user_id, total_orders, total_amount, avg_amount) AS
SELECT user_id, COUNT(*), SUM(total), AVG(total)
FROM orders
GROUP BY user_id;

-- ============================================================================
-- CREATE VIEW with WITH CHECK OPTION
-- ============================================================================

-- MySQL/PostgreSQL: WITH CHECK OPTION ensures inserts/updates through view respect WHERE clause
CREATE VIEW premium_customers AS
SELECT * FROM users
WHERE subscription_type = 'premium'
WITH CHECK OPTION;

-- ============================================================================
-- CREATE MATERIALIZED VIEW (PostgreSQL)
-- ============================================================================

-- Materialized views store the query result physically
CREATE MATERIALIZED VIEW sales_summary AS
SELECT product_id, SUM(amount) as total_sales, COUNT(*) as sale_count
FROM sales
GROUP BY product_id;

-- With IF NOT EXISTS
CREATE MATERIALIZED VIEW IF NOT EXISTS monthly_revenue AS
SELECT 
    DATE_TRUNC('month', order_date) as month,
    SUM(total) as revenue,
    COUNT(*) as order_count
FROM orders
GROUP BY DATE_TRUNC('month', order_date);

-- ============================================================================
-- CREATE OR REPLACE MATERIALIZED VIEW
-- ============================================================================

CREATE OR REPLACE MATERIALIZED VIEW top_products AS
SELECT product_id, product_name, SUM(quantity) as total_sold
FROM order_items
JOIN products ON order_items.product_id = products.id
GROUP BY product_id, product_name
ORDER BY total_sold DESC
LIMIT 100;

-- ============================================================================
-- Complex Views with Subqueries
-- ============================================================================

CREATE VIEW above_average_salaries AS
SELECT * FROM employees
WHERE salary > (SELECT AVG(salary) FROM employees);

CREATE VIEW department_stats AS
SELECT 
    d.name as department,
    COUNT(e.id) as employee_count,
    AVG(e.salary) as avg_salary,
    (SELECT MAX(salary) FROM employees WHERE department_id = d.id) as max_salary
FROM departments d
LEFT JOIN employees e ON d.id = e.department_id
GROUP BY d.id, d.name;

-- ============================================================================
-- Views with CTEs (Common Table Expressions)
-- ============================================================================

CREATE VIEW quarterly_sales AS
WITH quarterly_data AS (
    SELECT 
        EXTRACT(YEAR FROM order_date) as year,
        EXTRACT(QUARTER FROM order_date) as quarter,
        SUM(total) as revenue
    FROM orders
    GROUP BY year, quarter
)
SELECT * FROM quarterly_data
ORDER BY year DESC, quarter DESC;

-- ============================================================================
-- Views with Window Functions
-- ============================================================================

CREATE VIEW employee_rankings AS
SELECT 
    id,
    name,
    department_id,
    salary,
    ROW_NUMBER() OVER (PARTITION BY department_id ORDER BY salary DESC) as dept_rank,
    RANK() OVER (ORDER BY salary DESC) as company_rank
FROM employees;

-- ============================================================================
-- DROP VIEW
-- ============================================================================

-- Simple DROP
DROP VIEW active_users;

-- DROP IF EXISTS
DROP VIEW IF EXISTS user_stats;

-- DROP with schema
DROP VIEW myschema.customer_orders;

-- DROP with CASCADE (PostgreSQL - removes dependent views)
DROP VIEW IF EXISTS sales_summary CASCADE;

-- ============================================================================
-- DROP MATERIALIZED VIEW (PostgreSQL)
-- ============================================================================

DROP MATERIALIZED VIEW sales_summary;

DROP MATERIALIZED VIEW IF EXISTS monthly_revenue;

-- ============================================================================
-- Dialect-Specific Examples
-- ============================================================================

-- MySQL with backticks
CREATE VIEW `user_orders` AS
SELECT `users`.`id`, `users`.`name`, COUNT(`orders`.`id`) as `order_count`
FROM `users`
LEFT JOIN `orders` ON `users`.`id` = `orders`.`user_id`
GROUP BY `users`.`id`, `users`.`name`;

-- PostgreSQL with double quotes
CREATE VIEW "user_orders" AS
SELECT "users"."id", "users"."name", COUNT("orders"."id") as "order_count"
FROM "users"
LEFT JOIN "orders" ON "users"."id" = "orders"."user_id"
GROUP BY "users"."id", "users"."name";

-- SQL Server with brackets
CREATE VIEW [user_orders] AS
SELECT [users].[id], [users].[name], COUNT([orders].[id]) as [order_count]
FROM [users]
LEFT JOIN [orders] ON [users].[id] = [orders].[user_id]
GROUP BY [users].[id], [users].[name];

-- ============================================================================
-- Real-World Examples
-- ============================================================================

-- E-commerce: Active product catalog
CREATE VIEW active_products AS
SELECT p.id, p.name, p.price, c.name as category, i.quantity
FROM products p
JOIN categories c ON p.category_id = c.id
JOIN inventory i ON p.id = i.product_id
WHERE p.is_active = 1 AND i.quantity > 0;

-- Analytics: Daily order summary
CREATE VIEW daily_order_summary AS
SELECT 
    DATE(created_at) as order_date,
    COUNT(*) as total_orders,
    SUM(total) as total_revenue,
    AVG(total) as avg_order_value,
    COUNT(DISTINCT user_id) as unique_customers
FROM orders
WHERE status = 'completed'
GROUP BY DATE(created_at);

-- User management: User permissions view
CREATE VIEW user_permissions AS
SELECT 
    u.id as user_id,
    u.username,
    u.email,
    r.name as role,
    GROUP_CONCAT(p.name) as permissions
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
GROUP BY u.id, u.username, u.email, r.name;
