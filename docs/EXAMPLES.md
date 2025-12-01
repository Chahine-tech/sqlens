# SQL Parser Go - Usage Examples

This document provides comprehensive usage examples for SQL Parser Go across all supported features and SQL dialects.

## Table of Contents

1. [Basic Query Analysis](#basic-query-analysis)
2. [Multi-Dialect Support](#multi-dialect-support)
3. [DML Statements](#dml-statements)
4. [Subqueries](#subqueries)
5. [DDL Support](#ddl-support)
6. [Transaction Control](#transaction-control)
7. [Advanced SQL Features](#advanced-sql-features)
8. [Schema-Aware Parsing](#schema-aware-parsing)
9. [Execution Plan Analysis](#execution-plan-analysis)
10. [Optimization Suggestions](#optimization-suggestions)

For complete examples, see the [examples/queries/](../examples/queries/) directory.

---

## Basic Query Analysis

```bash
# From file
./bin/sqlparser -query examples/queries/complex_query.sql -output table

# From string
./bin/sqlparser -sql "SELECT * FROM users WHERE id > 100" -dialect mysql

# JSON output
./bin/sqlparser -sql "SELECT * FROM users" -output json

# Verbose mode
./bin/sqlparser -query file.sql -verbose
```

---

## Multi-Dialect Support

```bash
# MySQL with backticks
./bin/sqlparser -sql "SELECT \`user_id\` FROM \`users\`" -dialect mysql

# PostgreSQL with double quotes
./bin/sqlparser -sql "SELECT \"user_id\" FROM \"users\"" -dialect postgresql

# SQL Server with brackets
./bin/sqlparser -sql "SELECT [user_id] FROM [users]" -dialect sqlserver

# SQLite
./bin/sqlparser -sql "SELECT * FROM users LIMIT 10 OFFSET 20" -dialect sqlite

# Oracle
./bin/sqlparser -sql "SELECT * FROM users WHERE ROWNUM <= 10" -dialect oracle
```

See [DIALECT_SUPPORT.md](../DIALECT_SUPPORT.md) for complete dialect documentation.

---

## DML Statements

### INSERT

```bash
# Simple INSERT
./bin/sqlparser -sql "INSERT INTO users (name, email) VALUES ('John', 'john@test.com')" -dialect mysql

# Multiple rows
./bin/sqlparser -sql "INSERT INTO users (name, email) VALUES ('John', 'john@test.com'), ('Jane', 'jane@test.com')" -dialect mysql

# INSERT ... SELECT
./bin/sqlparser -sql "INSERT INTO archive SELECT * FROM users WHERE active = 0" -dialect postgresql
```

### UPDATE

```bash
# Simple UPDATE
./bin/sqlparser -sql "UPDATE users SET status = 'active' WHERE id > 100" -dialect postgresql

# Multiple columns
./bin/sqlparser -sql "UPDATE users SET name = 'Jane', email = 'jane@test.com' WHERE id = 1" -dialect mysql

# With ORDER BY and LIMIT (MySQL/SQLite)
./bin/sqlparser -sql "UPDATE users SET status = 'inactive' WHERE last_login < '2020-01-01' ORDER BY last_login LIMIT 100" -dialect mysql
```

### DELETE

```bash
# Simple DELETE
./bin/sqlparser -sql "DELETE FROM users WHERE id = 1" -dialect mysql

# Complex WHERE
./bin/sqlparser -sql "DELETE FROM logs WHERE created_at < '2020-01-01' AND level = 'debug'" -dialect postgresql

# With ORDER BY and LIMIT (MySQL/SQLite)
./bin/sqlparser -sql "DELETE FROM logs WHERE level = 'debug' ORDER BY created_at LIMIT 1000" -dialect mysql
```

---

## Subqueries

```bash
# Scalar subquery
./bin/sqlparser -sql "SELECT * FROM users WHERE salary > (SELECT AVG(salary) FROM employees)" -dialect postgresql

# EXISTS
./bin/sqlparser -sql "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id)" -dialect mysql

# NOT EXISTS
./bin/sqlparser -sql "DELETE FROM users WHERE NOT EXISTS (SELECT 1 FROM orders WHERE user_id = users.id)" -dialect mysql

# IN with subquery
./bin/sqlparser -sql "SELECT name FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > 1000)" -dialect postgresql

# Derived table (subquery in FROM)
./bin/sqlparser -sql "SELECT * FROM (SELECT id, name FROM users WHERE active = 1) AS active_users WHERE id > 100" -dialect postgresql

# JOIN with derived table
./bin/sqlparser -sql "SELECT u.name, o.total FROM users u JOIN (SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id) o ON u.id = o.user_id" -dialect mysql

# Nested subqueries (3 levels)
./bin/sqlparser -sql "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE product_id IN (SELECT id FROM products WHERE category_id IN (SELECT id FROM categories WHERE name = 'electronics')))" -dialect mysql

# Correlated subquery
./bin/sqlparser -sql "SELECT * FROM users u WHERE salary > (SELECT AVG(salary) FROM employees e WHERE e.department = u.department)" -dialect postgresql
```

---

## DDL Support

### CREATE TABLE

```bash
# Simple
./bin/sqlparser -sql "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255) UNIQUE)" -dialect mysql

# IF NOT EXISTS
./bin/sqlparser -sql "CREATE TABLE IF NOT EXISTS products (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL)" -dialect mysql

# With foreign keys
./bin/sqlparser -sql "CREATE TABLE orders (id INT PRIMARY KEY, user_id INT NOT NULL, FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE)" -dialect postgresql

# Composite primary key
./bin/sqlparser -sql "CREATE TABLE user_roles (user_id INT, role_id INT, PRIMARY KEY (user_id, role_id))" -dialect mysql
```

### DROP

```bash
# DROP TABLE
./bin/sqlparser -sql "DROP TABLE IF EXISTS users" -dialect postgresql

# DROP DATABASE
./bin/sqlparser -sql "DROP DATABASE IF EXISTS test_db" -dialect mysql

# DROP INDEX
./bin/sqlparser -sql "DROP INDEX IF EXISTS idx_users_email" -dialect postgresql
```

### ALTER TABLE

```bash
# ADD COLUMN
./bin/sqlparser -sql "ALTER TABLE users ADD COLUMN age INT NOT NULL" -dialect mysql

# DROP COLUMN
./bin/sqlparser -sql "ALTER TABLE users DROP COLUMN age" -dialect postgresql

# MODIFY COLUMN
./bin/sqlparser -sql "ALTER TABLE users MODIFY COLUMN name VARCHAR(150) NOT NULL" -dialect mysql

# ADD CONSTRAINT
./bin/sqlparser -sql "ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)" -dialect postgresql
```

### CREATE INDEX

```bash
# Simple index
./bin/sqlparser -sql "CREATE INDEX idx_users_email ON users (email)" -dialect mysql

# Unique index
./bin/sqlparser -sql "CREATE UNIQUE INDEX idx_users_email ON users (email)" -dialect postgresql

# Multi-column index
./bin/sqlparser -sql "CREATE INDEX idx_orders_user_product ON orders (user_id, product_id)" -dialect postgresql
```

---

## Transaction Control

```bash
# BEGIN/START TRANSACTION
./bin/sqlparser -sql "BEGIN TRANSACTION" -dialect sqlserver
./bin/sqlparser -sql "START TRANSACTION" -dialect mysql

# COMMIT
./bin/sqlparser -sql "COMMIT" -dialect postgresql
./bin/sqlparser -sql "COMMIT WORK" -dialect postgresql

# ROLLBACK
./bin/sqlparser -sql "ROLLBACK" -dialect mysql
./bin/sqlparser -sql "ROLLBACK WORK" -dialect mysql

# SAVEPOINT
./bin/sqlparser -sql "SAVEPOINT my_savepoint" -dialect postgresql

# ROLLBACK TO SAVEPOINT
./bin/sqlparser -sql "ROLLBACK TO SAVEPOINT my_savepoint" -dialect mysql

# RELEASE SAVEPOINT
./bin/sqlparser -sql "RELEASE SAVEPOINT my_savepoint" -dialect postgresql
```

---

## Advanced SQL Features

### CTEs (Common Table Expressions)

```bash
# Simple CTE
./bin/sqlparser -sql "WITH sales_summary AS (SELECT product_id, SUM(amount) as total FROM sales GROUP BY product_id) SELECT * FROM sales_summary WHERE total > 1000" -dialect postgresql

# Multiple CTEs
./bin/sqlparser -sql "WITH active_users AS (SELECT id, name FROM users WHERE active = 1), recent_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders WHERE created_at > '2024-01-01' GROUP BY user_id) SELECT u.name, COALESCE(o.order_count, 0) as orders FROM active_users u LEFT JOIN recent_orders o ON u.id = o.user_id" -dialect postgresql

# CTE with column list
./bin/sqlparser -sql "WITH user_stats (user_id, total_spent) AS (SELECT user_id, SUM(amount) FROM orders GROUP BY user_id) SELECT * FROM user_stats WHERE total_spent > 5000" -dialect mysql
```

### Window Functions

```bash
# ROW_NUMBER
./bin/sqlparser -sql "SELECT employee_id, salary, ROW_NUMBER() OVER (ORDER BY salary DESC) as rank FROM employees" -dialect postgresql

# PARTITION BY
./bin/sqlparser -sql "SELECT employee_id, department, salary, RANK() OVER (PARTITION BY department ORDER BY salary DESC) as dept_rank FROM employees" -dialect postgresql

# Window frames - ROWS
./bin/sqlparser -sql "SELECT employee_id, salary, AVG(salary) OVER (ORDER BY hire_date ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) as moving_avg FROM employees" -dialect postgresql

# Window frames - RANGE
./bin/sqlparser -sql "SELECT employee_id, salary, SUM(salary) OVER (ORDER BY hire_date RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) as running_total FROM employees" -dialect postgresql

# Multiple window functions
./bin/sqlparser -sql "SELECT employee_id, salary, ROW_NUMBER() OVER (ORDER BY salary DESC) as overall_rank, RANK() OVER (PARTITION BY department ORDER BY salary DESC) as dept_rank, AVG(salary) OVER (PARTITION BY department) as dept_avg FROM employees" -dialect postgresql
```

### Set Operations

```bash
# UNION
./bin/sqlparser -sql "SELECT id, name FROM customers UNION SELECT id, name FROM prospects" -dialect mysql

# UNION ALL
./bin/sqlparser -sql "SELECT id FROM active_users UNION ALL SELECT id FROM inactive_users" -dialect postgresql

# INTERSECT
./bin/sqlparser -sql "SELECT id FROM customers INTERSECT SELECT id FROM newsletter_subscribers" -dialect postgresql

# EXCEPT
./bin/sqlparser -sql "SELECT id FROM all_users EXCEPT SELECT id FROM banned_users" -dialect postgresql
```

---

## Schema-Aware Parsing

See [examples/schemas/](../examples/schemas/) for schema definition examples.

```go
package main

import (
	"context"
	"fmt"
	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
	"github.com/Chahine-tech/sql-parser-go/pkg/schema"
)

func main() {
	// Load schema from JSON file
	loader := schema.NewSchemaLoader()
	s, _ := loader.LoadFromFile("examples/schemas/test_schema.json")

	// Create validator
	validator := schema.NewValidator(s)

	// Parse SQL query
	sql := "SELECT id, name, invalid_column FROM users WHERE age > 18"
	p := parser.NewWithDialect(context.Background(), sql, dialect.GetDialect("mysql"))
	stmt, _ := p.ParseStatement()

	// Validate against schema
	errors := validator.ValidateStatement(stmt)
	for _, err := range errors {
		fmt.Printf("[%s] %s\n", err.Type, err.Message)
	}
}
```

---

## Execution Plan Analysis

```go
package main

import (
	"context"
	"fmt"
	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
	"github.com/Chahine-tech/sql-parser-go/pkg/plan"
)

func main() {
	// Parse EXPLAIN statement
	sql := "EXPLAIN ANALYZE SELECT u.name, COUNT(o.id) FROM users u JOIN orders o ON u.id = o.user_id GROUP BY u.name"
	p := parser.NewWithDialect(context.Background(), sql, dialect.GetDialect("postgresql"))
	stmt, _ := p.ParseStatement()

	explainStmt := stmt.(*parser.ExplainStatement)
	fmt.Printf("Analyzing: %s\n", explainStmt.Statement.Type())

	// Parse execution plan from database output
	jsonPlan := []byte(`{"Plan": {...}}`) // From EXPLAIN FORMAT=JSON
	executionPlan, _ := plan.ParseJSONPlan(jsonPlan, "postgresql")

	// Analyze plan for bottlenecks
	analyzer := plan.NewPlanAnalyzer("postgresql")
	analysis := analyzer.AnalyzePlan(executionPlan)

	fmt.Printf("Performance Score: %.2f/100\n", analysis.PerformanceScore)

	// Review issues and recommendations
	for _, issue := range analysis.Issues {
		fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Type, issue.Description)
	}
}
```

---

## Optimization Suggestions

```bash
# Basic optimization analysis
./bin/sqlparser -sql "SELECT * FROM users WHERE UPPER(email) = 'TEST'" -dialect mysql -output table

# MySQL LIMIT without ORDER BY
./bin/sqlparser -sql "SELECT name FROM users LIMIT 10" -dialect mysql

# SQL Server TOP suggestion
./bin/sqlparser -sql "SELECT name FROM users LIMIT 10" -dialect sqlserver

# PostgreSQL JSON optimization
./bin/sqlparser -sql "SELECT data FROM logs WHERE json_extract(data, '$.type') = 'error'" -dialect postgresql
```

---

## Additional Resources

- [README.md](../README.md) - Main documentation
- [PERFORMANCE.md](PERFORMANCE.md) - Performance benchmarks
- [DIALECT_SUPPORT.md](../DIALECT_SUPPORT.md) - Dialect-specific features
- [CLAUDE.md](../CLAUDE.md) - Developer guide
- [examples/queries/](../examples/queries/) - Example SQL files
