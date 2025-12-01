# SQL Parser Go - Usage Examples

This document provides comprehensive usage examples for SQL Parser Go across all supported features and SQL dialects.

## Table of Contents

1. [Basic Query Analysis](#basic-query-analysis)
2. [Multi-Dialect Support](#multi-dialect-support)
3. [DML Statements](#dml-statements-insert-update-delete)
4. [Subqueries](#subqueries)
5. [DDL Support](#ddl-support-create-drop-alter)
6. [Transaction Control](#transaction-control)
7. [Advanced SQL Features](#advanced-sql-features)
8. [Schema-Aware Parsing](#schema-aware-parsing)
9. [Execution Plan Analysis](#execution-plan-analysis)
10. [Optimization Suggestions](#optimization-suggestions)
11. [Log Parsing](#log-parsing)

---

## Basic Query Analysis

### Analyze from File

```bash
./bin/sqlparser -query examples/queries/complex_query.sql -output table
```

### Analyze from String

```bash
./bin/sqlparser -sql "SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id" -output json
```

### Output Formats

**JSON Output:**
```bash
./bin/sqlparser -sql "SELECT * FROM users" -output json
```

**Table Output:**
```bash
./bin/sqlparser -sql "SELECT * FROM users" -output table
```

**Verbose Mode:**
```bash
./bin/sqlparser -query file.sql -verbose
```

---

## Multi-Dialect Support

### MySQL

```bash
# Backtick identifiers
./bin/sqlparser -sql "SELECT \`user_id\`, \`email\` FROM \`users\`" -dialect mysql

# MySQL-specific LIMIT
./bin/sqlparser -sql "SELECT name FROM users ORDER BY created_at DESC LIMIT 10" -dialect mysql

# AUTO_INCREMENT
./bin/sqlparser -sql "CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100))" -dialect mysql
```

### PostgreSQL

```bash
# Double-quoted identifiers
./bin/sqlparser -sql "SELECT \"user_id\", \"email\" FROM \"users\"" -dialect postgresql

# PostgreSQL arrays
./bin/sqlparser -sql "SELECT tags FROM posts WHERE 'postgresql' = ANY(tags)" -dialect postgresql

# RETURNING clause
./bin/sqlparser -sql "INSERT INTO users (name) VALUES ('John') RETURNING id" -dialect postgresql
```

### SQL Server

```bash
# Bracket identifiers
./bin/sqlparser -sql "SELECT [user_id], [email] FROM [users]" -dialect sqlserver

# TOP clause
./bin/sqlparser -sql "SELECT TOP 10 name FROM users ORDER BY created_at DESC" -dialect sqlserver

# IDENTITY
./bin/sqlparser -sql "CREATE TABLE users (id INT IDENTITY(1,1) PRIMARY KEY, name VARCHAR(100))" -dialect sqlserver
```

### SQLite

```bash
# SQLite-specific syntax
./bin/sqlparser -sql "SELECT * FROM users LIMIT 10 OFFSET 20" -dialect sqlite

# AUTOINCREMENT
./bin/sqlparser -sql "CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)" -dialect sqlite
```

### Oracle

```bash
# Oracle-specific syntax
./bin/sqlparser -sql "SELECT * FROM users WHERE ROWNUM <= 10" -dialect oracle

# Dual table
./bin/sqlparser -sql "SELECT SYSDATE FROM DUAL" -dialect oracle
```

See [DIALECT_SUPPORT.md](DIALECT_SUPPORT.md) for complete dialect documentation.

---

## DML Statements (INSERT, UPDATE, DELETE)

### INSERT Statements

#### Simple INSERT with VALUES

```bash
./bin/sqlparser -sql "INSERT INTO users (name, email) VALUES ('John', 'john@test.com')" -dialect mysql -output table
```

#### INSERT with Multiple Rows

```bash
./bin/sqlparser -sql "INSERT INTO users (name, email) VALUES ('John', 'john@test.com'), ('Jane', 'jane@test.com'), ('Bob', 'bob@test.com')" -dialect mysql
```

#### INSERT with SELECT

```bash
./bin/sqlparser -sql "INSERT INTO archive SELECT * FROM users WHERE active = 0" -dialect postgresql

# With specific columns
./bin/sqlparser -sql "INSERT INTO user_stats (user_id, total_orders) SELECT user_id, COUNT(*) FROM orders GROUP BY user_id" -dialect mysql
```

### UPDATE Statements

#### Simple UPDATE

```bash
./bin/sqlparser -sql "UPDATE users SET status = 'active' WHERE id > 100" -dialect postgresql -output table
```

#### UPDATE Multiple Columns

```bash
./bin/sqlparser -sql "UPDATE users SET name = 'Jane', email = 'jane@test.com', status = 1 WHERE id = 1" -dialect mysql
```

#### UPDATE with ORDER BY and LIMIT (MySQL/SQLite)

```bash
./bin/sqlparser -sql "UPDATE users SET status = 'inactive' WHERE last_login < '2020-01-01' ORDER BY last_login LIMIT 100" -dialect mysql
```

#### UPDATE with Subquery

```bash
./bin/sqlparser -sql "UPDATE users SET status = (SELECT status FROM user_preferences WHERE user_id = users.id)" -dialect postgresql
```

### DELETE Statements

#### Simple DELETE

```bash
./bin/sqlparser -sql "DELETE FROM users WHERE id = 1" -dialect mysql -output table
```

#### DELETE with Complex WHERE

```bash
./bin/sqlparser -sql "DELETE FROM logs WHERE created_at < '2020-01-01' AND level = 'debug'" -dialect postgresql
```

#### DELETE with ORDER BY and LIMIT (MySQL/SQLite)

```bash
./bin/sqlparser -sql "DELETE FROM logs WHERE level = 'debug' ORDER BY created_at LIMIT 1000" -dialect mysql
```

---

## Subqueries

SQL Parser Go provides comprehensive subquery support across all clauses and statement types.

### Scalar Subqueries

#### In WHERE Clause

```bash
./bin/sqlparser -sql "SELECT * FROM users WHERE salary > (SELECT AVG(salary) FROM employees)" -dialect postgresql
```

#### In SELECT Clause

```bash
./bin/sqlparser -sql "SELECT id, name, (SELECT COUNT(*) FROM orders WHERE user_id = users.id) as order_count FROM users" -dialect postgresql
```

#### Multiple Scalar Subqueries

```bash
./bin/sqlparser -sql "SELECT id, (SELECT AVG(price) FROM products) as avg_price, (SELECT MAX(price) FROM products) as max_price FROM users" -dialect mysql
```

### EXISTS and NOT EXISTS

#### EXISTS Subquery

```bash
./bin/sqlparser -sql "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)" -dialect mysql
```

#### NOT EXISTS Subquery

```bash
./bin/sqlparser -sql "DELETE FROM users WHERE NOT EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)" -dialect mysql
```

#### Multiple EXISTS Clauses

```bash
./bin/sqlparser -sql "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id) AND EXISTS (SELECT 1 FROM payments WHERE user_id = users.id)" -dialect mysql
```

### IN and NOT IN with Subqueries

#### IN with Subquery

```bash
./bin/sqlparser -sql "SELECT name FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > 1000)" -dialect postgresql
```

#### NOT IN with Subquery

```bash
./bin/sqlparser -sql "SELECT * FROM users WHERE id NOT IN (SELECT user_id FROM banned_users)" -dialect mysql
```

### Derived Tables (Subqueries in FROM)

#### Simple Derived Table

```bash
./bin/sqlparser -sql "SELECT * FROM (SELECT id, name FROM users WHERE active = 1) AS active_users WHERE id > 100" -dialect postgresql
```

#### JOIN with Derived Table

```bash
./bin/sqlparser -sql "SELECT u.name, o.total FROM users u JOIN (SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id) o ON u.id = o.user_id" -dialect mysql
```

#### Multiple Derived Tables

```bash
./bin/sqlparser -sql "SELECT * FROM (SELECT id FROM users) u, (SELECT user_id FROM orders) o WHERE u.id = o.user_id" -dialect postgresql
```

### Nested and Correlated Subqueries

#### Triple Nested Subquery

```bash
./bin/sqlparser -sql "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE product_id IN (SELECT id FROM products WHERE category_id IN (SELECT id FROM categories WHERE name = 'electronics')))" -dialect mysql
```

#### Correlated Subquery

```bash
./bin/sqlparser -sql "SELECT * FROM users u WHERE salary > (SELECT AVG(salary) FROM employees e WHERE e.department = u.department)" -dialect postgresql
```

### Subqueries in DML Statements

#### INSERT with Subquery

```bash
./bin/sqlparser -sql "INSERT INTO user_stats (user_id, order_count) VALUES (1, (SELECT COUNT(*) FROM orders WHERE user_id = 1))" -dialect mysql
```

#### UPDATE with Subquery in WHERE

```bash
./bin/sqlparser -sql "UPDATE users SET active = 0 WHERE id IN (SELECT user_id FROM banned_users)" -dialect mysql
```

#### DELETE with EXISTS

```bash
./bin/sqlparser -sql "DELETE FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id AND total > 10000)" -dialect postgresql
```

---

## DDL Support (CREATE, DROP, ALTER)

### CREATE TABLE

#### Simple CREATE TABLE

```bash
./bin/sqlparser -sql "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255) UNIQUE)" -dialect mysql
```

#### CREATE TABLE with IF NOT EXISTS

```bash
./bin/sqlparser -sql "CREATE TABLE IF NOT EXISTS products (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, price DECIMAL(10,2) DEFAULT 0.00)" -dialect mysql
```

#### CREATE TABLE with Foreign Keys

```bash
./bin/sqlparser -sql "CREATE TABLE orders (id INT PRIMARY KEY, user_id INT NOT NULL, product_id INT, FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE, FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE SET NULL)" -dialect postgresql
```

#### CREATE TABLE with Composite Primary Key

```bash
./bin/sqlparser -sql "CREATE TABLE user_roles (user_id INT, role_id INT, assigned_at TIMESTAMP, PRIMARY KEY (user_id, role_id))" -dialect mysql
```

#### Dialect-Specific CREATE TABLE

**MySQL - AUTO_INCREMENT:**
```bash
./bin/sqlparser -sql "CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100))" -dialect mysql
```

**SQL Server - IDENTITY:**
```bash
./bin/sqlparser -sql "CREATE TABLE users (id INT IDENTITY(1,1) PRIMARY KEY, name VARCHAR(100))" -dialect sqlserver
```

**SQLite - AUTOINCREMENT:**
```bash
./bin/sqlparser -sql "CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)" -dialect sqlite
```

### DROP Statements

#### DROP TABLE

```bash
./bin/sqlparser -sql "DROP TABLE users" -dialect mysql

# With IF EXISTS
./bin/sqlparser -sql "DROP TABLE IF EXISTS users" -dialect postgresql
```

#### DROP DATABASE

```bash
./bin/sqlparser -sql "DROP DATABASE test_db" -dialect mysql

# With IF EXISTS
./bin/sqlparser -sql "DROP DATABASE IF EXISTS test_db" -dialect mysql
```

#### DROP INDEX

```bash
./bin/sqlparser -sql "DROP INDEX idx_users_email" -dialect mysql

# With IF EXISTS
./bin/sqlparser -sql "DROP INDEX IF EXISTS idx_users_email" -dialect postgresql
```

### ALTER TABLE

#### ADD COLUMN

```bash
./bin/sqlparser -sql "ALTER TABLE users ADD COLUMN age INT NOT NULL" -dialect mysql

# With default value
./bin/sqlparser -sql "ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active'" -dialect postgresql
```

#### DROP COLUMN

```bash
./bin/sqlparser -sql "ALTER TABLE users DROP COLUMN age" -dialect postgresql
```

#### MODIFY COLUMN

```bash
./bin/sqlparser -sql "ALTER TABLE users MODIFY COLUMN name VARCHAR(150) NOT NULL" -dialect mysql
```

#### CHANGE COLUMN (MySQL)

```bash
./bin/sqlparser -sql "ALTER TABLE users CHANGE COLUMN old_name new_name VARCHAR(100)" -dialect mysql
```

#### ADD CONSTRAINT

```bash
# Foreign key
./bin/sqlparser -sql "ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)" -dialect postgresql

# Unique constraint
./bin/sqlparser -sql "ALTER TABLE users ADD CONSTRAINT uk_email UNIQUE (email)" -dialect mysql
```

#### DROP CONSTRAINT

```bash
./bin/sqlparser -sql "ALTER TABLE orders DROP CONSTRAINT fk_user" -dialect postgresql
```

### CREATE INDEX

#### Simple Index

```bash
./bin/sqlparser -sql "CREATE INDEX idx_users_email ON users (email)" -dialect mysql
```

#### Unique Index

```bash
./bin/sqlparser -sql "CREATE UNIQUE INDEX idx_users_email ON users (email)" -dialect postgresql
```

#### Multi-Column Index

```bash
./bin/sqlparser -sql "CREATE INDEX idx_orders_user_product ON orders (user_id, product_id)" -dialect postgresql
```

#### CREATE INDEX with IF NOT EXISTS

```bash
./bin/sqlparser -sql "CREATE INDEX IF NOT EXISTS idx_products_category ON products (category)" -dialect mysql
```

---

## Transaction Control

### BEGIN/START TRANSACTION

```bash
# SQL Server
./bin/sqlparser -sql "BEGIN TRANSACTION" -dialect sqlserver

# MySQL/PostgreSQL
./bin/sqlparser -sql "START TRANSACTION" -dialect mysql
```

### COMMIT

```bash
./bin/sqlparser -sql "COMMIT" -dialect postgresql

# With WORK keyword
./bin/sqlparser -sql "COMMIT WORK" -dialect postgresql
```

### ROLLBACK

```bash
./bin/sqlparser -sql "ROLLBACK" -dialect mysql

# With WORK keyword
./bin/sqlparser -sql "ROLLBACK WORK" -dialect mysql
```

### SAVEPOINT

```bash
./bin/sqlparser -sql "SAVEPOINT my_savepoint" -dialect postgresql
```

### ROLLBACK TO SAVEPOINT

```bash
./bin/sqlparser -sql "ROLLBACK TO SAVEPOINT my_savepoint" -dialect mysql
```

### RELEASE SAVEPOINT (PostgreSQL/MySQL)

```bash
./bin/sqlparser -sql "RELEASE SAVEPOINT my_savepoint" -dialect postgresql
```

### Complete Transaction Example

```bash
# Start transaction
./bin/sqlparser -sql "BEGIN TRANSACTION" -dialect sqlserver

# Create savepoint
./bin/sqlparser -sql "SAVEPOINT before_update" -dialect postgresql

# If something goes wrong, rollback to savepoint
./bin/sqlparser -sql "ROLLBACK TO SAVEPOINT before_update" -dialect postgresql

# Or commit the transaction
./bin/sqlparser -sql "COMMIT" -dialect postgresql
```

---

## Advanced SQL Features

### CTEs (Common Table Expressions)

#### Simple CTE

```bash
./bin/sqlparser -sql "WITH sales_summary AS (SELECT product_id, SUM(amount) as total FROM sales GROUP BY product_id) SELECT * FROM sales_summary WHERE total > 1000" -dialect postgresql
```

#### Multiple CTEs

```bash
./bin/sqlparser -sql "WITH
  active_users AS (SELECT id, name FROM users WHERE active = 1),
  recent_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders WHERE created_at > '2024-01-01' GROUP BY user_id)
SELECT u.name, COALESCE(o.order_count, 0) as orders
FROM active_users u
LEFT JOIN recent_orders o ON u.id = o.user_id" -dialect postgresql
```

#### CTE with Column List

```bash
./bin/sqlparser -sql "WITH user_stats (user_id, total_spent) AS (SELECT user_id, SUM(amount) FROM orders GROUP BY user_id) SELECT * FROM user_stats WHERE total_spent > 5000" -dialect mysql
```

### Window Functions

#### ROW_NUMBER

```bash
./bin/sqlparser -sql "SELECT employee_id, salary, ROW_NUMBER() OVER (ORDER BY salary DESC) as rank FROM employees" -dialect postgresql
```

#### PARTITION BY

```bash
./bin/sqlparser -sql "SELECT employee_id, department, salary, RANK() OVER (PARTITION BY department ORDER BY salary DESC) as dept_rank FROM employees" -dialect postgresql
```

#### Window Frames - ROWS

```bash
./bin/sqlparser -sql "SELECT employee_id, salary, AVG(salary) OVER (ORDER BY hire_date ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) as moving_avg FROM employees" -dialect postgresql
```

#### Window Frames - RANGE

```bash
./bin/sqlparser -sql "SELECT employee_id, salary, SUM(salary) OVER (ORDER BY hire_date RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) as running_total FROM employees" -dialect postgresql
```

#### Multiple Window Functions

```bash
./bin/sqlparser -sql "SELECT
  employee_id,
  salary,
  ROW_NUMBER() OVER (ORDER BY salary DESC) as overall_rank,
  RANK() OVER (PARTITION BY department ORDER BY salary DESC) as dept_rank,
  AVG(salary) OVER (PARTITION BY department) as dept_avg
FROM employees" -dialect postgresql
```

### Set Operations

#### UNION

```bash
./bin/sqlparser -sql "SELECT id, name FROM customers UNION SELECT id, name FROM prospects" -dialect mysql
```

#### UNION ALL

```bash
./bin/sqlparser -sql "SELECT id FROM active_users UNION ALL SELECT id FROM inactive_users" -dialect postgresql
```

#### INTERSECT

```bash
./bin/sqlparser -sql "SELECT id FROM customers INTERSECT SELECT id FROM newsletter_subscribers" -dialect postgresql
```

#### EXCEPT (PostgreSQL/SQL Server)

```bash
./bin/sqlparser -sql "SELECT id FROM all_users EXCEPT SELECT id FROM banned_users" -dialect postgresql
```

#### Complex Set Operations

```bash
./bin/sqlparser -sql "SELECT id FROM customers UNION ALL SELECT id FROM prospects INTERSECT SELECT id FROM active_accounts" -dialect postgresql
```

---

## Schema-Aware Parsing

### Basic Schema Validation (Go Code)

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
		// Output: [COLUMN_NOT_FOUND] Column 'invalid_column' not found in table 'users'
		//         [COLUMN_NOT_FOUND] Column 'age' not found in table 'users'
	}
}
```

### Schema Definition (JSON)

```json
{
  "name": "my_database",
  "tables": [
    {
      "name": "users",
      "columns": [
        {"name": "id", "type": "INT", "primary_key": true},
        {"name": "name", "type": "VARCHAR", "length": 100, "nullable": false},
        {"name": "email", "type": "VARCHAR", "length": 255, "unique": true},
        {"name": "created_at", "type": "TIMESTAMP", "default": "CURRENT_TIMESTAMP"}
      ],
      "indexes": [
        {"name": "idx_email", "columns": ["email"], "unique": true}
      ]
    },
    {
      "name": "orders",
      "columns": [
        {"name": "id", "type": "INT", "primary_key": true},
        {"name": "user_id", "type": "INT", "nullable": false},
        {"name": "total", "type": "DECIMAL", "precision": 10, "scale": 2}
      ],
      "foreign_keys": [
        {
          "name": "fk_user",
          "columns": ["user_id"],
          "referenced_table": "users",
          "referenced_columns": ["id"],
          "on_delete": "CASCADE"
        }
      ]
    }
  ]
}
```

### Type Checking (Go Code)

```go
// Type compatibility checking
typeChecker := schema.NewTypeChecker(s)
errors := typeChecker.CheckStatement(stmt)

for _, err := range errors {
	fmt.Printf("[%s] %s\n", err.Type, err.Message)
	// Example: [TYPE_MISMATCH] Cannot compare VARCHAR with INT
}
```

---

## Execution Plan Analysis

### EXPLAIN Statement Parsing (Go Code)

```go
package main

import (
	"context"
	"fmt"
	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func main() {
	// Parse EXPLAIN statement
	sql := "EXPLAIN ANALYZE SELECT u.name, COUNT(o.id) FROM users u JOIN orders o ON u.id = o.user_id GROUP BY u.name"
	p := parser.NewWithDialect(context.Background(), sql, dialect.GetDialect("postgresql"))
	stmt, _ := p.ParseStatement()

	explainStmt := stmt.(*parser.ExplainStatement)
	fmt.Printf("Analyzing: %s\n", explainStmt.Statement.Type())
	fmt.Printf("With ANALYZE: %v\n", explainStmt.Analyze)
	fmt.Printf("Format: %s\n", explainStmt.Format)
}
```

### Execution Plan Analysis (Go Code)

```go
package main

import (
	"fmt"
	"github.com/Chahine-tech/sql-parser-go/pkg/plan"
)

func main() {
	// Parse execution plan from database output
	jsonPlan := []byte(`{
		"Plan": {
			"Node Type": "Hash Join",
			"Startup Cost": 10.5,
			"Total Cost": 125.75,
			"Plan Rows": 1000,
			"Plans": [
				{
					"Node Type": "Seq Scan",
					"Relation Name": "users",
					"Total Cost": 50.0
				}
			]
		}
	}`)

	executionPlan, _ := plan.ParseJSONPlan(jsonPlan, "postgresql")

	// Analyze plan for bottlenecks
	analyzer := plan.NewPlanAnalyzer("postgresql")
	analysis := analyzer.AnalyzePlan(executionPlan)

	// Print performance score
	fmt.Printf("Performance Score: %.2f/100\n", analysis.PerformanceScore)

	// Review issues
	for _, issue := range analysis.Issues {
		fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Type, issue.Description)
	}

	// Get recommendations
	for _, rec := range analysis.Recommendations {
		fmt.Printf("[%s] %s (Impact: %.2f)\n", rec.Priority, rec.Description, rec.ExpectedImpact)
	}

	// Find bottlenecks
	bottlenecks := executionPlan.FindBottlenecks()
	for _, b := range bottlenecks {
		fmt.Printf("Bottleneck: %s (Impact: %.2f)\n", b.Issue, b.ImpactScore)
		fmt.Printf("Fix: %s\n", b.Recommendation)
	}
}
```

---

## Optimization Suggestions

### Basic Optimization Analysis

```bash
./bin/sqlparser -sql "SELECT * FROM users WHERE UPPER(email) = 'TEST'" -dialect mysql -output table
```

**Output:**
```
┌─────────────────────────────────────┬──────────┬────────────────────────────────┬─────────────────────────┐
│ TYPE                                │ SEVERITY │ DESCRIPTION                    │ SUGGESTION              │
├─────────────────────────────────────┼──────────┼────────────────────────────────┼─────────────────────────┤
│ 🔍 SELECT_STAR                      │ WARNING  │ Avoid SELECT * for performance │ Specify explicit columns│
│ ⚡ FUNCTION_ON_COLUMN                │ WARNING  │ Function on indexed column     │ Use functional index    │
└─────────────────────────────────────┴──────────┴────────────────────────────────┴─────────────────────────┘
```

### Dialect-Specific Optimizations

#### MySQL - LIMIT without ORDER BY

```bash
./bin/sqlparser -sql "SELECT name FROM users LIMIT 10" -dialect mysql
```

**Optimization:** Warns about non-deterministic results, suggests adding ORDER BY.

#### SQL Server - Suggest TOP

```bash
./bin/sqlparser -sql "SELECT name FROM users LIMIT 10" -dialect sqlserver
```

**Optimization:** Suggests using SQL Server's TOP clause instead of LIMIT.

#### PostgreSQL - JSON vs JSONB

```bash
./bin/sqlparser -sql "SELECT data FROM logs WHERE json_extract(data, '$.type') = 'error'" -dialect postgresql
```

**Optimization:** Suggests using JSONB for better performance with JSON operations.

---

## Log Parsing

### Parse SQL Server Profiler Log

```bash
./bin/sqlparser -log examples/logs/sample_profiler.log -output table -verbose
```

### Parse Extended Events Log

```bash
./bin/sqlparser -log examples/logs/extended_events.xel -output json
```

### Parse Query Store Export

```bash
./bin/sqlparser -log examples/logs/query_store.csv -output table
```

---

## Command Line Reference

### Basic Options

- `-query FILE`: Analyze SQL query from file
- `-sql STRING`: Analyze SQL query from string
- `-log FILE`: Parse SQL Server log file
- `-output FORMAT`: Output format (json, table) - default: json
- `-dialect DIALECT`: SQL dialect (mysql, postgresql, sqlserver, sqlite, oracle) - default: sqlserver
- `-verbose`: Enable verbose output
- `-config FILE`: Configuration file path
- `-help`: Show help

### Examples

```bash
# Query from file with table output
./bin/sqlparser -query examples/queries/complex_query.sql -output table

# Query from string with MySQL dialect
./bin/sqlparser -sql "SELECT * FROM users" -dialect mysql -output json

# Parse log with verbose output
./bin/sqlparser -log examples/logs/sample.log -verbose

# Use custom config
./bin/sqlparser -query file.sql -config custom_config.yaml

# Show help
./bin/sqlparser -help
```

---

## Additional Resources

- [README.md](README.md) - Main documentation and getting started guide
- [DIALECT_SUPPORT.md](DIALECT_SUPPORT.md) - Complete dialect-specific documentation
- [PERFORMANCE.md](PERFORMANCE.md) - Detailed performance benchmarks
- [CLAUDE.md](CLAUDE.md) - Developer guide for working with Claude Code
- [examples/queries/](examples/queries/) - Example SQL files for all features
