package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// Test INSERT Statement
func TestInsertStatements(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "Simple INSERT with VALUES",
			sql:     `INSERT INTO users (id, name, email) VALUES (1, 'John', 'john@test.com')`,
			wantErr: false,
		},
		{
			name:    "INSERT with multiple rows",
			sql:     `INSERT INTO users (id, name) VALUES (1, 'John'), (2, 'Jane'), (3, 'Bob')`,
			wantErr: false,
		},
		{
			name:    "INSERT without column list",
			sql:     `INSERT INTO users VALUES (1, 'John', 'john@test.com')`,
			wantErr: false,
		},
		{
			name:    "INSERT with SELECT",
			sql:     `INSERT INTO users (id, name) SELECT id, name FROM temp_users`,
			wantErr: false,
		},
		{
			name:    "INSERT with expressions",
			sql:     `INSERT INTO orders (user_id, total, created_at) VALUES (1, 100 + 50, NOW())`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("INSERT parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				if insertStmt, ok := stmt.(*parser.InsertStatement); ok {
					t.Logf("Successfully parsed INSERT into %s", insertStmt.Table.Name)
				}
			}
		})
	}
}

// Test UPDATE Statement
func TestUpdateStatements(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "Simple UPDATE with WHERE",
			sql:     `UPDATE users SET name = 'Jane' WHERE id = 1`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "UPDATE multiple columns",
			sql:     `UPDATE users SET name = 'Jane', email = 'jane@test.com', active = 1 WHERE id = 1`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "UPDATE without WHERE",
			sql:     `UPDATE users SET active = 1`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "UPDATE with expressions",
			sql:     `UPDATE products SET price = price * 1.1 WHERE category = 'electronics'`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "UPDATE with ORDER BY and LIMIT (MySQL)",
			sql:     `UPDATE users SET status = 'inactive' WHERE last_login < '2020-01-01' ORDER BY last_login LIMIT 100`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "UPDATE with complex WHERE",
			sql:     `UPDATE users SET status = 'premium' WHERE credits > 1000 AND active = 1`,
			dialect: "postgresql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("UPDATE parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				if updateStmt, ok := stmt.(*parser.UpdateStatement); ok {
					t.Logf("Successfully parsed UPDATE on %s with %d SET clauses", updateStmt.Table.Name, len(updateStmt.Set))
				}
			}
		})
	}
}

// Test DELETE Statement
func TestDeleteStatements(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "Simple DELETE with WHERE",
			sql:     `DELETE FROM users WHERE id = 1`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "DELETE without WHERE",
			sql:     `DELETE FROM temp_data`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "DELETE with complex WHERE",
			sql:     `DELETE FROM users WHERE created_at < '2020-01-01' AND active = 0`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "DELETE with ORDER BY and LIMIT (MySQL)",
			sql:     `DELETE FROM logs WHERE created_at < '2020-01-01' ORDER BY created_at LIMIT 1000`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "DELETE with IN clause",
			sql:     `DELETE FROM users WHERE id IN (1, 2, 3, 4, 5)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "DELETE with LIKE",
			sql:     `DELETE FROM spam WHERE email LIKE '%@spam.com'`,
			dialect: "mysql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("DELETE parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				if deleteStmt, ok := stmt.(*parser.DeleteStatement); ok {
					t.Logf("Successfully parsed DELETE FROM %s", deleteStmt.From.Name)
				}
			}
		})
	}
}

// Test DML with different dialects
func TestDMLDialectSupport(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "MySQL INSERT with backticks",
			sql:     "INSERT INTO `users` (`id`, `name`) VALUES (1, 'John')",
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "PostgreSQL UPDATE with double quotes",
			sql:     `UPDATE "users" SET "name" = 'Jane' WHERE "id" = 1`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "SQL Server DELETE with brackets",
			sql:     `DELETE FROM [users] WHERE [id] = 1`,
			dialect: "sqlserver",
			wantErr: false,
		},
		{
			name:    "SQLite INSERT",
			sql:     `INSERT INTO users (id, name) VALUES (1, 'John')`,
			dialect: "sqlite",
			wantErr: false,
		},
		{
			name:    "Oracle UPDATE",
			sql:     `UPDATE users SET name = 'Jane' WHERE id = 1`,
			dialect: "oracle",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("DML dialect parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed %s statement with %s dialect", stmt.Type(), tt.dialect)
			}
		})
	}
}

// Test complex DML scenarios
func TestComplexDMLScenarios(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		// TODO: Subqueries in VALUES not yet supported - need to enhance parseExpression
		// {
		// 	name:    "INSERT with subquery in VALUES",
		// 	sql:     `INSERT INTO user_stats (user_id, order_count) VALUES (1, (SELECT COUNT(*) FROM orders WHERE user_id = 1))`,
		// 	wantErr: false,
		// },
		{
			name:    "UPDATE with CASE expression",
			sql:     `UPDATE users SET status = CASE WHEN credits > 1000 THEN 'premium' WHEN credits > 100 THEN 'standard' ELSE 'basic' END`,
			wantErr: false,
		},
		// TODO: EXISTS not yet supported in parseExpression
		// {
		// 	name:    "DELETE with EXISTS",
		// 	sql:     `DELETE FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND total > 1000)`,
		// 	wantErr: false,
		// },
		{
			name:    "INSERT with multiple value lists",
			sql:     `INSERT INTO products (name, price, category) VALUES ('Product A', 10.99, 'electronics'), ('Product B', 20.50, 'books'), ('Product C', 15.00, 'clothing')`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Complex DML parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed complex %s statement", stmt.Type())
			}
		})
	}
}
