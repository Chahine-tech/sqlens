-- SQL Parser Go - CTE (Common Table Expressions) Examples
-- These examples demonstrate WITH clause support across different dialects

-- ============================================================
-- SIMPLE CTE
-- ============================================================
-- Basic CTE with a single query
WITH sales_summary AS (
    SELECT
        product_id,
        SUM(quantity) as total_quantity,
        SUM(amount) as total_amount
    FROM sales
    WHERE sale_date >= '2024-01-01'
    GROUP BY product_id
)
SELECT
    p.product_name,
    s.total_quantity,
    s.total_amount
FROM sales_summary s
JOIN products p ON s.product_id = p.id
WHERE s.total_amount > 10000;

-- ============================================================
-- CTE WITH COLUMN LIST
-- ============================================================
-- Explicitly naming CTE columns
WITH employee_stats (emp_id, emp_name, dept, avg_score) AS (
    SELECT
        employee_id,
        name,
        department,
        AVG(performance_score)
    FROM employees
    GROUP BY employee_id, name, department
)
SELECT emp_name, dept, avg_score
FROM employee_stats
WHERE avg_score > 85;

-- ============================================================
-- MULTIPLE CTEs
-- ============================================================
-- Using multiple CTEs in a single query
WITH
    active_users AS (
        SELECT user_id, name, email
        FROM users
        WHERE status = 'active'
        AND last_login > DATEADD(day, -30, GETDATE())
    ),
    recent_orders AS (
        SELECT user_id, COUNT(*) as order_count, SUM(total) as total_spent
        FROM orders
        WHERE order_date > DATEADD(day, -30, GETDATE())
        GROUP BY user_id
    )
SELECT
    u.name,
    u.email,
    COALESCE(o.order_count, 0) as orders,
    COALESCE(o.total_spent, 0) as spending
FROM active_users u
LEFT JOIN recent_orders o ON u.user_id = o.user_id
ORDER BY o.total_spent DESC;

-- ============================================================
-- RECURSIVE CTE (Basic Structure)
-- ============================================================
-- Generating number sequence (note: full recursive parsing TBD)
WITH RECURSIVE number_sequence AS (
    SELECT 1 AS n
)
SELECT * FROM number_sequence;

-- ============================================================
-- CTE FOR DATA TRANSFORMATION
-- ============================================================
-- Using CTE to clean and transform data
WITH cleaned_data AS (
    SELECT
        customer_id,
        UPPER(TRIM(name)) as clean_name,
        LOWER(TRIM(email)) as clean_email,
        REPLACE(phone, '-', '') as clean_phone
    FROM raw_customers
    WHERE email IS NOT NULL
)
SELECT * FROM cleaned_data;

-- ============================================================
-- CTE WITH AGGREGATION
-- ============================================================
-- Complex aggregation with CTE
WITH monthly_metrics AS (
    SELECT
        YEAR(order_date) as year,
        MONTH(order_date) as month,
        COUNT(*) as order_count,
        SUM(total_amount) as revenue,
        AVG(total_amount) as avg_order_value
    FROM orders
    WHERE order_date >= '2024-01-01'
    GROUP BY YEAR(order_date), MONTH(order_date)
)
SELECT
    year,
    month,
    order_count,
    revenue,
    avg_order_value,
    revenue - LAG(revenue) OVER (ORDER BY year, month) as revenue_change
FROM monthly_metrics
ORDER BY year, month;

-- ============================================================
-- DIALECT-SPECIFIC EXAMPLES
-- ============================================================

-- MySQL: CTE with backtick identifiers
WITH `sales_data` AS (
    SELECT `product_id`, SUM(`amount`) as `total`
    FROM `sales`
    GROUP BY `product_id`
)
SELECT * FROM `sales_data` WHERE `total` > 1000;

-- PostgreSQL: CTE with double-quote identifiers
WITH "user_metrics" AS (
    SELECT "user_id", COUNT(*) as "activity_count"
    FROM "user_activities"
    GROUP BY "user_id"
)
SELECT * FROM "user_metrics";

-- SQL Server: CTE with bracket identifiers
WITH [department_summary] AS (
    SELECT [dept_id], AVG([salary]) as [avg_salary]
    FROM [employees]
    GROUP BY [dept_id]
)
SELECT * FROM [department_summary];
