-- ============================================================================
-- MERGE Statement Examples
-- ============================================================================
-- The MERGE statement (also known as UPSERT) performs INSERT, UPDATE, or DELETE
-- operations in a single statement based on matching conditions.
-- Supported by: SQL Server, Oracle, PostgreSQL (INSERT...ON CONFLICT), DB2
-- ============================================================================

-- ============================================================================
-- 1. BASIC MERGE OPERATIONS
-- ============================================================================

-- Example 1.1: Basic MERGE with UPDATE and INSERT
-- Synchronize customer data from a staging table
MERGE INTO customers AS target
USING customer_updates AS source
ON target.customer_id = source.customer_id
WHEN MATCHED THEN
    UPDATE SET
        target.name = source.name,
        target.email = source.email,
        target.updated_at = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN
    INSERT (customer_id, name, email, created_at)
    VALUES (source.customer_id, source.name, source.email, CURRENT_TIMESTAMP);

-- Example 1.2: MERGE without alias keywords
MERGE customers
USING updates
ON customers.id = updates.id
WHEN MATCHED THEN
    UPDATE SET status = updates.status;

-- Example 1.3: MERGE with INTO keyword
MERGE INTO products target
USING new_products source
ON target.product_id = source.product_id
WHEN MATCHED THEN
    UPDATE SET
        target.price = source.price,
        target.stock = source.stock;

-- ============================================================================
-- 2. MERGE WITH SUBQUERIES
-- ============================================================================

-- Example 2.1: Using subquery as source
MERGE INTO inventory
USING (
    SELECT product_id, SUM(quantity) AS total_quantity
    FROM daily_shipments
    WHERE shipment_date = '2024-01-01'
    GROUP BY product_id
) AS shipments
ON inventory.product_id = shipments.product_id
WHEN MATCHED THEN
    UPDATE SET inventory.quantity = inventory.quantity + shipments.total_quantity
WHEN NOT MATCHED THEN
    INSERT (product_id, quantity)
    VALUES (shipments.product_id, shipments.total_quantity);

-- Example 2.2: Complex subquery with JOIN
MERGE INTO user_stats
USING (
    SELECT
        u.user_id,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_amount) AS total_spent
    FROM users u
    LEFT JOIN orders o ON u.user_id = o.user_id
    WHERE o.order_date >= '2024-01-01'
    GROUP BY u.user_id
) AS stats
ON user_stats.user_id = stats.user_id
WHEN MATCHED THEN
    UPDATE SET
        user_stats.order_count = stats.order_count,
        user_stats.total_spent = stats.total_spent,
        user_stats.last_updated = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN
    INSERT (user_id, order_count, total_spent, last_updated)
    VALUES (stats.user_id, stats.order_count, stats.total_spent, CURRENT_TIMESTAMP);

-- ============================================================================
-- 3. MERGE WITH CONDITIONAL LOGIC
-- ============================================================================

-- Example 3.1: Multiple WHEN MATCHED clauses with conditions
-- Different actions based on source data
MERGE INTO inventory
USING product_updates
ON inventory.product_id = product_updates.product_id
WHEN MATCHED AND product_updates.quantity > 0 THEN
    UPDATE SET
        inventory.quantity = product_updates.quantity,
        inventory.last_restocked = CURRENT_TIMESTAMP
WHEN MATCHED AND product_updates.quantity = 0 THEN
    DELETE
WHEN NOT MATCHED THEN
    INSERT (product_id, quantity, last_restocked)
    VALUES (product_updates.product_id, product_updates.quantity, CURRENT_TIMESTAMP);

-- Example 3.2: Conditional UPDATE based on price change
MERGE INTO price_history
USING current_prices
ON price_history.product_id = current_prices.product_id
    AND price_history.is_current = 1
WHEN MATCHED AND price_history.price != current_prices.price THEN
    UPDATE SET
        price_history.is_current = 0,
        price_history.end_date = CURRENT_DATE
WHEN NOT MATCHED THEN
    INSERT (product_id, price, is_current, start_date)
    VALUES (current_prices.product_id, current_prices.price, 1, CURRENT_DATE);

-- Example 3.3: Complex business logic
MERGE INTO customer_tiers
USING (
    SELECT
        customer_id,
        SUM(order_amount) AS total_spent,
        COUNT(*) AS order_count
    FROM orders
    WHERE order_date >= DATEADD(year, -1, GETDATE())
    GROUP BY customer_id
) AS customer_activity
ON customer_tiers.customer_id = customer_activity.customer_id
WHEN MATCHED AND customer_activity.total_spent >= 10000 THEN
    UPDATE SET tier = 'PLATINUM', updated_at = GETDATE()
WHEN MATCHED AND customer_activity.total_spent >= 5000 THEN
    UPDATE SET tier = 'GOLD', updated_at = GETDATE()
WHEN MATCHED AND customer_activity.total_spent >= 1000 THEN
    UPDATE SET tier = 'SILVER', updated_at = GETDATE()
WHEN MATCHED THEN
    UPDATE SET tier = 'BRONZE', updated_at = GETDATE()
WHEN NOT MATCHED AND customer_activity.total_spent >= 1000 THEN
    INSERT (customer_id, tier, created_at, updated_at)
    VALUES (customer_activity.customer_id, 'SILVER', GETDATE(), GETDATE());

-- ============================================================================
-- 4. SQL SERVER SPECIFIC: WHEN NOT MATCHED BY SOURCE
-- ============================================================================

