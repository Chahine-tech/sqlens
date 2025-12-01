-- SQL Parser Go - TRIGGER Examples
-- Demonstrates CREATE TRIGGER and DROP TRIGGER parsing across different dialects

-- ============================================================================
-- Simple BEFORE/AFTER Triggers
-- ============================================================================

-- BEFORE INSERT trigger
CREATE TRIGGER audit_insert
BEFORE INSERT ON users
FOR EACH ROW
BEGIN
    INSERT INTO audit_log (user_id, action, timestamp)
    VALUES (NEW.id, 'INSERT', NOW());
END;

-- AFTER UPDATE trigger
CREATE TRIGGER update_timestamp
AFTER UPDATE ON products
FOR EACH ROW
BEGIN
    UPDATE products
    SET updated_at = NOW()
    WHERE id = NEW.id;
END;

-- BEFORE DELETE trigger
CREATE TRIGGER prevent_delete
BEFORE DELETE ON important_data
FOR EACH ROW
BEGIN
    SIGNAL SQLSTATE '45000'
    SET MESSAGE_TEXT = 'Cannot delete from important_data';
END;

-- AFTER INSERT trigger
CREATE TRIGGER welcome_email
AFTER INSERT ON users
FOR EACH ROW
BEGIN
    INSERT INTO email_queue (user_id, template, created_at)
    VALUES (NEW.id, 'welcome', NOW());
END;

-- ============================================================================
-- Multiple Event Triggers (INSERT OR UPDATE OR DELETE)
-- ============================================================================

-- Track all changes
CREATE TRIGGER track_changes
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW
BEGIN
    INSERT INTO changes_log (table_name, record_id, changed_at)
    VALUES ('users', NEW.id, NOW());
END;

-- Audit on INSERT or UPDATE
CREATE TRIGGER audit_changes
BEFORE INSERT OR UPDATE ON accounts
FOR EACH ROW
BEGIN
    INSERT INTO audit (table_name, action, timestamp)
    VALUES ('accounts', 'CHANGE', NOW());
END;

-- ============================================================================
-- Triggers with IF NOT EXISTS (MySQL)
-- ============================================================================

CREATE TRIGGER IF NOT EXISTS auto_timestamp
BEFORE INSERT ON orders
FOR EACH ROW
BEGIN
    SET NEW.created_at = NOW();
    SET NEW.updated_at = NOW();
END;

CREATE TRIGGER IF NOT EXISTS set_defaults
BEFORE INSERT ON products
FOR EACH ROW
BEGIN
    IF NEW.status IS NULL THEN
        SET NEW.status = 'active';
    END IF;
    IF NEW.quantity IS NULL THEN
        SET NEW.quantity = 0;
    END IF;
END;

-- ============================================================================
-- OR REPLACE Triggers (PostgreSQL)
-- ============================================================================

CREATE OR REPLACE TRIGGER update_audit
AFTER UPDATE ON accounts
FOR EACH ROW
BEGIN
    INSERT INTO audit_trail (account_id, old_balance, new_balance, changed_at)
    VALUES (NEW.id, OLD.balance, NEW.balance, NOW());
END;

CREATE OR REPLACE TRIGGER validate_email
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
BEGIN
    IF NEW.email NOT LIKE '%@%' THEN
        RAISE EXCEPTION 'Invalid email format';
    END IF;
END;

-- ============================================================================
-- INSTEAD OF Triggers (SQL Server, Oracle - for views)
-- ============================================================================

-- INSTEAD OF INSERT on view
CREATE TRIGGER view_insert_trigger
INSTEAD OF INSERT ON users_view
FOR EACH ROW
BEGIN
    INSERT INTO users (id, name, email, created_at)
    VALUES (NEW.id, NEW.name, NEW.email, NOW());

    INSERT INTO user_preferences (user_id, theme, language)
    VALUES (NEW.id, 'default', 'en');
END;

