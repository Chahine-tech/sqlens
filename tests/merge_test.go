package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func TestMergeBasic(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Basic MERGE with UPDATE and INSERT",
			sql: `MERGE INTO target_table t
				USING source_table s
				ON t.id = s.id
				WHEN MATCHED THEN
					UPDATE SET t.name = s.name
				WHEN NOT MATCHED THEN
					INSERT (id, name) VALUES (s.id, s.name)`,
			dialect: "sqlserver",
		},
		{
			name: "MERGE with INTO keyword",
			sql: `MERGE INTO customers t
				USING updates u
				ON t.customer_id = u.customer_id
				WHEN MATCHED THEN
					UPDATE SET t.status = u.status`,
			dialect: "postgresql",
		},
		{
			name: "MERGE without INTO keyword",
			sql: `MERGE customers t
				USING updates u
				ON t.customer_id = u.customer_id
				WHEN MATCHED THEN
					UPDATE SET t.status = u.status`,
			dialect: "sqlserver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			mergeStmt, ok := stmt.(*parser.MergeStatement)
			if !ok {
				t.Fatalf("Expected MergeStatement, got %T", stmt)
			}

			if mergeStmt == nil {
				t.Fatal("MergeStatement is nil")
			}
		})
	}
}

func TestMergeWithSubquery(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "MERGE with subquery source",
			sql: `MERGE INTO inventory i
				USING (SELECT product_id, quantity FROM updates WHERE date = '2024-01-01') u
				ON i.product_id = u.product_id
				WHEN MATCHED THEN
					UPDATE SET i.quantity = u.quantity`,
			dialect: "postgresql",
		},
		{
			name: "MERGE with aliased subquery",
			sql: `MERGE INTO products p
				USING (SELECT id, price FROM new_prices) AS src
				ON p.id = src.id
				WHEN MATCHED THEN
					UPDATE SET p.price = src.price`,
			dialect: "sqlserver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			mergeStmt, ok := stmt.(*parser.MergeStatement)
			if !ok {
				t.Fatalf("Expected MergeStatement, got %T", stmt)
			}

			// Check that source is a subquery
			_, ok = mergeStmt.SourceTable.(*parser.SelectStatement)
			if !ok {
				t.Fatalf("Expected source to be SelectStatement, got %T", mergeStmt.SourceTable)
			}
		})
	}
}

func TestMergeWithConditions(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "MERGE with AND condition in WHEN MATCHED",
			sql: `MERGE INTO customers c
				USING updates u
				ON c.id = u.id
				WHEN MATCHED AND u.status = 'active' THEN
					UPDATE SET c.status = u.status`,
			dialect: "postgresql",
		},
		{
			name: "MERGE with multiple WHEN clauses",
			sql: `MERGE INTO inventory i
				USING updates u
				ON i.product_id = u.product_id
				WHEN MATCHED AND u.quantity > 0 THEN
					UPDATE SET i.quantity = u.quantity
				WHEN MATCHED AND u.quantity = 0 THEN
					DELETE
				WHEN NOT MATCHED THEN
					INSERT (product_id, quantity) VALUES (u.product_id, u.quantity)`,
			dialect: "sqlserver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			mergeStmt, ok := stmt.(*parser.MergeStatement)
			if !ok {
				t.Fatalf("Expected MergeStatement, got %T", stmt)
			}

			if len(mergeStmt.WhenMatched) == 0 {
				t.Fatal("Expected at least one WHEN MATCHED clause")
			}
		})
	}
}

func TestMergeActions(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
		action  string
	}{
		{
			name: "MERGE with UPDATE action",
			sql: `MERGE INTO customers
				USING updates
				ON customers.id = updates.id
				WHEN MATCHED THEN
					UPDATE SET name = updates.name, status = updates.status`,
			dialect: "sqlserver",
			action:  "UPDATE",
		},
		{
			name: "MERGE with INSERT action",
			sql: `MERGE INTO customers
				USING updates
				ON customers.id = updates.id
				WHEN NOT MATCHED THEN
					INSERT (id, name, status) VALUES (updates.id, updates.name, updates.status)`,
			dialect: "postgresql",
			action:  "INSERT",
		},
		{
			name: "MERGE with DELETE action",
			sql: `MERGE INTO inventory
				USING updates
				ON inventory.id = updates.id
				WHEN MATCHED AND updates.quantity = 0 THEN
					DELETE`,
			dialect: "sqlserver",
			action:  "DELETE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := parser.NewWithDialect(ctx, tt.sql, dialect.GetDialect(tt.dialect))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			mergeStmt, ok := stmt.(*parser.MergeStatement)
			if !ok {
				t.Fatalf("Expected MergeStatement, got %T", stmt)
			}

			var foundAction bool
			for _, when := range mergeStmt.WhenMatched {
				if when.Action.ActionType == tt.action {
					foundAction = true
					break
				}
			}
			for _, when := range mergeStmt.WhenNotMatched {
				if when.Action.ActionType == tt.action {
					foundAction = true
					break
				}
			}

			if !foundAction {
				t.Fatalf("Expected action %s not found", tt.action)
			}
		})
	}
}

func TestMergeSQLServerBySource(t *testing.T) {
	sql := `MERGE INTO target
		USING source
		ON target.id = source.id
		WHEN MATCHED THEN
			UPDATE SET target.value = source.value
		WHEN NOT MATCHED BY SOURCE THEN
			DELETE`

	ctx := context.Background()
	p := parser.NewWithDialect(ctx, sql, dialect.GetDialect("sqlserver"))
	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	mergeStmt, ok := stmt.(*parser.MergeStatement)
	if !ok {
		t.Fatalf("Expected MergeStatement, got %T", stmt)
	}

	if len(mergeStmt.WhenNotMatchedBy) == 0 {
		t.Fatal("Expected WHEN NOT MATCHED BY SOURCE clause")
	}

	if !mergeStmt.WhenNotMatchedBy[0].BySource {
		t.Fatal("Expected BySource flag to be true")
	}
}
