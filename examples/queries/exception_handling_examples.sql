-- ====================================================================
-- Exception Handling Examples
-- SQL Parser Go - Comprehensive Examples
-- ====================================================================

-- This file contains real-world examples of exception handling statements
-- across multiple SQL dialects (SQL Server, PostgreSQL, MySQL, Oracle)

-- ====================================================================
-- SQL Server TRY...CATCH
-- ====================================================================

-- Example 1: Simple TRY...CATCH
BEGIN TRY
    SELECT 1/0;
END TRY
BEGIN CATCH
    SELECT ERROR_MESSAGE() AS ErrorMessage;
END CATCH

-- Example 2: TRY...CATCH with error handling
BEGIN TRY
    UPDATE accounts SET balance = 0 WHERE id = 1;
END TRY
BEGIN CATCH
    SELECT
        ERROR_NUMBER() AS ErrorNumber,
        ERROR_MESSAGE() AS ErrorMessage,
        ERROR_SEVERITY() AS ErrorSeverity,
        ERROR_STATE() AS ErrorState;
END CATCH

-- Example 3: Nested TRY...CATCH
BEGIN TRY
    BEGIN TRY
        INSERT INTO users (id, email) VALUES (1, 'test@example.com');
    END TRY
    BEGIN CATCH
        SELECT 'Inner error: ' + ERROR_MESSAGE();
    END CATCH
END TRY
BEGIN CATCH
    SELECT 'Outer error: ' + ERROR_MESSAGE();
END CATCH

-- Example 4: TRY...CATCH with THROW
BEGIN TRY
    UPDATE orders SET status = 'completed' WHERE id = 123;
END TRY
BEGIN CATCH
    THROW 50001, 'Failed to update order status', 1;
END CATCH

-- Example 5: TRY...CATCH with transaction rollback
BEGIN TRY
    UPDATE accounts SET balance = 0 WHERE id = 1;
    UPDATE accounts SET balance = 100 WHERE id = 2;
END TRY
BEGIN CATCH
    ROLLBACK;
    SELECT ERROR_MESSAGE();
END CATCH

-- Example 6: Re-throw exception
BEGIN TRY
    DELETE FROM users WHERE id = 1;
END TRY
BEGIN CATCH
    SELECT ERROR_MESSAGE();
    THROW;
END CATCH

-- ====================================================================
-- SQL Server THROW Statement
-- ====================================================================

-- Example 1: THROW with error details
THROW 50001, 'User not found', 1

-- Example 2: THROW with validation error
THROW 50002, 'Invalid email format', 1

-- Example 3: THROW to re-raise current exception
THROW

-- Example 4: THROW in stored procedure
BEGIN TRY
    SELECT * FROM users WHERE id = 999;
END TRY
BEGIN CATCH
    THROW 50003, 'Error retrieving user data', 1;
END CATCH

-- ====================================================================
-- PostgreSQL EXCEPTION Blocks
-- ====================================================================

-- Example 1: Simple EXCEPTION handler
CREATE OR REPLACE FUNCTION divide(a INT, b INT) RETURNS INT AS $$
BEGIN
    RETURN a / b;
EXCEPTION
    WHEN division_by_zero THEN
        RETURN 0;
END;
$$ LANGUAGE plpgsql;

-- Example 2: Multiple WHEN clauses
CREATE OR REPLACE FUNCTION insert_user(user_email TEXT) RETURNS TEXT AS $$
BEGIN
    INSERT INTO users (email) VALUES (user_email);
    RETURN 'Success';
EXCEPTION
    WHEN unique_violation THEN
        RETURN 'Email already exists';
    WHEN check_violation THEN
        RETURN 'Invalid email format';
    WHEN OTHERS THEN
        RETURN 'Unknown error occurred';
END;
$$ LANGUAGE plpgsql;

-- Example 3: EXCEPTION with RAISE
CREATE OR REPLACE FUNCTION validate_age(age INT) RETURNS VOID AS $$
BEGIN
    IF age < 0 THEN
        RAISE EXCEPTION 'Age cannot be negative';
    END IF;
    IF age > 150 THEN
        RAISE EXCEPTION 'Age too high';
    END IF;
EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE 'Validation error: %', SQLERRM;
END;
$$ LANGUAGE plpgsql;

-- Example 4: Nested exception handling
CREATE OR REPLACE FUNCTION process_order(order_id INT) RETURNS TEXT AS $$
BEGIN
    UPDATE orders SET status = 'processing' WHERE id = order_id;
    RETURN 'Order processed';
