-- ====================================================================
-- CASE Expression Examples
-- SQL Parser Go - Comprehensive Examples
-- ====================================================================

-- This file contains real-world examples of CASE expressions
-- across multiple SQL dialects (MySQL, PostgreSQL, SQL Server, Oracle)

-- ====================================================================
-- Simple CASE Expressions (CASE value WHEN...)
-- ====================================================================

-- Example 1: Basic status mapping
SELECT
    user_id,
    username,
    CASE status
        WHEN 'A' THEN 'Active'
        WHEN 'I' THEN 'Inactive'
        WHEN 'P' THEN 'Pending'
        ELSE 'Unknown'
    END AS status_description
FROM users;

-- Example 2: Numeric tier assignment
SELECT
    product_id,
    CASE category_id
        WHEN 1 THEN 'Electronics'
        WHEN 2 THEN 'Clothing'
        WHEN 3 THEN 'Food'
        WHEN 4 THEN 'Books'
        ELSE 'Other'
    END AS category_name
FROM products;

-- Example 3: Simple CASE in WHERE clause
SELECT * FROM orders
WHERE CASE payment_status
    WHEN 'paid' THEN 1
    WHEN 'processing' THEN 1
    ELSE 0
END = 1;

-- Example 4: Simple CASE with column alias
SELECT
    order_id,
    CASE order_type
        WHEN 'express' THEN 'Express Delivery'
        WHEN 'standard' THEN 'Standard Delivery'
        WHEN 'pickup' THEN 'Store Pickup'
    END AS delivery_method
FROM orders;

-- ====================================================================
-- Searched CASE Expressions (CASE WHEN condition...)
-- ====================================================================

-- Example 1: Age group classification
SELECT
    user_id,
    name,
    age,
    CASE
        WHEN age < 13 THEN 'Child'
        WHEN age < 18 THEN 'Teenager'
        WHEN age < 65 THEN 'Adult'
        ELSE 'Senior'
    END AS age_group
FROM users;

-- Example 2: Salary range categorization
SELECT
    employee_id,
    name,
    salary,
    CASE
        WHEN salary < 30000 THEN 'Entry Level'
        WHEN salary < 60000 THEN 'Mid Level'
        WHEN salary < 100000 THEN 'Senior Level'
        ELSE 'Executive Level'
    END AS salary_bracket
FROM employees;

-- Example 3: Grade calculation
SELECT
    student_id,
    name,
    score,
    CASE
        WHEN score >= 90 THEN 'A'
        WHEN score >= 80 THEN 'B'
        WHEN score >= 70 THEN 'C'
        WHEN score >= 60 THEN 'D'
        ELSE 'F'
    END AS grade
FROM student_scores;

-- Example 4: Multiple condition checks
SELECT
    order_id,
    CASE
        WHEN order_date > '2024-01-01' THEN 'Recent Order'
        WHEN order_date > '2023-01-01' THEN 'Last Year'
        WHEN order_date > '2022-01-01' THEN 'Old Order'
        ELSE 'Very Old'
    END AS order_age
FROM orders;

-- ====================================================================
-- CASE in Different SQL Clauses
-- ====================================================================

-- Example 1: CASE in SELECT with multiple columns
SELECT
    user_id,
    username,
    CASE status WHEN 'active' THEN 'Active User' ELSE 'Inactive' END AS status_desc,
    CASE type WHEN 'admin' THEN 'Administrator' WHEN 'user' THEN 'Regular User' ELSE 'Guest' END AS user_type
FROM users;

-- Example 2: CASE in WHERE clause for conditional filtering
SELECT * FROM products
WHERE CASE category
    WHEN 'electronics' THEN price > 100
    WHEN 'clothing' THEN price > 50
    ELSE price > 10
END;

-- Example 3: CASE in ORDER BY for custom sorting
SELECT product_name, priority
FROM products
ORDER BY CASE priority
    WHEN 'high' THEN 1
    WHEN 'medium' THEN 2
    WHEN 'low' THEN 3
    ELSE 4
END;

