package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// Benchmark CREATE TABLE statements
func BenchmarkCreateTableSimple(b *testing.B) {
	sql := `CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(255) UNIQUE)`
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

func BenchmarkCreateTableIfNotExists(b *testing.B) {
	sql := `CREATE TABLE IF NOT EXISTS products (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, price DECIMAL(10,2) DEFAULT 0.00)`
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

func BenchmarkCreateTableWithForeignKeys(b *testing.B) {
	sql := `CREATE TABLE orders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		product_id INT NOT NULL,
		quantity INT DEFAULT 1,
		total DECIMAL(10, 2),
		status VARCHAR(20) DEFAULT 'pending',
		created_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE SET NULL
	)`
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

func BenchmarkCreateTableCompositePK(b *testing.B) {
	sql := `CREATE TABLE user_roles (
		user_id INT,
		role_id INT,
		granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, role_id)
	)`
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

func BenchmarkCreateTableComplex(b *testing.B) {
	sql := `CREATE TABLE IF NOT EXISTS inventory (
		id INT AUTO_INCREMENT PRIMARY KEY,
		product_id INT NOT NULL,
		warehouse_id INT NOT NULL,
		quantity INT NOT NULL DEFAULT 0,
		reorder_level INT DEFAULT 10,
		last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		notes TEXT,
		CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT fk_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE RESTRICT,
		CONSTRAINT unique_product_warehouse UNIQUE (product_id, warehouse_id)
	)`
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

// Benchmark DROP statements
func BenchmarkDropTable(b *testing.B) {
	sql := `DROP TABLE users`
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

func BenchmarkDropTableIfExists(b *testing.B) {
	sql := `DROP TABLE IF EXISTS users`
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

func BenchmarkDropTableCascade(b *testing.B) {
	sql := `DROP TABLE IF EXISTS users CASCADE`
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

func BenchmarkDropDatabase(b *testing.B) {
	sql := `DROP DATABASE IF EXISTS test_db`
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

func BenchmarkDropIndex(b *testing.B) {
	sql := `DROP INDEX IF EXISTS idx_users_email`
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

// Benchmark ALTER TABLE statements
func BenchmarkAlterTableAddColumn(b *testing.B) {
	sql := `ALTER TABLE users ADD COLUMN age INT NOT NULL`
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

func BenchmarkAlterTableAddColumnWithConstraints(b *testing.B) {
	sql := `ALTER TABLE users ADD COLUMN email VARCHAR(100) NOT NULL UNIQUE`
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

func BenchmarkAlterTableDropColumn(b *testing.B) {
	sql := `ALTER TABLE users DROP COLUMN age`
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

func BenchmarkAlterTableModifyColumn(b *testing.B) {
	sql := `ALTER TABLE users MODIFY COLUMN name VARCHAR(150) NOT NULL`
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

func BenchmarkAlterTableChangeColumn(b *testing.B) {
	sql := `ALTER TABLE users CHANGE COLUMN old_name new_name VARCHAR(100)`
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

func BenchmarkAlterTableAddConstraint(b *testing.B) {
	sql := `ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)`
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

func BenchmarkAlterTableAddPrimaryKey(b *testing.B) {
	sql := `ALTER TABLE users ADD PRIMARY KEY (id)`
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

func BenchmarkAlterTableDropConstraint(b *testing.B) {
	sql := `ALTER TABLE orders DROP CONSTRAINT fk_user`
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

// Benchmark CREATE INDEX statements
func BenchmarkCreateIndex(b *testing.B) {
	sql := `CREATE INDEX idx_users_email ON users (email)`
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

func BenchmarkCreateUniqueIndex(b *testing.B) {
	sql := `CREATE UNIQUE INDEX idx_users_email ON users (email)`
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

func BenchmarkCreateIndexIfNotExists(b *testing.B) {
	sql := `CREATE INDEX IF NOT EXISTS idx_products_category ON products (category)`
	d := dialect.GetDialect("sqlite")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCreateIndexMultiColumn(b *testing.B) {
	sql := `CREATE INDEX idx_orders_user_product ON orders (user_id, product_id)`
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

// Benchmark DDL across different dialects
func BenchmarkDDLMySQL(b *testing.B) {
	sql := `CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100))`
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

func BenchmarkDDLPostgreSQL(b *testing.B) {
	sql := `CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100))`
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

func BenchmarkDDLSQLServer(b *testing.B) {
	sql := `CREATE TABLE users (id INT IDENTITY PRIMARY KEY, name VARCHAR(100))`
	d := dialect.GetDialect("sqlserver")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDDLSQLite(b *testing.B) {
	sql := `CREATE TABLE users (id INTEGER AUTOINCREMENT PRIMARY KEY, name TEXT)`
	d := dialect.GetDialect("sqlite")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}
