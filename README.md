# SQL Parser Go

A powerful multi-dialect SQL query analysis tool written in Go that provides comprehensive parsing, analysis, and optimization suggestions for SQL queries and log files.

## Features

- **Multi-Dialect Support**: Parse MySQL, PostgreSQL, SQL Server, SQLite, and Oracle queries
- **Full SQL Statement Support**:
  - **SELECT**: Complex queries with joins, subqueries, aggregations
  - **INSERT**: VALUES and SELECT variants, multiple rows
  - **UPDATE**: Multiple columns, WHERE clause, ORDER BY/LIMIT (MySQL/SQLite)
  - **DELETE**: WHERE clause, ORDER BY/LIMIT (MySQL/SQLite)
  - **CREATE TABLE**: Column definitions, constraints, foreign keys, IF NOT EXISTS
  - **DROP**: TABLE/DATABASE/INDEX with IF EXISTS and CASCADE
  - **ALTER TABLE**: ADD/DROP/MODIFY/CHANGE columns and constraints
  - **CREATE INDEX**: Simple and unique indexes with IF NOT EXISTS
  - **TRANSACTIONS**: BEGIN, START TRANSACTION, COMMIT, ROLLBACK, SAVEPOINT, RELEASE
  - **EXPLAIN**: Full support for EXPLAIN and EXPLAIN ANALYZE with dialect-specific options
  - **STORED PROCEDURES**: CREATE PROCEDURE with IN/OUT/INOUT parameters, variables, cursors
  - **FUNCTIONS**: CREATE FUNCTION with return types, DETERMINISTIC, dialect-specific options
- **Schema-Aware Parsing**: Validate SQL against database schemas
  - **Schema Loading**: Load schemas from JSON, YAML files
  - **Table/Column Validation**: Verify table and column existence
  - **Type Checking**: Validate data type compatibility
  - **Foreign Key Validation**: Check foreign key references
- **Execution Plan Analysis**: Analyze query performance and detect bottlenecks
  - **EXPLAIN Parsing**: Parse EXPLAIN statements across all dialects
  - **Plan Analysis**: Automatic detection of performance issues
  - **Optimization Suggestions**: Get actionable recommendations for query improvements
  - **Performance Scoring**: 0-100 score based on plan quality
- **SQL Query Parsing**: Parse and analyze complex SQL queries with dialect-specific syntax
- **Abstract Syntax Tree (AST)**: Generate detailed AST representations
- **Query Analysis**: Extract tables, columns, joins, and conditions
- **Advanced Optimization Suggestions**: Get intelligent, dialect-specific recommendations for query improvements
- **Comprehensive Subquery Support**: Full support for subqueries in all clauses
  - **Scalar Subqueries**: In WHERE, SELECT, INSERT VALUES, UPDATE SET
  - **EXISTS / NOT EXISTS**: Full implementation across all statement types
  - **IN / NOT IN with Subqueries**: Complete support
  - **Derived Tables**: Subqueries in FROM clause with JOIN support
  - **Nested & Correlated Subqueries**: Multiple levels of nesting
- **Advanced SQL Features**:
  - **CTEs (WITH clause)**: Common Table Expressions with recursive support
  - **Window Functions**: ROW_NUMBER, RANK, PARTITION BY, ORDER BY, window frames
  - **Set Operations**: UNION, UNION ALL, INTERSECT, EXCEPT
  - **CASE Expressions**: Searched and simple CASE statements
- **Log Parsing**: Parse SQL Server log files (Profiler, Extended Events, Query Store)
- **Multiple Output Formats**: JSON, table, and CSV output
- **CLI Interface**: Easy-to-use command-line interface with enhanced optimization output
- **Dialect-Specific Features**: Handle quoted identifiers, keywords, and features for each dialect

## Installation

### Prerequisites

- Go 1.21 or higher

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Chahine-tech/sql-parser-go.git
cd sql-parser-go

# Install dependencies
make deps

# Build the application
make build