EXCEPTION
    WHEN foreign_key_violation THEN
        RETURN 'Invalid order ID';
    WHEN check_violation THEN
        RETURN 'Order validation failed';
    WHEN OTHERS THEN
        RAISE WARNING 'Unexpected error: %', SQLERRM;
        RETURN 'Processing failed';
END;
$$ LANGUAGE plpgsql;

-- ====================================================================
-- PostgreSQL RAISE Statement
-- ====================================================================

-- Example 1: RAISE EXCEPTION
RAISE EXCEPTION 'User not found'

-- Example 2: RAISE NOTICE
RAISE NOTICE 'Processing user ID: %', user_id

-- Example 3: RAISE WARNING
RAISE WARNING 'Low balance detected for account %', account_id

-- Example 4: RAISE with SQLSTATE
RAISE EXCEPTION 'Custom error' USING ERRCODE = '45000'

-- Example 5: RAISE INFO
RAISE INFO 'Starting batch process'

-- Example 6: RAISE DEBUG
RAISE DEBUG 'Debug info: %', debug_value

-- Example 7: Re-raise current exception
RAISE

-- ====================================================================
-- MySQL DECLARE HANDLER
-- ====================================================================

-- Example 1: CONTINUE HANDLER for SQLEXCEPTION
CREATE PROCEDURE handle_errors()
BEGIN
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
    BEGIN
        SELECT 'Error occurred' AS message;
    END;

    INSERT INTO users (id, name) VALUES (1, 'John');
    SELECT 'Operation completed' AS message;
END

-- Example 2: EXIT HANDLER for NOT FOUND
CREATE PROCEDURE fetch_user(IN user_id INT)
BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE user_name VARCHAR(255);
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    SELECT name INTO user_name FROM users WHERE id = user_id;

    IF done = 1 THEN
        SELECT 'User not found' AS message;
    ELSE
        SELECT user_name AS name;
    END IF;
END

-- Example 3: HANDLER for SQLSTATE
CREATE PROCEDURE handle_duplicate()
BEGIN
    DECLARE CONTINUE HANDLER FOR SQLSTATE '23000'
    BEGIN
        SELECT 'Duplicate key error' AS message;
    END;

    INSERT INTO users (id, email) VALUES (1, 'test@example.com');
END

-- Example 4: HANDLER for SQLWARNING
CREATE PROCEDURE handle_warnings()
BEGIN
    DECLARE CONTINUE HANDLER FOR SQLWARNING
    BEGIN
        SELECT 'Warning occurred' AS message;
    END;

    UPDATE users SET name = 'test' WHERE id = 999;
END

-- Example 5: Multiple handlers in one procedure
CREATE PROCEDURE complex_handlers()
BEGIN
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
        SELECT 'SQL Exception' AS error_type;

    DECLARE CONTINUE HANDLER FOR SQLWARNING
        SELECT 'SQL Warning' AS error_type;

    DECLARE EXIT HANDLER FOR NOT FOUND
        SELECT 'Not Found' AS error_type;

    SELECT * FROM users WHERE id = 1;
    INSERT INTO logs (message) VALUES ('Processing complete');
END

-- Example 6: HANDLER with single statement
CREATE PROCEDURE simple_handler()
BEGIN
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
        ROLLBACK;

    INSERT INTO orders (id, status) VALUES (1, 'pending');
END

-- ====================================================================
-- MySQL SIGNAL Statement
-- ====================================================================

-- Example 1: Simple SIGNAL
SIGNAL SQLSTATE '45000'

-- Example 2: SIGNAL with MESSAGE_TEXT
SIGNAL SQLSTATE '45000'
SET MESSAGE_TEXT = 'User not found'

-- Example 3: SIGNAL with multiple properties
SIGNAL SQLSTATE '45000'
SET MESSAGE_TEXT = 'Invalid operation',
    MYSQL_ERRNO = 1525

-- Example 4: SIGNAL in procedure
CREATE PROCEDURE validate_email(IN email VARCHAR(255))
BEGIN
    IF email NOT LIKE '%@%.%' THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Invalid email format';
    END IF;
END

-- Example 5: SIGNAL with custom error code
SIGNAL SQLSTATE '45001'
SET MESSAGE_TEXT = 'Custom business rule violation',
    MYSQL_ERRNO = 9999

-- ====================================================================
-- Real-World Exception Handling Scenarios
-- ====================================================================

