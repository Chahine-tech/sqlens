package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// Test subqueries in WHERE clause
func TestSubqueriesInWhere(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "Simple subquery in WHERE",
			sql:     `SELECT * FROM users WHERE salary > (SELECT AVG(salary) FROM employees)`,
			wantErr: false,
		},
		{
			name:    "EXISTS subquery",
			sql:     `SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)`,
			wantErr: false,
		},
		{
			name:    "NOT EXISTS subquery",
			sql:     `SELECT * FROM users WHERE NOT EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)`,
			wantErr: false,
		},
		{
			name:    "IN with subquery",
			sql:     `SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > 1000)`,
			wantErr: false,
		},
		{
			name:    "NOT IN with subquery",
			sql:     `SELECT * FROM users WHERE id NOT IN (SELECT user_id FROM banned_users)`,
			wantErr: false,
		},
		{
			name:    "Multiple subqueries",
			sql:     `SELECT * FROM users WHERE salary > (SELECT AVG(salary) FROM employees) AND id IN (SELECT user_id FROM active_users)`,
			wantErr: false,
		},
		{
			name:    "Nested subqueries",
			sql:     `SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE product_id IN (SELECT id FROM products WHERE category = 'electronics'))`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Subquery in WHERE parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed query with subquery in WHERE")
			}
		})
	}
}

// Test subqueries in SELECT clause
func TestSubqueriesInSelect(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "Scalar subquery in SELECT",
			sql:     `SELECT id, name, (SELECT COUNT(*) FROM orders WHERE user_id = users.id) as order_count FROM users`,
			wantErr: false,
		},
		{
			name:    "Multiple scalar subqueries",
			sql:     `SELECT id, (SELECT AVG(price) FROM products) as avg_price, (SELECT MAX(price) FROM products) as max_price FROM users`,
			wantErr: false,
		},
		{
			name:    "Subquery with aggregation",
			sql:     `SELECT department, (SELECT AVG(salary) FROM employees e WHERE e.dept_id = d.id) as avg_salary FROM departments d`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("postgresql"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Subquery in SELECT parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed query with subquery in SELECT")
			}
		})
	}
}

// Test subqueries in FROM clause
func TestSubqueriesInFrom(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "Derived table (subquery in FROM)",
			sql:     `SELECT * FROM (SELECT id, name FROM users WHERE active = 1) AS active_users`,
			wantErr: false,
		},
		{
			name:    "JOIN with subquery",
			sql:     `SELECT u.name, o.total FROM users u JOIN (SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id) o ON u.id = o.user_id`,
			wantErr: false,
		},
		{
			name:    "Multiple derived tables",
			sql:     `SELECT * FROM (SELECT id FROM users) u, (SELECT user_id FROM orders) o WHERE u.id = o.user_id`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Subquery in FROM parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed query with subquery in FROM")
			}
		})
	}
}

// Test subqueries in INSERT statements
func TestSubqueriesInInsert(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "INSERT with SELECT",
			sql:     `INSERT INTO archive SELECT * FROM users WHERE created_at < '2020-01-01'`,
			wantErr: false,
		},
		{
			name:    "INSERT with column list and SELECT",
			sql:     `INSERT INTO user_stats (user_id, order_count) SELECT id, (SELECT COUNT(*) FROM orders WHERE user_id = users.id) FROM users`,
			wantErr: false,
		},
		{
			name:    "INSERT VALUES with subquery",
			sql:     `INSERT INTO user_stats (user_id, order_count) VALUES (1, (SELECT COUNT(*) FROM orders WHERE user_id = 1))`,
			wantErr: false,
		},
		{
			name:    "INSERT with multiple subqueries in VALUES",
			sql:     `INSERT INTO stats (user_id, order_count, total_spent) VALUES (1, (SELECT COUNT(*) FROM orders WHERE user_id = 1), (SELECT SUM(amount) FROM orders WHERE user_id = 1))`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Subquery in INSERT parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed INSERT with subquery")
			}
		})
	}
}