# Install to GOPATH/bin (optional)
make install
```

## Usage

### Analyze SQL Queries

#### From File
```bash
./bin/sqlparser -query examples/queries/complex_query.sql -output table
```

#### From String
```bash
./bin/sqlparser -sql "SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id" -output json
```

### Multi-Dialect Support

#### MySQL
```bash
./bin/sqlparser -sql "SELECT \`user_id\`, \`email\` FROM \`users\`" -dialect mysql
```

#### PostgreSQL
```bash
./bin/sqlparser -sql "SELECT \"user_id\", \"email\" FROM \"users\"" -dialect postgresql
```

#### SQL Server
```bash
./bin/sqlparser -sql "SELECT [user_id], [email] FROM [users]" -dialect sqlserver
```

See [DIALECT_SUPPORT.md](DIALECT_SUPPORT.md) for detailed information about dialect-specific features.

### DML Statement Support (INSERT, UPDATE, DELETE)

#### INSERT Statements
```bash
# Simple INSERT with VALUES
./bin/sqlparser -sql "INSERT INTO users (name, email) VALUES ('John', 'john@test.com')" -dialect mysql -output table

# INSERT with multiple rows
./bin/sqlparser -sql "INSERT INTO users (name, email) VALUES ('John', 'john@test.com'), ('Jane', 'jane@test.com')" -dialect mysql

# INSERT with SELECT
./bin/sqlparser -sql "INSERT INTO archive SELECT * FROM users WHERE active = 0" -dialect postgresql
```

#### UPDATE Statements
```bash
# Simple UPDATE
./bin/sqlparser -sql "UPDATE users SET status = 'active' WHERE id > 100" -dialect postgresql -output table

# UPDATE multiple columns
./bin/sqlparser -sql "UPDATE users SET name = 'Jane', email = 'jane@test.com', status = 1 WHERE id = 1" -dialect mysql

# UPDATE with ORDER BY and LIMIT (MySQL/SQLite)
./bin/sqlparser -sql "UPDATE users SET status = 'inactive' WHERE last_login < '2020-01-01' ORDER BY last_login LIMIT 100" -dialect mysql
```

#### DELETE Statements
```bash
# Simple DELETE
./bin/sqlparser -sql "DELETE FROM users WHERE id = 1" -dialect mysql -output table

# DELETE with complex WHERE
./bin/sqlparser -sql "DELETE FROM logs WHERE created_at < '2020-01-01' AND level = 'debug'" -dialect postgresql

# DELETE with ORDER BY and LIMIT (MySQL/SQLite)
./bin/sqlparser -sql "DELETE FROM logs WHERE level = 'debug' ORDER BY created_at LIMIT 1000" -dialect mysql
```

### Get Optimization Suggestions

#### Basic Optimization Analysis
```bash
./bin/sqlparser -sql "SELECT * FROM users WHERE UPPER(email) = 'TEST'" -dialect mysql -output table
```

#### Dialect-Specific Optimizations
```bash
# MySQL: LIMIT without ORDER BY warning
./bin/sqlparser -sql "SELECT name FROM users LIMIT 10" -dialect mysql

# SQL Server: Suggest TOP instead of LIMIT
./bin/sqlparser -sql "SELECT name FROM users LIMIT 10" -dialect sqlserver

# PostgreSQL: JSON vs JSONB suggestions
./bin/sqlparser -sql "SELECT data FROM logs WHERE json_extract(data, '$.type') = 'error'" -dialect postgresql
```

#### Advanced Subquery Support

SQLens provides comprehensive subquery support across all SQL dialects:

#### Subqueries in WHERE Clause
```bash
# Scalar subquery
./bin/sqlparser -sql "SELECT * FROM users WHERE salary > (SELECT AVG(salary) FROM employees)" -dialect postgresql

# EXISTS subquery
./bin/sqlparser -sql "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)" -dialect mysql

# NOT EXISTS subquery
./bin/sqlparser -sql "DELETE FROM users WHERE NOT EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)" -dialect mysql

# IN with subquery
./bin/sqlparser -sql "SELECT name FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > 1000)" -dialect postgresql

# NOT IN with subquery
./bin/sqlparser -sql "SELECT * FROM users WHERE id NOT IN (SELECT user_id FROM banned_users)" -dialect mysql
```

#### Subqueries in SELECT Clause
```bash
# Scalar subquery in SELECT
./bin/sqlparser -sql "SELECT id, name, (SELECT COUNT(*) FROM orders WHERE user_id = users.id) as order_count FROM users" -dialect postgresql

# Multiple scalar subqueries
./bin/sqlparser -sql "SELECT id, (SELECT AVG(price) FROM products) as avg_price, (SELECT MAX(price) FROM products) as max_price FROM users" -dialect mysql
```

#### Derived Tables (Subqueries in FROM Clause)
```bash
# Simple derived table
./bin/sqlparser -sql "SELECT * FROM (SELECT id, name FROM users WHERE active = 1) AS active_users WHERE id > 100" -dialect postgresql

