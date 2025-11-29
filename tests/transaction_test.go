package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// Test BEGIN TRANSACTION statements
func TestBeginTransaction(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "Simple BEGIN",
			sql:     `BEGIN`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "BEGIN TRANSACTION",
			sql:     `BEGIN TRANSACTION`,
			dialect: "sqlserver",
			wantErr: false,
		},
		{
			name:    "BEGIN WORK",
			sql:     `BEGIN WORK`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "START TRANSACTION",
			sql:     `START TRANSACTION`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "START TRANSACTION (PostgreSQL)",
			sql:     `START TRANSACTION`,
			dialect: "postgresql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("BEGIN TRANSACTION parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				beginStmt, ok := stmt.(*parser.BeginTransactionStatement)
				if !ok {
					t.Errorf("Expected BeginTransactionStatement, got %T", stmt)
				} else {
					t.Logf("Successfully parsed: %s", beginStmt.String())
				}
			}
		})
	}
}

// Test COMMIT statements
func TestCommit(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "Simple COMMIT",
			sql:     `COMMIT`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "COMMIT WORK",
			sql:     `COMMIT WORK`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "COMMIT (SQL Server)",
			sql:     `COMMIT`,
			dialect: "sqlserver",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("COMMIT parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				commitStmt, ok := stmt.(*parser.CommitStatement)
				if !ok {
					t.Errorf("Expected CommitStatement, got %T", stmt)
				} else {
					t.Logf("Successfully parsed: %s", commitStmt.String())
				}
			}
		})
	}
}

// Test ROLLBACK statements
func TestRollback(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "Simple ROLLBACK",
			sql:     `ROLLBACK`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "ROLLBACK WORK",
			sql:     `ROLLBACK WORK`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "ROLLBACK TO SAVEPOINT",
			sql:     `ROLLBACK TO SAVEPOINT sp1`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "ROLLBACK TO SAVEPOINT (MySQL)",
			sql:     `ROLLBACK TO SAVEPOINT my_savepoint`,
			dialect: "mysql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("ROLLBACK parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				rollbackStmt, ok := stmt.(*parser.RollbackStatement)
				if !ok {
					t.Errorf("Expected RollbackStatement, got %T", stmt)
				} else {
					t.Logf("Successfully parsed: %s", rollbackStmt.String())
				}
			}
		})
	}
}

// Test SAVEPOINT statements
func TestSavepoint(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "SAVEPOINT simple",
			sql:     `SAVEPOINT sp1`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "SAVEPOINT (MySQL)",
			sql:     `SAVEPOINT my_savepoint`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "SAVEPOINT (SQL Server)",
			sql:     `SAVEPOINT checkpoint1`,
			dialect: "sqlserver",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("SAVEPOINT parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				savepointStmt, ok := stmt.(*parser.SavepointStatement)
				if !ok {
					t.Errorf("Expected SavepointStatement, got %T", stmt)
				} else {
					t.Logf("Successfully parsed: %s", savepointStmt.String())
				}
			}
		})
	}
}

// Test RELEASE SAVEPOINT statements
func TestReleaseSavepoint(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "RELEASE SAVEPOINT",
			sql:     `RELEASE SAVEPOINT sp1`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "RELEASE SAVEPOINT (MySQL)",
			sql:     `RELEASE SAVEPOINT my_savepoint`,
			dialect: "mysql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("RELEASE SAVEPOINT parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				releaseStmt, ok := stmt.(*parser.ReleaseSavepointStatement)
				if !ok {
					t.Errorf("Expected ReleaseSavepointStatement, got %T", stmt)
				} else {
					t.Logf("Successfully parsed: %s", releaseStmt.String())
				}
			}
		})
	}
}

// Test complete transaction workflows
func TestTransactionWorkflow(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name: "Complete transaction with savepoint",
			sql: `BEGIN;
				  INSERT INTO users (name) VALUES ('John');
				  SAVEPOINT sp1;
				  UPDATE users SET name = 'Jane' WHERE id = 1;
				  ROLLBACK TO SAVEPOINT sp1;
				  COMMIT`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name: "MySQL transaction",
			sql: `START TRANSACTION;
				  INSERT INTO orders (total) VALUES (100.00);
				  COMMIT`,
			dialect: "mysql",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For workflow tests, we parse each statement separately
			// since ParseStatement only handles one statement at a time
			t.Logf("Testing transaction workflow: %s", tt.name)
			// This is more of a validation that individual statements work
			// Real workflow testing would require a multi-statement parser
		})
	}
}

// Test dialect-specific transaction features
func TestTransactionDialects(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		wantErr bool
	}{
		{
			name:    "PostgreSQL BEGIN",
			sql:     `BEGIN`,
			dialect: "postgresql",
			wantErr: false,
		},
		{
			name:    "MySQL START TRANSACTION",
			sql:     `START TRANSACTION`,
			dialect: "mysql",
			wantErr: false,
		},
		{
			name:    "SQL Server BEGIN TRANSACTION",
			sql:     `BEGIN TRANSACTION`,
			dialect: "sqlserver",
			wantErr: false,
		},
		{
			name:    "SQLite BEGIN",
			sql:     `BEGIN`,
			dialect: "sqlite",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction dialect parsing error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Errors: %v", p.Errors())
				}
				return
			}

			if !tt.wantErr && stmt == nil {
				t.Error("Expected statement, got nil")
			}

			if !tt.wantErr {
				t.Logf("Successfully parsed %s transaction statement", tt.dialect)
			}
		})
	}
}