// Test subqueries in UPDATE statements
func TestSubqueriesInUpdate(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "UPDATE SET with subquery",
			sql:     `UPDATE users SET status = (SELECT status FROM user_preferences WHERE user_id = users.id)`,
			wantErr: false,
		},
		{
			name:    "UPDATE with subquery in WHERE",
			sql:     `UPDATE users SET active = 0 WHERE id IN (SELECT user_id FROM banned_users)`,
			wantErr: false,
		},
		{
			name:    "UPDATE with EXISTS in WHERE",
			sql:     `UPDATE users SET premium = 1 WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id AND total > 10000)`,
			wantErr: false,
		},
		{
			name:    "UPDATE with multiple subqueries",
			sql:     `UPDATE products SET price = (SELECT AVG(price) FROM products WHERE category = products.category) WHERE id IN (SELECT product_id FROM discounted_products)`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("postgresql"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Subquery in UPDATE parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed UPDATE with subquery")
			}
		})
	}
}

// Test subqueries in DELETE statements
func TestSubqueriesInDelete(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "DELETE with IN subquery",
			sql:     `DELETE FROM users WHERE id IN (SELECT user_id FROM inactive_users)`,
			wantErr: false,
		},
		{
			name:    "DELETE with EXISTS",
			sql:     `DELETE FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND total > 1000)`,
			wantErr: false,
		},
		{
			name:    "DELETE with NOT EXISTS",
			sql:     `DELETE FROM users WHERE NOT EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)`,
			wantErr: false,
		},
		{
			name:    "DELETE with scalar subquery",
			sql:     `DELETE FROM products WHERE price < (SELECT AVG(price) FROM products)`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Subquery in DELETE parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed DELETE with subquery")
			}
		})
	}
}

// Test complex nested and correlated subqueries
func TestComplexSubqueries(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "Triple nested subquery",
			sql:     `SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE product_id IN (SELECT id FROM products WHERE category_id IN (SELECT id FROM categories WHERE name = 'electronics')))`,
			wantErr: false,
		},
		{
			name:    "Correlated subquery",
			sql:     `SELECT * FROM users u WHERE salary > (SELECT AVG(salary) FROM employees e WHERE e.department = u.department)`,
			wantErr: false,
		},
		{
			name:    "Multiple EXISTS clauses",
			sql:     `SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id) AND EXISTS (SELECT 1 FROM payments WHERE user_id = users.id)`,
			wantErr: false,
		},
		{
			name:    "Subquery with JOIN",
			sql:     `SELECT * FROM users WHERE id IN (SELECT o.user_id FROM orders o JOIN products p ON o.product_id = p.id WHERE p.category = 'electronics')`,
			wantErr: false,
		},
		{
			name:    "Subquery with aggregation and HAVING",
			sql:     `SELECT * FROM users WHERE id IN (SELECT user_id FROM orders GROUP BY user_id HAVING SUM(total) > 1000)`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("postgresql"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Complex subquery parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed complex subquery")
			}
		})
	}
}

// Test subqueries across different dialects
func TestSubqueriesDialectSupport(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "MySQL subquery",
			sql:     `SELECT * FROM users WHERE id IN (SELECT user_id FROM orders)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "PostgreSQL EXISTS",
			sql:     `SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "SQL Server subquery in UPDATE",
			sql:     `UPDATE users SET status = (SELECT TOP 1 status FROM user_status WHERE user_id = users.id)`,
			dialect: "sqlserver",
			wantErr: false,
		},
		{
			name:    "SQLite subquery",
			sql:     `SELECT * FROM users WHERE salary > (SELECT AVG(salary) FROM employees)`,
			dialect: "sqlite",
			wantErr: false,
		},
		{
			name:    "Oracle subquery",
			sql:     `SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE ROWNUM <= 10)`,
			dialect: "oracle",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Subquery dialect parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed subquery with %s dialect", tt.dialect)
			}
		})
	}
}