# JOIN with derived table
./bin/sqlparser -sql "SELECT u.name, o.total FROM users u JOIN (SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id) o ON u.id = o.user_id" -dialect mysql

# Multiple derived tables
./bin/sqlparser -sql "SELECT * FROM (SELECT id FROM users) u, (SELECT user_id FROM orders) o WHERE u.id = o.user_id" -dialect postgresql
```

#### Subqueries in INSERT/UPDATE/DELETE
```bash
# INSERT with subquery in VALUES
./bin/sqlparser -sql "INSERT INTO user_stats (user_id, order_count) VALUES (1, (SELECT COUNT(*) FROM orders WHERE user_id = 1))" -dialect mysql

# INSERT ... SELECT
./bin/sqlparser -sql "INSERT INTO archive SELECT * FROM users WHERE created_at < '2020-01-01'" -dialect postgresql

# UPDATE with subquery in SET
./bin/sqlparser -sql "UPDATE users SET status = (SELECT status FROM user_preferences WHERE user_id = users.id)" -dialect postgresql

# UPDATE with subquery in WHERE
./bin/sqlparser -sql "UPDATE users SET active = 0 WHERE id IN (SELECT user_id FROM banned_users)" -dialect mysql

# DELETE with EXISTS
./bin/sqlparser -sql "DELETE FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id AND total > 10000)" -dialect postgresql
```

#### Complex Nested and Correlated Subqueries
```bash
# Triple nested subquery
./bin/sqlparser -sql "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE product_id IN (SELECT id FROM products WHERE category_id IN (SELECT id FROM categories WHERE name = 'electronics')))" -dialect mysql

# Correlated subquery
./bin/sqlparser -sql "SELECT * FROM users u WHERE salary > (SELECT AVG(salary) FROM employees e WHERE e.department = u.department)" -dialect postgresql

# Multiple EXISTS clauses
./bin/sqlparser -sql "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id) AND EXISTS (SELECT 1 FROM payments WHERE user_id = users.id)" -dialect mysql
```

**Example Optimization Output:**
```
┌─────────────────────────────────────┬──────────┬────────────────────────────────┬─────────────────────────┐
│ TYPE                                │ SEVERITY │ DESCRIPTION                    │ SUGGESTION              │
├─────────────────────────────────────┼──────────┼────────────────────────────────┼─────────────────────────┤
│ 🔍 SELECT_STAR                      │ WARNING  │ Avoid SELECT * for performance │ Specify explicit columns│
│ ⚠️  INEFFICIENT_SUBQUERY            │ INFO     │ Subquery may be optimized      │ Consider JOIN instead   │
│ 🚀 MYSQL_LIMIT_WITHOUT_ORDER        │ WARNING  │ LIMIT without ORDER BY         │ Add ORDER BY clause     │
└─────────────────────────────────────┴──────────┴────────────────────────────────┴─────────────────────────┘
```

#### DDL Support (CREATE, DROP, ALTER, INDEX)

SQLens provides full DDL (Data Definition Language) support across all dialects:

```bash
# CREATE TABLE - Simple
./bin/sqlparser -sql "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255) UNIQUE)" -dialect mysql

# CREATE TABLE - IF NOT EXISTS
./bin/sqlparser -sql "CREATE TABLE IF NOT EXISTS products (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, price DECIMAL(10,2) DEFAULT 0.00)" -dialect mysql

# CREATE TABLE - Foreign Keys
./bin/sqlparser -sql "CREATE TABLE orders (id INT PRIMARY KEY, user_id INT NOT NULL, product_id INT, FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE, FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE SET NULL)" -dialect postgresql

# CREATE TABLE - Composite Primary Key
./bin/sqlparser -sql "CREATE TABLE user_roles (user_id INT, role_id INT, PRIMARY KEY (user_id, role_id))" -dialect mysql

# DROP TABLE
./bin/sqlparser -sql "DROP TABLE IF EXISTS users" -dialect postgresql

# DROP DATABASE
./bin/sqlparser -sql "DROP DATABASE IF EXISTS test_db" -dialect mysql

# DROP INDEX
./bin/sqlparser -sql "DROP INDEX IF EXISTS idx_users_email" -dialect postgresql

# ALTER TABLE - ADD COLUMN
./bin/sqlparser -sql "ALTER TABLE users ADD COLUMN age INT NOT NULL" -dialect mysql

# ALTER TABLE - DROP COLUMN
./bin/sqlparser -sql "ALTER TABLE users DROP COLUMN age" -dialect postgresql

