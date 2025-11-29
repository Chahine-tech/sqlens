package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
	"github.com/Chahine-tech/sql-parser-go/pkg/schema"
)

// Benchmark schema loading from JSON
func BenchmarkSchemaLoadFromJSON(b *testing.B) {
	loader := schema.NewSchemaLoader()

	jsonData := []byte(`{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT", "primary_key": true},
					{"name": "name", "type": "VARCHAR", "length": 100},
					{"name": "email", "type": "VARCHAR", "length": 255}
				]
			}
		]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := loader.LoadFromJSON(jsonData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark schema validation - SELECT
func BenchmarkSchemaValidateSelect(b *testing.B) {
	loader := schema.NewSchemaLoader()

	jsonData := `{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "name", "type": "VARCHAR", "length": 100},
					{"name": "email", "type": "VARCHAR", "length": 255}
				]
			}
		]
	}`

	s, _ := loader.LoadFromJSON([]byte(jsonData))
	validator := schema.NewValidator(s)

	sql := `SELECT id, name FROM users WHERE email = 'test@example.com'`
	p := parser.NewWithDialect(context.Background(), sql, dialect.GetDialect("mysql"))
	stmt, _ := p.ParseStatement()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.ValidateStatement(stmt)
	}
}

// Benchmark schema validation - INSERT
func BenchmarkSchemaValidateInsert(b *testing.B) {
	loader := schema.NewSchemaLoader()

	jsonData := `{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "name", "type": "VARCHAR", "length": 100},
					{"name": "email", "type": "VARCHAR", "length": 255}
				]
			}
		]
	}`

	s, _ := loader.LoadFromJSON([]byte(jsonData))
	validator := schema.NewValidator(s)

	sql := `INSERT INTO users (id, name, email) VALUES (1, 'John', 'john@example.com')`
	p := parser.NewWithDialect(context.Background(), sql, dialect.GetDialect("mysql"))
	stmt, _ := p.ParseStatement()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.ValidateStatement(stmt)
	}
}

// Benchmark schema validation - UPDATE
func BenchmarkSchemaValidateUpdate(b *testing.B) {
	loader := schema.NewSchemaLoader()

	jsonData := `{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "name", "type": "VARCHAR", "length": 100},
					{"name": "email", "type": "VARCHAR", "length": 255}
				]
			}
		]
	}`

	s, _ := loader.LoadFromJSON([]byte(jsonData))
	validator := schema.NewValidator(s)

	sql := `UPDATE users SET name = 'Jane' WHERE id = 1`
	p := parser.NewWithDialect(context.Background(), sql, dialect.GetDialect("mysql"))
	stmt, _ := p.ParseStatement()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.ValidateStatement(stmt)
	}
}

// Benchmark type checking
func BenchmarkSchemaTypeChecking(b *testing.B) {
	loader := schema.NewSchemaLoader()

	jsonData := `{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "name", "type": "VARCHAR", "length": 100},
					{"name": "age", "type": "INT"}
				]
			}
		]
	}`

	s, _ := loader.LoadFromJSON([]byte(jsonData))
	typeChecker := schema.NewTypeChecker(s)

	sql := `SELECT * FROM users WHERE age > 18`
	p := parser.NewWithDialect(context.Background(), sql, dialect.GetDialect("mysql"))
	stmt, _ := p.ParseStatement()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = typeChecker.CheckStatement(stmt)
	}
}

// Benchmark complex schema validation
func BenchmarkSchemaValidateComplex(b *testing.B) {
	loader := schema.NewSchemaLoader()

	jsonData := `{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "name", "type": "VARCHAR", "length": 100}
				]
			},
			{
				"name": "orders",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "user_id", "type": "INT"},
					{"name": "total", "type": "DECIMAL", "precision": 10, "scale": 2}
				]
			}
		]
	}`

	s, _ := loader.LoadFromJSON([]byte(jsonData))
	validator := schema.NewValidator(s)

	sql := `SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id WHERE o.total > 100`
	p := parser.NewWithDialect(context.Background(), sql, dialect.GetDialect("mysql"))
	stmt, _ := p.ParseStatement()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.ValidateStatement(stmt)
	}
}
