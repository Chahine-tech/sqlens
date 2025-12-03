-- ====================================================================
-- Control Flow Statements Examples
-- SQL Parser Go - Comprehensive Examples
-- ====================================================================

-- This file contains real-world examples of control flow statements
-- across multiple SQL dialects (MySQL, PostgreSQL, SQL Server, Oracle)

-- ====================================================================
-- IF...THEN...ELSE...END IF
-- ====================================================================

-- Example 1: Simple IF statement (MySQL)
IF user_status = 'active' THEN
    RETURN 'User is active';
END IF;

-- Example 2: IF with ELSE (MySQL)
IF balance > 1000 THEN
    RETURN 'Premium customer';
ELSE
    RETURN 'Standard customer';
END IF;

-- Example 3: IF with ELSEIF (MySQL)
IF score >= 90 THEN
    RETURN 'A';
ELSEIF score >= 80 THEN
    RETURN 'B';
ELSEIF score >= 70 THEN
    RETURN 'C';
ELSEIF score >= 60 THEN
    RETURN 'D';
ELSE
    RETURN 'F';
END IF;

-- Example 4: PostgreSQL ELSIF syntax
IF temperature > 30 THEN
    RETURN 'Hot';
ELSIF temperature > 20 THEN
    RETURN 'Warm';
ELSIF temperature > 10 THEN
    RETURN 'Cool';
ELSE
    RETURN 'Cold';
END IF;

-- Example 5: Nested IF statements
IF account_type = 'business' THEN
    IF revenue > 100000 THEN
        RETURN 'Enterprise';
    ELSE
        RETURN 'SMB';
    END IF;
ELSE
    RETURN 'Personal';
END IF;

-- Example 6: Complex business logic
IF subscription_tier = 'free' THEN
    IF usage_count >= free_limit THEN
        RETURN 'Upgrade required';
    ELSE
        RETURN 'Continue using';
    END IF;
ELSEIF subscription_tier = 'pro' THEN
    IF usage_count >= pro_limit THEN
        RETURN 'Approaching limit';
    ELSE
        RETURN 'Plenty of quota';
    END IF;
ELSE
    RETURN 'Unlimited access';
END IF;

-- ====================================================================
-- WHILE...DO...END WHILE
-- ====================================================================

-- Example 1: Simple WHILE loop (MySQL)
WHILE counter < 10 DO
    RETURN counter + 1;
END WHILE;

-- Example 2: WHILE with condition check
WHILE balance > minimum_balance DO
    RETURN balance - withdrawal_amount;
    EXIT WHEN balance <= minimum_balance;
END WHILE;

-- Example 3: Processing batches
WHILE batch_count < total_batches DO
    RETURN batch_count + 1;
    IF batch_count >= max_concurrent THEN
        EXIT;
    END IF;
END WHILE;

-- Example 4: Nested WHILE loops
WHILE row_index < max_rows DO
    WHILE col_index < max_cols DO
        RETURN col_index + 1;
    END WHILE;
    RETURN row_index + 1;
END WHILE;

-- Example 5: WHILE with complex condition
WHILE retry_count < max_retries AND status <> 'success' DO
    RETURN retry_count + 1;
    IF status = 'success' THEN
        EXIT;
    END IF;
END WHILE;

-- ====================================================================
-- FOR...LOOP (PostgreSQL style)
-- ====================================================================

-- Example 1: Simple FOR loop
FOR i IN 1..10 LOOP
    RETURN i * 2;
END LOOP;

-- Example 2: FOR with REVERSE
FOR i IN REVERSE 10..1 LOOP
    RETURN countdown_value - i;
END LOOP;

-- Example 3: FOR with step (BY clause)
FOR i IN 0..100 BY 10 LOOP
    RETURN percentage + i;
END LOOP;

-- Example 4: FOR loop with early exit
FOR month IN 1..12 LOOP
    IF month = target_month THEN
        RETURN month;
        EXIT;
    END IF;
END LOOP;