# ALTER TABLE - MODIFY COLUMN
./bin/sqlparser -sql "ALTER TABLE users MODIFY COLUMN name VARCHAR(150) NOT NULL" -dialect mysql

# ALTER TABLE - ADD CONSTRAINT
./bin/sqlparser -sql "ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)" -dialect postgresql

# ALTER TABLE - DROP CONSTRAINT
./bin/sqlparser -sql "ALTER TABLE orders DROP CONSTRAINT fk_user" -dialect postgresql

# CREATE INDEX
./bin/sqlparser -sql "CREATE INDEX idx_users_email ON users (email)" -dialect mysql

# CREATE UNIQUE INDEX
./bin/sqlparser -sql "CREATE UNIQUE INDEX idx_users_email ON users (email)" -dialect postgresql

# CREATE INDEX - IF NOT EXISTS
./bin/sqlparser -sql "CREATE INDEX IF NOT EXISTS idx_products_category ON products (category)" -dialect mysql

# CREATE INDEX - Multiple columns
./bin/sqlparser -sql "CREATE INDEX idx_orders_user_product ON orders (user_id, product_id)" -dialect postgresql
```

**Dialect-Specific DDL Features:**
- **MySQL**: `AUTO_INCREMENT`, `CHANGE COLUMN`
- **PostgreSQL**: `IF EXISTS`, `CASCADE` on DROP
- **SQLite**: `AUTOINCREMENT`, `IF NOT EXISTS`
- **SQL Server**: `IDENTITY`
- **Oracle**: Enterprise-grade DDL support

#### Transaction Support

SQLens provides full transaction control support across all dialects:

```bash
# BEGIN TRANSACTION
./bin/sqlparser -sql "BEGIN TRANSACTION" -dialect sqlserver

# START TRANSACTION (MySQL/PostgreSQL)
./bin/sqlparser -sql "START TRANSACTION" -dialect mysql

# COMMIT
./bin/sqlparser -sql "COMMIT" -dialect postgresql

# COMMIT WORK
./bin/sqlparser -sql "COMMIT WORK" -dialect postgresql

# ROLLBACK
./bin/sqlparser -sql "ROLLBACK" -dialect mysql

# SAVEPOINT
./bin/sqlparser -sql "SAVEPOINT my_savepoint" -dialect postgresql

# ROLLBACK TO SAVEPOINT
./bin/sqlparser -sql "ROLLBACK TO SAVEPOINT my_savepoint" -dialect mysql

# RELEASE SAVEPOINT
./bin/sqlparser -sql "RELEASE SAVEPOINT my_savepoint" -dialect postgresql
```

**Transaction Features:**
- **BEGIN/START TRANSACTION**: Start a new transaction
- **COMMIT**: Commit the current transaction
- **ROLLBACK**: Roll back the current transaction
- **SAVEPOINT**: Create a savepoint within a transaction
- **ROLLBACK TO SAVEPOINT**: Roll back to a specific savepoint
- **RELEASE SAVEPOINT**: Release a savepoint (PostgreSQL/MySQL)

#### Schema-Aware Parsing & Validation

SQLens can validate SQL queries against database schemas to catch errors before execution:

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

**Schema Features:**
- **Load schemas** from JSON or YAML files
- **Validate tables** - Detect non-existent tables
- **Validate columns** - Detect non-existent or mistyped column names
- **Type checking** - Validate data type compatibility in expressions
- **Foreign key validation** - Verify foreign key references

**Example schema (JSON):**
```json
{
  "name": "my_database",
  "tables": [
    {
      "name": "users",
      "columns": [
        {"name": "id", "type": "INT", "primary_key": true},
        {"name": "name", "type": "VARCHAR", "length": 100},
        {"name": "email", "type": "VARCHAR", "length": 255, "unique": true}
      ]
    }
  ]
}
```

#### Execution Plan Analysis

Analyze query execution plans to identify performance bottlenecks:

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
	fmt.Printf("With ANALYZE: %v\n", explainStmt.Analyze)

	// Analyze execution plan (from database output)
	jsonPlan := []byte(`{"Plan": {...}}`) // From EXPLAIN FORMAT=JSON
	executionPlan, _ := plan.ParseJSONPlan(jsonPlan, "postgresql")

	// Get performance analysis
	analyzer := plan.NewPlanAnalyzer("postgresql")
	analysis := analyzer.AnalyzePlan(executionPlan)

	fmt.Printf("Performance Score: %.2f/100\n", analysis.PerformanceScore)

	// Review issues
	for _, issue := range analysis.Issues {
		fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Type, issue.Description)
	}

	// Get recommendations
	for _, rec := range analysis.Recommendations {
		fmt.Printf("[%s] %s\n", rec.Priority, rec.Description)
	}
}
```

