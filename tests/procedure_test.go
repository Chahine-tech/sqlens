package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func TestSimpleProcedure(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "MySQL simple procedure",
			sql: `CREATE PROCEDURE get_user(IN user_id INT)
BEGIN
	RETURN;
END`,
			dialect: "mysql",
		},
		{
			name: "PostgreSQL simple procedure",
			sql: `CREATE PROCEDURE update_balance(user_id INT, amount DECIMAL)
LANGUAGE plpgsql
AS
BEGIN
	RETURN;
END`,
			dialect: "postgresql",
		},
		{
			name: "SQL Server simple procedure",
			sql: `CREATE PROCEDURE GetCustomers
AS
BEGIN
	RETURN;
END`,
			dialect: "sqlserver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dialect.GetDialect(tt.dialect)
			p := parser.NewWithDialect(context.Background(), tt.sql, d)
			stmt, err := p.ParseStatement()

			if err != nil {
				t.Fatalf("Failed to parse procedure: %v", err)
			}

			procStmt, ok := stmt.(*parser.CreateProcedureStatement)
			if !ok {
				t.Fatalf("Expected CreateProcedureStatement, got %T", stmt)
			}

			if procStmt.Name == "" {
				t.Error("Procedure name is empty")
			}

			t.Logf("✅ Successfully parsed procedure: %s", procStmt.String())
		})
	}
}

func TestProcedureWithParameters(t *testing.T) {
	sql := `CREATE PROCEDURE calculate_total(
		IN product_id INT,
		IN quantity INT,
		OUT total_price DECIMAL(10,2)
	)
BEGIN
	RETURN;
END`

	d := dialect.GetDialect("mysql")
	p := parser.NewWithDialect(context.Background(), sql, d)
	stmt, err := p.ParseStatement()

	if err != nil {
		t.Fatalf("Failed to parse procedure with parameters: %v", err)
	}

	procStmt, ok := stmt.(*parser.CreateProcedureStatement)
	if !ok {
		t.Fatalf("Expected CreateProcedureStatement, got %T", stmt)
	}

	if len(procStmt.Parameters) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(procStmt.Parameters))
	}

	// Check parameter modes
	expectedModes := []string{"IN", "IN", "OUT"}
	for i, param := range procStmt.Parameters {
		if param.Mode != expectedModes[i] {
			t.Errorf("Parameter %d: expected mode %s, got %s", i, expectedModes[i], param.Mode)
		}
	}

	t.Log("✅ Successfully parsed procedure with IN/OUT parameters")
}

func TestSimpleFunction(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "MySQL function",
			sql: `CREATE FUNCTION calculate_tax(amount DECIMAL(10,2))
RETURNS DECIMAL(10,2)
DETERMINISTIC
BEGIN
	RETURN amount * 0.20;
END`,
			dialect: "mysql",
		},
		// TODO: PostgreSQL uses $$ delimiters which requires parser changes
		// {
		// 	name: "PostgreSQL function",
		// 	sql: `CREATE FUNCTION add_numbers(a INT, b INT)
		// RETURNS INT
		// LANGUAGE plpgsql
		// AS $$
		// BEGIN
		// 	RETURN a + b;
		// END;
		// $$`,
		// 	dialect: "postgresql",
		// },
		// TODO: SQL Server uses @ for parameters which requires lexer changes
		// {
		// 	name: "SQL Server function",
		// 	sql: `CREATE FUNCTION GetFullName(@first VARCHAR(50), @last VARCHAR(50))
		// RETURNS VARCHAR(100)
		// AS
		// BEGIN
		// 	RETURN @first + ' ' + @last;
		// END`,
		// 	dialect: "sqlserver",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dialect.GetDialect(tt.dialect)
			p := parser.NewWithDialect(context.Background(), tt.sql, d)
			stmt, err := p.ParseStatement()

			if err != nil {
				t.Fatalf("Failed to parse function: %v", err)
			}

			funcStmt, ok := stmt.(*parser.CreateFunctionStatement)
			if !ok {
				t.Fatalf("Expected CreateFunctionStatement, got %T", stmt)
			}

			if funcStmt.Name == "" {
				t.Error("Function name is empty")
			}

			if funcStmt.ReturnType == nil {
				t.Error("Function must have a return type")
			}

			t.Logf("✅ Successfully parsed function: %s", funcStmt.String())
		})
	}
}

func TestFunctionWithVariousReturnTypes(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		returnType string
	}{
		{
			name: "INT return type",
			sql: `CREATE FUNCTION get_count()
RETURNS INT
BEGIN
	RETURN 42;
END`,
			returnType: "INT",
		},
		{
			name: "VARCHAR return type",
			sql: `CREATE FUNCTION get_name()
RETURNS VARCHAR(100)
BEGIN
	RETURN 'John Doe';
END`,
			returnType: "VARCHAR",
		},
		{
			name: "DECIMAL return type",
			sql: `CREATE FUNCTION get_price()
RETURNS DECIMAL(10,2)
BEGIN
	RETURN 99.99;
END`,
			returnType: "DECIMAL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dialect.GetDialect("mysql")
			p := parser.NewWithDialect(context.Background(), tt.sql, d)
			stmt, err := p.ParseStatement()

			if err != nil {
				t.Fatalf("Failed to parse function: %v", err)
			}

			funcStmt, ok := stmt.(*parser.CreateFunctionStatement)
			if !ok {
				t.Fatalf("Expected CreateFunctionStatement, got %T", stmt)
			}

			if funcStmt.ReturnType.Name != tt.returnType {
				t.Errorf("Expected return type %s, got %s", tt.returnType, funcStmt.ReturnType.Name)
			}

			t.Logf("✅ Function has correct return type: %s", funcStmt.ReturnType.String())
		})
	}
}

