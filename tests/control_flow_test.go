package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// TestIfStatement tests IF...THEN...ELSE...END IF parsing
func TestIfStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple IF THEN END IF",
			sql: `
				IF x > 10 THEN
					RETURN 1;
				END IF
			`,
			dialect: "mysql",
		},
		{
			name: "IF THEN ELSE END IF",
			sql: `
				IF x > 10 THEN
					RETURN 1;
				ELSE
					RETURN 0;
				END IF
			`,
			dialect: "mysql",
		},
		{
			name: "IF THEN ELSEIF ELSE END IF",
			sql: `
				IF x > 10 THEN
					RETURN 1;
				ELSEIF x > 5 THEN
					RETURN 2;
				ELSE
					RETURN 0;
				END IF;
			`,
			dialect: "mysql",
		},
		{
			name: "PostgreSQL ELSIF syntax",
			sql: `
				IF x > 10 THEN
					RETURN 1;
				ELSIF x > 5 THEN
					RETURN 2;
				ELSE
					RETURN 0;
				END IF
			`,
			dialect: "postgresql",
		},
		{
			name: "Multiple ELSEIF blocks",
			sql: `
				IF status = 'active' THEN
					RETURN 1;
				ELSEIF status = 'pending' THEN
					RETURN 2;
				ELSEIF status = 'inactive' THEN
					RETURN 3;
				ELSE
					RETURN 0;
				END IF;
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
				t.Fatalf("Failed to parse IF statement: %v", err)
			}

			ifStmt, ok := stmt.(*parser.IfStatement)
			if !ok {
				t.Fatalf("Expected *parser.IfStatement, got %T", stmt)
			}

			t.Logf("✅ Successfully parsed IF statement: %s", ifStmt.String())
		})
	}
}

// TestWhileStatement tests WHILE...DO...END WHILE parsing
func TestWhileStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple WHILE loop",
			sql: `
				WHILE x < 10 DO
					RETURN x + 1;
				END WHILE;
			`,
			dialect: "mysql",
		},
		{
			name: "WHILE with multiple statements",
			sql: `
				WHILE counter < 100 DO
					RETURN counter + 1;
					RETURN total + counter;
				END WHILE;
			`,
			dialect: "mysql",
		},
		{
			name: "Nested WHILE loops",
			sql: `
				WHILE i < 10 DO
					RETURN 0;
					WHILE j < 5 DO
						RETURN j + 1;
					END WHILE;
					RETURN i + 1;
				END WHILE;
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
				t.Fatalf("Failed to parse WHILE statement: %v", err)
			}

			whileStmt, ok := stmt.(*parser.WhileStatement)
			if !ok {
				t.Fatalf("Expected *parser.WhileStatement, got %T", stmt)
			}

			t.Logf("✅ Successfully parsed WHILE statement: %s", whileStmt.String())
		})
	}
}

// TestLoopStatement tests LOOP...END LOOP parsing
func TestLoopStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple LOOP with EXIT",
			sql: `
				LOOP
					RETURN x + 1;
					EXIT WHEN x > 10;
				END LOOP;
			`,
			dialect: "postgresql",
		},
		{
			name: "LOOP with unconditional EXIT",
			sql: `
				LOOP
					RETURN counter + 1;
					IF counter > 100 THEN
						EXIT;
					END IF;
				END LOOP;
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
				t.Fatalf("Failed to parse LOOP statement: %v", err)
			}

			loopStmt, ok := stmt.(*parser.LoopStatement)
			if !ok {
				t.Fatalf("Expected *parser.LoopStatement, got %T", stmt)
			}

			t.Logf("✅ Successfully parsed LOOP statement: %s", loopStmt.String())
		})
	}
}

// TestForStatement tests FOR...LOOP parsing
func TestForStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple FOR loop",
			sql: `
				FOR i IN 1..10 LOOP
					RETURN total + i;
				END LOOP;
			`,
			dialect: "postgresql",
		},
		{
			name: "FOR REVERSE loop",
			sql: `
				FOR i IN REVERSE 10..1 LOOP
					RETURN total + i;
				END LOOP;
			`,
			dialect: "postgresql",
		},
		{
			name: "FOR loop with BY step",
			sql: `
				FOR i IN 0..100 BY 10 LOOP
					RETURN result + i;
				END LOOP;
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
				t.Fatalf("Failed to parse FOR statement: %v", err)
			}

			forStmt, ok := stmt.(*parser.ForStatement)
			if !ok {
				t.Fatalf("Expected *parser.ForStatement, got %T", stmt)
			}

			t.Logf("✅ Successfully parsed FOR statement: %s", forStmt.String())
		})
	}
}

