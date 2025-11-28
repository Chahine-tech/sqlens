package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// Test CTE (Common Table Expressions) - WITH clause
func TestCTEParsing(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name: "Simple CTE",
			sql: `WITH cte AS (
				SELECT id, name FROM users
			)
			SELECT * FROM cte`,
			wantErr: false,
		},
		{
			name: "CTE with column list",
			sql: `WITH employee_cte (emp_id, emp_name) AS (
				SELECT id, name FROM employees
			)
			SELECT emp_id, emp_name FROM employee_cte`,
			wantErr: false,
		},
		{
			name: "Multiple CTEs",
			sql: `WITH
				users_cte AS (SELECT id, name FROM users),
				orders_cte AS (SELECT user_id, total FROM orders)
			SELECT u.name, o.total
			FROM users_cte u
			JOIN orders_cte o ON u.id = o.user_id`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("sqlserver"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("CTE parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			// Type assertion to check it's a WithStatement
			if !tt.wantErr {
				if withStmt, ok := stmt.(*parser.WithStatement); ok {
					if len(withStmt.CTEs) == 0 {
						t.Error("Expected at least one CTE")
					}
					t.Logf("Successfully parsed CTE with %d CTEs", len(withStmt.CTEs))
				} else {
					t.Logf("Statement type: %T", stmt)
				}
			}
		})
	}
}

// Test Window Functions
func TestWindowFunctionParsing(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name: "ROW_NUMBER with ORDER BY",
			sql: `SELECT
				id,
				name,
				ROW_NUMBER() OVER (ORDER BY id) as row_num
			FROM users`,
			wantErr: false,
		},
		{
			name: "RANK with PARTITION BY and ORDER BY",
			sql: `SELECT
				department,
				employee,
				salary,
				RANK() OVER (PARTITION BY department ORDER BY salary DESC) as rank
			FROM employees`,
			wantErr: false,
		},
		{
			name: "Multiple window functions",
			sql: `SELECT
				name,
				ROW_NUMBER() OVER (ORDER BY id) as rn,
				RANK() OVER (ORDER BY score DESC) as rank,
				DENSE_RANK() OVER (ORDER BY score DESC) as dense_rank
			FROM students`,
			wantErr: false,
		},
		{
			name: "Window function with frame clause",
			sql: `SELECT
				date,
				amount,
				SUM(amount) OVER (
					ORDER BY date
					ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
				) as running_total
			FROM transactions`,
			wantErr: false,
		},
		{
			name: "RANGE frame clause",
			sql: `SELECT
				id,
				value,
				AVG(value) OVER (
					ORDER BY id
					RANGE BETWEEN 2 PRECEDING AND 2 FOLLOWING
				) as avg_value
			FROM data`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("sqlserver"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Window function parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed window function query")
			}
		})
	}
}

// Test Set Operations (UNION, INTERSECT, EXCEPT)
func TestSetOperations(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name: "Simple UNION",
			sql: `SELECT id, name FROM users
			UNION
			SELECT id, name FROM customers`,
			wantErr: false,
		},
		{
			name: "UNION ALL",
			sql: `SELECT id FROM table1
			UNION ALL
			SELECT id FROM table2`,
			wantErr: false,
		},
		{
			name: "INTERSECT",
			sql: `SELECT id FROM users
			INTERSECT
			SELECT user_id FROM orders`,
			wantErr: false,
		},
		{
			name: "EXCEPT",
			sql: `SELECT id FROM all_users
			EXCEPT
			SELECT id FROM inactive_users`,
			wantErr: false,
		},
		{
			name: "Chained set operations",
			sql: `SELECT id FROM table1
			UNION
			SELECT id FROM table2
			UNION
			SELECT id FROM table3`,
			wantErr: false,
		},
		{
			name: "Mixed set operations",
			sql: `SELECT id FROM table1
			UNION ALL
			SELECT id FROM table2
			INTERSECT
			SELECT id FROM table3`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("sqlserver"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Set operation parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				if setOp, ok := stmt.(*parser.SetOperation); ok {
					t.Logf("Successfully parsed %s operation", setOp.Operator)
				}
			}
		})
	}
}