**Execution Plan Features:**
- **EXPLAIN parsing** - Full support for EXPLAIN and EXPLAIN ANALYZE
- **Multi-dialect** - MySQL, PostgreSQL, SQL Server, SQLite
- **Bottleneck detection** - Automatic identification of performance issues
- **Optimization suggestions** - Actionable recommendations with priority levels
- **Performance scoring** - 0-100 score based on plan quality
- **Cost analysis** - Startup and total cost estimation

### Parse SQL Server Logs

```bash
./bin/sqlparser -log examples/logs/sample_profiler.log -output table -verbose
```

### Command Line Options

- `-query FILE`: Analyze SQL query from file
- `-sql STRING`: Analyze SQL query from string
- `-log FILE`: Parse SQL Server log file
- `-output FORMAT`: Output format (json, table) - default: json
- `-dialect DIALECT`: SQL dialect (mysql, postgresql, sqlserver, sqlite, oracle) - default: sqlserver
- `-verbose`: Enable verbose output
- `-config FILE`: Configuration file path
- `-help`: Show help

## Example Output

### Query Analysis (JSON)
```json
{
  "analysis": {
    "tables": [
      {
        "name": "users",
        "alias": "u",
        "usage": "SELECT"
      },
      {
        "name": "orders", 
        "alias": "o",
        "usage": "SELECT"
      }
    ],
    "columns": [
      {
        "table": "u",
        "name": "name",
        "usage": "SELECT"
      },
      {
        "table": "o", 
        "name": "total",
        "usage": "SELECT"
      }
    ],
    "joins": [
      {
        "type": "INNER",
        "right_table": "orders",
        "condition": "(u.id = o.user_id)"
      }
    ],
    "query_type": "SELECT",
    "complexity": 4
  },
  "suggestions": [
    {
      "type": "COMPLEX_QUERY",
      "description": "Query involves many tables. Consider breaking into smaller queries.",
      "severity": "INFO"
    }
  ]
}
```

### Query Analysis (Table)
```
=== SQL Query Analysis ===
Query Type: SELECT
Complexity: 4

Tables:
Name                 Schema     Alias      Usage
------------------------------------------------------------
users                           u          SELECT
orders                          o          SELECT

Columns:
Name                 Table      Usage
----------------------------------------
name                 u          SELECT
total                o          SELECT

Joins:
Type       Left Table      Right Table     Condition
------------------------------------------------------------
INNER                      orders          (u.id = o.user_id)
```

## Configuration

You can use a configuration file to customize the behavior:

```yaml
# config.yaml
parser:
  strict_mode: false
  max_query_size: 1000000
  dialect: "sqlserver"

analyzer:
  enable_optimizations: true
  complexity_threshold: 10
  detailed_analysis: true

logger:
  default_format: "profiler"
  max_file_size_mb: 100
  filters:
    min_duration_ms: 0
    exclude_system: true

output:
  format: "json"
  pretty_json: true
  include_timestamps: true
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Running Examples

```bash
# Analyze complex query
make dev-query

# Analyze simple query  
make dev-simple

# Parse log file
make dev-log
```

### Code Quality

```bash
# Format code
make fmt

# Lint code
make lint

