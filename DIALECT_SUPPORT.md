# SQL Dialect Support

This SQL Parser supports multiple SQL dialects with dialect-specific features and syntax parsing.

## Supported Dialects

### MySQL
- **Identifier quoting**: Backticks (`` ` ``)
- **Features**: JSON support, Common Table Expressions (CTE), Window Functions, Full-text search
- **Usage**: `--dialect mysql`

### PostgreSQL
- **Identifier quoting**: Double quotes (`"`)
- **Features**: Array support, JSON/JSONB, XML, RETURNING clause, Advanced window functions
- **Usage**: `--dialect postgresql` or `--dialect postgres`

### SQL Server
- **Identifier quoting**: Square brackets (`[]`)
- **Features**: XML support, TOP clause, MERGE statement, Advanced analytics
- **Usage**: `--dialect sqlserver` or `--dialect mssql`

### SQLite
- **Identifier quoting**: Double quotes (`"`)
- **Features**: JSON support (3.38+), Window functions (3.25+), Full-text search
- **Usage**: `--dialect sqlite`

### Oracle
- **Identifier quoting**: Double quotes (`"`)
- **Features**: XML support, Advanced partitioning, Hierarchical queries
- **Usage**: `--dialect oracle`

## Usage Examples

### Command Line
```bash
# MySQL with backtick identifiers
./sqlparser -sql "SELECT \`user_id\` FROM \`users\`" -dialect mysql

# PostgreSQL with double quotes
./sqlparser -sql "SELECT \"user_id\" FROM \"users\"" -dialect postgresql

# SQL Server with brackets
./sqlparser -sql "SELECT [user_id] FROM [users]" -dialect sqlserver
```

### Configuration File
```yaml
parser:
  dialect: mysql
  strict_mode: false
  max_query_size: 1000000
```

### Programmatic Usage
```go
import (
    "github.com/Chahine-tech/sql-parser-go/pkg/parser"
    "github.com/Chahine-tech/sql-parser-go/pkg/dialect"
)

// Create parser with specific dialect
d := dialect.GetDialect("mysql")
p := parser.NewWithDialect(ctx, sql, d)

// Parse with dialect-specific features
stmt, err := p.ParseStatement()
```

## Dialect-Specific Features

### Identifier Quoting
Each dialect has its own quoting style for identifiers:
- **MySQL**: `SELECT \`table\`.\`column\` FROM \`database\`.\`table\``
- **PostgreSQL**: `SELECT "table"."column" FROM "schema"."table"`
- **SQL Server**: `SELECT [table].[column] FROM [database].[schema].[table]`
- **SQLite**: `SELECT "table"."column" FROM "table"`
- **Oracle**: `SELECT "table"."column" FROM "schema"."table"`

### Advanced SQL Features (Newly Supported! ✨)

#### Common Table Expressions (CTEs)
All dialects now support WITH clause:
```sql
-- MySQL
WITH `sales_cte` AS (
    SELECT `product_id`, SUM(`amount`) as `total`
    FROM `sales` GROUP BY `product_id`
)
SELECT * FROM `sales_cte`;

-- PostgreSQL
WITH "revenue_cte" AS (
    SELECT "dept_id", SUM("revenue") as "total"
    FROM "departments" GROUP BY "dept_id"
)
SELECT * FROM "revenue_cte";

-- SQL Server
WITH [employee_cte] AS (
    SELECT [emp_id], [name], [salary]
    FROM [employees]
)
SELECT * FROM [employee_cte];
```

**Features:**
- Simple CTEs
- Multiple CTEs (comma-separated)
- CTEs with explicit column lists
- Recursive CTE support (keyword recognized)

#### Window Functions
Full support for window functions across all dialects:
```sql
-- Partition and ordering
SELECT
    employee_id,
    ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) as rank
FROM employees;

-- Window frames
SELECT
    date,
    SUM(amount) OVER (
        ORDER BY date
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) as running_total
FROM transactions;
```

**Supported window features:**
- `ROW_NUMBER()`, `RANK()`, `DENSE_RANK()`, and all aggregate functions
- `PARTITION BY` with multiple expressions
- `ORDER BY` within window
- Frame specifications: `ROWS` and `RANGE`
- Frame boundaries: `UNBOUNDED PRECEDING/FOLLOWING`, `CURRENT ROW`, expression-based offsets

#### Set Operations
All dialects support set operations:
- `UNION` - Combines and removes duplicates
- `UNION ALL` - Combines keeping duplicates
- `INTERSECT` - Returns common records
- `EXCEPT` - Returns records in first set but not in second

```sql
SELECT id FROM customers
UNION ALL
SELECT id FROM prospects
INTERSECT
SELECT id FROM active_users;
```

### Other Feature Support
- **JSON Support**: MySQL 5.7+, PostgreSQL, SQL Server 2016+, SQLite 3.38+
- **Array Support**: PostgreSQL only
- **XML Support**: PostgreSQL, SQL Server, Oracle

### Keyword Recognition
Each dialect has its own set of reserved keywords and functions that are recognized during parsing.

**Newly added keywords:** `WITH`, `RECURSIVE`, `OVER`, `PARTITION`, `ROWS`, `RANGE`, `UNBOUNDED`, `PRECEDING`, `FOLLOWING`, `CURRENT`, `ROW`, `INTERSECT`, `EXCEPT`, `CASE`, `WHEN`, `THEN`, `ELSE`, `END`

## Default Behavior
- **Default dialect**: SQL Server
- **Configuration**: Can be overridden via CLI flag (`-dialect`) or config file
- **Fallback**: Unknown dialects default to SQL Server syntax

## Error Handling
The parser will produce appropriate error messages when:
- Dialect-specific syntax is used with wrong dialect
- Unsupported features are used for a dialect
- Invalid identifier quoting is detected
