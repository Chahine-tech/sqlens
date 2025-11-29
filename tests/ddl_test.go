package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// Test CREATE TABLE statements
func TestCreateTable(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "Simple CREATE TABLE",
			sql:     `CREATE TABLE users (id INT, name VARCHAR(100))`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with IF NOT EXISTS",
			sql:     `CREATE TABLE IF NOT EXISTS products (id INT, name VARCHAR(255))`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with PRIMARY KEY",
			sql:     `CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with NOT NULL",
			sql:     `CREATE TABLE users (id INT NOT NULL, email VARCHAR(100) NOT NULL)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with AUTO_INCREMENT (MySQL)",
			sql:     `CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(100))`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with AUTOINCREMENT (SQLite)",
			sql:     `CREATE TABLE users (id INTEGER AUTOINCREMENT PRIMARY KEY, name TEXT)`,
			dialect: "sqlite",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with DEFAULT value",
			sql:     `CREATE TABLE users (id INT, status VARCHAR(20) DEFAULT 'active')`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with UNIQUE constraint",
			sql:     `CREATE TABLE users (id INT, email VARCHAR(100) UNIQUE)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with inline FOREIGN KEY",
			sql:     `CREATE TABLE orders (id INT, user_id INT REFERENCES users(id))`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with FOREIGN KEY and ON DELETE",
			sql:     `CREATE TABLE orders (id INT, user_id INT REFERENCES users(id) ON DELETE CASCADE)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with FOREIGN KEY and ON UPDATE",
			sql:     `CREATE TABLE orders (id INT, user_id INT REFERENCES users(id) ON UPDATE SET NULL)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with data type lengths",
			sql:     `CREATE TABLE products (id INT, name VARCHAR(255), price DECIMAL(10, 2))`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with table-level PRIMARY KEY",
			sql:     `CREATE TABLE users (id INT, email VARCHAR(100), PRIMARY KEY (id))`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with composite PRIMARY KEY",
			sql:     `CREATE TABLE user_roles (user_id INT, role_id INT, PRIMARY KEY (user_id, role_id))`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with table-level FOREIGN KEY",
			sql:     `CREATE TABLE orders (id INT, user_id INT, FOREIGN KEY (user_id) REFERENCES users(id))`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with named constraint",
			sql:     `CREATE TABLE orders (id INT, user_id INT, CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id))`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with UNIQUE constraint",
			sql:     `CREATE TABLE users (id INT, email VARCHAR(100), UNIQUE (email))`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with multiple constraints",
			sql:     `CREATE TABLE users (id INT PRIMARY KEY, email VARCHAR(100) NOT NULL UNIQUE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE with schema name",
			sql:     `CREATE TABLE myschema.users (id INT, name VARCHAR(100))`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE TABLE complex example",
			sql: `CREATE TABLE IF NOT EXISTS orders (
				id INT AUTO_INCREMENT PRIMARY KEY,
				user_id INT NOT NULL,
				product_id INT NOT NULL,
				quantity INT DEFAULT 1,
				total DECIMAL(10, 2),
				status VARCHAR(20) DEFAULT 'pending',
				created_at TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
				FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE SET NULL
			)`,
			dialect: "mysql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("CREATE TABLE parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				createStmt, ok := stmt.(*parser.CreateTableStatement)
				if !ok {
					t.Errorf("Expected CreateTableStatement, got %T", stmt)
				} else {
					t.Logf("Successfully parsed CREATE TABLE: %s", createStmt.Table.Name)
				}
			}
		})
	}
}

// Test DROP statements
func TestDropStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "DROP TABLE",
			sql:     `DROP TABLE users`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "DROP TABLE IF EXISTS",
			sql:     `DROP TABLE IF EXISTS users`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "DROP TABLE CASCADE",
			sql:     `DROP TABLE users CASCADE`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "DROP TABLE IF EXISTS CASCADE",
			sql:     `DROP TABLE IF EXISTS users CASCADE`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "DROP DATABASE",
			sql:     `DROP DATABASE mydb`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "DROP DATABASE IF EXISTS",
			sql:     `DROP DATABASE IF EXISTS mydb`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "DROP SCHEMA",
			sql:     `DROP SCHEMA myschema`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "DROP INDEX",
			sql:     `DROP INDEX idx_users_email`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "DROP INDEX IF EXISTS",
			sql:     `DROP INDEX IF EXISTS idx_users_email`,
			dialect: "postgresql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("DROP statement parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				dropStmt, ok := stmt.(*parser.DropStatement)
				if !ok {
					t.Errorf("Expected DropStatement, got %T", stmt)
				} else {
					t.Logf("Successfully parsed DROP %s: %s", dropStmt.ObjectType, dropStmt.ObjectName)
				}
			}
		})
	}
}