# Security check
make security
```

## Architecture

The project follows a modular architecture:

```
sql-parser-go/
├── cmd/sqlparser/          # CLI application
├── pkg/
│   ├── lexer/             # SQL tokenization
│   ├── parser/            # SQL parsing and AST
│   ├── analyzer/          # Query analysis
│   └── logger/            # Log parsing
├── internal/config/        # Configuration management
├── examples/              # Example queries and logs
└── tests/                 # Test files
```

### Key Components

1. **Lexer**: Tokenizes SQL text into tokens
2. **Parser**: Builds Abstract Syntax Tree from tokens  
3. **Analyzer**: Extracts metadata and provides insights
4. **Logger**: Parses various SQL Server log formats

## Supported SQL Features

### Core SQL Statements
- SELECT statements with complex joins
- WHERE, GROUP BY, HAVING, ORDER BY clauses
- **Full Subquery Support**: WHERE, SELECT, FROM (derived tables), INSERT, UPDATE, DELETE
- **Subquery Types**: Scalar subqueries, EXISTS, NOT EXISTS, IN, NOT IN, correlated subqueries
- **Derived Tables**: Subqueries in FROM clause with JOIN support
- Functions and expressions
- INSERT, UPDATE, DELETE statements with full DML support

### DDL (Data Definition Language) ✨
- **CREATE TABLE**
  ```sql
  CREATE TABLE IF NOT EXISTS users (
      id INT AUTO_INCREMENT PRIMARY KEY,
      name VARCHAR(100) NOT NULL,
      email VARCHAR(255) UNIQUE,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );
  ```
- **DROP Statements**
  ```sql
  DROP TABLE IF EXISTS users;
  DROP DATABASE IF EXISTS test_db;
  DROP INDEX IF EXISTS idx_users_email;
  ```
- **ALTER TABLE**
  ```sql
  ALTER TABLE users ADD COLUMN age INT;
  ALTER TABLE users DROP COLUMN age;
  ALTER TABLE users MODIFY COLUMN name VARCHAR(150);
  ALTER TABLE users ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id);
  ```
- **CREATE INDEX**
  ```sql
  CREATE UNIQUE INDEX idx_users_email ON users (email);
  CREATE INDEX IF NOT EXISTS idx_products_category ON products (category);
  ```

### Transaction Control ✨
- **Transaction Statements**
  ```sql
  -- Begin a transaction
  BEGIN TRANSACTION;
  START TRANSACTION;  -- MySQL/PostgreSQL

  -- Commit transaction
  COMMIT;
  COMMIT WORK;  -- Optional WORK keyword

  -- Rollback transaction
  ROLLBACK;
  ROLLBACK WORK;

  -- Savepoints
  SAVEPOINT my_savepoint;
  ROLLBACK TO SAVEPOINT my_savepoint;
  RELEASE SAVEPOINT my_savepoint;  -- PostgreSQL/MySQL
  ```

### Advanced Features ✨
- **CTEs (Common Table Expressions)**
  ```sql
  WITH sales_summary AS (
      SELECT product_id, SUM(amount) as total
      FROM sales GROUP BY product_id
  )
  SELECT * FROM sales_summary WHERE total > 1000;
  ```

- **Window Functions**
  ```sql
  SELECT
      employee_id,
      salary,
      ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) as rank,
      AVG(salary) OVER (ORDER BY hire_date ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) as moving_avg
  FROM employees;
  ```

- **Set Operations**
  ```sql
  SELECT id FROM customers
  UNION ALL
  SELECT id FROM prospects
  INTERSECT
  SELECT id FROM active_accounts;
  ```

- **Comprehensive Subquery Support** ✨
  ```sql
  -- Scalar subquery
  SELECT * FROM users WHERE salary > (SELECT AVG(salary) FROM employees);

  -- EXISTS / NOT EXISTS
  SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id);
  DELETE FROM users WHERE NOT EXISTS (SELECT 1 FROM orders WHERE user_id = users.id);

  -- IN / NOT IN with subquery
  SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > 1000);

  -- Derived tables (subqueries in FROM)
  SELECT * FROM (SELECT id, name FROM users WHERE active = 1) AS active_users;

  -- JOIN with derived table
  SELECT u.name, o.total
  FROM users u
  JOIN (SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id) o
  ON u.id = o.user_id;

  -- Nested subqueries
  SELECT * FROM users
  WHERE id IN (
      SELECT user_id FROM orders
      WHERE product_id IN (SELECT id FROM products WHERE category = 'electronics')
  );

  -- Correlated subquery
  SELECT * FROM users u
  WHERE salary > (SELECT AVG(salary) FROM employees e WHERE e.department = u.department);
  ```

See [examples/queries/](examples/queries/) for more comprehensive examples.

## Supported Log Formats

- SQL Server Profiler traces
- Extended Events
- Query Store exports
- SQL Server Error Logs
- Performance Counter logs

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make all` to ensure code quality
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Roadmap

- [x] Support for more SQL dialects (MySQL, PostgreSQL, SQLite, Oracle, SQL Server)
- [x] Dialect-specific identifier quoting and keyword recognition
- [x] Multi-dialect CLI interface with dialect selection
- [x] **Advanced optimization suggestions** ✅ **COMPLETED!**
- [x] **Dialect-specific optimization recommendations** ✅ **COMPLETED!**
- [x] **Comprehensive Subquery Support** ✅ **COMPLETED!**
  - [x] **Scalar Subqueries** - In WHERE, SELECT, INSERT VALUES, UPDATE SET
  - [x] **EXISTS / NOT EXISTS** - Full support in all statement types
  - [x] **IN / NOT IN with Subqueries** - Complete implementation
  - [x] **Derived Tables** - Subqueries in FROM clause with JOIN support
  - [x] **Nested & Correlated Subqueries** - Multiple levels of nesting
