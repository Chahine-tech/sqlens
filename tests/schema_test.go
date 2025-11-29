package tests

import (
	"context"
	"os"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
	"github.com/Chahine-tech/sql-parser-go/pkg/schema"
)

// Test schema loading from JSON
func TestSchemaLoadFromJSON(t *testing.T) {
	loader := schema.NewSchemaLoader()

	jsonData := []byte(`{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{
						"name": "id",
						"type": "INT",
						"primary_key": true
					},
					{
						"name": "name",
						"type": "VARCHAR",
						"length": 100
					}
				]
			}
		]
	}`)

	s, err := loader.LoadFromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to load schema from JSON: %v", err)
	}

	if s.Name != "test_db" {
		t.Errorf("Expected schema name 'test_db', got '%s'", s.Name)
	}

	if !s.HasTable("users") {
		t.Error("Expected table 'users' to exist")
	}

	table, _ := s.GetTable("users")
	if !table.HasColumn("id") {
		t.Error("Expected column 'id' in table 'users'")
	}

	if !table.HasColumn("name") {
		t.Error("Expected column 'name' in table 'users'")
	}
}

// Test schema loading from file
func TestSchemaLoadFromFile(t *testing.T) {
	loader := schema.NewSchemaLoader()

	schemaPath := "../examples/schemas/test_schema.json"
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Skip("Schema file not found, skipping test")
	}

	s, err := loader.LoadFromFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to load schema from file: %v", err)
	}

	if s.Name != "test_db" {
		t.Errorf("Expected schema name 'test_db', got '%s'", s.Name)
	}

	// Check all tables
	expectedTables := []string{"users", "orders", "products"}
	for _, tableName := range expectedTables {
		if !s.HasTable(tableName) {
			t.Errorf("Expected table '%s' to exist", tableName)
		}
	}

	// Check users table columns
	usersTable, _ := s.GetTable("users")
	expectedColumns := []string{"id", "name", "email", "age", "created_at"}
	for _, colName := range expectedColumns {
		if !usersTable.HasColumn(colName) {
			t.Errorf("Expected column '%s' in table 'users'", colName)
		}
	}

	// Check foreign key in orders table
	ordersTable, _ := s.GetTable("orders")
	userIdCol, _ := ordersTable.GetColumn("user_id")
	if !userIdCol.IsForeignKey {
		t.Error("Expected user_id to be a foreign key")
	}
	if userIdCol.ForeignKey.Table != "users" {
		t.Errorf("Expected foreign key to reference 'users', got '%s'", userIdCol.ForeignKey.Table)
	}
}

// Test table validation
func TestTableValidation(t *testing.T) {
	loader := schema.NewSchemaLoader()

	s, err := loader.LoadFromJSON([]byte(`{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{
						"name": "id",
						"type": "INT",
						"primary_key": true
					}
				]
			},
			{
				"name": "orders",
				"columns": [
					{
						"name": "id",
						"type": "INT",
						"primary_key": true
					},
					{
						"name": "user_id",
						"type": "INT",
						"foreign_key": true,
						"fk_table": "users",
						"fk_column": "id"
					}
				]
			}
		]
	}`))
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	// Schema should be valid
	if err := s.Validate(); err != nil {
		t.Errorf("Schema validation failed: %v", err)
	}
}

// Test invalid foreign key
func TestInvalidForeignKey(t *testing.T) {
	loader := schema.NewSchemaLoader()

	s, err := loader.LoadFromJSON([]byte(`{
		"name": "test_db",
		"tables": [
			{
				"name": "orders",
				"columns": [
					{
						"name": "user_id",
						"type": "INT",
						"foreign_key": true,
						"fk_table": "users",
						"fk_column": "id"
					}
				]
			}
		]
	}`))
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	// Schema should be invalid (users table doesn't exist)
	if err := s.Validate(); err == nil {
		t.Error("Expected validation error for missing foreign key table")
	}
}

// Test SELECT statement validation
func TestValidateSelectStatement(t *testing.T) {
	loader := schema.NewSchemaLoader()

	jsonData := `{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "name", "type": "VARCHAR", "length": 100},
					{"name": "email", "type": "VARCHAR", "length": 255}
				]
			}
		]
	}`

	s, err := loader.LoadFromJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	validator := schema.NewValidator(s)

	tests := []struct {
		name      string
		sql       string
		wantErr   bool
		errorType string
	}{
		{
			name:    "Valid SELECT",
			sql:     `SELECT id, name FROM users`,
			wantErr: false,
		},
		{
			name:      "Invalid table",
			sql:       `SELECT id FROM orders`,
			wantErr:   true,
			errorType: "TABLE_NOT_FOUND",
		},
		{
			name:      "Invalid column",
			sql:       `SELECT id, age FROM users`,
			wantErr:   true,
			errorType: "COLUMN_NOT_FOUND",
		},
		{
			name:    "Valid WHERE clause",
			sql:     `SELECT name FROM users WHERE id = 1`,
			wantErr: false,
		},
		{
			name:      "Invalid column in WHERE",
			sql:       `SELECT name FROM users WHERE age > 18`,
			wantErr:   true,
			errorType: "COLUMN_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			errors := validator.ValidateStatement(stmt)

			if tt.wantErr && len(errors) == 0 {
				t.Error("Expected validation errors, got none")
			}

			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("Expected no validation errors, got: %v", errors)
			}

			if tt.wantErr && len(errors) > 0 {
				if errors[0].Type != tt.errorType {
					t.Errorf("Expected error type '%s', got '%s'", tt.errorType, errors[0].Type)
				}
			}
		})
	}
}

