package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// TestSimpleCaseExpression tests simple CASE expressions (CASE value WHEN...)
func TestSimpleCaseExpression(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name:    "Simple CASE with one WHEN",
			sql:     `SELECT CASE status WHEN 'active' THEN 1 END FROM users`,
			dialect: "mysql",
		},
		{
			name:    "Simple CASE with multiple WHEN",
			sql:     `SELECT CASE status WHEN 'active' THEN 1 WHEN 'inactive' THEN 0 END FROM users`,
			dialect: "mysql",
		},
		{
			name:    "Simple CASE with ELSE",
			sql:     `SELECT CASE status WHEN 'active' THEN 1 WHEN 'inactive' THEN 0 ELSE 2 END FROM users`,
			dialect: "mysql",
		},
		{
			name:    "Simple CASE with string results",
			sql:     `SELECT CASE type WHEN 'admin' THEN 'Administrator' WHEN 'user' THEN 'Regular User' ELSE 'Guest' END FROM accounts`,
			dialect: "postgresql",
		},
		{
			name:    "Simple CASE in WHERE clause",
			sql:     `SELECT * FROM users WHERE CASE status WHEN 'active' THEN 1 ELSE 0 END = 1`,
			dialect: "mysql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse simple CASE expression: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement, got nil")
			}

			t.Logf("✅ Successfully parsed simple CASE: %s", stmt.String())
		})
	}
}

// TestSearchedCaseExpression tests searched CASE expressions (CASE WHEN condition...)
func TestSearchedCaseExpression(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name:    "Searched CASE with one WHEN",
			sql:     `SELECT CASE WHEN age > 18 THEN 'adult' END FROM users`,
			dialect: "postgresql",
		},
		{
			name:    "Searched CASE with multiple WHEN",
			sql:     `SELECT CASE WHEN age < 13 THEN 'child' WHEN age < 18 THEN 'teen' ELSE 'adult' END FROM users`,
			dialect: "mysql",
		},
		{
			name:    "Searched CASE with complex conditions",
			sql:     `SELECT CASE WHEN status = 'active' THEN 1 WHEN status = 'pending' THEN 2 ELSE 0 END FROM orders`,
			dialect: "sqlserver",
		},
		{
			name:    "Searched CASE with AND condition",
			sql:     `SELECT CASE WHEN age > 18 THEN 'adult' ELSE 'minor' END FROM people`,
			dialect: "mysql",
		},
		{
			name:    "Searched CASE in ORDER BY",
			sql:     `SELECT name FROM users ORDER BY CASE WHEN priority = 'high' THEN 1 WHEN priority = 'medium' THEN 2 ELSE 3 END`,
			dialect: "postgresql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse searched CASE expression: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement, got nil")
			}

			t.Logf("✅ Successfully parsed searched CASE: %s", stmt.String())
		})
	}
}

// TestNestedCaseExpressions tests nested CASE expressions
func TestNestedCaseExpressions(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "CASE inside CASE",
			sql: `SELECT CASE
					WHEN type = 'premium' THEN CASE status WHEN 'active' THEN 100 ELSE 50 END
					ELSE 10
				END FROM users`,
			dialect: "mysql",
		},
		{
			name: "Multiple CASE in SELECT",
			sql: `SELECT
					CASE status WHEN 'active' THEN 1 ELSE 0 END,
					CASE type WHEN 'admin' THEN 'A' ELSE 'U' END
				FROM users`,
			dialect: "postgresql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse nested CASE expression: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement, got nil")
			}

			t.Logf("✅ Successfully parsed nested CASE: %s", stmt.String())
		})
	}
}

// TestCaseInDifferentClauses tests CASE in various SQL clauses
func TestCaseInDifferentClauses(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name:    "CASE in SELECT",
			sql:     `SELECT CASE status WHEN 'active' THEN 'Active User' ELSE 'Inactive' END AS user_status FROM users`,
			dialect: "mysql",
		},
		{
			name:    "CASE in WHERE",
			sql:     `SELECT * FROM orders WHERE CASE type WHEN 'express' THEN 1 ELSE 0 END = 1`,
			dialect: "postgresql",
		},
		{
			name:    "CASE in GROUP BY",
			sql:     `SELECT CASE status WHEN 'active' THEN 'Active' ELSE 'Other' END, COUNT(*) FROM users GROUP BY CASE status WHEN 'active' THEN 'Active' ELSE 'Other' END`,
			dialect: "mysql",
		},
		{
			name:    "CASE in ORDER BY",
			sql:     `SELECT name FROM users ORDER BY CASE status WHEN 'premium' THEN 1 WHEN 'active' THEN 2 ELSE 3 END`,
			dialect: "sqlserver",
		},
		{
			name:    "CASE in HAVING",
			sql:     `SELECT status, COUNT(*) FROM users GROUP BY status HAVING COUNT(*) > CASE status WHEN 'active' THEN 10 ELSE 5 END`,
			dialect: "mysql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse CASE in %s: %v", tt.name, err)
			}

			if stmt == nil {
				t.Fatal("Expected statement, got nil")
			}

			t.Logf("✅ Successfully parsed CASE in clause: %s", stmt.String())
		})
	}
}

// TestCaseWithFunctions tests CASE with function calls
func TestCaseWithFunctions(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name:    "CASE with function in WHEN",
			sql:     `SELECT CASE WHEN LENGTH(name) > 10 THEN 'long' ELSE 'short' END FROM users`,
			dialect: "mysql",
		},
		{
			name:    "CASE with function in THEN",
			sql:     `SELECT CASE status WHEN 'active' THEN UPPER(name) ELSE LOWER(name) END FROM users`,
			dialect: "postgresql",
		},
		{
			name:    "Function wrapping CASE",
			sql:     `SELECT UPPER(CASE status WHEN 'active' THEN 'yes' ELSE 'no' END) FROM users`,
			dialect: "mysql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse CASE with functions: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement, got nil")
			}

			t.Logf("✅ Successfully parsed CASE with functions: %s", stmt.String())
		})
	}
}

// TestComplexCaseScenarios tests real-world complex CASE scenarios
func TestComplexCaseScenarios(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "User tier calculation",
			sql: `SELECT
					user_id,
					CASE
						WHEN total_purchases > 10000 THEN 'platinum'
						WHEN total_purchases > 5000 THEN 'gold'
						WHEN total_purchases > 1000 THEN 'silver'
						ELSE 'bronze'
					END AS tier
				FROM user_stats`,
			dialect: "mysql",
		},
		{
			name: "Grade calculation",
			sql: `SELECT
					student_name,
					CASE
						WHEN score >= 90 THEN 'A'
						WHEN score >= 80 THEN 'B'
						WHEN score >= 70 THEN 'C'
						WHEN score >= 60 THEN 'D'
						ELSE 'F'
					END AS grade
				FROM student_scores`,
			dialect: "postgresql",
		},
		{
			name: "Status mapping",
			sql: `SELECT
					order_id,
					CASE order_status
						WHEN 'P' THEN 'Pending'
						WHEN 'S' THEN 'Shipped'
						WHEN 'D' THEN 'Delivered'
						WHEN 'C' THEN 'Cancelled'
						ELSE 'Unknown'
					END AS status_description
				FROM orders`,
			dialect: "sqlserver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse complex CASE scenario: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement, got nil")
			}

			t.Logf("✅ Successfully parsed complex CASE: %s", stmt.String())
		})
	}
}
