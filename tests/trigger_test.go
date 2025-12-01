package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func TestCreateTrigger(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		dialect     string
		expectError bool
	}{
		{
			name:        "Simple BEFORE INSERT trigger",
			sql:         "CREATE TRIGGER audit_insert BEFORE INSERT ON users FOR EACH ROW BEGIN INSERT INTO audit_log VALUES (NEW.id, NOW()); END",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "AFTER UPDATE trigger",
			sql:         "CREATE TRIGGER update_timestamp AFTER UPDATE ON products FOR EACH ROW BEGIN UPDATE products SET updated_at = NOW() WHERE id = NEW.id; END",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "BEFORE DELETE trigger",
			sql:         "CREATE TRIGGER prevent_delete BEFORE DELETE ON important_data FOR EACH ROW BEGIN SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Cannot delete'; END",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "Multiple events (INSERT OR UPDATE)",
			sql:         "CREATE TRIGGER track_changes AFTER INSERT OR UPDATE ON users FOR EACH ROW BEGIN INSERT INTO changes_log VALUES (NEW.id); END",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "Trigger with IF NOT EXISTS (MySQL)",
			sql:         "CREATE TRIGGER IF NOT EXISTS auto_timestamp BEFORE INSERT ON orders FOR EACH ROW BEGIN SET NEW.created_at = NOW(); END",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "OR REPLACE trigger (PostgreSQL)",
			sql:         "CREATE OR REPLACE TRIGGER update_audit AFTER UPDATE ON accounts FOR EACH ROW BEGIN INSERT INTO audit VALUES (NEW.id); END",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "Trigger on schema.table",
			sql:         "CREATE TRIGGER log_changes AFTER INSERT ON myschema.users FOR EACH ROW BEGIN INSERT INTO log VALUES (NEW.id); END",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "INSTEAD OF trigger (SQL Server/Oracle)",
			sql:         "CREATE TRIGGER view_insert INSTEAD OF INSERT ON users_view FOR EACH ROW BEGIN INSERT INTO users VALUES (NEW.id, NEW.name); END",
			dialect:     "sqlserver",
			expectError: false,
		},
		{
			name:        "Trigger with WHEN condition",
			sql:         "CREATE TRIGGER conditional_trigger BEFORE UPDATE ON products FOR EACH ROW WHEN (NEW.price > 1000) BEGIN UPDATE products SET requires_approval = 1 WHERE id = NEW.id; END",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "FOR EACH STATEMENT trigger",
			sql:         "CREATE TRIGGER statement_trigger AFTER DELETE ON users FOR EACH STATEMENT BEGIN INSERT INTO deletion_log VALUES (NOW()); END",
			dialect:     "postgresql",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			triggerStmt, ok := stmt.(*parser.CreateTriggerStatement)
			if !ok {
				t.Errorf("Expected *CreateTriggerStatement, got %T", stmt)
				return
			}

			if triggerStmt.TriggerName == "" {
				t.Errorf("Trigger name is empty")
				return
			}

			if triggerStmt.Timing == "" {
				t.Errorf("Trigger timing is empty")
				return
			}

			if len(triggerStmt.Events) == 0 {
				t.Errorf("Trigger events list is empty")
				return
			}

			if triggerStmt.TableName.Name == "" {
				t.Errorf("Table name is empty")
				return
			}

			t.Logf("✅ Successfully parsed CREATE TRIGGER: %s %s %v ON %s",
				triggerStmt.TriggerName,
				triggerStmt.Timing,
				triggerStmt.Events,
				triggerStmt.TableName.Name)
		})
	}
}

func TestDropTrigger(t *testing.T) {
	tests := []struct {
		name        string
		sql         string
		dialect     string
		expectError bool
	}{
		{
			name:        "Simple DROP TRIGGER",
			sql:         "DROP TRIGGER audit_insert",
			dialect:     "mysql",
			expectError: false,
		},
		{
			name:        "DROP TRIGGER IF EXISTS",
			sql:         "DROP TRIGGER IF EXISTS update_timestamp",
			dialect:     "postgresql",
			expectError: false,
		},
		{
			name:        "DROP TRIGGER with schema",
			sql:         "DROP TRIGGER myschema.log_changes",
			dialect:     "postgresql",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			dropStmt, ok := stmt.(*parser.DropStatement)
			if !ok {
				t.Errorf("Expected *DropStatement, got %T", stmt)
				return
			}

			if dropStmt.ObjectType != "TRIGGER" {
				t.Errorf("Expected ObjectType TRIGGER, got %s", dropStmt.ObjectType)
				return
			}

			if dropStmt.ObjectName == "" {
				t.Errorf("Trigger name is empty")
				return
			}

			t.Logf("✅ Successfully parsed DROP TRIGGER: %s", dropStmt.ObjectName)
		})
	}
}

func TestCreateTriggerDialects(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name:    "MySQL BEFORE INSERT trigger",
			sql:     "CREATE TRIGGER `audit_trigger` BEFORE INSERT ON `users` FOR EACH ROW BEGIN INSERT INTO audit VALUES (NEW.id); END",
			dialect: "mysql",
		},
		{
			name:    "PostgreSQL AFTER UPDATE trigger",
			sql:     "CREATE TRIGGER \"audit_trigger\" AFTER UPDATE ON \"users\" FOR EACH ROW BEGIN INSERT INTO audit VALUES (NEW.id); END",
			dialect: "postgresql",
		},
		{
			name:    "SQL Server INSTEAD OF trigger",
			sql:     "CREATE TRIGGER [audit_trigger] INSTEAD OF DELETE ON [users] FOR EACH ROW BEGIN INSERT INTO audit VALUES (OLD.id); END",
			dialect: "sqlserver",
		},
		{
			name:    "SQLite trigger",
			sql:     "CREATE TRIGGER audit_trigger AFTER INSERT ON users FOR EACH ROW BEGIN INSERT INTO audit VALUES (NEW.id); END",
			dialect: "sqlite",
		},
		{
			name:    "Oracle trigger",
			sql:     "CREATE TRIGGER audit_trigger BEFORE UPDATE ON users FOR EACH ROW BEGIN INSERT INTO audit VALUES (NEW.id); END",
			dialect: "oracle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewWithDialect(context.Background(), tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.dialect, err)
				return
			}

			triggerStmt, ok := stmt.(*parser.CreateTriggerStatement)
			if !ok {
				t.Errorf("Expected *CreateTriggerStatement for %s, got %T", tt.dialect, stmt)
				return
			}

			t.Logf("✅ %s: Successfully parsed CREATE TRIGGER", tt.dialect)
			_ = triggerStmt
		})
	}
}