-- Example 1: SQL Server - Safe division function
BEGIN TRY
    SELECT 100 / 0 AS result;
END TRY
BEGIN CATCH
    SELECT 0 AS result;
END CATCH

-- Example 2: PostgreSQL - Safe user insertion
CREATE OR REPLACE FUNCTION safe_insert_user(
    user_email TEXT,
    user_name TEXT
) RETURNS TEXT AS $$
BEGIN
    INSERT INTO users (email, name) VALUES (user_email, user_name);
    RETURN 'User created successfully';
EXCEPTION
    WHEN unique_violation THEN
        RETURN 'Email already registered';
    WHEN check_violation THEN
        RETURN 'Invalid user data';
    WHEN foreign_key_violation THEN
        RETURN 'Invalid reference';
    WHEN OTHERS THEN
        RAISE NOTICE 'Unexpected error: %', SQLERRM;
        RETURN 'Failed to create user';
END;
$$ LANGUAGE plpgsql;

-- Example 3: MySQL - Safe update with handler
CREATE PROCEDURE safe_update_balance(
    IN account_id INT,
    IN amount DECIMAL
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Failed to update balance';
    END;

    UPDATE accounts SET balance = amount WHERE id = account_id;
END

-- Example 4: SQL Server - Batch processing with error handling
BEGIN TRY
    UPDATE orders SET status = 'processed' WHERE status = 'pending';
    UPDATE inventory SET quantity = 0 WHERE quantity < 0;
END TRY
BEGIN CATCH
    THROW 50010, 'Batch processing failed', 1;
END CATCH

-- Example 5: PostgreSQL - Multi-step operation with rollback
CREATE OR REPLACE FUNCTION transfer_funds(
    from_account INT,
    to_account INT,
    amount DECIMAL
) RETURNS TEXT AS $$
BEGIN
    UPDATE accounts SET balance = 0 WHERE id = from_account;
    UPDATE accounts SET balance = amount WHERE id = to_account;
    RETURN 'Transfer successful';
EXCEPTION
    WHEN check_violation THEN
        RETURN 'Insufficient funds';
    WHEN foreign_key_violation THEN
        RETURN 'Invalid account';
    WHEN OTHERS THEN
        RAISE WARNING 'Transfer failed: %', SQLERRM;
        RETURN 'Transfer failed';
END;
$$ LANGUAGE plpgsql;

-- Example 6: MySQL - Cursor with NOT FOUND handler
CREATE PROCEDURE process_pending_orders()
BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE order_id INT;

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        SELECT 'Error processing orders' AS error_message;
    END;

    WHILE done = 0 DO
        SELECT id INTO order_id FROM orders WHERE status = 'pending' LIMIT 1;
        IF done = 0 THEN
            UPDATE orders SET status = 'processed' WHERE id = order_id;
        END IF;
    END WHILE;
END

-- Example 7: SQL Server - Validation with custom errors
BEGIN TRY
    SELECT * FROM users WHERE id = 999;
END TRY
BEGIN CATCH
    THROW 50100, 'User validation failed', 1;
END CATCH

-- Example 8: PostgreSQL - Cascading exception handlers
CREATE OR REPLACE FUNCTION delete_user_safe(user_id INT) RETURNS TEXT AS $$
BEGIN
    DELETE FROM users WHERE id = user_id;
    RETURN 'User deleted';
EXCEPTION
    WHEN foreign_key_violation THEN
        RAISE NOTICE 'Cannot delete user with related records';
        RETURN 'Delete failed: has related data';
    WHEN OTHERS THEN
        RAISE WARNING 'Delete operation failed: %', SQLERRM;
        RETURN 'Delete failed';
END;
$$ LANGUAGE plpgsql;

-- Example 9: MySQL - Conditional error signaling
CREATE PROCEDURE validate_order(IN order_total DECIMAL)
BEGIN
    IF order_total < 0 THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Order total cannot be negative';
    END IF;

    IF order_total > 1000000 THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Order total exceeds maximum limit';
    END IF;
END

-- Example 10: SQL Server - Complex error handling workflow
BEGIN TRY
    BEGIN TRY
        UPDATE products SET stock = 0 WHERE id = 123;
    END TRY
    BEGIN CATCH
        SELECT 'Product update failed';
        THROW;
    END CATCH
END TRY
BEGIN CATCH
    SELECT 'Workflow failed: ' + ERROR_MESSAGE();
END CATCH
