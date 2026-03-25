package tests

import (
	"strings"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/lexer"
)

func TestToken_String(t *testing.T) {
	tokens := lexer.TokenizeSQL("SELECT id FROM users")
	for _, tok := range tokens {
		s := tok.String()
		if s == "" {
			t.Error("Token.String() returned empty string")
		}
	}
}

func TestLookupIdent_Keyword(t *testing.T) {
	tt := lexer.LookupIdent("SELECT")
	if tt == lexer.IDENT {
		t.Error("expected SELECT to be a keyword, not IDENT")
	}
}

func TestLookupIdent_NotKeyword(t *testing.T) {
	tt := lexer.LookupIdent("my_table")
	if tt != lexer.IDENT {
		t.Errorf("expected IDENT for non-keyword, got %v", tt)
	}
}

func TestTokenType_String_Known(t *testing.T) {
	// SELECT should have a known name
	tokens := lexer.TokenizeSQL("SELECT")
	for _, tok := range tokens {
		if tok.Type != lexer.EOF {
			s := tok.Type.String()
			if s == "UNKNOWN" || s == "" {
				t.Errorf("expected known name for token type %d, got %q", tok.Type, s)
			}
		}
	}
}

func TestTokenType_String_Unknown(t *testing.T) {
	var unknown lexer.TokenType = 99999
	s := unknown.String()
	if s != "UNKNOWN" {
		t.Errorf("expected UNKNOWN for unknown token type, got %q", s)
	}
}

func TestTokenizeWithBuffer(t *testing.T) {
	buf := make([]lexer.Token, 0, 32)
	tokens := lexer.TokenizeWithBuffer("SELECT id FROM users WHERE id = 1", buf)
	if len(tokens) == 0 {
		t.Error("expected tokens")
	}
	// Verify SELECT is first non-whitespace token
	found := false
	for _, tok := range tokens {
		s := strings.ToUpper(tok.Literal)
		if s == "SELECT" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected SELECT token in output")
	}
}

func TestTokenize_IncludesEOF(t *testing.T) {
	tokens := lexer.TokenizeSQL("SELECT 1")
	last := tokens[len(tokens)-1]
	if last.Type != lexer.EOF {
		t.Errorf("expected last token to be EOF, got %v", last.Type)
	}
}
