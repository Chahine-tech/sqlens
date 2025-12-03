package tests

import (
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/lexer"
)

func TestDollarQuotedStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple dollar-quoted string",
			input:    `$$Hello World$$`,
			expected: "Hello World",
		},
		{
			name:     "Dollar-quoted with tag",
			input:    `$body$Function body here$body$`,
			expected: "Function body here",
		},
		{
			name:     "Dollar-quoted with numbers in tag",
			input:    `$tag123$Content with $$ inside$tag123$`,
			expected: "Content with $$ inside",
		},
		{
			name:     "Dollar-quoted with SQL",
			input:    `$$SELECT * FROM users WHERE id = 1$$`,
			expected: "SELECT * FROM users WHERE id = 1",
		},
		{
			name:     "Dollar-quoted with newlines",
			input:    "$$Line 1\nLine 2\nLine 3$$",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "Dollar-quoted with single quotes",
			input:    `$$It's a test with 'quotes'$$`,
			expected: "It's a test with 'quotes'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewWithDialect(tt.input, dialect.GetDialect("postgresql"))
			tok := l.NextToken()

			if tok.Type != lexer.STRING {
				t.Fatalf("Expected STRING token, got %v", tok.Type)
			}

			if tok.Literal != tt.expected {
				t.Fatalf("Expected %q, got %q", tt.expected, tok.Literal)
			}
		})
	}
}

func TestDollarQuotedInProcedure(t *testing.T) {
	// Test dollar-quoted strings in PostgreSQL function
	sql := `CREATE FUNCTION test_func() RETURNS text AS $$
BEGIN
    RETURN 'Hello from function';
END;
$$ LANGUAGE plpgsql;`

	l := lexer.NewWithDialect(sql, dialect.GetDialect("postgresql"))

	// Skip tokens until we find the first dollar-quoted string
	for {
		tok := l.NextToken()
		if tok.Type == lexer.EOF {
			t.Fatal("Did not find dollar-quoted string")
		}
		if tok.Type == lexer.STRING {
			// Found the function body
			expected := "\nBEGIN\n    RETURN 'Hello from function';\nEND;\n"
			if tok.Literal != expected {
				t.Fatalf("Expected function body:\n%q\nGot:\n%q", expected, tok.Literal)
			}
			break
		}
	}
}

func TestDollarQuotedNested(t *testing.T) {
	// Test nested dollar-quoted strings with different tags
	input := `$outer$Text with $inner$nested$$inner$ content$outer$`

	l := lexer.NewWithDialect(input, dialect.GetDialect("postgresql"))
	tok := l.NextToken()

	if tok.Type != lexer.STRING {
		t.Fatalf("Expected STRING token, got %v", tok.Type)
	}

	expected := "Text with $inner$nested$$inner$ content"
	if tok.Literal != expected {
		t.Fatalf("Expected %q, got %q", expected, tok.Literal)
	}
}

func TestDollarQuotedNotPostgreSQL(t *testing.T) {
	// Dollar signs should not be treated as delimiters in non-PostgreSQL dialects
	input := `$$test$$`

	dialects := []string{"mysql", "sqlserver", "sqlite", "oracle"}

	for _, dialectName := range dialects {
		t.Run(dialectName, func(t *testing.T) {
			l := lexer.NewWithDialect(input, dialect.GetDialect(dialectName))
			tok := l.NextToken()

			// Should get ILLEGAL token for dollar sign
			if tok.Type != lexer.ILLEGAL {
				t.Fatalf("Expected ILLEGAL token for $ in %s, got %v", dialectName, tok.Type)
			}
		})
	}
}