- [x] **Extended SQL Features** ✅ **COMPLETED!**
  - [x] **CTEs (WITH clause)** - Simple and multiple CTEs with column lists
  - [x] **Window Functions** - ROW_NUMBER, RANK, PARTITION BY, ORDER BY, frame clauses
  - [x] **Set Operations** - UNION, UNION ALL, INTERSECT, EXCEPT
  - [x] **CASE Expressions** - Searched and simple CASE statements
- [x] **DML Statements** ✅ **COMPLETED!**
  - [x] **INSERT** - VALUES, multiple rows, INSERT...SELECT
  - [x] **UPDATE** - Multiple columns, WHERE, ORDER BY/LIMIT (MySQL/SQLite)
  - [x] **DELETE** - WHERE, ORDER BY/LIMIT (MySQL/SQLite)
- [x] **DDL Support** ✅ **COMPLETED!**
  - [x] **CREATE TABLE** - Columns, constraints, foreign keys, IF NOT EXISTS
  - [x] **DROP** - TABLE/DATABASE/INDEX with IF EXISTS and CASCADE
  - [x] **ALTER TABLE** - ADD/DROP/MODIFY/CHANGE columns and constraints
  - [x] **CREATE INDEX** - Simple and unique indexes with IF NOT EXISTS
  - [x] **Dialect-specific features** - AUTO_INCREMENT, IDENTITY, AUTOINCREMENT
- [x] **Transaction Support** ✅ **COMPLETED!**
  - [x] **BEGIN/START TRANSACTION** - Start transactions
  - [x] **COMMIT/ROLLBACK** - Commit or rollback transactions
  - [x] **SAVEPOINT** - Create savepoints within transactions
  - [x] **ROLLBACK TO SAVEPOINT** - Roll back to specific savepoints
  - [x] **RELEASE SAVEPOINT** - Release savepoints (PostgreSQL/MySQL)
- [x] Performance benchmarking
- [ ] Query execution plan analysis
- [ ] Real-time log monitoring
- [ ] Integration with monitoring tools
- [ ] Schema-aware parsing and validation
- [ ] Stored procedure parsing
- [ ] Trigger parsing
- [ ] View definitions (CREATE VIEW)

## Acknowledgments

- Inspired by various SQL parsing libraries
- Built with Go's excellent standard library
- Uses minimal external dependencies for better maintainability

## 🚀 Performance Optimizations

This project has been heavily optimized for production use with Go's strengths in mind:

### Key Performance Features

- **Sub-millisecond parsing**: Parse queries in <1ms
- **Multi-dialect optimization**: Optimized lexing and parsing for each SQL dialect
- **Memory efficient**: Uses only ~200KB-7KB depending on dialect complexity
- **Concurrent processing**: Multi-core analysis support
- **Zero-allocation paths**: Optimized hot paths for identifier quoting

### Multi-Dialect Performance Benchmarks

**Tested on Apple M2 Pro:**

#### Lexing Performance (ns/op | MB/s)
```
SQL Server:   2,492 ns/op  | 200.24 MB/s   (bracket parsing - fastest!)
SQLite:       2,900 ns/op  | 163.44 MB/s   (lightweight parsing)
Oracle:       3,620 ns/op  | 137.85 MB/s   (enterprise parsing)
PostgreSQL:   8,736 ns/op  |  56.32 MB/s   (double quote parsing)
MySQL:       16,708 ns/op  |  28.55 MB/s   (complex backtick parsing)
```

#### Parsing Performance (ns/op | MB/s)
```
SQL Server:    375.9 ns/op |1327.54 MB/s   (🚀 ultra-fast!)
Oracle:      1,315 ns/op   | 379.61 MB/s   
SQLite:      1,248 ns/op   | 379.77 MB/s   
PostgreSQL:  2,753 ns/op   | 178.71 MB/s   
MySQL:       4,887 ns/op   |  97.60 MB/s   
```

#### Memory Usage (per operation)
```
SQL Server:   704 B/op,  8 allocs/op   (most efficient)
SQLite:     3,302 B/op, 25 allocs/op   
Oracle:     3,302 B/op, 25 allocs/op   
PostgreSQL: 4,495 B/op, 27 allocs/op   
MySQL:      7,569 B/op, 27 allocs/op   (complex syntax overhead)
```

