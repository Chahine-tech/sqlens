# Claude Code Guide - SQL Parser Go

This guide helps you work with the SQL Parser Go project using Claude Code.

## 📚 Project Overview

**SQL Parser Go** is a high-performance, multi-dialect SQL query analysis tool that provides:
- Parsing and analysis for MySQL, PostgreSQL, SQL Server, SQLite, and Oracle
- Advanced optimization suggestions (dialect-specific)
- **Extended SQL features**: CTEs (WITH clause), Window Functions, Set Operations
- **Schema-aware parsing and validation** - Validate SQL against database schemas
- SQL Server log parsing (Profiler, Extended Events, Query Store)
- Sub-microsecond query parsing with intelligent caching
- Comprehensive CLI interface

**Tech Stack**: Go 1.25, minimal dependencies (yaml.v3 only)

## 🚀 Quick Start

### Initial Setup
```bash
# Install dependencies
make deps

# Build the project
make build

# Run tests
make test

# Run benchmarks
make bench
```

### Common Development Commands
```bash
# Analyze a query from file
./bin/sqlparser -query examples/queries/complex_query.sql -output table

# Analyze inline SQL with specific dialect
./bin/sqlparser -sql "SELECT * FROM users" -dialect mysql

# Parse SQL Server logs
./bin/sqlparser -log examples/logs/sample_profiler.log -output table -verbose

# Run performance benchmarks
make bench

# Run with verbose output
./bin/sqlparser -query file.sql -verbose
```

## 📁 Project Structure

```
sql-parser-go/
├── cmd/sqlparser/          # CLI entry point (main.go)
├── pkg/
│   ├── lexer/             # SQL tokenization (lexer.go, tokens.go)
│   ├── parser/            # SQL parsing & AST (parser.go, ast.go, procedure_parser.go, errors.go, pool.go)
│   ├── analyzer/          # Query analysis (analyzer.go, extractor.go, optimization*.go, concurrent.go)
│   ├── dialect/           # Dialect support (mysql.go, postgresql.go, sqlserver.go, sqlite.go, oracle.go)
│   ├── schema/            # Schema definitions and validation (schema.go, loader.go, validator.go, type_checker.go)
│   ├── plan/              # Execution plan analysis (plan.go, analyzer.go)
│   └── logger/            # Log parsing (parser.go, formats.go)
├── internal/
│   ├── config/            # Configuration (config.go)
│   └── performance/       # Performance monitoring (monitor.go)
├── tests/                 # All test files (*_test.go)
└── examples/              # Example queries, logs, and schemas
```

## 🔧 Key Components

### 1. Lexer (`pkg/lexer/`)
- **Purpose**: Tokenizes SQL text into tokens
- **Performance**: ~1826 ns/op
- **Files**:
  - [lexer.go](pkg/lexer/lexer.go) - Main lexer logic
  - [tokens.go](pkg/lexer/tokens.go) - Token definitions

### 2. Parser (`pkg/parser/`)
- **Purpose**: Builds Abstract Syntax Tree (AST) from tokens
- **Performance**: ~1141 ns/op (sub-microsecond!)
- **Files**:
  - [parser.go](pkg/parser/parser.go) - Main parser
  - [ast.go](pkg/parser/ast.go) - AST node definitions
  - [pool.go](pkg/parser/pool.go) - Object pooling for performance
  - [errors.go](pkg/parser/errors.go) - Error handling

### 3. Analyzer (`pkg/analyzer/`)
- **Purpose**: Extracts metadata and provides optimization suggestions
- **Performance**: 1786 ns/op (cold) / 26.42 ns/op (cached) - 67x speedup with cache!
- **Files**:
  - [analyzer.go](pkg/analyzer/analyzer.go) - Main analyzer
  - [extractor.go](pkg/analyzer/extractor.go) - Table/column extraction
  - [optimization.go](pkg/analyzer/optimization.go) - Core optimization logic
  - [optimization_dialect.go](pkg/analyzer/optimization_dialect.go) - Dialect-specific optimizations
  - [optimization_rules.go](pkg/analyzer/optimization_rules.go) - Optimization rules
  - [concurrent.go](pkg/analyzer/concurrent.go) - Multi-core analysis

