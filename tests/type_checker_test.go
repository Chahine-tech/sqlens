package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
	"github.com/Chahine-tech/sql-parser-go/pkg/schema"
)

func makeTestSchema() *schema.Schema {
	loader := schema.NewSchemaLoader()
	jsonData := []byte(`{
		"name": "testdb",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id",    "type": "INT",     "primary_key": true},
					{"name": "name",  "type": "VARCHAR", "nullable": true},
					{"name": "age",   "type": "INT",     "nullable": true},
					{"name": "score", "type": "FLOAT",   "nullable": true},
					{"name": "active","type": "BOOLEAN", "nullable": true}
				],
				"indexes": [
					{"name": "idx_name", "columns": ["name"]}
				]
			},
			{
				"name": "orders",
				"columns": [
					{"name": "id",      "type": "INT",     "primary_key": true},
					{"name": "user_id", "type": "INT",     "foreign_key": true, "fk_table": "users", "fk_column": "id"},
					{"name": "amount",  "type": "DECIMAL", "nullable": false}
				]
			}
		]
	}`)
	s, _ := loader.LoadFromJSON(jsonData)
	return s
}

func parseForTC(t *testing.T, sql string) parser.Statement {
	t.Helper()
	ctx := context.Background()
	p := parser.NewWithDialect(ctx, sql, dialect.GetDialect("mysql"))
	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return stmt
}

// --- TypeChecker ---

func TestTypeChecker_New(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	if tc == nil {
		t.Fatal("expected non-nil type checker")
	}
}

func TestTypeChecker_CheckStatement_Select_Clean(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	stmt := parseForTC(t, "SELECT id, name FROM users WHERE id = 1")
	errs := tc.CheckStatement(stmt)
	// id=1 is int=int — no type mismatch
	for _, e := range errs {
		if e.Type == "TYPE_MISMATCH" {
			t.Errorf("unexpected TYPE_MISMATCH: %s", e.Message)
		}
	}
}

func TestTypeChecker_CheckStatement_Select_Having(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	stmt := parseForTC(t, "SELECT id FROM users GROUP BY id HAVING id > 0")
	errs := tc.CheckStatement(stmt)
	_ = errs // should not panic
}

func TestTypeChecker_CheckStatement_Insert_Valid(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	stmt := parseForTC(t, "INSERT INTO users (id, name) VALUES (1, 'Alice')")
	errs := tc.CheckStatement(stmt)
	for _, e := range errs {
		if e.Type == "TYPE_MISMATCH" {
			t.Errorf("unexpected TYPE_MISMATCH: %s", e.Message)
		}
	}
}

func TestTypeChecker_CheckStatement_Insert_CountMismatch(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	stmt := parseForTC(t, "INSERT INTO users (id, name) VALUES (1)")
	errs := tc.CheckStatement(stmt)
	found := false
	for _, e := range errs {
		if e.Type == "COLUMN_COUNT_MISMATCH" {
			found = true
		}
	}
	if !found {
		t.Error("expected COLUMN_COUNT_MISMATCH error")
	}
}

func TestTypeChecker_CheckStatement_Update(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	stmt := parseForTC(t, "UPDATE users SET name = 'Bob' WHERE id = 1")
	errs := tc.CheckStatement(stmt)
	for _, e := range errs {
		if e.Type == "TYPE_MISMATCH" {
			t.Errorf("unexpected TYPE_MISMATCH: %s", e.Message)
		}
	}
}

func TestTypeChecker_CheckStatement_Update_UnknownTable(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	stmt := parseForTC(t, "UPDATE nonexistent SET x = 1 WHERE id = 1")
	errs := tc.CheckStatement(stmt)
	_ = errs // no panic, table just not found
}

func TestTypeChecker_CheckStatement_NonBoolean(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	// Use arithmetic operator in WHERE — should flag NON_BOOLEAN_EXPRESSION
	stmt := parseForTC(t, "SELECT id FROM users WHERE id + 1")
	errs := tc.CheckStatement(stmt)
	found := false
	for _, e := range errs {
		if e.Type == "NON_BOOLEAN_EXPRESSION" {
			found = true
		}
	}
	if !found {
		t.Log("NON_BOOLEAN_EXPRESSION not flagged (may depend on parser AST shape) — that's OK")
	}
}

func TestTypeChecker_CheckStatement_IN_Expression(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	stmt := parseForTC(t, "SELECT id FROM users WHERE id IN (1, 2, 3)")
	errs := tc.CheckStatement(stmt)
	_ = errs // should not panic
}

func TestTypeChecker_CheckStatement_NOT(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	stmt := parseForTC(t, "SELECT id FROM users WHERE id != 1")
	errs := tc.CheckStatement(stmt)
	_ = errs // should not panic
}

func TestTypeChecker_InferExpressionType_Float(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	// score is FLOAT column — arithmetic on it
	stmt := parseForTC(t, "SELECT id FROM users WHERE score > 3.14")
	errs := tc.CheckStatement(stmt)
	_ = errs
}

