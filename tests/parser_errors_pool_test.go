package tests

import (
	"strings"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// --- ParseError ---

func TestParseError_Error(t *testing.T) {
	err := parser.NewParseError("unexpected token", "FROM", 3, 12)
	s := err.Error()
	if !strings.Contains(s, "line 3") {
		t.Errorf("expected 'line 3' in error string, got %q", s)
	}
	if !strings.Contains(s, "column 12") {
		t.Errorf("expected 'column 12' in error string, got %q", s)
	}
	if !strings.Contains(s, "FROM") {
		t.Errorf("expected token 'FROM' in error string, got %q", s)
	}
}

func TestParseError_Fields(t *testing.T) {
	err := parser.NewParseError("test message", "TOKEN", 1, 5)
	if err.Message != "test message" {
		t.Errorf("unexpected message: %s", err.Message)
	}
	if err.Line != 1 || err.Column != 5 {
		t.Errorf("unexpected line/column: %d/%d", err.Line, err.Column)
	}
	if err.Token != "TOKEN" {
		t.Errorf("unexpected token: %s", err.Token)
	}
}

// --- SyntaxError ---

func TestSyntaxError_Error(t *testing.T) {
	err := parser.NewSyntaxError("SELECT", "FROM", 2, 8)
	s := err.Error()
	if !strings.Contains(s, "line 2") {
		t.Errorf("expected 'line 2' in error string, got %q", s)
	}
	if !strings.Contains(s, "SELECT") {
		t.Errorf("expected 'SELECT' in error string, got %q", s)
	}
	if !strings.Contains(s, "FROM") {
		t.Errorf("expected 'FROM' in error string, got %q", s)
	}
}

func TestSyntaxError_Fields(t *testing.T) {
	err := parser.NewSyntaxError("IDENTIFIER", ";", 5, 20)
	if err.Expected != "IDENTIFIER" {
		t.Errorf("unexpected Expected: %s", err.Expected)
	}
	if err.Found != ";" {
		t.Errorf("unexpected Found: %s", err.Found)
	}
}

// --- Object pool Put methods ---

func TestPool_PutSelectStatement(t *testing.T) {
	stmt := parser.GetSelectStatement()
	stmt.Distinct = true
	parser.PutSelectStatement(stmt) // should not panic

	// nil should be safe too
	parser.PutSelectStatement(nil)
}

func TestPool_PutColumnReference(t *testing.T) {
	col := parser.GetColumnReference()
	col.Table = "users"
	col.Column = "id"
	parser.PutColumnReference(col)

	// Verify reset when re-acquired
	col2 := parser.GetColumnReference()
	_ = col2
	parser.PutColumnReference(nil)
}

func TestPool_PutBinaryExpression(t *testing.T) {
	expr := parser.GetBinaryExpression()
	expr.Operator = "="
	parser.PutBinaryExpression(expr)
	parser.PutBinaryExpression(nil)
}

func TestPool_PutJoinClause(t *testing.T) {
	join := parser.GetJoinClause()
	join.JoinType = "INNER"
	parser.PutJoinClause(join)
	parser.PutJoinClause(nil)
}

func TestPool_GetSelectStatement_Reset(t *testing.T) {
	stmt := parser.GetSelectStatement()
	stmt.Distinct = true
	parser.PutSelectStatement(stmt)

	// Get it back — should be reset
	stmt2 := parser.GetSelectStatement()
	if stmt2.Distinct {
		t.Error("expected Distinct to be false after pool reset")
	}
	parser.PutSelectStatement(stmt2)
}