func TestCreateOrReplaceProcedure(t *testing.T) {
	sql := `CREATE OR REPLACE PROCEDURE update_stats()
LANGUAGE plpgsql
AS
BEGIN
	RETURN;
END`

	d := dialect.GetDialect("postgresql")
	p := parser.NewWithDialect(context.Background(), sql, d)
	stmt, err := p.ParseStatement()

	if err != nil {
		t.Fatalf("Failed to parse CREATE OR REPLACE PROCEDURE: %v", err)
	}

	procStmt, ok := stmt.(*parser.CreateProcedureStatement)
	if !ok {
		t.Fatalf("Expected CreateProcedureStatement, got %T", stmt)
	}

	if !procStmt.OrReplace {
		t.Error("Expected OrReplace to be true")
	}

	if procStmt.Language != "plpgsql" {
		t.Errorf("Expected language plpgsql, got %s", procStmt.Language)
	}

	t.Log("✅ Successfully parsed CREATE OR REPLACE PROCEDURE")
}

func TestCreateOrReplaceFunction(t *testing.T) {
	sql := `CREATE OR REPLACE FUNCTION square(x INT)
RETURNS INT
LANGUAGE plpgsql
AS
BEGIN
	RETURN x;
END`

	d := dialect.GetDialect("postgresql")
	p := parser.NewWithDialect(context.Background(), sql, d)
	stmt, err := p.ParseStatement()

	if err != nil {
		t.Fatalf("Failed to parse CREATE OR REPLACE FUNCTION: %v", err)
	}

	funcStmt, ok := stmt.(*parser.CreateFunctionStatement)
	if !ok {
		t.Fatalf("Expected CreateFunctionStatement, got %T", stmt)
	}

	if !funcStmt.OrReplace {
		t.Error("Expected OrReplace to be true")
	}

	t.Log("✅ Successfully parsed CREATE OR REPLACE FUNCTION")
}

func TestProcedureWithVariableDeclarations(t *testing.T) {
	sql := `CREATE PROCEDURE process_order(order_id INT)
BEGIN
	DECLARE total DECIMAL(10,2);
	DECLARE status VARCHAR(20);
	RETURN;
END`

	d := dialect.GetDialect("mysql")
	p := parser.NewWithDialect(context.Background(), sql, d)
	stmt, err := p.ParseStatement()

	if err != nil {
		t.Fatalf("Failed to parse procedure with variables: %v", err)
	}

	procStmt, ok := stmt.(*parser.CreateProcedureStatement)
	if !ok {
		t.Fatalf("Expected CreateProcedureStatement, got %T", stmt)
	}

	if procStmt.Body == nil {
		t.Fatal("Procedure body is nil")
	}

	if len(procStmt.Body.Variables) != 2 {
		t.Errorf("Expected 2 variable declarations, got %d", len(procStmt.Body.Variables))
	}

	t.Logf("✅ Successfully parsed procedure with %d variable declarations", len(procStmt.Body.Variables))
}

func TestProcedureWithCursor(t *testing.T) {
	sql := `CREATE PROCEDURE list_users()
BEGIN
	DECLARE user_cursor CURSOR FOR SELECT id, name FROM users;
	RETURN;
END`

	d := dialect.GetDialect("mysql")
	p := parser.NewWithDialect(context.Background(), sql, d)
	stmt, err := p.ParseStatement()

	if err != nil {
		t.Fatalf("Failed to parse procedure with cursor: %v", err)
	}

	procStmt, ok := stmt.(*parser.CreateProcedureStatement)
	if !ok {
		t.Fatalf("Expected CreateProcedureStatement, got %T", stmt)
	}

	if procStmt.Body == nil {
		t.Fatal("Procedure body is nil")
	}

	if len(procStmt.Body.Cursors) != 1 {
		t.Errorf("Expected 1 cursor declaration, got %d", len(procStmt.Body.Cursors))
	}

	if len(procStmt.Body.Statements) < 1 {
		t.Errorf("Expected at least 1 statement (RETURN), got %d", len(procStmt.Body.Statements))
	}

	t.Log("✅ Successfully parsed procedure with cursor")
}

func TestFunctionDeterministic(t *testing.T) {
	sql := `CREATE FUNCTION double_value(x INT)
RETURNS INT
DETERMINISTIC
BEGIN
	RETURN x * 2;
END`

	d := dialect.GetDialect("mysql")
	p := parser.NewWithDialect(context.Background(), sql, d)
	stmt, err := p.ParseStatement()

	if err != nil {
		t.Fatalf("Failed to parse DETERMINISTIC function: %v", err)
	}

	funcStmt, ok := stmt.(*parser.CreateFunctionStatement)
	if !ok {
		t.Fatalf("Expected CreateFunctionStatement, got %T", stmt)
	}

	if !funcStmt.Deterministic {
		t.Error("Expected Deterministic to be true")
	}

	t.Log("✅ Successfully parsed DETERMINISTIC function")
}