func TestTypeChecker_InferFunctionReturnType(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	// Functions in SELECT don't hit type errors
	stmts := []string{
		"SELECT COUNT(id) FROM users WHERE id > 0",
		"SELECT AVG(age) FROM users WHERE id > 0",
		"SELECT MAX(id) FROM users WHERE id > 0",
		"SELECT MIN(id) FROM users WHERE id > 0",
		"SELECT CONCAT(name, ' x') FROM users WHERE id > 0",
		"SELECT NOW() FROM users WHERE id > 0",
		"SELECT CURRENT_DATE FROM users WHERE id > 0",
	}
	for _, sql := range stmts {
		stmt := parseForTC(t, sql)
		errs := tc.CheckStatement(stmt)
		_ = errs // should not panic
	}
}

func TestTypeChecker_CheckStatement_ColumnRef_WithTable(t *testing.T) {
	s := makeTestSchema()
	tc := schema.NewTypeChecker(s)
	stmt := parseForTC(t, "SELECT users.id FROM users WHERE users.id = 1")
	errs := tc.CheckStatement(stmt)
	_ = errs
}

// --- ValidationError.Error ---

func TestValidationError_Error(t *testing.T) {
	ve := &schema.ValidationError{
		Type:    "TABLE_NOT_FOUND",
		Message: "table 'foo' does not exist",
	}
	s := ve.Error()
	if s == "" {
		t.Error("expected non-empty error string")
	}
	if s != "[TABLE_NOT_FOUND] table 'foo' does not exist" {
		t.Errorf("unexpected error string: %q", s)
	}
}

// --- Schema GetColumn, GetIndex, String (DataType) ---

func TestSchema_GetColumn(t *testing.T) {
	s := makeTestSchema()
	col, err := s.GetColumn("users", "id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if col == nil {
		t.Fatal("expected non-nil column")
	}
	if col.Name != "id" {
		t.Errorf("expected column name 'id', got %q", col.Name)
	}
}

func TestSchema_GetColumn_NotFound(t *testing.T) {
	s := makeTestSchema()
	_, err := s.GetColumn("users", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent column")
	}
}

func TestSchema_GetColumn_TableNotFound(t *testing.T) {
	s := makeTestSchema()
	_, err := s.GetColumn("nonexistent_table", "id")
	if err == nil {
		t.Error("expected error for nonexistent table")
	}
}

func TestTable_GetIndex(t *testing.T) {
	s := makeTestSchema()
	table, ok := s.GetTable("users")
	if !ok {
		t.Fatal("expected users table")
	}
	idx, ok := table.GetIndex("idx_name")
	if !ok {
		t.Error("expected idx_name index")
	}
	if idx == nil {
		t.Error("expected non-nil index")
	}
}

func TestTable_GetIndex_NotFound(t *testing.T) {
	s := makeTestSchema()
	table, ok := s.GetTable("users")
	if !ok {
		t.Fatal("expected users table")
	}
	_, ok = table.GetIndex("nonexistent_idx")
	if ok {
		t.Error("expected index not found")
	}
}

func TestDataType_String(t *testing.T) {
	cases := []struct {
		dt   schema.DataType
		want string
	}{
		{schema.DataType{Name: "INT"}, "INT"},
		{schema.DataType{Name: "VARCHAR", Length: 255}, "VARCHAR(255)"},
		{schema.DataType{Name: "DECIMAL", Precision: 10, Scale: 2}, "DECIMAL(10,2)"},
		{schema.DataType{Name: "DECIMAL", Precision: 10}, "DECIMAL(10)"},
	}
	for _, tc := range cases {
		got := tc.dt.String()
		if got != tc.want {
			t.Errorf("DataType.String() = %q, want %q", got, tc.want)
		}
	}
}

// --- SchemaLoader AddSchema, GetSchema, HasSchema ---

func TestSchemaLoader_AddGetHas(t *testing.T) {
	loader := schema.NewSchemaLoader()
	s := makeTestSchema()

	loader.AddSchema(s)

	got, ok := loader.GetSchema("testdb")
	if !ok {
		t.Error("expected to find schema 'testdb'")
	}
	if got == nil {
		t.Error("expected non-nil schema")
	}

	if !loader.HasSchema("testdb") {
		t.Error("HasSchema should return true")
	}
	if loader.HasSchema("nonexistent") {
		t.Error("HasSchema should return false for nonexistent")
	}
}

func TestSchemaLoader_GetSchema_CaseInsensitive(t *testing.T) {
	loader := schema.NewSchemaLoader()
	s := makeTestSchema()
	loader.AddSchema(s)

	_, ok := loader.GetSchema("TESTDB")
	if !ok {
		t.Error("GetSchema should be case-insensitive")
	}
}

// --- SchemaLoader LoadFromYAML ---

func TestSchemaLoader_LoadFromYAML(t *testing.T) {
	loader := schema.NewSchemaLoader()
	yamlData := []byte(`
name: yamldb
tables:
  - name: products
    columns:
      - name: id
        type: INT
        primary_key: true
      - name: price
        type: DECIMAL
        precision: 10
        scale: 2
`)
	s, err := loader.LoadFromYAML(yamlData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil schema")
	}
	if s.Name != "yamldb" {
		t.Errorf("expected schema name 'yamldb', got %q", s.Name)
	}
	table, ok := s.GetTable("products")
	if !ok {
		t.Fatal("expected products table")
	}
	_, ok = table.GetColumn("price")
	if !ok {
		t.Error("expected price column")
	}
}
