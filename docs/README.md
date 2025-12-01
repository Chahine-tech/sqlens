# SQL Parser Go - Documentation

Welcome to the SQL Parser Go documentation! This directory contains comprehensive guides for using and understanding SQL Parser Go.

## 📚 Documentation Files

### [EXAMPLES.md](EXAMPLES.md)
**Comprehensive usage examples for all features**

Learn how to use SQL Parser Go with practical examples covering:
- Multi-dialect support (MySQL, PostgreSQL, SQL Server, SQLite, Oracle)
- DML statements (INSERT, UPDATE, DELETE)
- DDL statements (CREATE, DROP, ALTER, INDEX)
- Advanced features (CTEs, Window Functions, Set Operations)
- Subqueries and complex queries
- Transaction control
- Schema-aware parsing
- Execution plan analysis
- Command-line reference

**Start here if:** You want to see practical examples of how to use specific features.

---

### [PERFORMANCE.md](PERFORMANCE.md)
**Detailed performance benchmarks and optimization strategies**

Explore the performance characteristics of SQL Parser Go:
- Multi-dialect performance comparison (lexing, parsing, analysis)
- Advanced features benchmarks (subqueries, CTEs, window functions)
- DDL/DML operations performance
- Schema validation performance (zero-allocation!)
- Transaction performance (sub-microsecond!)
- Memory usage analysis
- Optimization strategies applied (object pooling, caching, zero-allocation paths)

**Start here if:** You want to understand performance characteristics or optimize your usage.

---

## 🎯 Quick Navigation

### I want to...

- **Get started** → [Main README](../README.md)
- **See usage examples** → [EXAMPLES.md](EXAMPLES.md)
- **Check performance** → [PERFORMANCE.md](PERFORMANCE.md)
- **Learn about dialects** → [DIALECT_SUPPORT.md](../DIALECT_SUPPORT.md)
- **Contribute code** → [CLAUDE.md](../CLAUDE.md)
- **See example files** → [examples/queries/](../examples/queries/)

---

## 📖 Additional Resources

- **[Main README](../README.md)** - Project overview, installation, quick start
- **[DIALECT_SUPPORT.md](../DIALECT_SUPPORT.md)** - Dialect-specific syntax and features
- **[CLAUDE.md](../CLAUDE.md)** - Developer guide for working with Claude Code
- **[examples/queries/](../examples/queries/)** - Example SQL files for all features
- **[examples/schemas/](../examples/schemas/)** - Example schema definitions

---

## 🚀 Quick Start Examples

### Basic Usage

```bash
# Analyze query from file
./bin/sqlparser -query examples/queries/complex_query.sql -output table

# Analyze query from string
./bin/sqlparser -sql "SELECT * FROM users WHERE id > 100" -dialect mysql

# Get optimization suggestions
./bin/sqlparser -sql "SELECT * FROM users" -dialect postgresql -output table
```

### Multi-Dialect

```bash
# MySQL with backticks
./bin/sqlparser -sql "SELECT \`user_id\` FROM \`users\`" -dialect mysql

# PostgreSQL with double quotes
./bin/sqlparser -sql "SELECT \"user_id\" FROM \"users\"" -dialect postgresql

# SQL Server with brackets
./bin/sqlparser -sql "SELECT [user_id] FROM [users]" -dialect sqlserver
```

For more examples, see [EXAMPLES.md](EXAMPLES.md).

---

## 📊 Performance Highlights

**Tested on Apple M2 Pro**

| Metric                   | Performance      | Notes                |
|--------------------------|------------------|----------------------|
| SQL Server Parsing       | 375 ns/op        | Sub-microsecond!     |
| Cached Analysis          | 26 ns/op         | 67x speedup          |
| Schema Validation        | 155-264 ns/op    | Zero-allocation      |
| Transaction COMMIT       | 149 ns/op        | Lightning-fast       |
| Scalar Subqueries        | 8-10 μs          | Sub-10 microseconds  |
| Window Functions         | 12-32 μs         | Production-ready     |

For complete benchmarks, see [PERFORMANCE.md](PERFORMANCE.md).

---

## 🤝 Contributing

See [CLAUDE.md](../CLAUDE.md) for:
- Project architecture and structure
- Development workflow
- Common development tasks
- Testing strategies
- Code conventions

---

## 📞 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/Chahine-tech/sql-parser-go/issues)
- 📖 **Documentation**: You're reading it!
- 💬 **Discussions**: [GitHub Discussions](https://github.com/Chahine-tech/sql-parser-go/discussions)

---

**Happy parsing!** 🚀
