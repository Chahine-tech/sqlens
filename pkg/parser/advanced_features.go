// Package parser - Advanced SQL features (CTEs, Window Functions, Set Operations, CASE)
package parser

import (
	"fmt"

	"github.com/Chahine-tech/sql-parser-go/pkg/lexer"
)

// parseWithStatement parses a WITH (CTE) statement
// Syntax: WITH [RECURSIVE] cte_name [(columns)] AS (query) [, ...] main_query
func (p *Parser) parseWithStatement() (*WithStatement, error) {
	stmt := &WithStatement{}

	if !p.curTokenIs(lexer.WITH) {
		return nil, fmt.Errorf("expected WITH, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Check for RECURSIVE keyword
	if p.curTokenIs(lexer.RECURSIVE) {
		stmt.Recursive = true
		p.nextToken()
	}

	// Parse CTEs (one or more)
	for {
		cte, err := p.parseCommonTableExpression()
		if err != nil {
			return nil, err
		}
		stmt.CTEs = append(stmt.CTEs, cte)

		// Check if there are more CTEs (comma-separated)
		if !p.curTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma and move to next CTE name
	}

	// Parse the main query (usually a SELECT)
	// At this point, curToken should be the start of the main query
	mainQuery, err := p.ParseStatement()
	if err != nil {
		return nil, fmt.Errorf("failed to parse main query after WITH clause: %v", err)
	}
	stmt.Query = mainQuery

	return stmt, nil
}

// parseCommonTableExpression parses a single CTE
// Syntax: cte_name [(col1, col2, ...)] AS (SELECT ...)
func (p *Parser) parseCommonTableExpression() (*CommonTableExpression, error) {
	cte := &CommonTableExpression{}

	// Parse CTE name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected CTE name, got %s", p.curToken.Literal)
	}
	cte.Name = p.curToken.Literal
	p.nextToken()

	// Optional: parse column list
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()

		// Check if it's a column list or the AS clause
		if p.curTokenIs(lexer.IDENT) {
			// Parse column names
			for {
				if !p.curTokenIs(lexer.IDENT) {
					return nil, fmt.Errorf("expected column name in CTE column list")
				}
				cte.Columns = append(cte.Columns, p.curToken.Literal)
				p.nextToken()

				if !p.curTokenIs(lexer.COMMA) {
					break
				}
				p.nextToken() // consume comma
			}

			if !p.curTokenIs(lexer.RPAREN) {
				return nil, fmt.Errorf("expected ')' after CTE column list, got %s", p.curToken.Literal)
			}
			p.nextToken()
		}
	}

	// Expect AS keyword
	if !p.curTokenIs(lexer.AS) {
		return nil, fmt.Errorf("expected AS after CTE name, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Expect opening parenthesis
	if !p.curTokenIs(lexer.LPAREN) {
		return nil, fmt.Errorf("expected '(' after AS, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse the SELECT statement
	if !p.curTokenIs(lexer.SELECT) {
		return nil, fmt.Errorf("expected SELECT in CTE, got %s", p.curToken.Literal)
	}

	selectStmt, err := p.parseSelectStatement()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CTE query: %v", err)
	}
	cte.Query = selectStmt

	// Expect closing parenthesis
	if !p.curTokenIs(lexer.RPAREN) {
		return nil, fmt.Errorf("expected ')' after CTE query, got %s", p.curToken.Literal)
	}
	p.nextToken()

	return cte, nil
}

// parseSetOperation parses UNION, INTERSECT, EXCEPT operations
func (p *Parser) parseSetOperation(left Statement) (Statement, error) {
	// Check if next token is a set operator
	if !p.peekTokenIs(lexer.UNION) && !p.peekTokenIs(lexer.INTERSECT) && !p.peekTokenIs(lexer.EXCEPT) {
		// No set operation, return the original statement
		return left, nil
	}

	p.nextToken() // move to UNION/INTERSECT/EXCEPT

	setOp := &SetOperation{
		Left:     left,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	// Check for ALL keyword
	if p.curTokenIs(lexer.ALL) {
		setOp.All = true
		p.nextToken()
	}

	// Parse right side (must be SELECT)
	if !p.curTokenIs(lexer.SELECT) {
		return nil, fmt.Errorf("expected SELECT after %s, got %s", setOp.Operator, p.curToken.Literal)
	}

	right, err := p.parseSelectStatement()
	if err != nil {
		return nil, fmt.Errorf("failed to parse right side of %s: %v", setOp.Operator, err)
	}
	setOp.Right = right

	// Check for chained set operations
	return p.parseSetOperation(setOp)
}

// parseWindowFunction parses a window function
// Syntax: function_name(args) OVER (...)
func (p *Parser) parseWindowFunction(funcCall *FunctionCall) (*WindowFunction, error) {
	wf := &WindowFunction{
		Function: funcCall,
	}

	// curToken is already OVER (from parseFunctionCall)
	// Parse OVER clause
	overClause, err := p.parseOverClause()
	if err != nil {
		return nil, err
	}
	wf.OverClause = overClause

	return wf, nil
}

// parseOverClause parses the OVER clause of a window function
// Syntax: OVER (PARTITION BY ... ORDER BY ... frame_clause)
func (p *Parser) parseOverClause() (*OverClause, error) {
	oc := &OverClause{}

	p.nextToken() // move past OVER

	// Expect opening parenthesis
	if !p.curTokenIs(lexer.LPAREN) {
		return nil, fmt.Errorf("expected '(' after OVER, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Optional: PARTITION BY clause
	if p.curTokenIs(lexer.PARTITION) {
		p.nextToken()
		if !p.curTokenIs(lexer.BY) {
			return nil, fmt.Errorf("expected BY after PARTITION, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Parse partition expressions
		for {
			expr, err := p.parseExpression()
			if err != nil {
				return nil, fmt.Errorf("failed to parse PARTITION BY expression: %v", err)
			}
			oc.PartitionBy = append(oc.PartitionBy, expr)

			if !p.curTokenIs(lexer.COMMA) {
				break
			}
			p.nextToken() // consume comma
		}
	}

	// Optional: ORDER BY clause
	if p.curTokenIs(lexer.ORDER) {
		p.nextToken()
		if !p.curTokenIs(lexer.BY) {
			return nil, fmt.Errorf("expected BY after ORDER, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Parse order by items manually (we've already consumed ORDER BY)
		for {
			item, err := p.parseOrderByItem()
			if err != nil {
				return nil, fmt.Errorf("failed to parse ORDER BY item in window: %v", err)
			}
			oc.OrderBy = append(oc.OrderBy, item)

			if !p.curTokenIs(lexer.COMMA) {
				break
			}
			p.nextToken() // consume comma
		}
	}

	// Optional: Window frame (ROWS/RANGE)
	if p.curTokenIs(lexer.ROWS) || p.curTokenIs(lexer.RANGE) {
		frame, err := p.parseWindowFrame()
		if err != nil {
			return nil, err
		}
		oc.Frame = frame
	}

	// Expect closing parenthesis
	if !p.curTokenIs(lexer.RPAREN) {
		return nil, fmt.Errorf("expected ')' to close OVER clause, got %s", p.curToken.Literal)
	}
	p.nextToken() // consume the closing paren

	return oc, nil
}

// parseWindowFrame parses window frame specification
// Syntax: ROWS|RANGE BETWEEN start AND end
func (p *Parser) parseWindowFrame() (*WindowFrame, error) {
	wf := &WindowFrame{}

	// ROWS or RANGE
	wf.FrameType = p.curToken.Literal
	p.nextToken()

	// BETWEEN keyword (or single bound)
	if p.curTokenIs(lexer.BETWEEN) {
		p.nextToken()

		// Parse start bound
		start, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		wf.Start = start

		// Expect AND
		if !p.curTokenIs(lexer.AND) {
			return nil, fmt.Errorf("expected AND in window frame, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Parse end bound
		end, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		wf.End = end
	} else {
		// Single bound (e.g., "ROWS UNBOUNDED PRECEDING")
		bound, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		wf.Start = bound
		wf.End = bound
	}

	return wf, nil
}

// parseFrameBound parses a frame boundary
// Syntax: UNBOUNDED PRECEDING|FOLLOWING | CURRENT ROW | <expr> PRECEDING|FOLLOWING
func (p *Parser) parseFrameBound() (*FrameBound, error) {
	fb := &FrameBound{}

	if p.curTokenIs(lexer.UNBOUNDED) {
		fb.BoundType = "UNBOUNDED"
		p.nextToken()
		if p.curTokenIs(lexer.PRECEDING) {
			fb.Direction = "PRECEDING"
		} else if p.curTokenIs(lexer.FOLLOWING) {
			fb.Direction = "FOLLOWING"
		} else {
			return nil, fmt.Errorf("expected PRECEDING or FOLLOWING after UNBOUNDED")
		}
		p.nextToken()
	} else if p.curTokenIs(lexer.CURRENT) {
		fb.BoundType = "CURRENT"
		p.nextToken()
		if !p.curTokenIs(lexer.ROW) {
			return nil, fmt.Errorf("expected ROW after CURRENT, got %s", p.curToken.Literal)
		}
		p.nextToken()
	} else if p.curTokenIs(lexer.NUMBER) {
		// Expression-based bound
		offset, err := p.parseNumberLiteral()
		if err != nil {
			return nil, err
		}
		fb.Offset = offset
		fb.BoundType = "EXPRESSION"

		if p.curTokenIs(lexer.PRECEDING) {
			fb.Direction = "PRECEDING"
		} else if p.curTokenIs(lexer.FOLLOWING) {
			fb.Direction = "FOLLOWING"
		} else {
			return nil, fmt.Errorf("expected PRECEDING or FOLLOWING after frame offset")
		}
		p.nextToken()
	}

	return fb, nil
}

// parseCaseExpression parses a CASE expression
// Syntax: CASE [input] WHEN condition THEN result [...] [ELSE result] END
func (p *Parser) parseCaseExpression() (*CaseExpression, error) {
	ce := &CaseExpression{}

	if !p.curTokenIs(lexer.CASE) {
		return nil, fmt.Errorf("expected CASE, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Check if it's a simple CASE (has input expression)
	if !p.curTokenIs(lexer.WHEN) {
		// Parse input expression for simple CASE
		input, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse CASE input: %v", err)
		}
		ce.Input = input
		p.nextToken()
	}

	// Parse WHEN clauses
	for p.curTokenIs(lexer.WHEN) {
		whenClause, err := p.parseWhenClause()
		if err != nil {
			return nil, err
		}
		ce.WhenClauses = append(ce.WhenClauses, whenClause)
	}

	if len(ce.WhenClauses) == 0 {
		return nil, fmt.Errorf("CASE expression must have at least one WHEN clause")
	}

	// Optional: ELSE clause
	if p.curTokenIs(lexer.ELSE) {
		p.nextToken()
		elseResult, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse ELSE result: %v", err)
		}
		ce.ElseResult = elseResult
		p.nextToken()
	}

	// Expect END keyword
	if !p.curTokenIs(lexer.END) {
		return nil, fmt.Errorf("expected END to close CASE expression, got %s", p.curToken.Literal)
	}
	p.nextToken()

	return ce, nil
}

// parseWhenClause parses a WHEN clause in a CASE expression
func (p *Parser) parseWhenClause() (*WhenClause, error) {
	wc := &WhenClause{}

	if !p.curTokenIs(lexer.WHEN) {
		return nil, fmt.Errorf("expected WHEN, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse condition
	condition, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to parse WHEN condition: %v", err)
	}
	wc.Condition = condition
	p.nextToken()

	// Expect THEN keyword
	if !p.curTokenIs(lexer.THEN) {
		return nil, fmt.Errorf("expected THEN after WHEN condition, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse result
	result, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to parse THEN result: %v", err)
	}
	wc.Result = result
	p.nextToken()

	return wc, nil
}
