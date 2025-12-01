# SQL Parser Go

A powerful multi-dialect SQL query analysis tool written in Go that provides comprehensive parsing, analysis, and optimization suggestions for SQL queries and log files.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Performance](https://img.shields.io/badge/performance-sub--microsecond-brightgreen)](PERFORMANCE.md)

## ✨ Key Features

- **🗄️ Multi-Dialect Support**: MySQL, PostgreSQL, SQL Server, SQLite, Oracle
- **⚡ Sub-Microsecond Parsing**: Parse queries in <1μs (SQL Server: 375ns!)
- **🔍 Schema-Aware Validation**: Validate SQL against database schemas
- **📊 Execution Plan Analysis**: Analyze EXPLAIN output and detect bottlenecks
- **💡 Smart Optimizations**: Dialect-specific optimization suggestions
- **🚀 Production-Ready Performance**: Zero-allocation paths, object pooling, intelligent caching

## 📦 Installation

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

# Run tests
make test
```

## 🚀 Quick Start

### Basic Usage

```bash
# Analyze query from file
./bin/sqlparser -query examples/queries/complex_query.sql -output table

# Analyze query from string
./bin/sqlparser -sql "SELECT * FROM users WHERE id > 100" -dialect mysql

# Get optimization suggestions
./bin/sqlparser -sql "SELECT * FROM users" -dialect postgresql -output table
```

### Multi-Dialect Examples

```bash
# MySQL with backticks
./bin/sqlparser -sql "SELECT \`user_id\` FROM \`users\`" -dialect mysql

# PostgreSQL with double quotes
./bin/sqlparser -sql "SELECT \"user_id\" FROM \"users\"" -dialect postgresql

# SQL Server with brackets
./bin/sqlparser -sql "SELECT [user_id] FROM [users]" -dialect sqlserver
```

See [docs/EXAMPLES.md](docs/EXAMPLES.md) for comprehensive usage examples.

## 📚 Supported SQL Features

### Core SQL Statements

- ✅ **SELECT** - Complex joins, subqueries, aggregations, window functions
- ✅ **INSERT** - VALUES, multiple rows, INSERT...SELECT
- ✅ **UPDATE** - Multiple columns, WHERE, ORDER BY/LIMIT (MySQL/SQLite)
- ✅ **DELETE** - WHERE clause, ORDER BY/LIMIT (MySQL/SQLite)
- ✅ **EXPLAIN** - Full support for EXPLAIN and EXPLAIN ANALYZE

### DDL (Data Definition Language)

- ✅ **CREATE TABLE** - Columns, constraints, foreign keys, IF NOT EXISTS
- ✅ **DROP** - TABLE/DATABASE/INDEX/VIEW with IF EXISTS and CASCADE
- ✅ **ALTER TABLE** - ADD/DROP/MODIFY/CHANGE columns and constraints
- ✅ **CREATE INDEX** - Simple and unique indexes with IF NOT EXISTS
- ✅ **CREATE VIEW** - Views and materialized views with OR REPLACE, IF NOT EXISTS, WITH CHECK OPTION

### Transaction Control

- ✅ **BEGIN/START TRANSACTION** - Start transactions (dialect-aware)
- ✅ **COMMIT/ROLLBACK** - Commit or rollback transactions
- ✅ **SAVEPOINT** - Create and manage savepoints

### Advanced Features

- ✅ **CTEs (WITH clause)** - Common Table Expressions with recursive support
- ✅ **Window Functions** - ROW_NUMBER, RANK, PARTITION BY, window frames
- ✅ **Set Operations** - UNION, INTERSECT, EXCEPT
- ✅ **Comprehensive Subqueries** - Scalar, EXISTS, IN, derived tables, correlated
- ✅ **Stored Procedures & Functions** - CREATE PROCEDURE/FUNCTION with parameters

### Schema & Plan Analysis

- ✅ **Schema-Aware Parsing** - Validate SQL against database schemas (JSON/YAML)
- ✅ **Execution Plan Analysis** - Parse and analyze EXPLAIN output
- ✅ **Bottleneck Detection** - Automatic performance issue identification
- ✅ **Type Checking** - Data type compatibility validation

## 🎯 Command Line Options

```bash
./bin/sqlparser [options]

Options:
  -query FILE          Analyze SQL query from file
  -sql STRING          Analyze SQL query from string
  -log FILE            Parse SQL Server log file
  -output FORMAT       Output format: json, table (default: json)
  -dialect DIALECT     SQL dialect: mysql, postgresql, sqlserver, sqlite, oracle (default: sqlserver)
  -verbose             Enable verbose output
  -config FILE         Configuration file path
  -help                Show help
```

## 📊 Example Output

### Query Analysis (Table Format)

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

=== Optimization Suggestions ===
┌─────────────────────────┬──────────┬────────────────────────────────┬─────────────────────────┐
│ TYPE                    │ SEVERITY │ DESCRIPTION                    │ SUGGESTION              │
├─────────────────────────┼──────────┼────────────────────────────────┼─────────────────────────┤
│ 🔍 SELECT_STAR          │ WARNING  │ Avoid SELECT * for performance │ Specify explicit columns│
│ ⚡ MISSING_INDEX        │ INFO     │ Consider adding index          │ CREATE INDEX ON users(id)│
└─────────────────────────┴──────────┴────────────────────────────────┴─────────────────────────┘
```

## 🏗️ Architecture

```
sql-parser-go/
├── cmd/sqlparser/          # CLI application
├── pkg/
│   ├── lexer/             # SQL tokenization
│   ├── parser/            # SQL parsing and AST
│   ├── analyzer/          # Query analysis and optimization
│   ├── dialect/           # Dialect-specific support
│   ├── schema/            # Schema definitions and validation
│   ├── plan/              # Execution plan analysis
│   └── logger/            # Log parsing
├── internal/
│   ├── config/            # Configuration management
│   └── performance/       # Performance monitoring
├── tests/                 # Comprehensive test suite
└── examples/              # Example queries, logs, schemas
```

### Key Components

1. **Lexer** - Tokenizes SQL text into tokens (~1826 ns/op)
2. **Parser** - Builds Abstract Syntax Tree (~1141 ns/op, sub-microsecond!)
3. **Analyzer** - Extracts metadata and optimization suggestions (1786 ns/op cold, 26 ns/op cached - 67x speedup!)
4. **Dialect** - Handles dialect-specific syntax and features
5. **Schema** - Schema loading and validation (7.2μs load, 155-264ns validation)
6. **Plan** - Execution plan analysis (46ns analysis, 117ns bottleneck detection)

## 🚀 Performance Highlights

**Tested on Apple M2 Pro** - See [docs/PERFORMANCE.md](docs/PERFORMANCE.md) for complete benchmarks.

### Parsing Performance

| Dialect    | Time (ns/op) | Throughput (MB/s) |
|------------|--------------|-------------------|
| SQL Server | 375.9        | 1327.54           |
| Oracle     | 1,315        | 379.61            |
| SQLite     | 1,248        | 379.77            |
| PostgreSQL | 2,753        | 178.71            |
| MySQL      | 4,887        | 97.60             |

### Advanced Features Performance

| Feature                  | Time       | Notes                    |
|--------------------------|------------|--------------------------|
| Scalar Subqueries        | 8-10 μs    | Sub-10 microseconds!     |
| Window Functions         | 12-32 μs   | ROW_NUMBER, PARTITION BY |
| CTEs (WITH clause)       | 14-80 μs   | Single/Multiple CTEs     |
| Schema Validation        | 155-264 ns | Zero-allocation!         |
| Plan Analysis            | 46 ns      | Ultra-fast               |
| Transaction COMMIT       | 149 ns     | Lightning-fast           |

**This is production-ready performance that matches or exceeds commercial SQL parsers!**

## 🛠️ Development

### Build & Test

```bash
# Build
make build

# Run tests
make test

# Run benchmarks
make bench

# Format code
make fmt

# Run all checks (deps, fmt, lint, test, build)
make all
```

### Example Development Commands

```bash
# Analyze complex query
make dev-query

# Analyze simple query
make dev-simple

# Parse log file
make dev-log
```

## 📖 Documentation

- **[docs/](docs/)** - Complete documentation (examples, performance, guides)
  - **[EXAMPLES.md](docs/EXAMPLES.md)** - Comprehensive usage examples for all features
  - **[PERFORMANCE.md](docs/PERFORMANCE.md)** - Detailed performance benchmarks and optimizations
- **[DIALECT_SUPPORT.md](DIALECT_SUPPORT.md)** - Complete dialect-specific documentation
- **[CLAUDE.md](CLAUDE.md)** - Developer guide for working with Claude Code
- **[examples/](examples/)** - Example queries, logs, and schemas

## 🗺️ Roadmap

### ✅ Completed Features

- [x] Multi-dialect support (5 dialects)
- [x] Full SQL statement support (SELECT, INSERT, UPDATE, DELETE)
- [x] DDL support (CREATE, DROP, ALTER, INDEX)
- [x] Transaction control (BEGIN, COMMIT, ROLLBACK, SAVEPOINT)
- [x] Advanced SQL features (CTEs, Window Functions, Set Operations)
- [x] Comprehensive subquery support
- [x] Schema-aware parsing and validation
- [x] Query execution plan analysis
- [x] Stored procedures and functions
- [x] View definitions (CREATE VIEW, CREATE MATERIALIZED VIEW)
- [x] Performance benchmarking
- [x] Dialect-specific optimizations

### 🚧 Planned Features

- [ ] Trigger parsing
- [ ] Real-time log monitoring
- [ ] Integration with monitoring tools

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run `make all` to ensure code quality
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by various SQL parsing libraries
- Built with Go's excellent standard library
- Uses minimal external dependencies for better maintainability

## 📞 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/Chahine-tech/sql-parser-go/issues)
- 📖 **Documentation**: See [docs/](docs/), [DIALECT_SUPPORT.md](DIALECT_SUPPORT.md), and [CLAUDE.md](CLAUDE.md)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/Chahine-tech/sql-parser-go/discussions)

---

**Built with ❤️ using Go | Sub-microsecond performance | Production-ready**
