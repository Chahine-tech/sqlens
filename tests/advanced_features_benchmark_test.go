package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// Benchmark Subqueries
func BenchmarkSubqueryScalar(b *testing.B) {
	sql := `SELECT * FROM users WHERE salary > (SELECT AVG(salary) FROM employees)`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSubqueryExists(b *testing.B) {
	sql := `SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)`
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSubqueryIN(b *testing.B) {
	sql := `SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > 1000)`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSubqueryDerivedTable(b *testing.B) {
	sql := `SELECT * FROM (SELECT id, name FROM users WHERE active = 1) AS active_users WHERE id > 100`
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSubqueryNested(b *testing.B) {
	sql := `SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE product_id IN (SELECT id FROM products WHERE category_id IN (SELECT id FROM categories WHERE name = 'electronics')))`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSubqueryCorrelated(b *testing.B) {
	sql := `SELECT * FROM users u WHERE salary > (SELECT AVG(salary) FROM employees e WHERE e.department = u.department)`
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark CTEs (WITH clause)
func BenchmarkCTESimple(b *testing.B) {
	sql := `WITH sales_summary AS (SELECT product_id, SUM(amount) as total FROM sales GROUP BY product_id) SELECT * FROM sales_summary WHERE total > 1000`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCTEMultiple(b *testing.B) {
	sql := `WITH users_cte AS (SELECT id, name FROM users WHERE active = 1), orders_cte AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) SELECT u.name, o.order_count FROM users_cte u JOIN orders_cte o ON u.id = o.user_id`
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCTERecursive(b *testing.B) {
	sql := `WITH RECURSIVE employee_hierarchy AS (SELECT id, name, manager_id FROM employees WHERE manager_id IS NULL UNION ALL SELECT e.id, e.name, e.manager_id FROM employees e INNER JOIN employee_hierarchy eh ON e.manager_id = eh.id) SELECT * FROM employee_hierarchy`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Window Functions
func BenchmarkWindowFunctionSimple(b *testing.B) {
	sql := `SELECT id, name, salary, ROW_NUMBER() OVER (ORDER BY salary DESC) as rank FROM employees`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWindowFunctionPartition(b *testing.B) {
	sql := `SELECT id, department, salary, RANK() OVER (PARTITION BY department ORDER BY salary DESC) as dept_rank FROM employees`
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWindowFunctionFrame(b *testing.B) {
	sql := `SELECT id, date, amount, SUM(amount) OVER (ORDER BY date ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) as rolling_sum FROM sales`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Set Operations
func BenchmarkSetOperationUnion(b *testing.B) {
	sql := `SELECT id, name FROM customers UNION SELECT id, name FROM prospects`
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetOperationUnionAll(b *testing.B) {
	sql := `SELECT id FROM customers UNION ALL SELECT id FROM prospects UNION ALL SELECT id FROM partners`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetOperationIntersect(b *testing.B) {
	sql := `SELECT id FROM customers INTERSECT SELECT id FROM premium_users`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetOperationExcept(b *testing.B) {
	sql := `SELECT id FROM all_users EXCEPT SELECT id FROM banned_users`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark DML with Subqueries
func BenchmarkInsertWithSubquery(b *testing.B) {
	sql := `INSERT INTO user_stats (user_id, order_count) VALUES (1, (SELECT COUNT(*) FROM orders WHERE user_id = 1))`
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInsertSelect(b *testing.B) {
	sql := `INSERT INTO archive SELECT * FROM users WHERE created_at < '2020-01-01'`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUpdateWithSubquery(b *testing.B) {
	sql := `UPDATE users SET status = (SELECT status FROM user_preferences WHERE user_id = users.id) WHERE id IN (SELECT user_id FROM active_sessions)`
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeleteWithSubquery(b *testing.B) {
	sql := `DELETE FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = users.id AND total > 10000)`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Comprehensive benchmark combining multiple features
func BenchmarkComplexQueryWithAllFeatures(b *testing.B) {
	sql := `WITH ranked_sales AS (
		SELECT
			product_id,
			category,
			amount,
			ROW_NUMBER() OVER (PARTITION BY category ORDER BY amount DESC) as rank
		FROM sales
		WHERE amount > (SELECT AVG(amount) FROM sales)
	)
	SELECT
		r.category,
		r.product_id,
		r.amount,
		(SELECT COUNT(*) FROM products WHERE category = r.category) as total_products
	FROM ranked_sales r
	WHERE r.rank <= 10
		AND EXISTS (SELECT 1 FROM products p WHERE p.id = r.product_id AND p.active = 1)
	ORDER BY r.category, r.rank`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}
