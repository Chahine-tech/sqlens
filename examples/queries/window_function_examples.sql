-- SQL Parser Go - Window Function Examples
-- These examples demonstrate OVER clause and window functions

-- ============================================================
-- ROW_NUMBER() - Sequential numbering
-- ============================================================
-- Simple row numbering
SELECT
    employee_id,
    name,
    department,
    ROW_NUMBER() OVER (ORDER BY employee_id) as row_num
FROM employees;

-- Row numbering within partitions
SELECT
    employee_id,
    name,
    department,
    salary,
    ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) as dept_rank
FROM employees;

-- ============================================================
-- RANK() and DENSE_RANK()
-- ============================================================
-- Ranking with gaps (RANK)
SELECT
    student_id,
    name,
    score,
    RANK() OVER (ORDER BY score DESC) as rank,
    DENSE_RANK() OVER (ORDER BY score DESC) as dense_rank
FROM students;

-- Ranking within groups
SELECT
    product_id,
    category,
    sales_amount,
    RANK() OVER (PARTITION BY category ORDER BY sales_amount DESC) as category_rank
FROM product_sales;

-- ============================================================
-- AGGREGATE WINDOW FUNCTIONS
-- ============================================================
-- Running totals
SELECT
    order_date,
    order_amount,
    SUM(order_amount) OVER (ORDER BY order_date) as running_total
FROM orders
ORDER BY order_date;

-- Moving averages
SELECT
    sale_date,
    daily_revenue,
    AVG(daily_revenue) OVER (
        ORDER BY sale_date
        ROWS BETWEEN 6 PRECEDING AND CURRENT ROW
    ) as seven_day_avg
FROM daily_sales;

-- ============================================================
-- PARTITION BY with ORDER BY
-- ============================================================
-- Department-wise rankings
SELECT
    employee_id,
    name,
    department,
    salary,
    RANK() OVER (PARTITION BY department ORDER BY salary DESC) as dept_salary_rank,
    AVG(salary) OVER (PARTITION BY department) as dept_avg_salary
FROM employees;

-- Multiple window functions
SELECT
    product_id,
    sale_date,
    quantity,
    ROW_NUMBER() OVER (PARTITION BY product_id ORDER BY sale_date) as sale_sequence,
    SUM(quantity) OVER (PARTITION BY product_id ORDER BY sale_date) as cumulative_qty,
    AVG(quantity) OVER (PARTITION BY product_id) as avg_qty
FROM product_sales;

-- ============================================================
-- WINDOW FRAMES - ROWS
-- ============================================================
-- Fixed window size
SELECT
    date,
    value,
    AVG(value) OVER (
        ORDER BY date
        ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING
    ) as moving_avg_5
FROM measurements;

-- From start to current
SELECT
    month,
    revenue,
    SUM(revenue) OVER (
        ORDER BY month
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) as ytd_revenue
FROM monthly_revenue;

-- Custom frame bounds
SELECT
    timestamp,
    sensor_value,
    MIN(sensor_value) OVER (
        ORDER BY timestamp
        ROWS BETWEEN 10 PRECEDING AND 5 PRECEDING
    ) as recent_min
FROM sensor_data;

-- ============================================================
-- WINDOW FRAMES - RANGE
-- ============================================================
-- Range-based window
SELECT
    employee_id,
    hire_date,
    salary,
    AVG(salary) OVER (
        ORDER BY hire_date
        RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) as avg_salary_to_date
FROM employees;

-- ============================================================
-- ADVANCED EXAMPLES
-- ============================================================
-- Finding top N per group
SELECT * FROM (
    SELECT
        department,
        employee_name,
        salary,
        RANK() OVER (PARTITION BY department ORDER BY salary DESC) as rank
    FROM employees
) ranked
WHERE rank <= 3;

-- Calculating percentiles
SELECT
    product_id,
    price,
    PERCENT_RANK() OVER (ORDER BY price) as price_percentile,
    NTILE(4) OVER (ORDER BY price) as price_quartile
FROM products;

-- Gap and island detection
SELECT
    user_id,
    activity_date,
    activity_date - ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY activity_date) as island_id
FROM user_activities;

-- ============================================================
-- MULTIPLE WINDOW SPECIFICATIONS
-- ============================================================
SELECT
    employee_id,
    department,
    salary,
    -- Overall ranking
    ROW_NUMBER() OVER (ORDER BY salary DESC) as overall_rank,
    -- Department ranking
    ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) as dept_rank,
    -- Department statistics
    AVG(salary) OVER (PARTITION BY department) as dept_avg,
    MAX(salary) OVER (PARTITION BY department) as dept_max,
    -- Percentile within department
    PERCENT_RANK() OVER (PARTITION BY department ORDER BY salary) as dept_percentile
FROM employees;

-- ============================================================
-- PRACTICAL USE CASES
-- ============================================================
-- Customer lifetime value analysis
SELECT
    customer_id,
    order_date,
    order_amount,
    SUM(order_amount) OVER (
        PARTITION BY customer_id
        ORDER BY order_date
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) as customer_ltv,
    AVG(order_amount) OVER (
        PARTITION BY customer_id
        ORDER BY order_date
        ROWS BETWEEN 3 PRECEDING AND CURRENT ROW
    ) as recent_avg_order
FROM orders;

-- Time-series analysis with window functions
SELECT
    metric_date,
    metric_value,
    metric_value - LAG(metric_value) OVER (ORDER BY metric_date) as daily_change,
    (metric_value - LAG(metric_value) OVER (ORDER BY metric_date)) * 100.0 /
        LAG(metric_value) OVER (ORDER BY metric_date) as pct_change,
    AVG(metric_value) OVER (
        ORDER BY metric_date
        ROWS BETWEEN 6 PRECEDING AND CURRENT ROW
    ) as seven_day_ma
FROM daily_metrics;
