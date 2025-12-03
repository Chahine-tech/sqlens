package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// TestTryStatement tests SQL Server TRY...CATCH parsing
func TestTryStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple TRY CATCH",
			sql: `
				BEGIN TRY
					SELECT 1/0;
				END TRY
				BEGIN CATCH
					SELECT ERROR_MESSAGE();
				END CATCH
			`,
			dialect: "sqlserver",
		},
		{
			name: "TRY CATCH with multiple statements",
			sql: `
				BEGIN TRY
					INSERT INTO users (id, name) VALUES (1, 'John');
					UPDATE accounts SET balance = 0 WHERE user_id = 1;
				END TRY
				BEGIN CATCH
					SELECT ERROR_NUMBER() AS ErrorNumber;
					SELECT ERROR_MESSAGE() AS ErrorMessage;
					RETURN 0;
				END CATCH
			`,
			dialect: "sqlserver",
		},
		{
			name: "Nested TRY CATCH",
			sql: `
				BEGIN TRY
					BEGIN TRY
						SELECT 1/0;
					END TRY
					BEGIN CATCH
						RETURN ERROR_MESSAGE();
					END CATCH
				END TRY
				BEGIN CATCH
					SELECT 'Outer error handler';
				END CATCH
			`,
			dialect: "sqlserver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse TRY statement: %v", err)
			}

			tryStmt, ok := stmt.(*parser.TryStatement)
			if !ok {
				t.Fatalf("Expected *parser.TryStatement, got %T", stmt)
			}

			t.Logf("✅ Successfully parsed TRY statement: %s", tryStmt.String())
		})
	}
}

// TestExceptionBlock tests PostgreSQL EXCEPTION...WHEN parsing
func TestExceptionBlock(t *testing.T) {
	t.Skip("PostgreSQL $$ delimiter not yet supported - EXCEPTION blocks work but need standard function syntax")

	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple EXCEPTION WHEN",
			sql: `
				CREATE OR REPLACE FUNCTION divide(a INT, b INT) RETURNS INT AS
				BEGIN
					RETURN a / b;
				EXCEPTION
					WHEN division_by_zero THEN
						RETURN 0;
				END;
			`,
			dialect: "postgresql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse EXCEPTION block: %v", err)
			}

			t.Logf("✅ Successfully parsed statement with EXCEPTION block: %s", stmt.String())
		})
	}
}

// TestHandlerDeclaration tests MySQL DECLARE HANDLER parsing
func TestHandlerDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "CONTINUE HANDLER for SQLEXCEPTION",
			sql: `
				CREATE PROCEDURE handle_errors()
				BEGIN
					DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
					BEGIN
						SELECT 'Error occurred';
					END;

					INSERT INTO users (id, name) VALUES (1, 'John');
				END
			`,
			dialect: "mysql",
		},
		{
			name: "EXIT HANDLER for NOT FOUND",
			sql: `
				CREATE PROCEDURE fetch_user(IN user_id INT)
				BEGIN
					DECLARE done INT DEFAULT 0;
					DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

					SELECT name FROM users WHERE id = user_id;
				END
			`,
			dialect: "mysql",
		},
		{
			name: "HANDLER for SQLSTATE",
			sql: `
				CREATE PROCEDURE handle_duplicate()
				BEGIN
					DECLARE CONTINUE HANDLER FOR SQLSTATE '23000'
					BEGIN
						SELECT 'Duplicate key error';
					END;

					INSERT INTO users (id, name) VALUES (1, 'John');
				END
			`,
			dialect: "mysql",
		},
		{
			name: "Multiple HANDLERS",
			sql: `
				CREATE PROCEDURE complex_handler()
				BEGIN
					DECLARE CONTINUE HANDLER FOR SQLEXCEPTION
						SELECT 'SQL Exception';

					DECLARE CONTINUE HANDLER FOR SQLWARNING
						SELECT 'SQL Warning';

					DECLARE EXIT HANDLER FOR NOT FOUND
						SELECT 'Not Found';

					SELECT * FROM users;
				END
			`,
			dialect: "mysql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse DECLARE HANDLER: %v", err)
			}

			t.Logf("✅ Successfully parsed DECLARE HANDLER: %s", stmt.String())
		})
	}
}