#### Feature Operations (ultra-fast)
```
Identifier Quoting:    ~154-160 ns/op (all dialects)
Feature Support:     ~18-27 ns/op    (all dialects)
Keyword Lookup:   2,877-43,984 ns/op (varies by dialect complexity)
```

### Advanced SQL Features Performance

**Tested on Apple M2 Pro:**

#### Subqueries & Advanced Features (μs/op)
```
Simple Scalar Subqueries:        8-10 μs    ✅ Sub-10 microseconds!
EXISTS/NOT EXISTS:               22 μs      ✅ Fast predicate checks
Nested Subqueries (3 levels):    19 μs      ✅ Excellent scaling
Correlated Subqueries:           39 μs      ✅ Production-ready
Derived Tables (FROM):           22 μs      ✅ Efficient JOIN alternatives

CTEs (WITH clause):              14-80 μs   ✅ Single/Multiple CTEs
Window Functions:                12-32 μs   ✅ ROW_NUMBER, RANK, PARTITION BY
Set Operations:                  3-11 μs    ✅ UNION, INTERSECT, EXCEPT
```

#### DDL Operations Performance (μs/op)
```
CREATE TABLE (simple):           24 μs      ✅ Fast schema creation
CREATE TABLE (complex FK):       78-111 μs  ✅ Multi-constraint support
DROP TABLE/DATABASE/INDEX:       1.7-3.7 μs ✅ Blazing fast!
ALTER TABLE (ADD/DROP):          3-14 μs    ✅ Quick schema changes
CREATE INDEX:                    1.4-11 μs  ✅ Efficient indexing
```

#### DML with Subqueries (μs/op)
```
INSERT ... SELECT:               5 μs       ✅ Very fast bulk operations
INSERT with Subquery:            22 μs      ✅ Dynamic value insertion
UPDATE with Subquery:            38 μs      ✅ Complex updates
DELETE with EXISTS:              10 μs      ✅ Conditional deletion
```

#### Transaction Operations (ns/op)
```
BEGIN/START TRANSACTION:         200 ns     ✅ Ultra-fast transaction start
COMMIT:                          149 ns     ✅ Lightning-fast commits
ROLLBACK:                        173 ns     ✅ Fast rollbacks
SAVEPOINT:                       3.6 μs     ✅ Efficient savepoint creation
ROLLBACK TO SAVEPOINT:           3.0 μs     ✅ Quick savepoint rollback
RELEASE SAVEPOINT:               1.7 μs     ✅ Fast savepoint release
```

#### Schema-Aware Validation (ns/op)
```
Schema Loading (JSON):           7.2 μs     ✅ Fast schema loading
Validate SELECT:                 264 ns     ✅ Ultra-fast validation
Validate INSERT:                 155 ns     ✅ Lightning-fast checks
Validate UPDATE:                 170 ns     ✅ Quick validation
Type Checking:                   590 ns     ✅ Sub-microsecond type checks
Complex Validation (JOIN):       1.1 μs     ✅ Fast multi-table validation
```

**Memory Efficiency:**
- Simple queries: **8-20 KB** per operation
- Complex queries: **40-80 KB** per operation
- DDL operations: **4-200 KB** depending on complexity
- Transaction operations: **337-7360 B** per operation
- Schema validation: **0-504 B** per operation (zero-allocation for simple queries!)

### Real-world Performance

- **🏆 Best overall**: SQL Server (375ns parsing, 1.3GB/s throughput)
- **🥇 Best lexing**: SQL Server bracket parsing at 200MB/s
- **🥈 Most balanced**: PostgreSQL (fast + memory efficient)
- **🥉 Most features**: MySQL (comprehensive but slower due to complexity)

### Performance Insights

1. **SQL Server dominance**: Bracket parsing is extremely efficient
2. **PostgreSQL efficiency**: Great balance of speed and memory usage
3. **MySQL complexity**: Feature-rich but higher memory overhead
4. **SQLite optimization**: Lightweight and fast for embedded use
5. **Oracle enterprise**: Robust performance for complex queries
6. **🆕 Advanced features**: Sub-millisecond parsing for 95%+ of queries
7. **🆕 DDL operations**: Ultra-fast with DROP < 2μs, CREATE TABLE < 25μs
8. **🆕 Subqueries**: Excellent scaling even with deep nesting (3+ levels)
9. **🆕 Transactions**: Sub-microsecond COMMIT/ROLLBACK, ~3μs for savepoints

**This is production-ready performance that matches or exceeds commercial SQL parsers across all major dialects!**