### 3.5. Advanced Features (`pkg/parser/advanced_features.go`)
- **Purpose**: Parse modern SQL features (CTEs, Window Functions, Set Operations)
- **Features**:
  - WITH clause (CTEs) - simple and multiple
  - Window functions with OVER, PARTITION BY, ORDER BY, frames
  - Set operations (UNION, INTERSECT, EXCEPT)
  - CASE expressions (AST only - parsing TBD)
- **Files**:
  - [advanced_features.go](pkg/parser/advanced_features.go) - Advanced SQL parsing (448 lines)

### 4. Dialect Support (`pkg/dialect/`)
- **Purpose**: Handle dialect-specific syntax and features
- **Supported**: MySQL, PostgreSQL, SQL Server, SQLite, Oracle
- **Files**: One file per dialect ([mysql.go](pkg/dialect/mysql.go), [postgresql.go](pkg/dialect/postgresql.go), etc.)

### 5. Schema (`pkg/schema/`)
- **Purpose**: Define database schemas and validate SQL against them
- **Performance**: 7.2μs schema loading, 155-264ns validation (zero-allocation!)
- **Files**:
  - [schema.go](pkg/schema/schema.go) - Schema, Table, Column, DataType definitions
  - [loader.go](pkg/schema/loader.go) - Load schemas from JSON/YAML
  - [validator.go](pkg/schema/validator.go) - Validate SQL statements against schema
  - [type_checker.go](pkg/schema/type_checker.go) - Type compatibility checking

### 6. Execution Plan (`pkg/plan/`)
- **Purpose**: Analyze query execution plans and provide optimization suggestions
- **Performance**: 46ns plan analysis, 117ns bottleneck detection
- **Files**:
  - [plan.go](pkg/plan/plan.go) - ExecutionPlan, PlanNode, cost structures
  - [analyzer.go](pkg/plan/analyzer.go) - Plan analysis and JSON/XML parsing

### 7. Logger (`pkg/logger/`)
- **Purpose**: Parse SQL Server log files
- **Formats**: Profiler, Extended Events, Query Store
- **Files**:
  - [parser.go](pkg/logger/parser.go) - Log parsing logic
  - [formats.go](pkg/logger/formats.go) - Format definitions

## 🎯 Common Tasks for Claude

### Task 1: Add Support for a New SQL Statement Type
**Example**: Adding support for `CREATE TABLE`

1. **Update Token Types** in [pkg/lexer/tokens.go](pkg/lexer/tokens.go)
   - Add new keywords: `CREATE`, `TABLE`, `PRIMARY`, `FOREIGN`, etc.

2. **Update Lexer** in [pkg/lexer/lexer.go](pkg/lexer/lexer.go)
   - Add keyword mappings in `keywords` map

3. **Create AST Node** in [pkg/parser/ast.go](pkg/parser/ast.go)
   - Add new struct for `CreateTableStatement`

4. **Implement Parser** in [pkg/parser/parser.go](pkg/parser/parser.go)
   - Add `parseCreateTableStatement()` method
   - Update `ParseStatement()` to handle `CREATE` keyword

5. **Add Tests** in [tests/parser_test.go](tests/parser_test.go)
   - Test basic CREATE TABLE
   - Test with constraints, foreign keys, etc.

### Task 2: Add New Optimization Rule
**Example**: Detecting missing indexes

1. **Add Rule Logic** in [pkg/analyzer/optimization_rules.go](pkg/analyzer/optimization_rules.go)
   - Create `detectMissingIndexes()` function

2. **Integrate in Analyzer** in [pkg/analyzer/optimization.go](pkg/analyzer/optimization.go)
   - Call new rule in `SuggestOptimizations()`

3. **Add Dialect-Specific Logic** (if needed) in [pkg/analyzer/optimization_dialect.go](pkg/analyzer/optimization_dialect.go)

4. **Add Tests** in [tests/optimization_test.go](tests/optimization_test.go)

### Task 3: Fix a Parser Bug
**Example**: Parser fails on certain JOIN syntax

1. **Write Failing Test First** in [tests/parser_test.go](tests/parser_test.go)
2. **Debug with Verbose Mode**:
   ```bash
   ./bin/sqlparser -query problem.sql -verbose
   ```
3. **Fix in Parser** in [pkg/parser/parser.go](pkg/parser/parser.go)
4. **Verify Test Passes**: `make test`

### Task 4: Add New Dialect Support
**Example**: Adding Snowflake dialect

1. **Create Dialect File**: [pkg/dialect/snowflake.go](pkg/dialect/snowflake.go)
2. **Implement Interface** from [pkg/dialect/dialect.go](pkg/dialect/dialect.go)
3. **Register Dialect** in `GetDialect()` function
4. **Add Tests**: [tests/dialect_test.go](tests/dialect_test.go)
5. **Update Documentation**: [DIALECT_SUPPORT.md](DIALECT_SUPPORT.md)

### Task 5: Add Support for Advanced SQL Features
**Example**: Adding MERGE statement or materialized views

Since we now have advanced features support, here's the pattern:

1. **Add Tokens** in [pkg/lexer/tokens.go](pkg/lexer/tokens.go)
2. **Create AST Nodes** in [pkg/parser/ast.go](pkg/parser/ast.go)
3. **Implement Parsing** in [pkg/parser/advanced_features.go](pkg/parser/advanced_features.go)
4. **Add Tests** in [tests/advanced_features_test.go](tests/advanced_features_test.go)
5. **Add Examples** in [examples/queries/](examples/queries/)

**Currently Supported Advanced Features:**
- ✅ CTEs (WITH clause) - see [cte_examples.sql](examples/queries/cte_examples.sql)
- ✅ Window Functions - see [window_function_examples.sql](examples/queries/window_function_examples.sql)
- ✅ Set Operations - see [set_operations_examples.sql](examples/queries/set_operations_examples.sql)
- ✅ CASE expressions - see [case_expression_examples.sql](examples/queries/case_expression_examples.sql)
- ✅ **MERGE Statement** - see [merge_examples.sql](examples/queries/merge_examples.sql) 🆕
- ✅ **Advanced Cursor Operations** - FETCH NEXT/PRIOR/FIRST/LAST/ABSOLUTE/RELATIVE, DEALLOCATE 🆕
- ✅ **PostgreSQL Dollar-Quoted Strings** - $$...$$, $tag$...$tag$ for function bodies 🆕

### Task 6: Improve Performance
**Target**: Lexer, Parser, or Analyzer

1. **Run Benchmarks**:
   ```bash
   make bench
   ```

2. **Profile with CPU profiling**:
   ```bash
   make bench-cpu
   go tool pprof cpu.prof
   ```

3. **Common Optimizations**:
   - Use object pooling (see [pkg/parser/pool.go](pkg/parser/pool.go))
   - Pre-allocate slices with capacity
   - Reduce string allocations
   - Add caching where appropriate

4. **Verify Improvement**:
   ```bash
   make bench > before.txt
   # Make changes
   make bench > after.txt
   # Compare results
   ```

### Task 7: Add Schema-Aware Validation
**Example**: Validating SQL queries against a database schema

1. **Create Schema File** in [examples/schemas/](examples/schemas/)
   - Define tables, columns, data types in JSON/YAML format

2. **Load Schema** using [pkg/schema/loader.go](pkg/schema/loader.go)
   ```go
   loader := schema.NewSchemaLoader()
   s, err := loader.LoadFromFile("schema.json")
   ```

3. **Validate SQL** using [pkg/schema/validator.go](pkg/schema/validator.go)
   ```go
   validator := schema.NewValidator(s)
   errors := validator.ValidateStatement(stmt)
   ```

4. **Check Types** using [pkg/schema/type_checker.go](pkg/schema/type_checker.go)
   ```go
   typeChecker := schema.NewTypeChecker(s)
   errors := typeChecker.CheckStatement(stmt)
   ```

5. **Add Tests** in [tests/schema_test.go](tests/schema_test.go)
   - Test table/column existence validation
   - Test type compatibility checking
   - Test foreign key validation

### Task 8: Analyze Execution Plans
**Example**: Analyzing query execution plans for performance bottlenecks

1. **Parse EXPLAIN statement**
   ```go
   sql := "EXPLAIN ANALYZE SELECT * FROM users WHERE id = 1"
   p := parser.NewWithDialect(ctx, sql, dialect.GetDialect("postgresql"))
   stmt, err := p.ParseStatement()
   explainStmt := stmt.(*parser.ExplainStatement)
   ```

2. **Parse JSON execution plan** (from database output)
   ```go
   jsonPlan := []byte(`{"Plan": {...}}`) // From EXPLAIN FORMAT=JSON
   executionPlan, err := plan.ParseJSONPlan(jsonPlan, "postgresql")
   ```