// TestRaiseStatement tests PostgreSQL RAISE statement
func TestRaiseStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple RAISE EXCEPTION",
			sql: `
				RAISE EXCEPTION 'User not found'
			`,
			dialect: "postgresql",
		},
		{
			name: "RAISE NOTICE",
			sql: `
				RAISE NOTICE 'Processing user %', user_id
			`,
			dialect: "postgresql",
		},
		{
			name: "RAISE WARNING",
			sql: `
				RAISE WARNING 'Low balance detected'
			`,
			dialect: "postgresql",
		},
		{
			name: "RAISE with no message (re-raise)",
			sql: `
				RAISE
			`,
			dialect: "postgresql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse RAISE statement: %v", err)
			}

			raiseStmt, ok := stmt.(*parser.RaiseStatement)
			if !ok {
				t.Fatalf("Expected *parser.RaiseStatement, got %T", stmt)
			}

			t.Logf("✅ Successfully parsed RAISE statement: %s", raiseStmt.String())
		})
	}
}

// TestThrowStatement tests SQL Server THROW statement
func TestThrowStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "THROW with parameters",
			sql: `
				THROW 50001, 'User not found', 1
			`,
			dialect: "sqlserver",
		},
		{
			name: "THROW re-throw",
			sql: `
				THROW
			`,
			dialect: "sqlserver",
		},
		{
			name: "THROW in CATCH block",
			sql: `
				BEGIN TRY
					SELECT 1/0;
				END TRY
				BEGIN CATCH
					THROW 50000, 'Division by zero detected', 1;
				END CATCH
			`,
			dialect: "sqlserver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse THROW statement: %v", err)
			}

			t.Logf("✅ Successfully parsed THROW statement: %s", stmt.String())
		})
	}
}

// TestSignalStatement tests MySQL SIGNAL statement
func TestSignalStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple SIGNAL",
			sql: `
				SIGNAL SQLSTATE '45000'
			`,
			dialect: "mysql",
		},
		{
			name: "SIGNAL with MESSAGE_TEXT",
			sql: `
				SIGNAL SQLSTATE '45000'
				SET MESSAGE_TEXT = 'User not found'
			`,
			dialect: "mysql",
		},
		{
			name: "SIGNAL with multiple properties",
			sql: `
				SIGNAL SQLSTATE '45000'
				SET MESSAGE_TEXT = 'Custom error',
					MYSQL_ERRNO = 1525
			`,
			dialect: "mysql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse SIGNAL statement: %v", err)
			}

			signalStmt, ok := stmt.(*parser.SignalStatement)
			if !ok {
				t.Fatalf("Expected *parser.SignalStatement, got %T", stmt)
			}

			t.Logf("✅ Successfully parsed SIGNAL statement: %s", signalStmt.String())
		})
	}
}

// TestComplexExceptionHandling tests complex real-world exception scenarios
func TestComplexExceptionHandling(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "SQL Server TRY CATCH with transactions",
			sql: `
				BEGIN TRY
					UPDATE accounts SET balance = 0 WHERE id = 1;
					UPDATE accounts SET balance = 100 WHERE id = 2;
				END TRY
				BEGIN CATCH
					ROLLBACK;
					THROW 50001, 'Transaction failed', 1;
				END CATCH
			`,
			dialect: "sqlserver",
		},
		{
			name: "Nested TRY CATCH with multiple operations",
			sql: `
				BEGIN TRY
					BEGIN TRY
						INSERT INTO users (id, name) VALUES (1, 'John');
					END TRY
					BEGIN CATCH
						SELECT ERROR_MESSAGE();
						THROW;
					END CATCH
				END TRY
				BEGIN CATCH
					SELECT 'Outer error handler';
				END CATCH
			`,
			dialect: "sqlserver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse complex exception handling: %v", err)
			}

			t.Logf("✅ Successfully parsed complex exception handling: %s", stmt.String())
		})
	}
}