// Test INSERT statement validation
func TestValidateInsertStatement(t *testing.T) {
	loader := schema.NewSchemaLoader()

	jsonData := `{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "name", "type": "VARCHAR", "length": 100},
					{"name": "email", "type": "VARCHAR", "length": 255}
				]
			}
		]
	}`

	s, err := loader.LoadFromJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	validator := schema.NewValidator(s)

	tests := []struct {
		name      string
		sql       string
		wantErr   bool
		errorType string
	}{
		{
			name:    "Valid INSERT",
			sql:     `INSERT INTO users (id, name, email) VALUES (1, 'John', 'john@example.com')`,
			wantErr: false,
		},
		{
			name:      "Invalid table",
			sql:       `INSERT INTO orders (id) VALUES (1)`,
			wantErr:   true,
			errorType: "TABLE_NOT_FOUND",
		},
		{
			name:      "Invalid column",
			sql:       `INSERT INTO users (id, age) VALUES (1, 25)`,
			wantErr:   true,
			errorType: "COLUMN_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			errors := validator.ValidateStatement(stmt)

			if tt.wantErr && len(errors) == 0 {
				t.Error("Expected validation errors, got none")
			}

			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("Expected no validation errors, got: %v", errors)
			}

			if tt.wantErr && len(errors) > 0 {
				if errors[0].Type != tt.errorType {
					t.Errorf("Expected error type '%s', got '%s'", tt.errorType, errors[0].Type)
				}
			}
		})
	}
}

// Test UPDATE statement validation
func TestValidateUpdateStatement(t *testing.T) {
	loader := schema.NewSchemaLoader()

	jsonData := `{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "name", "type": "VARCHAR", "length": 100},
					{"name": "email", "type": "VARCHAR", "length": 255}
				]
			}
		]
	}`

	s, err := loader.LoadFromJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	validator := schema.NewValidator(s)

	tests := []struct {
		name      string
		sql       string
		wantErr   bool
		errorType string
	}{
		{
			name:    "Valid UPDATE",
			sql:     `UPDATE users SET name = 'Jane' WHERE id = 1`,
			wantErr: false,
		},
		{
			name:      "Invalid table",
			sql:       `UPDATE orders SET status = 'shipped'`,
			wantErr:   true,
			errorType: "TABLE_NOT_FOUND",
		},
		{
			name:      "Invalid column in SET",
			sql:       `UPDATE users SET age = 25 WHERE id = 1`,
			wantErr:   true,
			errorType: "COLUMN_NOT_FOUND",
		},
		{
			name:      "Invalid column in WHERE",
			sql:       `UPDATE users SET name = 'Jane' WHERE age > 18`,
			wantErr:   true,
			errorType: "COLUMN_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			errors := validator.ValidateStatement(stmt)

			if tt.wantErr && len(errors) == 0 {
				t.Error("Expected validation errors, got none")
			}

			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("Expected no validation errors, got: %v", errors)
			}

			if tt.wantErr && len(errors) > 0 {
				if errors[0].Type != tt.errorType {
					t.Errorf("Expected error type '%s', got '%s'", tt.errorType, errors[0].Type)
				}
			}
		})
	}
}

// Test DELETE statement validation
func TestValidateDeleteStatement(t *testing.T) {
	loader := schema.NewSchemaLoader()

	jsonData := `{
		"name": "test_db",
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "INT"},
					{"name": "name", "type": "VARCHAR", "length": 100}
				]
			}
		]
	}`

	s, err := loader.LoadFromJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	validator := schema.NewValidator(s)

	tests := []struct {
		name      string
		sql       string
		wantErr   bool
		errorType string
	}{
		{
			name:    "Valid DELETE",
			sql:     `DELETE FROM users WHERE id = 1`,
			wantErr: false,
		},
		{
			name:      "Invalid table",
			sql:       `DELETE FROM orders WHERE id = 1`,
			wantErr:   true,
			errorType: "TABLE_NOT_FOUND",
		},
		{
			name:      "Invalid column in WHERE",
			sql:       `DELETE FROM users WHERE age > 18`,
			wantErr:   true,
			errorType: "COLUMN_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect("mysql"))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			errors := validator.ValidateStatement(stmt)

			if tt.wantErr && len(errors) == 0 {
				t.Error("Expected validation errors, got none")
			}

			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("Expected no validation errors, got: %v", errors)
			}

			if tt.wantErr && len(errors) > 0 {
				if errors[0].Type != tt.errorType {
					t.Errorf("Expected error type '%s', got '%s'", tt.errorType, errors[0].Type)
				}
			}
		})
	}
}

// Test data type compatibility
func TestDataTypeCompatibility(t *testing.T) {
	tests := []struct {
		type1      *schema.DataType
		type2      *schema.DataType
		compatible bool
	}{
		{
			type1:      &schema.DataType{Name: "INT"},
			type2:      &schema.DataType{Name: "INT"},
			compatible: true,
		},
		{
			type1:      &schema.DataType{Name: "INT"},
			type2:      &schema.DataType{Name: "BIGINT"},
			compatible: true, // Numeric types are compatible
		},
		{
			type1:      &schema.DataType{Name: "VARCHAR"},
			type2:      &schema.DataType{Name: "TEXT"},
			compatible: true, // String types are compatible
		},
		{
			type1:      &schema.DataType{Name: "INT"},
			type2:      &schema.DataType{Name: "VARCHAR"},
			compatible: false, // Different type categories
		},
	}

	for i, tt := range tests {
		result := tt.type1.IsCompatibleWith(tt.type2)
		if result != tt.compatible {
			t.Errorf("Test %d: Expected compatibility %v, got %v for %s vs %s",
				i, tt.compatible, result, tt.type1.Name, tt.type2.Name)
		}
	}
}
