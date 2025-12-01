# Documentation Structure

This document explains the organization of SQL Parser Go documentation.

## 📚 Documentation Files

### Core Documentation

1. **[README.md](README.md)** (310 lines) - **START HERE**
   - Project overview and key features
   - Quick start guide
   - Supported SQL features summary
   - Performance highlights
   - Development commands
   - Roadmap and contributing guidelines
   
2. **[EXAMPLES.md](EXAMPLES.md)** - **Comprehensive Usage Examples**
   - Complete usage examples for all features
   - Multi-dialect examples (MySQL, PostgreSQL, SQL Server, SQLite, Oracle)
   - DML statements (INSERT, UPDATE, DELETE)
   - DDL statements (CREATE, DROP, ALTER, INDEX)
   - Advanced features (CTEs, Window Functions, Set Operations)
   - Subqueries and complex queries
   - Transaction control
   - Schema-aware parsing
   - Execution plan analysis
   - Command-line reference

3. **[PERFORMANCE.md](PERFORMANCE.md)** - **Detailed Performance Benchmarks**
   - Multi-dialect performance comparison
   - Lexing, parsing, and analysis benchmarks
   - Advanced features performance
   - DDL/DML operations performance
   - Schema validation performance
   - Transaction performance
   - Memory usage analysis
   - Optimization strategies applied

4. **[DIALECT_SUPPORT.md](DIALECT_SUPPORT.md)** - **Dialect-Specific Details**
   - MySQL syntax and features
   - PostgreSQL syntax and features
   - SQL Server syntax and features
   - SQLite syntax and features
   - Oracle syntax and features
   - Dialect comparison matrix

5. **[CLAUDE.md](CLAUDE.md)** - **Developer Guide**
   - Project architecture and structure
   - Component documentation
   - Development workflow
   - Common development tasks
   - Performance tuning guide
   - Testing strategies
   - Code conventions

## 🎯 Which Document to Read?

### I want to...

- **Get started quickly** → [README.md](README.md)
- **See usage examples** → [EXAMPLES.md](EXAMPLES.md)
- **Check performance** → [PERFORMANCE.md](PERFORMANCE.md)
- **Learn about dialects** → [DIALECT_SUPPORT.md](DIALECT_SUPPORT.md)
- **Contribute code** → [CLAUDE.md](CLAUDE.md)

## 📊 Documentation Reorganization (December 2024)

### Before
- **README.md**: 974 lines (too long, hard to navigate)
- Performance benchmarks scattered
- Examples mixed with documentation

### After
- **README.md**: 310 lines (concise, focused)
- **EXAMPLES.md**: Dedicated examples file
- **PERFORMANCE.md**: Dedicated performance documentation
- **DIALECT_SUPPORT.md**: Dialect-specific information
- **CLAUDE.md**: Developer guide

### Benefits
✅ Easier to find information  
✅ Better organization  
✅ More maintainable  
✅ Clearer separation of concerns  
✅ Faster onboarding for new users  

## 🔗 Quick Links

- [Main README](README.md)
- [Usage Examples](EXAMPLES.md)
- [Performance Benchmarks](PERFORMANCE.md)
- [Dialect Support](DIALECT_SUPPORT.md)
- [Developer Guide](CLAUDE.md)
- [Example Queries](examples/queries/)
- [Example Schemas](examples/schemas/)

---

**Need help?** Start with [README.md](README.md) and follow the links to detailed documentation!