-- INSTEAD OF UPDATE on view
CREATE TRIGGER view_update_trigger
INSTEAD OF UPDATE ON users_view
FOR EACH ROW
BEGIN
    UPDATE users
    SET name = NEW.name, email = NEW.email, updated_at = NOW()
    WHERE id = NEW.id;
END;

-- INSTEAD OF DELETE on view
CREATE TRIGGER view_delete_trigger
INSTEAD OF DELETE ON users_view
FOR EACH ROW
BEGIN
    UPDATE users
    SET deleted = 1, deleted_at = NOW()
    WHERE id = OLD.id;
END;

-- ============================================================================
-- Triggers with WHEN Conditions (PostgreSQL)
-- ============================================================================

-- Only trigger when price changes
CREATE TRIGGER price_change_audit
AFTER UPDATE ON products
FOR EACH ROW
WHEN (NEW.price <> OLD.price)
BEGIN
    INSERT INTO price_history (product_id, old_price, new_price, changed_at)
    VALUES (NEW.id, OLD.price, NEW.price, NOW());
END;

-- Only trigger for high-value orders
CREATE TRIGGER high_value_alert
AFTER INSERT ON orders
FOR EACH ROW
WHEN (NEW.total > 10000)
BEGIN
    INSERT INTO alerts (type, message, created_at)
    VALUES ('high_value_order', 'Order #' || NEW.id || ' exceeds $10,000', NOW());
END;

-- Only trigger for status changes
CREATE TRIGGER status_change_log
AFTER UPDATE ON orders
FOR EACH ROW
WHEN (NEW.status <> OLD.status)
BEGIN
    INSERT INTO status_log (order_id, old_status, new_status, changed_at)
    VALUES (NEW.id, OLD.status, NEW.status, NOW());
END;

-- ============================================================================
-- FOR EACH STATEMENT Triggers (vs FOR EACH ROW)
-- ============================================================================

-- Statement-level trigger (executes once per SQL statement)
CREATE TRIGGER statement_audit
AFTER DELETE ON users
FOR EACH STATEMENT
BEGIN
    INSERT INTO deletion_log (table_name, deleted_at, deleted_count)
    VALUES ('users', NOW(), (SELECT COUNT(*) FROM deleted_temp));
END;

-- Track bulk operations
CREATE TRIGGER bulk_insert_log
AFTER INSERT ON products
FOR EACH STATEMENT
BEGIN
    INSERT INTO operation_log (operation, timestamp)
    VALUES ('BULK_INSERT_PRODUCTS', NOW());
END;

-- ============================================================================
-- Triggers on Schema-Qualified Tables
-- ============================================================================

CREATE TRIGGER log_changes
AFTER INSERT ON myschema.users
FOR EACH ROW
BEGIN
    INSERT INTO myschema.audit_log (user_id, action, timestamp)
    VALUES (NEW.id, 'INSERT', NOW());
END;

CREATE TRIGGER validate_data
BEFORE UPDATE ON production.orders
FOR EACH ROW
BEGIN
    IF NEW.status NOT IN ('pending', 'processing', 'completed', 'cancelled') THEN
        RAISE EXCEPTION 'Invalid order status';
    END IF;
END;

-- ============================================================================
-- Complex Business Logic Triggers
-- ============================================================================

-- Inventory management
CREATE TRIGGER update_inventory
AFTER INSERT ON order_items
FOR EACH ROW
BEGIN
    UPDATE inventory
    SET quantity = quantity - NEW.quantity,
        last_updated = NOW()
    WHERE product_id = NEW.product_id;

    -- Alert if low stock
    IF (SELECT quantity FROM inventory WHERE product_id = NEW.product_id) < 10 THEN
        INSERT INTO alerts (type, product_id, message)
        VALUES ('low_stock', NEW.product_id, 'Stock below 10 units');
    END IF;
END;