-- Example 4: CASE in GROUP BY
SELECT
    CASE
        WHEN age < 25 THEN 'Young'
        WHEN age < 50 THEN 'Middle Aged'
        ELSE 'Senior'
    END AS age_group,
    COUNT(*) as user_count
FROM users
GROUP BY CASE
    WHEN age < 25 THEN 'Young'
    WHEN age < 50 THEN 'Middle Aged'
    ELSE 'Senior'
END;

-- Example 5: CASE in HAVING clause
SELECT
    category,
    COUNT(*) as product_count
FROM products
GROUP BY category
HAVING COUNT(*) > CASE category
    WHEN 'popular' THEN 100
    WHEN 'standard' THEN 50
    ELSE 10
END;

-- ====================================================================
-- Nested CASE Expressions
-- ====================================================================

-- Example 1: CASE within CASE for complex logic
SELECT
    user_id,
    CASE subscription_type
        WHEN 'premium' THEN
            CASE status
                WHEN 'active' THEN 'Premium Active'
                WHEN 'trial' THEN 'Premium Trial'
                ELSE 'Premium Inactive'
            END
        WHEN 'basic' THEN
            CASE status
                WHEN 'active' THEN 'Basic Active'
                ELSE 'Basic Inactive'
            END
        ELSE 'No Subscription'
    END AS user_tier
FROM users;

-- Example 2: Multiple CASE expressions in one SELECT
SELECT
    order_id,
    CASE status WHEN 'completed' THEN 'Done' ELSE 'Pending' END AS order_status,
    CASE payment_method WHEN 'credit' THEN 'CC' WHEN 'debit' THEN 'DC' ELSE 'Other' END AS payment_type,
    CASE
        WHEN total > 1000 THEN 'Large'
        WHEN total > 100 THEN 'Medium'
        ELSE 'Small'
    END AS order_size
FROM orders;

-- ====================================================================
-- CASE with Functions
-- ====================================================================

-- Example 1: Function in WHEN condition
SELECT
    name,
    CASE
        WHEN LENGTH(name) > 20 THEN 'Very Long Name'
        WHEN LENGTH(name) > 10 THEN 'Long Name'
        ELSE 'Short Name'
    END AS name_length_category
FROM users;

-- Example 2: Function in THEN result
SELECT
    user_id,
    CASE status
        WHEN 'active' THEN UPPER(username)
        WHEN 'inactive' THEN LOWER(username)
        ELSE username
    END AS formatted_username
FROM users;

-- Example 3: CASE wrapped in function
SELECT
    order_id,
    UPPER(CASE status
        WHEN 'pending' THEN 'waiting'
        WHEN 'shipped' THEN 'in transit'
        WHEN 'delivered' THEN 'completed'
        ELSE 'unknown'
    END) AS status_uppercase
FROM orders;

-- Example 4: CASE with aggregate functions
SELECT
    category,
    SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) AS active_count,
    SUM(CASE WHEN status = 'inactive' THEN 1 ELSE 0 END) AS inactive_count
FROM products
GROUP BY category;

-- ====================================================================
-- Real-World Complex Scenarios
-- ====================================================================

-- Example 1: E-commerce shipping cost calculation
SELECT
    order_id,
    total_amount,
    CASE
        WHEN total_amount > 100 THEN 0
        WHEN shipping_method = 'express' THEN 15
        WHEN shipping_method = 'standard' THEN 5
        ELSE 10
    END AS shipping_cost
FROM orders;

-- Example 2: Customer segmentation for marketing
SELECT
    customer_id,
    name,
    total_purchases,
    last_purchase_date,
    CASE
        WHEN total_purchases > 10000 THEN 'VIP'
        WHEN total_purchases > 5000 THEN 'Gold'
        WHEN total_purchases > 1000 THEN 'Silver'
        WHEN last_purchase_date > '2024-01-01' THEN 'Active Bronze'
        ELSE 'Inactive'
    END AS customer_tier
FROM customers;

-- Example 3: Employee bonus calculation
SELECT
    employee_id,
    name,
    performance_score,
    years_of_service,
    salary,
    CASE
        WHEN performance_score >= 9 THEN salary * 0.20
        WHEN performance_score >= 8 THEN salary * 0.15
        WHEN performance_score >= 7 THEN salary * 0.10
        WHEN performance_score >= 6 THEN salary * 0.05
        ELSE 0
    END AS bonus_amount
FROM employees;

-- Example 4: Product availability status
SELECT
    product_id,
    product_name,
    stock_quantity,
    reserved_quantity,
    CASE
        WHEN stock_quantity = 0 THEN 'Out of Stock'
        WHEN stock_quantity <= reserved_quantity THEN 'Reserved'
        WHEN stock_quantity < 10 THEN 'Low Stock'
        WHEN stock_quantity < 50 THEN 'In Stock'
        ELSE 'Well Stocked'
    END AS availability_status
FROM inventory;

-- Example 5: Order priority assignment
SELECT
    order_id,
    customer_tier,
    order_value,
    order_date,
    CASE
        WHEN customer_tier = 'VIP' THEN 'Urgent'
        WHEN order_value > 1000 THEN 'High'
        WHEN order_date > CURRENT_DATE THEN 'Normal'
        ELSE 'Low'
    END AS processing_priority
FROM orders;

-- Example 6: Dynamic discount calculation
SELECT
    product_id,
    price,
    category,
    CASE category
        WHEN 'electronics' THEN
            CASE
                WHEN price > 1000 THEN price * 0.85
                WHEN price > 500 THEN price * 0.90
                ELSE price * 0.95
            END
        WHEN 'clothing' THEN
            CASE
                WHEN price > 100 THEN price * 0.80
                ELSE price * 0.90
            END
        ELSE price * 0.95
    END AS discounted_price
FROM products;

-- Example 7: SaaS subscription tier determination
SELECT
    user_id,
    plan_type,
    usage_count,
    CASE plan_type
        WHEN 'free' THEN
            CASE
                WHEN usage_count >= 100 THEN 'Upgrade Required'
                WHEN usage_count >= 80 THEN 'Approaching Limit'
                ELSE 'Within Limit'
            END
        WHEN 'pro' THEN
            CASE
                WHEN usage_count >= 1000 THEN 'Consider Enterprise'
                ELSE 'Within Limit'
            END
        WHEN 'enterprise' THEN 'Unlimited'
        ELSE 'Unknown Plan'
    END AS usage_status
FROM subscriptions;

-- Example 8: Data quality classification
SELECT
    record_id,
    CASE
        WHEN email IS NULL THEN 'Missing Email'
        WHEN phone IS NULL THEN 'Missing Phone'
        WHEN address IS NULL THEN 'Missing Address'
        WHEN email IS NULL THEN 'Complete'
        ELSE 'Complete'
    END AS data_completeness
FROM contacts;

-- Example 9: Tax rate calculation by region
SELECT
    order_id,
    region,
    subtotal,
    CASE region
        WHEN 'CA' THEN subtotal * 0.0725
        WHEN 'NY' THEN subtotal * 0.08
        WHEN 'TX' THEN subtotal * 0.0625
        WHEN 'FL' THEN 0
        ELSE subtotal * 0.05
    END AS tax_amount
FROM orders;

-- Example 10: Performance metric categorization
SELECT
    metric_name,
    metric_value,
    CASE metric_name
        WHEN 'response_time' THEN
            CASE
                WHEN metric_value < 100 THEN 'Excellent'
                WHEN metric_value < 500 THEN 'Good'
                WHEN metric_value < 1000 THEN 'Fair'
                ELSE 'Poor'
            END
        WHEN 'error_rate' THEN
            CASE
                WHEN metric_value < 0.01 THEN 'Excellent'
                WHEN metric_value < 0.05 THEN 'Good'
                WHEN metric_value < 0.10 THEN 'Fair'
                ELSE 'Poor'
            END
        ELSE 'Unknown Metric'
    END AS performance_rating
FROM system_metrics;