3. **Analyze plan for bottlenecks**
   ```go
   analyzer := plan.NewPlanAnalyzer("postgresql")
   analysis := analyzer.AnalyzePlan(executionPlan)

   // Check performance score
   fmt.Printf("Performance Score: %.2f/100\n", analysis.PerformanceScore)

   // Review issues
   for _, issue := range analysis.Issues {
       fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Type, issue.Description)
   }

   // Get recommendations
   for _, rec := range analysis.Recommendations {
       fmt.Printf("[%s] %s\n", rec.Priority, rec.Description)
   }
   ```

4. **Find specific bottlenecks**
   ```go
   bottlenecks := executionPlan.FindBottlenecks()
   for _, b := range bottlenecks {
       fmt.Printf("Issue: %s (Impact: %.2f)\n", b.Issue, b.ImpactScore)
       fmt.Printf("Fix: %s\n", b.Recommendation)
   }
   ```

## 🧪 Testing Strategy

### Run All Tests
```bash
make test
```

### Run Specific Test
```bash
go test -v ./tests -run TestParserSimpleSelect
```

### Run Benchmarks
```bash
make bench
make bench-cpu    # With CPU profiling
make bench-mem    # With memory profiling
```

### Test Coverage
```bash
go test -cover ./...
```

## 🐛 Debugging Tips

### Enable Verbose Output
```bash
./bin/sqlparser -query file.sql -verbose
```

### Check Token Stream
Add debug print in [pkg/lexer/lexer.go](pkg/lexer/lexer.go:177) `NextToken()`:
```go
fmt.Printf("Token: %v\n", tok)
```

### Check AST Structure
Add debug print in [cmd/sqlparser/main.go](cmd/sqlparser/main.go) after parsing:
```go
fmt.Printf("AST: %+v\n", stmt)
```

### Profile Performance
```bash
go test -cpuprofile=cpu.prof -bench=BenchmarkParser ./tests
go tool pprof cpu.prof
# Then: top10, list functionName
```

## 📝 Code Conventions

### Naming
- **Files**: lowercase with underscores (`optimization_rules.go`)
- **Types**: PascalCase (`SelectStatement`)
- **Functions**: camelCase for private, PascalCase for public
- **Constants**: UPPER_CASE for token types

### Error Handling
- Always return errors, don't panic
- Use meaningful error messages
- Include position information when parsing

### Performance
- Pre-allocate slices when size is known
- Use object pooling for frequently allocated objects
- Benchmark before and after changes

### Testing
- Test file names: `*_test.go`
- Benchmark names: `Benchmark*`
- Use table-driven tests for multiple cases

## 🔍 Important Files to Know

### Core Logic
- [cmd/sqlparser/main.go](cmd/sqlparser/main.go) - CLI entry point
- [pkg/parser/parser.go](pkg/parser/parser.go) - Main parser logic (~800 lines)
- [pkg/analyzer/analyzer.go](pkg/analyzer/analyzer.go) - Analysis engine

### Performance Critical
- [pkg/parser/pool.go](pkg/parser/pool.go) - Object pooling (60% allocation reduction)
- [pkg/lexer/lexer.go](pkg/lexer/lexer.go) - Tokenization hot path
- [pkg/analyzer/concurrent.go](pkg/analyzer/concurrent.go) - Multi-core processing

### Configuration
- [config.yaml](config.yaml) - Default configuration
- [internal/config/config.go](internal/config/config.go) - Config loader

### Documentation
- [README.md](README.md) - Main documentation
- [DIALECT_SUPPORT.md](DIALECT_SUPPORT.md) - Dialect details
- [PERFORMANCE.md](PERFORMANCE.md) - Performance notes
- [todo.md](todo.md) - Development roadmap

## 🎨 Roadmap & Future Features

### ✅ Completed
- Multi-dialect support (5 dialects)
- Advanced optimization suggestions
- Performance benchmarking
- Dialect-specific identifier quoting
- **Extended SQL features** ✨
  - **CTEs (WITH clause)** - Simple, multiple, with column lists
  - **Window Functions** - OVER, PARTITION BY, ORDER BY, ROWS/RANGE frames
  - **Set Operations** - UNION, UNION ALL, INTERSECT, EXCEPT
  - **CASE Expressions** - Searched and simple CASE statements
- **DML Statement Support** ✅
  - **INSERT** - VALUES, multiple rows, INSERT...SELECT
  - **UPDATE** - Multiple columns, WHERE, ORDER BY/LIMIT (MySQL/SQLite)
  - **DELETE** - WHERE, ORDER BY/LIMIT (MySQL/SQLite)
- **Comprehensive Subquery Support** ✅
  - **Scalar Subqueries** - In WHERE, SELECT, INSERT VALUES, UPDATE SET
  - **EXISTS / NOT EXISTS** - Full support in all statement types
  - **IN / NOT IN with Subqueries** - Complete implementation
  - **Derived Tables** - Subqueries in FROM clause with JOIN support
  - **Nested & Correlated Subqueries** - Multiple levels of nesting
  - **40+ comprehensive tests** - All passing
- **DDL Support** ✅
  - **CREATE TABLE** - Columns, constraints, foreign keys, IF NOT EXISTS
  - **DROP** - TABLE/DATABASE/INDEX with IF EXISTS and CASCADE
  - **ALTER TABLE** - ADD/DROP/MODIFY/CHANGE columns and constraints
  - **CREATE INDEX** - Simple and unique indexes with IF NOT EXISTS
  - **Dialect-specific features** - AUTO_INCREMENT (MySQL), IDENTITY (SQL Server), AUTOINCREMENT (SQLite)
  - **Foreign key references** - ON DELETE/UPDATE actions (CASCADE, SET NULL, SET DEFAULT, NO ACTION)
  - **50+ comprehensive tests** - All passing
- **Transaction Support** ✅
  - **BEGIN/START TRANSACTION** - Start transactions (dialect-aware)
  - **COMMIT** - Commit transactions (with optional WORK keyword)
  - **ROLLBACK** - Roll back transactions (with optional WORK keyword)
  - **SAVEPOINT** - Create savepoints within transactions
  - **ROLLBACK TO SAVEPOINT** - Roll back to specific savepoints
  - **RELEASE SAVEPOINT** - Release savepoints (PostgreSQL/MySQL)
  - **16+ comprehensive tests** - All passing
  - **Ultra-fast performance** - Sub-microsecond COMMIT/ROLLBACK
- **Schema-Aware Parsing** ✅
  - **Schema Definition** - Tables, columns, data types, constraints, indexes, foreign keys
  - **Schema Loading** - JSON and YAML format support
  - **SQL Validation** - Validate SELECT, INSERT, UPDATE, DELETE against schema
  - **Table/Column Validation** - Check existence of tables and columns
  - **Type Checking** - Data type compatibility validation in expressions
  - **Foreign Key Support** - Validate foreign key references
  - **9+ comprehensive tests** - All passing (67 total tests)
  - **Zero-allocation validation** - 155-264ns per statement
  - **Fast schema loading** - 7.2μs from JSON
- **Query Execution Plan Analysis** ✅ 🆕
  - **EXPLAIN Statement Support** - Full multi-dialect EXPLAIN parsing
  - **EXPLAIN ANALYZE** - Support for actual execution statistics
  - **Execution Plan Structures** - 20+ node types (Scans, Joins, Aggregates, Sorts)
  - **Plan Analyzer** - Automatic bottleneck detection and optimization suggestions
  - **Performance Score** - 0-100 score calculation based on plan quality
  - **Dialect-Specific Parsing** - MySQL JSON, PostgreSQL JSON, SQL Server XML formats
  - **14+ comprehensive tests** - All passing (81 total tests)
  - **Ultra-fast analysis** - 46ns plan analysis, 117ns bottleneck detection
  - **8 benchmarks** - Sub-microsecond execution plan parsing
- **Stored Procedures and Functions** ✅ 🆕
  - **CREATE PROCEDURE** - Full support with parameters (IN/OUT/INOUT)
  - **CREATE FUNCTION** - With return types and DETERMINISTIC
  - **OR REPLACE** - PostgreSQL/Oracle style CREATE OR REPLACE
  - **Parameter Modes** - IN, OUT, INOUT with default values
  - **Data Types** - VARCHAR(n), DECIMAL(p,s), INT, etc.
  - **Variable Declarations** - DECLARE variables with types
  - **Cursor Support** - DECLARE CURSOR FOR SELECT
  - **Procedural Statements** - RETURN, assignments, cursor operations
  - **Dialect-Specific Options** - LANGUAGE, DETERMINISTIC, SECURITY DEFINER/INVOKER
  - **10+ comprehensive tests** - All passing (MySQL, PostgreSQL, SQL Server)
  - **8 benchmarks** - 10-54μs procedure parsing
  - **650+ lines** - Complete procedure parser implementation