-- Cascade updates
CREATE TRIGGER cascade_email_update
AFTER UPDATE ON users
FOR EACH ROW
WHEN (NEW.email <> OLD.email)
BEGIN
    -- Update all related tables
    UPDATE orders SET customer_email = NEW.email WHERE user_id = NEW.id;
    UPDATE subscriptions SET email = NEW.email WHERE user_id = NEW.id;

    -- Log the change
    INSERT INTO email_change_log (user_id, old_email, new_email, changed_at)
    VALUES (NEW.id, OLD.email, NEW.email, NOW());
END;

-- Referential integrity enforcement
CREATE TRIGGER prevent_orphan_orders
BEFORE DELETE ON customers
FOR EACH ROW
BEGIN
    IF EXISTS (SELECT 1 FROM orders WHERE customer_id = OLD.id AND status = 'pending') THEN
        RAISE EXCEPTION 'Cannot delete customer with pending orders';
    END IF;

    -- Archive completed orders
    INSERT INTO archived_orders
    SELECT * FROM orders WHERE customer_id = OLD.id AND status = 'completed';

    -- Delete completed orders
    DELETE FROM orders WHERE customer_id = OLD.id AND status = 'completed';
END;

-- ============================================================================
-- Audit Trail Triggers
-- ============================================================================

-- Complete audit trail
CREATE TRIGGER audit_trail
AFTER INSERT OR UPDATE OR DELETE ON sensitive_data
FOR EACH ROW
BEGIN
    DECLARE v_action VARCHAR(10);
    DECLARE v_user_id INT;

    -- Determine action
    IF INSERTING THEN
        SET v_action = 'INSERT';
    ELSIF UPDATING THEN
        SET v_action = 'UPDATE';
    ELSIF DELETING THEN
        SET v_action = 'DELETE';
    END IF;

    -- Get current user
    SET v_user_id = (SELECT user_id FROM session WHERE session_id = CURRENT_SESSION);

    -- Insert audit record
    INSERT INTO audit_trail (
        table_name,
        record_id,
        action,
        old_value,
        new_value,
        changed_by,
        changed_at
    ) VALUES (
        'sensitive_data',
        COALESCE(NEW.id, OLD.id),
        v_action,
        OLD.value,
        NEW.value,
        v_user_id,
        NOW()
    );
END;

-- ============================================================================
-- Validation Triggers
-- ============================================================================

-- Email validation
CREATE TRIGGER validate_email_format
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
BEGIN
    IF NEW.email IS NOT NULL AND NEW.email NOT LIKE '%_@_%._%' THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Invalid email format';
    END IF;
END;

-- Business rules validation
CREATE TRIGGER validate_order
BEFORE INSERT OR UPDATE ON orders
FOR EACH ROW
BEGIN
    -- Validate order total
    IF NEW.total < 0 THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Order total cannot be negative';
    END IF;

    -- Validate customer exists
    IF NOT EXISTS (SELECT 1 FROM customers WHERE id = NEW.customer_id) THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Customer does not exist';
    END IF;

    -- Validate order date
    IF NEW.order_date > NOW() THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Order date cannot be in the future';
    END IF;
END;

-- ============================================================================
-- Calculated Fields Triggers
-- ============================================================================

-- Auto-calculate order total
CREATE TRIGGER calculate_order_total
BEFORE INSERT OR UPDATE ON orders
FOR EACH ROW
BEGIN
    SET NEW.total = (
        SELECT SUM(quantity * unit_price)
        FROM order_items
        WHERE order_id = NEW.id
    );

    -- Apply discount if eligible
    IF NEW.total > 1000 THEN
        SET NEW.discount = NEW.total * 0.10;
        SET NEW.total = NEW.total - NEW.discount;
    END IF;
END;

-- Auto-calculate age from birthdate
CREATE TRIGGER calculate_age
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
BEGIN
    SET NEW.age = YEAR(NOW()) - YEAR(NEW.birthdate);
END;

-- ============================================================================
-- DROP TRIGGER Examples
-- ============================================================================

-- Simple DROP
DROP TRIGGER audit_insert;

-- DROP IF EXISTS
DROP TRIGGER IF EXISTS update_timestamp;