-- Example 5: Nested FOR loops (matrix traversal)
FOR row IN 1..grid_rows LOOP
    FOR col IN 1..grid_cols LOOP
        RETURN row * grid_cols + col;
    END LOOP;
END LOOP;

-- Example 6: FOR with complex calculation
FOR year IN start_year..end_year LOOP
    IF year % 4 = 0 THEN
        RETURN 'Leap year';
    ELSE
        RETURN 'Regular year';
    END IF;
END LOOP;

-- ====================================================================
-- LOOP...END LOOP (infinite loop with EXIT)
-- ====================================================================

-- Example 1: Simple LOOP with EXIT WHEN
LOOP
    RETURN counter + 1;
    EXIT WHEN counter >= 10;
END LOOP;

-- Example 2: LOOP with unconditional EXIT
LOOP
    RETURN attempts + 1;
    IF attempts > max_attempts THEN
        EXIT;
    END IF;
END LOOP;

-- Example 3: LOOP with multiple exit conditions
LOOP
    RETURN iteration + 1;
    EXIT WHEN iteration >= max_iterations;
    EXIT WHEN status = 'complete';
    EXIT WHEN error_count > threshold;
END LOOP;

-- Example 4: Nested LOOPs
LOOP
    LOOP
        RETURN inner_count + 1;
        EXIT WHEN inner_count >= 5;
    END LOOP;
    EXIT WHEN outer_count >= 3;
END LOOP;

-- Example 5: LOOP for retry logic
LOOP
    RETURN retry_attempt + 1;
    IF operation_successful THEN
        EXIT;
    END IF;
    IF retry_attempt >= max_retries THEN
        RETURN 'Failed after max retries';
        EXIT;
    END IF;
END LOOP;

-- ====================================================================
-- REPEAT...UNTIL (MySQL)
-- ====================================================================

-- Example 1: Simple REPEAT UNTIL
REPEAT
    RETURN counter + 1;
UNTIL counter >= 10;

-- Example 2: REPEAT with multiple statements
REPEAT
    RETURN current_value + increment;
    IF current_value > threshold THEN
        EXIT;
    END IF;
UNTIL current_value >= max_value;

-- Example 3: REPEAT for batch processing
REPEAT
    RETURN batch_processed + 1;
    IF batch_processed % 100 = 0 THEN
        RETURN 'Checkpoint reached';
    END IF;
UNTIL batch_processed >= total_batches;

-- Example 4: REPEAT with validation
REPEAT
    RETURN validation_attempt + 1;
    IF data_valid THEN
        RETURN 'Validation successful';
        EXIT;
    END IF;
UNTIL validation_attempt >= max_validation_attempts;

-- ====================================================================
-- EXIT and CONTINUE statements
-- ====================================================================

-- Example 1: EXIT with condition
LOOP
    RETURN counter + 1;
    EXIT WHEN counter > 100;
END LOOP;

-- Example 2: Unconditional EXIT
WHILE processing DO
    IF error_occurred THEN
        EXIT;
    END IF;
END WHILE;

-- Example 3: CONTINUE (skip to next iteration)
LOOP
    RETURN item_index + 1;
    CONTINUE WHEN item_index % 2 = 0;  -- Skip even numbers
    RETURN 'Processing odd number';
    EXIT WHEN item_index >= 100;
END LOOP;

-- Example 4: ITERATE (MySQL synonym for CONTINUE)
WHILE batch_index < total_batches DO
    RETURN batch_index + 1;
    IF batch_index % skip_interval = 0 THEN
        ITERATE;  -- Skip this batch
    END IF;
    RETURN 'Processing batch';
END WHILE;

-- ====================================================================
-- Real-World Complex Examples
-- ====================================================================

-- Example 1: User subscription tier calculation
IF account_age_days < 30 THEN
    IF trial_used THEN
        RETURN 'Trial expired - upgrade required';
    ELSE
        RETURN 'Trial active';
    END IF;
ELSEIF monthly_spend > 1000 THEN
    RETURN 'Enterprise tier';