- **View Support** ✅ 🆕
  - **CREATE VIEW** - Standard view creation with column lists
  - **CREATE OR REPLACE VIEW** - Idempotent view definitions
  - **CREATE VIEW IF NOT EXISTS** - Safe view creation
  - **CREATE MATERIALIZED VIEW** - PostgreSQL materialized views
  - **WITH CHECK OPTION** - Enforce view constraints
  - **DROP VIEW** - Simple and conditional (IF EXISTS) view dropping
  - **DROP MATERIALIZED VIEW** - PostgreSQL materialized view dropping
  - **Multi-dialect support** - MySQL (backticks), PostgreSQL (double quotes, MATERIALIZED), SQL Server (brackets)
  - **Schema-qualified views** - Support for schema.view_name syntax
  - **27 comprehensive tests** - All passing across all dialects
  - **231 lines of examples** - Complete examples in view_examples.sql
- **Trigger Support** ✅ 🆕
  - **CREATE TRIGGER** - Full trigger creation with timing, events, and body
  - **BEFORE/AFTER/INSTEAD OF** - All trigger timings supported
  - **Multiple events** - INSERT, UPDATE, DELETE with OR combinations
  - **FOR EACH ROW/STATEMENT** - Row-level and statement-level triggers
  - **WHEN conditions** - Optional trigger conditions (PostgreSQL)
  - **IF NOT EXISTS** - Safe trigger creation (MySQL)
  - **OR REPLACE** - Idempotent trigger definitions (PostgreSQL)
  - **DROP TRIGGER** - Simple and conditional (IF EXISTS) trigger dropping
  - **BEGIN...END blocks** - Trigger body parsing
  - **Multi-dialect support** - MySQL, PostgreSQL, SQL Server, SQLite, Oracle
  - **Schema-qualified tables** - Support for triggers on schema.table
  - **23 comprehensive tests** - All passing across all dialects
  - **450 lines of examples** - Complete real-world examples in trigger_examples.sql
- **Control Flow Statements** ✅ 🆕
  - **IF...THEN...ELSE...END IF** - Full conditional logic with ELSEIF/ELSIF support
  - **WHILE...DO...END WHILE** - Conditional loops with pre-condition check
  - **FOR...LOOP** - Range-based iteration (PostgreSQL style) with REVERSE and BY step
  - **LOOP...END LOOP** - Infinite loops with EXIT conditions
  - **REPEAT...UNTIL** - Post-condition loops (MySQL)
  - **EXIT / CONTINUE** - Loop control with optional WHEN conditions and labels
  - **ITERATE** - MySQL synonym for CONTINUE
  - **Nested control flow** - IF inside WHILE, FOR inside LOOP, etc.
  - **Multi-dialect support** - MySQL, PostgreSQL, SQL Server compatible syntax
  - **28+ comprehensive tests** - All passing across all control flow types
  - **450+ lines of examples** - Real-world scenarios in control_flow_examples.sql
- **Exception Handling** ✅ 🆕
  - **TRY...CATCH** - SQL Server exception handling with BEGIN TRY/END TRY/BEGIN CATCH/END CATCH
  - **EXCEPTION...WHEN** - PostgreSQL/Oracle exception blocks with WHEN clauses
  - **DECLARE HANDLER** - MySQL handler declarations (CONTINUE, EXIT, UNDO)
  - **RAISE** - PostgreSQL error raising (EXCEPTION, NOTICE, WARNING, INFO, LOG, DEBUG)
  - **THROW** - SQL Server error throwing with error number, message, and state
  - **SIGNAL** - MySQL error signaling with SQLSTATE and properties
  - **Exception conditions** - SQLEXCEPTION, SQLWARNING, NOT FOUND, SQLSTATE, OTHERS
  - **Nested exception handling** - TRY within TRY, multiple WHEN clauses
  - **Re-raise support** - THROW and RAISE without parameters
  - **Multi-dialect support** - SQL Server, PostgreSQL, MySQL, Oracle syntax
  - **23+ comprehensive tests** - All passing across all exception types
  - **400+ lines of examples** - Real-world error handling scenarios in exception_handling_examples.sql