// Test Bug Fixes - Column Aliases
func TestColumnAliases(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "Simple column alias",
			sql:     `SELECT id, name as employee_name FROM employees`,
			wantErr: false,
		},
		{
			name:    "Aggregate function with alias",
			sql:     `SELECT product_id, SUM(amount) as total FROM sales GROUP BY product_id`,
			wantErr: false,
		},
		{
			name:    "Multiple aliases",
			sql:     `SELECT id as emp_id, name as emp_name, salary as emp_salary FROM employees`,
			wantErr: false,
		},
		{
			name:    "Alias in CTE",
			sql:     `WITH test AS (SELECT id, name as employee_name FROM employees) SELECT * FROM test`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("sqlserver"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Column alias parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed query with column aliases")
			}
		})
	}
}

// Test Bug Fixes - CTEs with GROUP BY/HAVING
func TestCTEWithGroupBy(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name: "CTE with GROUP BY",
			sql: `WITH sales_summary AS (
				SELECT product_id, SUM(amount) as total
				FROM sales
				GROUP BY product_id
			)
			SELECT * FROM sales_summary`,
			wantErr: false,
		},
		{
			name: "CTE with GROUP BY and HAVING",
			sql: `WITH top_products AS (
				SELECT product_id, SUM(amount) as total
				FROM sales
				GROUP BY product_id
				HAVING SUM(amount) > 1000
			)
			SELECT * FROM top_products`,
			wantErr: false,
		},
		{
			name: "CTE with aggregate alias and GROUP BY",
			sql: `WITH revenue AS (
				SELECT dept_id, AVG(salary) as avg_salary
				FROM employees
				GROUP BY dept_id
			)
			SELECT * FROM revenue WHERE avg_salary > 50000`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("sqlserver"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("CTE with GROUP BY parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				if withStmt, ok := stmt.(*parser.WithStatement); ok {
					t.Logf("Successfully parsed CTE with %d CTEs", len(withStmt.CTEs))
				}
			}
		})
	}
}

// Test Bug Fixes - Window Functions with Aliases in CTEs
func TestWindowFunctionsInCTEs(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name: "Window function with alias in CTE",
			sql: `WITH ranked AS (
				SELECT employee_id, ROW_NUMBER() OVER (ORDER BY salary DESC) as rank
				FROM employees
			)
			SELECT * FROM ranked`,
			wantErr: false,
		},
		{
			name: "Multiple window functions in CTE",
			sql: `WITH analytics AS (
				SELECT
					employee_id,
					ROW_NUMBER() OVER (ORDER BY salary DESC) as rn,
					RANK() OVER (PARTITION BY department ORDER BY salary DESC) as dept_rank
				FROM employees
			)
			SELECT * FROM analytics`,
			wantErr: false,
		},
		{
			name: "Window function with frame and alias in CTE",
			sql: `WITH running_totals AS (
				SELECT
					date,
					SUM(amount) OVER (
						ORDER BY date
						ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
					) as running_total
				FROM transactions
			)
			SELECT * FROM running_totals`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("sqlserver"))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Window function in CTE parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				if withStmt, ok := stmt.(*parser.WithStatement); ok {
					t.Logf("Successfully parsed window functions in CTE with %d CTEs", len(withStmt.CTEs))
				}
			}
		})
	}
}

// Test CASE Expressions - TODO: Requires refactoring expression parser
// CASE expressions consume THEN keyword during expression parsing
// This needs a more sophisticated expression parser that understands context
/*
func TestCaseExpressions(t *testing.T) {
	// TODO: Implement proper CASE expression parsing
	// The challenge is that parseExpression() consumes THEN as part of the expression
	// Need to implement expression parsing with boundary detection
}
*/