// TestRepeatStatement tests REPEAT...UNTIL parsing (MySQL)
func TestRepeatStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple REPEAT UNTIL",
			sql: `
				REPEAT
					RETURN x + 1;
				UNTIL x > 10;
			`,
			dialect: "mysql",
		},
		{
			name: "REPEAT with multiple statements",
			sql: `
				REPEAT
					RETURN counter + 1;
					RETURN total + counter;
				UNTIL counter >= 100;
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
				t.Fatalf("Failed to parse REPEAT statement: %v", err)
			}

			repeatStmt, ok := stmt.(*parser.RepeatStatement)
			if !ok {
				t.Fatalf("Expected *parser.RepeatStatement, got %T", stmt)
			}

			t.Logf("✅ Successfully parsed REPEAT statement: %s", repeatStmt.String())
		})
	}
}

// TestExitStatement tests EXIT and EXIT WHEN parsing
func TestExitStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Simple EXIT",
			sql: `
				LOOP
					RETURN x + 1;
					EXIT;
				END LOOP;
			`,
			dialect: "postgresql",
		},
		{
			name: "EXIT WHEN condition",
			sql: `
				LOOP
					RETURN x + 1;
					EXIT WHEN x > 10;
				END LOOP;
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
				t.Fatalf("Failed to parse EXIT statement: %v", err)
			}

			t.Logf("✅ Successfully parsed statement with EXIT: %s", stmt.String())
		})
	}
}

// TestContinueStatement tests CONTINUE and ITERATE parsing
func TestContinueStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "CONTINUE in LOOP",
			sql: `
				LOOP
					RETURN x + 1;
					CONTINUE WHEN x < 5;
					RETURN y + 1;
					EXIT WHEN x > 10;
				END LOOP;
			`,
			dialect: "postgresql",
		},
		{
			name: "ITERATE in WHILE (MySQL)",
			sql: `
				WHILE x < 10 DO
					RETURN x + 1;
					IF x = 5 THEN
						ITERATE;
					END IF;
					RETURN total + x;
				END WHILE;
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
				t.Fatalf("Failed to parse CONTINUE/ITERATE statement: %v", err)
			}

			t.Logf("✅ Successfully parsed statement with CONTINUE/ITERATE: %s", stmt.String())
		})
	}
}

// TestNestedControlFlow tests nested control flow structures
func TestNestedControlFlow(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "IF inside WHILE",
			sql: `
				WHILE x < 100 DO
					IF x > 50 THEN
						RETURN 1;
					ELSE
						RETURN 0;
					END IF;
					RETURN x + 1;
				END WHILE;
			`,
			dialect: "mysql",
		},
		{
			name: "WHILE inside IF",
			sql: `
				IF mode = 'batch' THEN
					WHILE counter < 100 DO
						RETURN counter + 1;
					END WHILE;
				END IF
			`,
			dialect: "mysql",
		},
		{
			name: "FOR inside LOOP",
			sql: `
				LOOP
					FOR i IN 1..10 LOOP
						RETURN total + i;
					END LOOP;
					EXIT WHEN total > 1000;
				END LOOP;
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
				t.Fatalf("Failed to parse nested control flow: %v", err)
			}

			t.Logf("✅ Successfully parsed nested control flow statement: %s", stmt.String())
		})
	}
}

// TestComplexControlFlow tests complex real-world scenarios
func TestComplexControlFlow(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		dialect string
	}{
		{
			name: "Complex business logic with multiple control flows",
			sql: `
				IF account_type = 'premium' THEN
					FOR i IN 1..premium_limit LOOP
						IF i > threshold THEN
							RETURN bonus + 10;
						ELSE
							RETURN bonus + 5;
						END IF;
					END LOOP;
				ELSEIF account_type = 'standard' THEN
					WHILE counter < standard_limit DO
						RETURN counter + 1;
						RETURN reward + 1;
					END WHILE;
				ELSE
					RETURN 0;
				END IF;
			`,
			dialect: "postgresql",
		},
		{
			name: "MySQL batch processing with REPEAT",
			sql: `
				REPEAT
					RETURN batch_count + 1;
					IF batch_count > 10 THEN
						EXIT;
					END IF;
				UNTIL batch_count >= 100
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
				t.Fatalf("Failed to parse complex control flow: %v", err)
			}

			t.Logf("✅ Successfully parsed complex control flow statement: %s", stmt.String())
		})
	}
}
