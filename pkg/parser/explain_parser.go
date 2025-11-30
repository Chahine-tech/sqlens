package parser

import (
	"fmt"

	"github.com/Chahine-tech/sql-parser-go/pkg/lexer"
)

// parseExplainStatement parses EXPLAIN statements
// Supports:
// - EXPLAIN SELECT ...
// - EXPLAIN ANALYZE SELECT ...
// - EXPLAIN FORMAT=JSON SELECT ... (MySQL, PostgreSQL)
// - EXPLAIN (ANALYZE, BUFFERS) SELECT ... (PostgreSQL)
// - EXPLAIN QUERY PLAN SELECT ... (SQLite)
// - EXPLAIN EXTENDED SELECT ... (MySQL)
func (p *Parser) parseExplainStatement() (Statement, error) {
	stmt := &ExplainStatement{
		Options: make(map[string]string),
	}

	// Consume EXPLAIN keyword
	p.nextToken()

	// Check for ANALYZE keyword (EXPLAIN ANALYZE)
	if p.curTokenIs(lexer.ANALYZE) {
		stmt.Analyze = true
		p.nextToken()
	}

	// Check for EXTENDED keyword (MySQL: EXPLAIN EXTENDED)
	if p.curTokenIs(lexer.EXTENDED) {
		stmt.Options["extended"] = "true"
		p.nextToken()
	}

	// Check for FORMAT option (MySQL/PostgreSQL: EXPLAIN FORMAT=JSON)
	if p.curTokenIs(lexer.FORMAT) {
		p.nextToken()

		// Expect = or just the format name
		if p.curTokenIs(lexer.ASSIGN) {
			p.nextToken()
		}

		if p.curTokenIs(lexer.IDENT) {
			stmt.Format = p.curToken.Literal
			p.nextToken()
		} else {
			return nil, fmt.Errorf("expected format type after FORMAT, got %s", p.curToken.Type)
		}
	}

	// Check for PostgreSQL style options: EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()

		for !p.curTokenIs(lexer.RPAREN) && !p.curTokenIs(lexer.EOF) {
			optionName := ""
			optionValue := ""

			// Get option name
			if p.curTokenIs(lexer.IDENT) {
				optionName = p.curToken.Literal
				p.nextToken()
			} else if p.curTokenIs(lexer.ANALYZE) {
				optionName = "ANALYZE"
				stmt.Analyze = true
				p.nextToken()
			} else if p.curTokenIs(lexer.FORMAT) {
				optionName = "FORMAT"
				p.nextToken()
			} else {
				return nil, fmt.Errorf("expected option name in EXPLAIN options, got %s", p.curToken.Type)
			}

			// Check for option value (e.g., FORMAT JSON)
			if !p.curTokenIs(lexer.COMMA) && !p.curTokenIs(lexer.RPAREN) {
				if p.curTokenIs(lexer.IDENT) {
					optionValue = p.curToken.Literal
					p.nextToken()
				}
			}

			// Store the option
			if optionName != "" {
				if optionName == "FORMAT" && optionValue != "" {
					stmt.Format = optionValue
				} else if optionValue != "" {
					stmt.Options[optionName] = optionValue
				} else {
					stmt.Options[optionName] = "true"
				}
			}

			// Skip comma if present
			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			}
		}

		// Expect closing parenthesis
		if !p.curTokenIs(lexer.RPAREN) {
			return nil, fmt.Errorf("expected ) after EXPLAIN options")
		}
		p.nextToken()
	}

	// Check for SQLite style: EXPLAIN QUERY PLAN
	if p.curTokenIs(lexer.QUERY) {
		p.nextToken()
		if p.curTokenIs(lexer.PLAN) {
			stmt.Options["query_plan"] = "true"
			p.nextToken()
		} else {
			return nil, fmt.Errorf("expected PLAN after QUERY in EXPLAIN QUERY PLAN")
		}
	}

	// Now parse the actual statement to explain
	innerStatement, err := p.ParseStatement()
	if err != nil {
		return nil, fmt.Errorf("failed to parse statement after EXPLAIN: %w", err)
	}

	stmt.Statement = innerStatement

	return stmt, nil
}
