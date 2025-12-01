package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func TestCreateView(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		dialect     string
		expectError bool
	}{
		{
			name:        "Simple CREATE VIEW",
			sql:         "CREATE VIEW active_users AS SELECT * FROM users WHERE active = 1",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "CREATE VIEW with schema",
			sql:         "CREATE VIEW myschema.user_stats AS SELECT user_id, COUNT(*) as total FROM orders GROUP BY user_id",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "CREATE OR REPLACE VIEW",
			sql:         "CREATE OR REPLACE VIEW high_value_orders AS SELECT * FROM orders WHERE total > 1000",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "CREATE VIEW IF NOT EXISTS",
			sql:         "CREATE VIEW IF NOT EXISTS recent_orders AS SELECT * FROM orders WHERE created_at > '2024-01-01'",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "CREATE VIEW with column list",
			sql:         "CREATE VIEW user_summary (user_id, order_count, total_spent) AS SELECT user_id, COUNT(*), SUM(total) FROM orders GROUP BY user_id",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "CREATE VIEW with WITH CHECK OPTION",
			sql:         "CREATE VIEW premium_users AS SELECT * FROM users WHERE subscription = 'premium' WITH CHECK OPTION",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "CREATE MATERIALIZED VIEW (PostgreSQL)",
			sql:         "CREATE MATERIALIZED VIEW sales_summary AS SELECT product_id, SUM(amount) as total FROM sales GROUP BY product_id",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "CREATE OR REPLACE MATERIALIZED VIEW",
			sql:         "CREATE OR REPLACE MATERIALIZED VIEW monthly_stats AS SELECT DATE_TRUNC('month', created_at) as month, COUNT(*) FROM orders GROUP BY month",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "CREATE VIEW with complex SELECT",
			sql:         "CREATE VIEW customer_orders AS SELECT u.name, o.order_id, o.total FROM users u JOIN orders o ON u.id = o.user_id WHERE o.status = 'completed'",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "CREATE VIEW with subquery",
			sql:         "CREATE VIEW above_average AS SELECT * FROM users WHERE salary > (SELECT AVG(salary) FROM users)",
			dialect:     "postgresql",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			viewStmt, ok := stmt.(*parser.CreateViewStatement)
			if !ok {
				t.Errorf("Expected *CreateViewStatement, got %T", stmt)
				return
			}

			if viewStmt.SelectStmt == nil {
				t.Errorf("View SELECT statement is nil")
				return
			}

			t.Logf("✅ Successfully parsed CREATE VIEW: %s", viewStmt.ViewName.Name)
		})
	}
}

func TestDropView(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		dialect     string
		expectError bool
	}{
		{
			name:        "Simple DROP VIEW",
			sql:         "DROP VIEW active_users",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "DROP VIEW IF EXISTS",
			sql:         "DROP VIEW IF EXISTS user_stats",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "DROP VIEW with schema",
			sql:         "DROP VIEW myschema.sales_summary",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "DROP VIEW IF EXISTS CASCADE",
			sql:         "DROP VIEW IF EXISTS old_view CASCADE",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "DROP MATERIALIZED VIEW",
			sql:         "DROP MATERIALIZED VIEW sales_summary",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "DROP MATERIALIZED VIEW IF EXISTS",
			sql:         "DROP MATERIALIZED VIEW IF EXISTS monthly_stats",
			dialect:     "postgresql",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			dropStmt, ok := stmt.(*parser.DropStatement)
			if !ok {
				t.Errorf("Expected *DropStatement, got %T", stmt)
				return
			}

			if dropStmt.ObjectType != "VIEW" && dropStmt.ObjectType != "MATERIALIZED VIEW" {
				t.Errorf("Expected ObjectType VIEW or MATERIALIZED VIEW, got %s", dropStmt.ObjectType)
				return
			}

			t.Logf("✅ Successfully parsed DROP VIEW: %s", dropStmt.ObjectName)
		})
	}
}

func TestCreateViewDialects(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name:    "MySQL CREATE VIEW",
			sql:     "CREATE VIEW `user_orders` AS SELECT `id`, `name` FROM `users`",
			dialect: "mysql",
		},
		{
			name:    "PostgreSQL CREATE VIEW",
			sql:     "CREATE VIEW \"user_orders\" AS SELECT \"id\", \"name\" FROM \"users\"",
			dialect: "postgresql",
		},
		{
			name:    "SQL Server CREATE VIEW",
			sql:     "CREATE VIEW [user_orders] AS SELECT [id], [name] FROM [users]",
			dialect: "sqlserver",
		},
		{
			name:    "SQLite CREATE VIEW",
			sql:     "CREATE VIEW user_orders AS SELECT id, name FROM users",
			dialect: "sqlite",
		},
		{
			name:    "Oracle CREATE VIEW",
			sql:     "CREATE VIEW user_orders AS SELECT id, name FROM users",
			dialect: "oracle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.dialect, err)
				return
			}

			viewStmt, ok := stmt.(*parser.CreateViewStatement)
			if !ok {
				t.Errorf("Expected *CreateViewStatement for %s, got %T", tt.dialect, stmt)
				return
			}

			t.Logf("✅ %s: Successfully parsed CREATE VIEW", tt.dialect)
			_ = viewStmt
		})
	}
}