ELSEIF monthly_spend > 100 THEN
    RETURN 'Pro tier';
ELSE
    RETURN 'Free tier';
END IF;

-- Example 2: Inventory replenishment logic
WHILE stock_level < reorder_point DO
    IF supplier_available THEN
        RETURN order_quantity;
        RETURN stock_level + order_quantity;
    ELSE
        RETURN 'Waiting for supplier';
        EXIT;
    END IF;
END WHILE;

-- Example 3: Data migration with batching
FOR batch_num IN 1..total_batches LOOP
    IF batch_num % 10 = 0 THEN
        RETURN 'Checkpoint - committing batch';
    END IF;

    IF error_count > error_threshold THEN
        RETURN 'Too many errors - stopping migration';
        EXIT;
    END IF;
END LOOP;

-- Example 4: Rate limiting with retry
REPEAT
    RETURN api_call_attempt + 1;

    IF rate_limit_exceeded THEN
        IF retry_count < max_retries THEN
            RETURN retry_count + 1;
            CONTINUE;
        ELSE
            RETURN 'Rate limit exceeded - max retries reached';
            EXIT;
        END IF;
    END IF;

    IF api_call_successful THEN
        RETURN 'API call succeeded';
        EXIT;
    END IF;
UNTIL retry_count >= max_retries;

-- Example 5: SaaS billing calculation
IF subscription_status = 'active' THEN
    FOR month IN billing_start_month..billing_end_month LOOP
        IF month = current_month THEN
            IF proration_needed THEN
                RETURN base_price * proration_days / days_in_month;
            ELSE
                RETURN base_price;
            END IF;
        END IF;
    END LOOP;
ELSIF subscription_status = 'cancelled' THEN
    RETURN 0;
ELSE
    RETURN 'Invalid subscription status';
END IF;

-- Example 6: Multi-level approval workflow
IF request_amount <= manager_limit THEN
    RETURN 'Manager approval required';
ELSIF request_amount <= director_limit THEN
    RETURN 'Director approval required';
ELSIF request_amount <= vp_limit THEN
    RETURN 'VP approval required';
ELSE
    RETURN 'Board approval required';
END IF;

-- Example 7: Data validation pipeline
LOOP
    RETURN validation_step + 1;

    IF validation_step = 1 THEN
        IF NOT format_valid THEN
            RETURN 'Format validation failed';
            EXIT;
        END IF;
    ELSIF validation_step = 2 THEN
        IF NOT schema_valid THEN
            RETURN 'Schema validation failed';
            EXIT;
        END IF;
    ELSIF validation_step = 3 THEN
        IF NOT business_rules_valid THEN
            RETURN 'Business rules validation failed';
            EXIT;
        END IF;
    END IF;

    EXIT WHEN validation_step >= 3;
END LOOP;

-- Example 8: A/B testing distribution
IF experiment_enabled THEN
    IF user_id % 2 = 0 THEN
        RETURN 'Variant A';
    ELSE
        RETURN 'Variant B';
    END IF;
ELSE
    RETURN 'Control';
END IF;

-- Example 9: Cache warm-up with backoff
REPEAT
    RETURN cache_load_attempt + 1;

    IF cache_hit_rate >= target_hit_rate THEN
        RETURN 'Cache warmed up successfully';
        EXIT;
    END IF;

    IF cache_load_attempt > 1 THEN
        RETURN backoff_seconds * 2;  -- Exponential backoff
    END IF;
UNTIL cache_load_attempt >= max_cache_attempts;

-- Example 10: Feature flag evaluation
IF feature_rollout_percentage >= 100 THEN
    RETURN 'Feature enabled for all users';
ELSIF feature_rollout_percentage > 0 THEN
    IF user_id % 100 < feature_rollout_percentage THEN
        RETURN 'Feature enabled';
    ELSE
        RETURN 'Feature disabled';
    END IF;
ELSE
    RETURN 'Feature disabled for all users';
END IF;