-- Example 4.1: Delete records not in source
-- Remove products that no longer exist in the master catalog
MERGE INTO local_products AS target
USING master_catalog AS source
ON target.product_id = source.product_id
WHEN MATCHED THEN
    UPDATE SET
        target.product_name = source.product_name,
        target.price = source.price
WHEN NOT MATCHED BY TARGET THEN
    INSERT (product_id, product_name, price)
    VALUES (source.product_id, source.product_name, source.price)
WHEN NOT MATCHED BY SOURCE THEN
    DELETE;

-- Example 4.2: Archive instead of delete
MERGE INTO active_users AS target
USING current_active_users AS source
ON target.user_id = source.user_id
WHEN MATCHED THEN
    UPDATE SET target.last_active = source.last_active
WHEN NOT MATCHED BY SOURCE THEN
    UPDATE SET target.is_archived = 1, target.archived_at = GETDATE();

-- ============================================================================
-- 5. ETL AND DATA SYNCHRONIZATION
-- ============================================================================

-- Example 5.1: Daily sales data synchronization
MERGE INTO sales_summary AS target
USING (
    SELECT
        DATE(order_date) AS sale_date,
        product_id,
        SUM(quantity) AS total_quantity,
        SUM(amount) AS total_amount
    FROM orders
    WHERE DATE(order_date) = CURRENT_DATE
    GROUP BY DATE(order_date), product_id
) AS daily_sales
ON target.sale_date = daily_sales.sale_date
    AND target.product_id = daily_sales.product_id
WHEN MATCHED THEN
    UPDATE SET
        target.total_quantity = daily_sales.total_quantity,
        target.total_amount = daily_sales.total_amount,
        target.last_updated = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN
    INSERT (sale_date, product_id, total_quantity, total_amount, last_updated)
    VALUES (
        daily_sales.sale_date,
        daily_sales.product_id,
        daily_sales.total_quantity,
        daily_sales.total_amount,
        CURRENT_TIMESTAMP
    );

-- Example 5.2: Slowly Changing Dimension (Type 2)
-- Track historical changes by creating new records
MERGE INTO dim_customer AS target
USING stg_customer AS source
ON target.customer_id = source.customer_id
    AND target.is_current = 1
WHEN MATCHED AND (
    target.name != source.name OR
    target.address != source.address OR
    target.email != source.email
) THEN
    UPDATE SET
        target.is_current = 0,
        target.end_date = CURRENT_DATE
WHEN NOT MATCHED BY TARGET THEN
    INSERT (customer_id, name, address, email, is_current, start_date, end_date)
    VALUES (
        source.customer_id,
        source.name,
        source.address,
        source.email,
        1,
        CURRENT_DATE,
        '9999-12-31'
    );

-- Example 5.3: Incremental data warehouse load
MERGE INTO fact_orders AS target
USING (
    SELECT
        order_id,
        customer_id,
        product_id,
        order_date,
        quantity,
        amount,
        CHECKSUM(customer_id, product_id, order_date, quantity, amount) AS row_hash
    FROM staging_orders
) AS source
ON target.order_id = source.order_id
WHEN MATCHED AND target.row_hash != source.row_hash THEN
    UPDATE SET
        target.customer_id = source.customer_id,
        target.product_id = source.product_id,
        target.order_date = source.order_date,
        target.quantity = source.quantity,
        target.amount = source.amount,
        target.row_hash = source.row_hash,
        target.updated_at = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN
    INSERT (order_id, customer_id, product_id, order_date, quantity, amount, row_hash, created_at)
    VALUES (
        source.order_id,
        source.customer_id,
        source.product_id,
        source.order_date,
        source.quantity,
        source.amount,
        source.row_hash,
        CURRENT_TIMESTAMP
    );

-- ============================================================================
-- 6. REAL-WORLD SCENARIOS
-- ============================================================================

-- Example 6.1: Product inventory reconciliation
MERGE INTO inventory AS target
USING physical_count AS source
ON target.warehouse_id = source.warehouse_id
    AND target.product_id = source.product_id
WHEN MATCHED AND target.quantity != source.counted_quantity THEN
    UPDATE SET
        target.quantity = source.counted_quantity,
        target.last_count_date = CURRENT_DATE,
        target.adjustment_amount = source.counted_quantity - target.quantity
WHEN NOT MATCHED BY TARGET THEN
    INSERT (warehouse_id, product_id, quantity, last_count_date)
    VALUES (source.warehouse_id, source.product_id, source.counted_quantity, CURRENT_DATE)
WHEN NOT MATCHED BY SOURCE THEN
    DELETE;

-- Example 6.2: Employee salary updates
MERGE INTO employee_salaries AS target
USING salary_adjustments AS source
ON target.employee_id = source.employee_id
    AND target.is_current = 1
WHEN MATCHED AND source.new_salary != target.salary THEN
    UPDATE SET
        target.is_current = 0,
        target.end_date = source.effective_date
WHEN NOT MATCHED BY TARGET THEN
    INSERT (employee_id, salary, is_current, start_date, end_date)
    VALUES (source.employee_id, source.new_salary, 1, source.effective_date, '9999-12-31');

-- Example 6.3: Customer preferences update
MERGE INTO customer_preferences AS target
USING new_preferences AS source
ON target.customer_id = source.customer_id
    AND target.preference_key = source.preference_key
WHEN MATCHED THEN
    UPDATE SET
        target.preference_value = source.preference_value,
        target.updated_at = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN
    INSERT (customer_id, preference_key, preference_value, created_at, updated_at)
    VALUES (
        source.customer_id,
        source.preference_key,
        source.preference_value,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );
