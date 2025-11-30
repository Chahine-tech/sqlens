package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func TestExplainSimpleSelect(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		analyze bool
	}{
		{
			name:    "Simple EXPLAIN SELECT",
			sql:     "EXPLAIN SELECT * FROM users",
			dialect: "mysql",
			analyze: false,
		},
		{
			name:    "EXPLAIN ANALYZE SELECT",
			sql:     "EXPLAIN ANALYZE SELECT id, name FROM users WHERE id = 1",
			dialect: "postgresql",
			analyze: true,
		},
		{
			name:    "EXPLAIN with FORMAT",
			sql:     "EXPLAIN FORMAT=JSON SELECT * FROM orders",
			dialect: "mysql",
			analyze: false,
		},
		{
			name:    "PostgreSQL EXPLAIN with options",
			sql:     "EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM users",
			dialect: "postgresql",
			analyze: true,
		},
		{
			name:    "SQLite EXPLAIN QUERY PLAN",
			sql:     "EXPLAIN QUERY PLAN SELECT * FROM users WHERE email = 'test@example.com'",
			dialect: "sqlite",
			analyze: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dialect.GetDialect(tt.dialect)
			p := parser.NewWithDialect(context.Background(), tt.sql, d)
			stmt, err := p.ParseStatement()

			if err != nil {
				t.Fatalf("Failed to parse EXPLAIN statement: %v", err)
			}

			explainStmt, ok := stmt.(*parser.ExplainStatement)
			if !ok {
				t.Fatalf("Expected ExplainStatement, got %T", stmt)
			}

			if explainStmt.Analyze != tt.analyze {
				t.Errorf("Expected Analyze=%v, got %v", tt.analyze, explainStmt.Analyze)
			}

			if explainStmt.Statement == nil {
				t.Fatal("Inner statement is nil")
			}

			t.Logf("✅ Successfully parsed EXPLAIN statement: %s", explainStmt.String())
		})
	}
}

func TestExplainDMLStatements(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		stmtType string
	}{
		{
			name:     "EXPLAIN INSERT",
			sql:      "EXPLAIN INSERT INTO users (name, email) VALUES ('John', 'john@example.com')",
			stmtType: "InsertStatement",
		},
		{
			name:     "EXPLAIN UPDATE",
			sql:      "EXPLAIN UPDATE users SET name = 'Jane' WHERE id = 1",
			stmtType: "UpdateStatement",
		},
		{
			name:     "EXPLAIN DELETE",
			sql:      "EXPLAIN DELETE FROM users WHERE id > 100",
			stmtType: "DeleteStatement",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dialect.GetDialect("mysql")
			p := parser.NewWithDialect(context.Background(), tt.sql, d)
			stmt, err := p.ParseStatement()

			if err != nil {
				t.Fatalf("Failed to parse EXPLAIN statement: %v", err)
			}

			explainStmt, ok := stmt.(*parser.ExplainStatement)
			if !ok {
				t.Fatalf("Expected ExplainStatement, got %T", stmt)
			}

			if explainStmt.Statement.Type() != tt.stmtType {
				t.Errorf("Expected inner statement type %s, got %s", tt.stmtType, explainStmt.Statement.Type())
			}

			t.Logf("✅ Successfully parsed EXPLAIN %s", tt.stmtType)
		})
	}
}

func TestExplainExtended(t *testing.T) {
	sql := "EXPLAIN EXTENDED SELECT * FROM users"
	d := dialect.GetDialect("mysql")
	p := parser.NewWithDialect(context.Background(), sql, d)
	stmt, err := p.ParseStatement()

	if err != nil {
		t.Fatalf("Failed to parse EXPLAIN EXTENDED: %v", err)
	}

	explainStmt, ok := stmt.(*parser.ExplainStatement)
	if !ok {
		t.Fatalf("Expected ExplainStatement, got %T", stmt)
	}

	if extended, ok := explainStmt.Options["extended"]; !ok || extended != "true" {
		t.Errorf("Expected extended option to be true")
	}

	t.Log("✅ Successfully parsed EXPLAIN EXTENDED")
}

func TestExplainPostgreSQLOptions(t *testing.T) {
	sql := "EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT * FROM users"
	d := dialect.GetDialect("postgresql")
	p := parser.NewWithDialect(context.Background(), sql, d)
	stmt, err := p.ParseStatement()

	if err != nil {
		t.Fatalf("Failed to parse EXPLAIN with PostgreSQL options: %v", err)
	}

	explainStmt, ok := stmt.(*parser.ExplainStatement)
	if !ok {
		t.Fatalf("Expected ExplainStatement, got %T", stmt)
	}

	if !explainStmt.Analyze {
		t.Error("Expected Analyze to be true")
	}

	if explainStmt.Format != "JSON" {
		t.Errorf("Expected FORMAT=JSON, got %s", explainStmt.Format)
	}

	if buffers, ok := explainStmt.Options["BUFFERS"]; !ok || buffers != "true" {
		t.Error("Expected BUFFERS option to be true")
	}

	t.Log("✅ Successfully parsed EXPLAIN with PostgreSQL options")
}

func TestExplainComplexQuery(t *testing.T) {
	sql := `EXPLAIN ANALYZE
		SELECT u.id, u.name, COUNT(o.id) as order_count
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
		WHERE u.created_at > '2023-01-01'
		GROUP BY u.id, u.name
		HAVING COUNT(o.id) > 5
		ORDER BY order_count DESC
		LIMIT 10`

	d := dialect.GetDialect("postgresql")
	p := parser.NewWithDialect(context.Background(), sql, d)
	stmt, err := p.ParseStatement()

	if err != nil {
		t.Fatalf("Failed to parse complex EXPLAIN query: %v", err)
	}

	explainStmt, ok := stmt.(*parser.ExplainStatement)
	if !ok {
		t.Fatalf("Expected ExplainStatement, got %T", stmt)
	}

	if !explainStmt.Analyze {
		t.Error("Expected Analyze to be true")
	}

	selectStmt, ok := explainStmt.Statement.(*parser.SelectStatement)
	if !ok {
		t.Fatalf("Expected inner SelectStatement, got %T", explainStmt.Statement)
	}

	if len(selectStmt.Joins) != 1 {
		t.Errorf("Expected 1 join, got %d", len(selectStmt.Joins))
	}

	if len(selectStmt.GroupBy) != 2 {
		t.Errorf("Expected 2 GROUP BY columns, got %d", len(selectStmt.GroupBy))
	}

	t.Log("✅ Successfully parsed complex EXPLAIN ANALYZE query")
}