- **MERGE Statement Support** ✅ 🆕
  - **MERGE INTO...USING...ON** - Full MERGE statement syntax
  - **WHEN MATCHED / NOT MATCHED** - All conditional clauses
  - **WHEN NOT MATCHED BY SOURCE** - SQL Server specific clause
  - **Multiple WHEN clauses** - Multiple conditions per statement
  - **UPDATE / INSERT / DELETE actions** - All DML actions in MERGE
  - **Subquery sources** - Support for subqueries as data sources
  - **Qualified column names** - table.column syntax in SET clauses
  - **18+ comprehensive tests** - All passing across SQL Server, PostgreSQL, Oracle
  - **350+ lines of examples** - Real-world ETL and synchronization scenarios in merge_examples.sql
- **Advanced Cursor Operations** ✅ 🆕
  - **FETCH directions** - NEXT, PRIOR, FIRST, LAST, ABSOLUTE n, RELATIVE n
  - **FETCH...INTO** - Fetch cursor values into variables
  - **DEALLOCATE** - Deallocate cursor or prepared statement
  - **FROM/IN keywords** - PostgreSQL-style FETCH FROM cursor
  - **12+ comprehensive tests** - All passing for MySQL, PostgreSQL, SQL Server
  - **Full cursor lifecycle** - DECLARE, OPEN, FETCH, CLOSE, DEALLOCATE
- **PostgreSQL Dollar-Quoted Strings** ✅ 🆕
  - **$$...$$ syntax** - Simple dollar quoting for strings
  - **$tag$...$tag$ syntax** - Tagged dollar quoting with custom delimiters
  - **Nested content** - Supports $$ inside tagged strings
  - **Function bodies** - Perfect for PostgreSQL function/procedure bodies
  - **Multi-line support** - Preserves newlines and formatting
  - **Dialect-specific** - Only enabled for PostgreSQL dialect
  - **13+ comprehensive tests** - All passing including nested and edge cases

### 🚧 In Progress / Planned
- [ ] Real-time log monitoring
- [ ] Integration with monitoring tools

### ❌ Not Planned
- Web interface (project stays CLI-focused)

## 💡 Tips for Working with Claude Code

### When Asking for Help
Be specific about:
1. Which component you're working on (lexer/parser/analyzer)
2. What SQL dialect you're targeting
3. Expected vs actual behavior
4. Include example SQL that fails/succeeds

### Example Requests
Good:
- "Add support for MySQL's `LIMIT` clause in the parser"
- "The PostgreSQL dialect doesn't recognize double-quoted identifiers in joins"
- "Optimize the analyzer's table extraction - it's too slow for queries with 10+ tables"

Less Good:
- "Fix the parser" (too vague)
- "Add more features" (what features?)

### Before Asking Claude to Code
1. Run tests to see current state: `make test`
2. Check if similar code exists elsewhere in the project
3. Review existing tests for patterns
4. Look at [todo.md](todo.md) for planned features

## 🚀 Performance Targets

Current performance (Apple M2 Pro):
- **Lexer**: ~1826 ns/op, ~260 MB/s
- **Parser**: ~1141 ns/op, **sub-microsecond!**
- **Analyzer**: 1786 ns/op (cold) / 26.42 ns/op (cached)
- **Memory**: Very efficient with object pooling

When optimizing, aim to maintain or improve these metrics.

## 📚 Useful Commands

```bash
# Development
make build              # Build binary
make test               # Run all tests
make bench              # Run benchmarks
make fmt                # Format code
make lint               # Lint code

# Examples
make dev-query          # Analyze complex_query.sql
make dev-simple         # Analyze simple_query.sql
make dev-log            # Parse sample log

# Performance
make bench-cpu          # CPU profiling
make bench-mem          # Memory profiling
make perf-compare       # Compare performance

# Release
make build-release      # Optimized build
make build-all          # Multi-platform build
```

## 🤝 Contributing Guidelines

1. Write tests first (TDD approach)
2. Run `make all` before committing (deps, fmt, lint, test, build)
3. Update documentation if adding new features
4. Keep performance in mind - benchmark if changing hot paths
5. Follow existing code style and conventions

## 📞 Getting Help

- Check [README.md](README.md) for usage examples
- Review [todo.md](todo.md) for known issues and planned features
- Look at tests in [tests/](tests/) for usage patterns
- Use `make help` to see all available commands

---

**Happy coding with Claude!** 🚀
