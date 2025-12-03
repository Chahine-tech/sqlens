package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func TestFetchDirections(t *testing.T) {
	tests := []struct {
		name      string
		fetchSQL  string
		dialect   string
		direction string
	}{
		{
			name:      "FETCH NEXT",
			fetchSQL:  `FETCH NEXT FROM my_cursor`,
			dialect:   "postgresql",
			direction: "NEXT",
		},
		{
			name:      "FETCH PRIOR",
			fetchSQL:  `FETCH PRIOR FROM my_cursor`,
			dialect:   "postgresql",
			direction: "PRIOR",
		},
		{
			name:      "FETCH FIRST",
			fetchSQL:  `FETCH FIRST FROM my_cursor`,
			dialect:   "postgresql",
			direction: "FIRST",
		},
		{
			name:      "FETCH LAST",
			fetchSQL:  `FETCH LAST FROM my_cursor`,
			dialect:   "postgresql",
			direction: "LAST",
		},
		{
			name:      "FETCH ABSOLUTE",
			fetchSQL:  `FETCH ABSOLUTE 10 FROM my_cursor`,
			dialect:   "postgresql",
			direction: "ABSOLUTE",
		},
		{
			name:      "FETCH RELATIVE",
			fetchSQL:  `FETCH RELATIVE 5 FROM my_cursor`,
			dialect:   "postgresql",
			direction: "RELATIVE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wrap FETCH in a procedure
			sql := `CREATE PROCEDURE test_proc() AS BEGIN ` + tt.fetchSQL + `; END`
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, sql, dialect.GetDialect(tt.dialect))

			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			procStmt, ok := stmt.(*parser.CreateProcedureStatement)
			if !ok {
				t.Fatalf("Expected CreateProcedureStatement, got %T", stmt)
			}

			if len(procStmt.Body.Statements) == 0 {
				t.Fatal("Procedure body has no statements")
			}

			fetchStmt, ok := procStmt.Body.Statements[0].(*parser.FetchStatement)
			if !ok {
				t.Fatalf("Expected FetchStatement, got %T", procStmt.Body.Statements[0])
			}

			if fetchStmt.Direction != tt.direction {
				t.Fatalf("Expected direction %s, got %s", tt.direction, fetchStmt.Direction)
			}
		})
	}
}

func TestFetchWithInto(t *testing.T) {
	tests := []struct {
		name     string
		fetchSQL string
		dialect  string
		varCount int
	}{
		{
			name:     "FETCH with single variable",
			fetchSQL: `FETCH my_cursor INTO v_name`,
			dialect:  "mysql",
			varCount: 1,
		},
		{
			name:     "FETCH with multiple variables",
			fetchSQL: `FETCH my_cursor INTO v_id, v_name, v_email`,
			dialect:  "mysql",
			varCount: 3,
		},
		{
			name:     "FETCH NEXT with INTO",
			fetchSQL: `FETCH NEXT FROM my_cursor INTO v_data`,
			dialect:  "postgresql",
			varCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := `CREATE PROCEDURE test_proc() AS BEGIN ` + tt.fetchSQL + `; END`
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, sql, dialect.GetDialect(tt.dialect))

			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			procStmt, ok := stmt.(*parser.CreateProcedureStatement)
			if !ok {
				t.Fatalf("Expected CreateProcedureStatement, got %T", stmt)
			}

			fetchStmt, ok := procStmt.Body.Statements[0].(*parser.FetchStatement)
			if !ok {
				t.Fatalf("Expected FetchStatement, got %T", procStmt.Body.Statements[0])
			}

			if len(fetchStmt.Variables) != tt.varCount {
				t.Fatalf("Expected %d variables, got %d", tt.varCount, len(fetchStmt.Variables))
			}
		})
	}
}

func TestDeallocateStatement(t *testing.T) {
	tests := []struct {
		name        string
		deallocSQL  string
		dialect     string
		cursorName  string
	}{
		{
			name:        "Simple DEALLOCATE",
			deallocSQL:  `DEALLOCATE my_cursor`,
			dialect:     "postgresql",
			cursorName:  "my_cursor",
		},
		{
			name:        "DEALLOCATE PREPARE (MySQL)",
			deallocSQL:  `DEALLOCATE PREPARE stmt1`,
			dialect:     "mysql",
			cursorName:  "stmt1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := `CREATE PROCEDURE test_proc() AS BEGIN ` + tt.deallocSQL + `; END`
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, sql, dialect.GetDialect(tt.dialect))

			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			procStmt, ok := stmt.(*parser.CreateProcedureStatement)
			if !ok {
				t.Fatalf("Expected CreateProcedureStatement, got %T", stmt)
			}

			deallocStmt, ok := procStmt.Body.Statements[0].(*parser.DeallocateStatement)
			if !ok {
				t.Fatalf("Expected DeallocateStatement, got %T", procStmt.Body.Statements[0])
			}

			if deallocStmt.CursorName != tt.cursorName {
				t.Fatalf("Expected cursor name %s, got %s", tt.cursorName, deallocStmt.CursorName)
			}
		})
	}
}

func TestCursorLifecycle(t *testing.T) {
	// Test a complete cursor lifecycle: DECLARE, OPEN, FETCH, CLOSE, DEALLOCATE
	tests := []struct {
		name       string
		cursorSQL  string
		dialect    string
		stmtType   string
	}{
		{
			name:       "OPEN cursor",
			cursorSQL:  `OPEN my_cursor`,
			dialect:    "mysql",
			stmtType:   "OpenCursorStatement",
		},
		{
			name:       "FETCH simple",
			cursorSQL:  `FETCH my_cursor`,
			dialect:    "mysql",
			stmtType:   "FetchStatement",
		},
		{
			name:       "CLOSE cursor",
			cursorSQL:  `CLOSE my_cursor`,
			dialect:    "postgresql",
			stmtType:   "CloseStatement",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := `CREATE PROCEDURE test_proc() AS BEGIN ` + tt.cursorSQL + `; END`
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, sql, dialect.GetDialect(tt.dialect))

			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			procStmt, ok := stmt.(*parser.CreateProcedureStatement)
			if !ok {
				t.Fatalf("Expected CreateProcedureStatement, got %T", stmt)
			}

			if len(procStmt.Body.Statements) == 0 {
				t.Fatal("Expected procedure body to have statements")
			}

			innerStmt := procStmt.Body.Statements[0]
			actualType := innerStmt.Type()
			if actualType != tt.stmtType {
				t.Fatalf("Expected %s, got %s", tt.stmtType, actualType)
			}
		})
	}
}
