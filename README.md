# SQLens

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Performance](https://img.shields.io/badge/performance-sub--microsecond-brightgreen)](PERFORMANCE.md)

A multi-dialect SQL parser and analyzer written in Go. Parses queries into an AST, extracts metadata, suggests optimizations, validates against schemas, and monitors SQL log files in real time.

## Installation

**With Nix:**

```bash
nix run github:Chahine-tech/sql-parser-go -- --help
nix develop   # development environment
```

**From source** (requires Go 1.21+):

```bash
git clone https://github.com/Chahine-tech/sql-parser-go.git
cd sql-parser-go
make build
```

## Usage

```bash
# Analyze a query file
./bin/sqlparser -query examples/queries/complex_query.sql -output table

# Analyze an inline query
./bin/sqlparser -sql "SELECT * FROM users WHERE id > 100" -dialect mysql

# Parse a SQL Server log
./bin/sqlparser -log examples/logs/sample_profiler.log -output table

# Monitor a log file in real time
./bin/sqlparser -monitor /var/log/postgresql/postgresql.log -dialect postgresql
```

**Options:**

```
-query FILE      Read SQL from file
-sql STRING      Read SQL from string
-log FILE        Parse SQL Server log file
-monitor FILE    Watch log file for new entries
-dialect NAME    mysql | postgresql | sqlserver | sqlite | oracle  (default: sqlserver)
-output FORMAT   json | table  (default: json)
-verbose         Verbose output
-config FILE     Config file path
```

## What it supports

**Statements:** SELECT, INSERT, UPDATE, DELETE, EXPLAIN, all DDL (CREATE/DROP/ALTER TABLE, INDEX, VIEW, TRIGGER), transactions (BEGIN, COMMIT, ROLLBACK, SAVEPOINT), stored procedures and functions, MERGE, cursors, control flow (IF, WHILE, FOR, LOOP, REPEAT), exception handling (TRY/CATCH, RAISE, SIGNAL).

**Advanced SQL:** CTEs, window functions, set operations (UNION/INTERSECT/EXCEPT), scalar/correlated/derived subqueries, EXISTS/IN subqueries, PostgreSQL dollar-quoted strings.

**Analysis:** schema-aware validation (JSON/YAML schema files), execution plan parsing (EXPLAIN JSON/XML), bottleneck detection, dialect-specific optimization suggestions, real-time log monitoring with alert rules.

## Performance

Benchmarked on Apple M2 Pro:

| Component       | Cold         | Cached    |
|-----------------|--------------|-----------|
| Lexer           | ~1826 ns/op  | —         |
| Parser          | ~1141 ns/op  | —         |
| Analyzer        | ~1786 ns/op  | ~26 ns/op |
| Schema validate | 155–264 ns   | —         |
| Plan analysis   | ~46 ns       | —         |

## Development

```bash
make test        # run tests
make bench       # run benchmarks
make fmt         # format code
make all         # deps, fmt, lint, test, build
```

See [EXAMPLES.md](EXAMPLES.md) for usage examples and [PERFORMANCE.md](PERFORMANCE.md) for benchmark details.

## License

MIT — see [LICENSE](LICENSE).
