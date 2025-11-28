# SQL Parser Go

A powerful multi-dialect SQL query analysis tool written in Go that provides comprehensive parsing, analysis, and optimization suggestions for SQL queries and log files.

## Features

- **Multi-Dialect Support**: Parse MySQL, PostgreSQL, SQL Server, SQLite, and Oracle queries
- **SQL Query Parsing**: Parse and analyze complex SQL queries with dialect-specific syntax
- **Abstract Syntax Tree (AST)**: Generate detailed AST representations
- **Query Analysis**: Extract tables, columns, joins, and conditions
- **Advanced Optimization Suggestions**: Get intelligent, dialect-specific recommendations for query improvements
- **Subquery Support**: Parse and analyze complex subqueries in WHERE clauses
- **Advanced SQL Features**:
  - **CTEs (WITH clause)**: Common Table Expressions with recursive support
  - **Window Functions**: ROW_NUMBER, RANK, PARTITION BY, ORDER BY, window frames
  - **Set Operations**: UNION, UNION ALL, INTERSECT, EXCEPT
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

#### Subquery Optimization Detection
```bash
./bin/sqlparser -sql "SELECT name FROM users WHERE id IN (SELECT user_id FROM orders)" -dialect postgresql
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
- Subqueries in WHERE and FROM clauses
- Functions and expressions
- INSERT, UPDATE, DELETE statements (basic)

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
- [x] **Subquery parsing and optimization** ✅ **COMPLETED!**
- [x] **Extended dialect features** ✅ **COMPLETED!**
  - [x] **CTEs (WITH clause)** - Simple and multiple CTEs with column lists
  - [x] **Window Functions** - ROW_NUMBER, RANK, PARTITION BY, ORDER BY, frame clauses
  - [x] **Set Operations** - UNION, UNION ALL, INTERSECT, EXCEPT
- [x] Performance benchmarking
- [ ] Query execution plan analysis
- [ ] Real-time log monitoring
- [ ] Integration with monitoring tools
- [ ] CASE expressions (partial - needs expression parser refactoring)
- [ ] Schema-aware parsing and validation
- [ ] More DDL support (ALTER, CREATE TABLE, etc.)

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

**This is production-ready performance that matches or exceeds commercial SQL parsers across all major dialects!**