// Test ALTER TABLE statements
func TestAlterTable(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "ALTER TABLE ADD COLUMN",
			sql:     `ALTER TABLE users ADD COLUMN age INT`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE ADD COLUMN with constraints",
			sql:     `ALTER TABLE users ADD COLUMN email VARCHAR(100) NOT NULL UNIQUE`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE ADD without COLUMN keyword",
			sql:     `ALTER TABLE users ADD age INT`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE DROP COLUMN",
			sql:     `ALTER TABLE users DROP COLUMN age`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE DROP without COLUMN keyword",
			sql:     `ALTER TABLE users DROP age`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE MODIFY COLUMN",
			sql:     `ALTER TABLE users MODIFY COLUMN age BIGINT`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE MODIFY without COLUMN keyword",
			sql:     `ALTER TABLE users MODIFY age BIGINT NOT NULL`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE CHANGE COLUMN (MySQL)",
			sql:     `ALTER TABLE users CHANGE COLUMN old_name new_name VARCHAR(100)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE CHANGE without COLUMN keyword",
			sql:     `ALTER TABLE users CHANGE old_name new_name VARCHAR(100)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE ADD CONSTRAINT",
			sql:     `ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE ADD PRIMARY KEY",
			sql:     `ALTER TABLE users ADD PRIMARY KEY (id)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE ADD UNIQUE",
			sql:     `ALTER TABLE users ADD UNIQUE (email)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "ALTER TABLE DROP CONSTRAINT",
			sql:     `ALTER TABLE orders DROP CONSTRAINT fk_user`,
			dialect: "postgresql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("ALTER TABLE parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				alterStmt, ok := stmt.(*parser.AlterTableStatement)
				if !ok {
					t.Errorf("Expected AlterTableStatement, got %T", stmt)
				} else {
					t.Logf("Successfully parsed ALTER TABLE: %s (%s)", alterStmt.Table.Name, alterStmt.Action.ActionType)
				}
			}
		})
	}
}

// Test CREATE INDEX statements
func TestCreateIndex(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "CREATE INDEX",
			sql:     `CREATE INDEX idx_users_email ON users (email)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE INDEX IF NOT EXISTS",
			sql:     `CREATE INDEX IF NOT EXISTS idx_users_email ON users (email)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE UNIQUE INDEX",
			sql:     `CREATE UNIQUE INDEX idx_users_email ON users (email)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "CREATE UNIQUE INDEX IF NOT EXISTS",
			sql:     `CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email)`,
			dialect: "sqlite",
			wantErr: false,
		},
		{
			name:    "CREATE INDEX on multiple columns",
			sql:     `CREATE INDEX idx_users_name_email ON users (name, email)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "CREATE INDEX with schema",
			sql:     `CREATE INDEX idx_users_email ON myschema.users (email)`,
			dialect: "postgresql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("CREATE INDEX parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				indexStmt, ok := stmt.(*parser.CreateIndexStatement)
				if !ok {
					t.Errorf("Expected CreateIndexStatement, got %T", stmt)
				} else {
					t.Logf("Successfully parsed CREATE INDEX: %s ON %s", indexStmt.IndexName, indexStmt.Table.Name)
				}
			}
		})
	}
}

// Test DDL across different dialects
func TestDDLDialects(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "MySQL AUTO_INCREMENT",
			sql:     `CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "SQLite AUTOINCREMENT",
			sql:     `CREATE TABLE users (id INTEGER AUTOINCREMENT PRIMARY KEY)`,
			dialect: "sqlite",
			wantErr: false,
		},
		{
			name:    "PostgreSQL SERIAL",
			sql:     `CREATE TABLE users (id SERIAL PRIMARY KEY)`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "SQL Server IDENTITY",
			sql:     `CREATE TABLE users (id INT IDENTITY PRIMARY KEY)`,
			dialect: "sqlserver",
			wantErr: false,
		},
		{
			name:    "MySQL CHANGE COLUMN",
			sql:     `ALTER TABLE users CHANGE COLUMN old_name new_name VARCHAR(100)`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "PostgreSQL CASCADE",
			sql:     `DROP TABLE users CASCADE`,
			dialect: "postgresql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("DDL dialect parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed %s dialect DDL", tt.dialect)
			}
		})
	}
}