-- DROP with schema
DROP TRIGGER myschema.log_changes;

-- ============================================================================
-- Dialect-Specific Examples
-- ============================================================================

-- MySQL with backticks
CREATE TRIGGER `audit_trigger`
BEFORE INSERT ON `users`
FOR EACH ROW
BEGIN
    INSERT INTO `audit_log` (`user_id`, `action`)
    VALUES (NEW.`id`, 'INSERT');
END;

-- PostgreSQL with double quotes
CREATE TRIGGER "audit_trigger"
AFTER UPDATE ON "users"
FOR EACH ROW
BEGIN
    INSERT INTO "audit_log" ("user_id", "action")
    VALUES (NEW."id", 'UPDATE');
END;

-- SQL Server with brackets
CREATE TRIGGER [audit_trigger]
INSTEAD OF DELETE ON [users]
FOR EACH ROW
BEGIN
    INSERT INTO [audit_log] ([user_id], [action])
    VALUES (OLD.[id], 'DELETE');
END;

-- ============================================================================
-- Real-World Use Cases
-- ============================================================================

-- E-commerce: Inventory tracking
CREATE TRIGGER track_inventory_changes
AFTER UPDATE ON products
FOR EACH ROW
WHEN (NEW.stock_quantity <> OLD.stock_quantity)
BEGIN
    INSERT INTO inventory_history (
        product_id,
        old_quantity,
        new_quantity,
        change_amount,
        change_date,
        change_reason
    ) VALUES (
        NEW.id,
        OLD.stock_quantity,
        NEW.stock_quantity,
        NEW.stock_quantity - OLD.stock_quantity,
        NOW(),
        'AUTO_UPDATE'
    );

    -- Send alert if out of stock
    IF NEW.stock_quantity = 0 THEN
        INSERT INTO stock_alerts (product_id, alert_type, created_at)
        VALUES (NEW.id, 'OUT_OF_STOCK', NOW());
    END IF;
END;

-- Banking: Transaction logging
CREATE TRIGGER log_transactions
AFTER INSERT ON transactions
FOR EACH ROW
BEGIN
    -- Update account balance
    UPDATE accounts
    SET balance = balance + NEW.amount,
        last_transaction_date = NOW()
    WHERE account_id = NEW.account_id;

    -- Log for compliance
    INSERT INTO transaction_audit (
        transaction_id,
        account_id,
        amount,
        transaction_type,
        timestamp,
        compliance_flag
    ) VALUES (
        NEW.id,
        NEW.account_id,
        NEW.amount,
        NEW.transaction_type,
        NOW(),
        CASE WHEN ABS(NEW.amount) > 10000 THEN 'REVIEW_REQUIRED' ELSE 'NORMAL' END
    );
END;

-- SaaS: Subscription management
CREATE TRIGGER manage_subscription
AFTER UPDATE ON subscriptions
FOR EACH ROW
WHEN (NEW.status <> OLD.status)
BEGIN
    -- Log status change
    INSERT INTO subscription_history (
        subscription_id,
        old_status,
        new_status,
        changed_at
    ) VALUES (
        NEW.id,
        OLD.status,
        NEW.status,
        NOW()
    );

    -- Handle cancellation
    IF NEW.status = 'cancelled' THEN
        UPDATE users
        SET subscription_active = 0,
            subscription_end_date = NOW()
        WHERE id = NEW.user_id;

        -- Queue exit survey email
        INSERT INTO email_queue (user_id, template, scheduled_for)
        VALUES (NEW.user_id, 'exit_survey', NOW() + INTERVAL 1 DAY);
    END IF;

    -- Handle activation
    IF NEW.status = 'active' AND OLD.status <> 'active' THEN
        UPDATE users
        SET subscription_active = 1,
            subscription_start_date = NOW()
        WHERE id = NEW.user_id;

        -- Queue welcome email
        INSERT INTO email_queue (user_id, template)
        VALUES (NEW.user_id, 'subscription_activated');
    END IF;
END;
