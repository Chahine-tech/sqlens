// Package parser provides SQL parsing functionality for SQL queries.
package parser

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/lexer"
)

type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	errors []string

	parseStartTime time.Time
	tokenCount     int

	ctx     context.Context
	dialect dialect.Dialect
}

func New(input string) *Parser {
	return NewWithContext(context.Background(), input)
}

func NewWithContext(ctx context.Context, input string) *Parser {
	return NewWithDialect(ctx, input, dialect.GetDialect("sqlserver"))
}

func NewWithDialect(ctx context.Context, input string, d dialect.Dialect) *Parser {
	l := lexer.NewWithDialect(input, d)
	p := &Parser{
		l:              l,
		errors:         make([]string, 0, 4),
		parseStartTime: time.Now(),
		ctx:            ctx,
		dialect:        d,
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	select {
	case <-p.ctx.Done():
		p.errors = append(p.errors, "parsing cancelled due to timeout")
		return
	default:
		p.curToken = p.peekToken
		p.peekToken = p.l.NextToken()
		p.tokenCount++
	}
}

// GetDialect returns the dialect used by this parser
func (p *Parser) GetDialect() dialect.Dialect {
	return p.dialect
}

// SetDialect sets the dialect for this parser
func (p *Parser) SetDialect(d dialect.Dialect) {
	p.dialect = d
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	syntaxErr := NewSyntaxError(
		t.String(),
		p.peekToken.Type.String(),
		p.peekToken.Line,
		p.peekToken.Column,
	)
	p.errors = append(p.errors, syntaxErr.Error())
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) GetParseMetrics() map[string]interface{} {
	duration := time.Since(p.parseStartTime)
	return map[string]interface{}{
		"parse_duration_ms": duration.Milliseconds(),
		"tokens_processed":  p.tokenCount,
		"tokens_per_second": float64(p.tokenCount) / duration.Seconds(),
		"error_count":       len(p.errors),
	}
}

func (p *Parser) ParseStatement() (Statement, error) {
	switch p.curToken.Type {
	case lexer.WITH:
		return p.parseWithStatement()
	case lexer.SELECT:
		stmt, err := p.parseSelectStatement()
		if err != nil {
			return nil, err
		}
		// Check for set operations (UNION, INTERSECT, EXCEPT)
		return p.parseSetOperation(stmt)
	case lexer.INSERT:
		return p.parseInsertStatement()
	case lexer.UPDATE:
		return p.parseUpdateStatement()
	case lexer.DELETE:
		return p.parseDeleteStatement()
	case lexer.CREATE:
		return p.parseCreateStatement()
	case lexer.DROP:
		return p.parseDropStatement()
	case lexer.ALTER:
		return p.parseAlterStatement()
	case lexer.BEGIN, lexer.START:
		return p.parseBeginTransaction()
	case lexer.COMMIT:
		return p.parseCommit()
	case lexer.ROLLBACK:
		return p.parseRollback()
	case lexer.SAVEPOINT:
		return p.parseSavepoint()
	case lexer.RELEASE:
		return p.parseReleaseSavepoint()
	default:
		return nil, fmt.Errorf("unsupported statement type: %s", p.curToken.Literal)
	}
}

// Parse SELECT statement
func (p *Parser) parseSelectStatement() (*SelectStatement, error) {
	stmt := GetSelectStatement() // Use object pool

	if !p.curTokenIs(lexer.SELECT) {
		PutSelectStatement(stmt)
		return nil, fmt.Errorf("expected SELECT, got %s", p.curToken.Literal)
	}

	p.nextToken()

	if p.curTokenIs(lexer.DISTINCT) {
		stmt.Distinct = true
		p.nextToken()
	}

	if p.curTokenIs(lexer.TOP) {
		topClause, err := p.parseTopClause()
		if err != nil {
			return nil, err
		}
		stmt.Top = topClause
	}

	columns, err := p.parseSelectList()
	if err != nil {
		return nil, err
	}
	stmt.Columns = columns

	if p.curTokenIs(lexer.FROM) {
		fromClause, err := p.parseFromClause()
		if err != nil {
			return nil, err
		}
		stmt.From = fromClause
	}

	for p.curTokenIs(lexer.JOIN) || p.curTokenIs(lexer.INNER) || p.curTokenIs(lexer.LEFT) || p.curTokenIs(lexer.RIGHT) || p.curTokenIs(lexer.FULL) {
		joinClause, err := p.parseJoinClause()
		if err != nil {
			return nil, err
		}
		stmt.Joins = append(stmt.Joins, joinClause)
	}

	if p.curTokenIs(lexer.WHERE) {
		p.nextToken()
		whereExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = whereExpr
	}

	if p.curTokenIs(lexer.GROUP) {
		groupBy, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		stmt.GroupBy = groupBy
	}

	if p.curTokenIs(lexer.HAVING) {
		p.nextToken()
		havingExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Having = havingExpr
	}

	if p.curTokenIs(lexer.ORDER) {
		orderBy, err := p.parseOrderByClause()
		if err != nil {
			return nil, err
		}
		stmt.OrderBy = orderBy
	}

	// Parse LIMIT clause
	if p.curTokenIs(lexer.LIMIT) {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = limit
	}

	return stmt, nil
}

func (p *Parser) parseTopClause() (*TopClause, error) {
	if !p.curTokenIs(lexer.TOP) {
		return nil, fmt.Errorf("expected TOP, got %s", p.curToken.Literal)
	}

	p.nextToken()

	if !p.curTokenIs(lexer.NUMBER) {
		return nil, fmt.Errorf("expected number after TOP, got %s", p.curToken.Literal)
	}

	count, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil, fmt.Errorf("invalid number in TOP clause: %s", p.curToken.Literal)
	}

	topClause := &TopClause{Count: count}
	p.nextToken()

	// Check for PERCENT
	if p.curTokenIs(lexer.IDENT) && strings.ToUpper(p.curToken.Literal) == "PERCENT" {
		topClause.Percent = true
		p.nextToken()
	}

	return topClause, nil
}

func (p *Parser) parseSelectList() ([]Expression, error) {
	var columns []Expression

	if p.curTokenIs(lexer.ASTERISK) {
		columns = append(columns, &StarExpression{})
		p.nextToken()
		return columns, nil
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Check for alias (AS keyword or implicit alias)
	if p.curTokenIs(lexer.AS) {
		p.nextToken()
		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected identifier after AS, got %s", p.curToken.Literal)
		}
		aliasExpr := &AliasedExpression{
			Expression: expr,
			Alias:      p.curToken.Literal,
		}
		columns = append(columns, aliasExpr)
		p.nextToken()
	} else {
		columns = append(columns, expr)
	}

	for p.curTokenIs(lexer.COMMA) {
		p.nextToken()

		if p.curTokenIs(lexer.ASTERISK) {
			columns = append(columns, &StarExpression{})
			p.nextToken()
		} else {
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			// Check for alias
			if p.curTokenIs(lexer.AS) {
				p.nextToken()
				if !p.curTokenIs(lexer.IDENT) {
					return nil, fmt.Errorf("expected identifier after AS, got %s", p.curToken.Literal)
				}
				aliasExpr := &AliasedExpression{
					Expression: expr,
					Alias:      p.curToken.Literal,
				}
				columns = append(columns, aliasExpr)
				p.nextToken()
			} else {
				columns = append(columns, expr)
			}
		}
	}

	return columns, nil
}

func (p *Parser) parseFromClause() (*FromClause, error) {
	if !p.curTokenIs(lexer.FROM) {
		return nil, fmt.Errorf("expected FROM, got %s", p.curToken.Literal)
	}

	p.nextToken()

	fromClause := &FromClause{}

	table, err := p.parseTableReference()
	if err != nil {
		return nil, err
	}
	fromClause.Tables = append(fromClause.Tables, *table)

	for p.curTokenIs(lexer.COMMA) {
		p.nextToken()
		table, err := p.parseTableReference()
		if err != nil {
			return nil, err
		}
		fromClause.Tables = append(fromClause.Tables, *table)
	}

	return fromClause, nil
}

func (p *Parser) parseTableReference() (*TableReference, error) {
	table := &TableReference{}

	// Check for derived table: (SELECT ...) AS alias
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()

		// Check if this is a subquery
		if p.curTokenIs(lexer.SELECT) {
			subquery, err := p.parseSelectStatement()
			if err != nil {
				return nil, fmt.Errorf("failed to parse derived table subquery: %v", err)
			}

			if !p.curTokenIs(lexer.RPAREN) {
				return nil, fmt.Errorf("expected ')' to close derived table, got %s", p.curToken.Literal)
			}
			p.nextToken()

			table.Subquery = subquery

			// Alias is required for derived tables (in most SQL dialects)
			if p.curTokenIs(lexer.AS) {
				p.nextToken()
			}

			if !p.curTokenIs(lexer.IDENT) {
				return nil, fmt.Errorf("derived table requires an alias, got %s", p.curToken.Literal)
			}
			table.Alias = p.curToken.Literal
			p.nextToken()

			return table, nil
		} else {
			return nil, fmt.Errorf("expected SELECT in derived table, got %s", p.curToken.Literal)
		}
	}

	// Regular table reference
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected table name, got %s", p.curToken.Literal)
	}

	firstIdent := p.curToken.Literal
	p.nextToken()

	if p.curTokenIs(lexer.DOT) {
		p.nextToken()
		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected table name after dot, got %s", p.curToken.Literal)
		}
		table.Schema = firstIdent
		table.Name = p.curToken.Literal
		p.nextToken()
	} else {
		table.Name = firstIdent
	}

	if p.curTokenIs(lexer.AS) {
		p.nextToken()
		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected alias after AS, got %s", p.curToken.Literal)
		}
		table.Alias = p.curToken.Literal
		p.nextToken()
	} else if p.curTokenIs(lexer.IDENT) {
		// Implicit alias (no AS keyword)
		table.Alias = p.curToken.Literal
		p.nextToken()
	}

	return table, nil
}

func (p *Parser) parseJoinClause() (*JoinClause, error) {
	joinClause := GetJoinClause()

	if p.curTokenIs(lexer.INNER) {
		joinClause.JoinType = "INNER"
		p.nextToken()
		if !p.curTokenIs(lexer.JOIN) {
			PutJoinClause(joinClause)
			return nil, fmt.Errorf("expected JOIN after INNER, got %s", p.curToken.Literal)
		}
		p.nextToken()
	} else if p.curTokenIs(lexer.LEFT) {
		joinClause.JoinType = "LEFT"
		p.nextToken()
		if !p.curTokenIs(lexer.JOIN) {
			PutJoinClause(joinClause)
			return nil, fmt.Errorf("expected JOIN after LEFT, got %s", p.curToken.Literal)
		}
		p.nextToken()
	} else if p.curTokenIs(lexer.RIGHT) {
		joinClause.JoinType = "RIGHT"
		p.nextToken()
		if !p.curTokenIs(lexer.JOIN) {
			PutJoinClause(joinClause)
			return nil, fmt.Errorf("expected JOIN after RIGHT, got %s", p.curToken.Literal)
		}
		p.nextToken()
	} else if p.curTokenIs(lexer.FULL) {
		joinClause.JoinType = "FULL"
		p.nextToken()
		if !p.curTokenIs(lexer.JOIN) {
			PutJoinClause(joinClause)
			return nil, fmt.Errorf("expected JOIN after FULL, got %s", p.curToken.Literal)
		}
		p.nextToken()
	} else if p.curTokenIs(lexer.JOIN) {
		joinClause.JoinType = "INNER"
		p.nextToken()
	}

	table, err := p.parseTableReference()
	if err != nil {
		PutJoinClause(joinClause)
		return nil, err
	}
	joinClause.Table = *table

	// Parse ON condition
	if !p.curTokenIs(lexer.ON) {
		PutJoinClause(joinClause) // Return to pool on error
		return nil, fmt.Errorf("expected ON after JOIN table, got %s", p.curToken.Literal)
	}

	p.nextToken()
	condition, err := p.parseExpression()
	if err != nil {
		PutJoinClause(joinClause) // Return to pool on error
		return nil, err
	}
	joinClause.Condition = condition

	return joinClause, nil
}

func (p *Parser) parseGroupByClause() ([]Expression, error) {
	if !p.curTokenIs(lexer.GROUP) {
		return nil, fmt.Errorf("expected GROUP, got %s", p.curToken.Literal)
	}

	p.nextToken()
	if !p.curTokenIs(lexer.BY) {
		return nil, fmt.Errorf("expected BY after GROUP, got %s", p.curToken.Literal)
	}

	p.nextToken()

	var expressions []Expression

	// Parse first expression
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	expressions = append(expressions, expr)

	// Parse additional expressions
	for p.curTokenIs(lexer.COMMA) {
		p.nextToken()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}

func (p *Parser) parseOrderByClause() ([]*OrderByClause, error) {
	if !p.curTokenIs(lexer.ORDER) {
		return nil, fmt.Errorf("expected ORDER, got %s", p.curToken.Literal)
	}

	p.nextToken()
	if !p.curTokenIs(lexer.BY) {
		return nil, fmt.Errorf("expected BY after ORDER, got %s", p.curToken.Literal)
	}
	p.nextToken()

	var clauses []*OrderByClause

	// Parse first clause
	clause, err := p.parseOrderByItem()
	if err != nil {
		return nil, err
	}
	clauses = append(clauses, clause)

	// Parse additional clauses
	for p.curTokenIs(lexer.COMMA) {
		p.nextToken()
		clause, err := p.parseOrderByItem()
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, clause)
	}

	return clauses, nil
}

func (p *Parser) parseOrderByItem() (*OrderByClause, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	clause := &OrderByClause{
		Expression: expr,
		Direction:  "ASC", // Default
	}

	// Check for ASC/DESC
	if p.curTokenIs(lexer.IDENT) {
		direction := strings.ToUpper(p.curToken.Literal)
		if direction == "ASC" || direction == "DESC" {
			clause.Direction = direction
			p.nextToken()
		}
	}

	return clause, nil
}

func (p *Parser) parseLimitClause() (*LimitClause, error) {
	if !p.curTokenIs(lexer.LIMIT) {
		return nil, fmt.Errorf("expected LIMIT, got %s", p.curToken.Literal)
	}

	p.nextToken()

	// Parse count
	if !p.curTokenIs(lexer.NUMBER) {
		return nil, fmt.Errorf("expected number after LIMIT, got %s", p.curToken.Literal)
	}

	count, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil, fmt.Errorf("invalid LIMIT count: %s", p.curToken.Literal)
	}

	clause := &LimitClause{
		Count:  count,
		Offset: 0, // Default
	}

	p.nextToken()

	// Check for OFFSET (MySQL/PostgreSQL style: LIMIT count OFFSET offset)
	if p.curTokenIs(lexer.OFFSET) {
		p.nextToken()
		if !p.curTokenIs(lexer.NUMBER) {
			return nil, fmt.Errorf("expected number after OFFSET, got %s", p.curToken.Literal)
		}

		offset, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			return nil, fmt.Errorf("invalid OFFSET value: %s", p.curToken.Literal)
		}

		clause.Offset = offset
		p.nextToken()
	}

	return clause, nil
}

// Basic expression parsing (simplified for now)
func (p *Parser) parseExpression() (Expression, error) {
	return p.parseInfixExpression()
}

func (p *Parser) parseInfixExpression() (Expression, error) {
	left, err := p.parsePrimaryExpression()
	if err != nil {
		return nil, err
	}

	for p.isInfixOperator(p.curToken.Type) {
		if p.curToken.Type == lexer.IN {
			// Special handling for IN expressions
			inExpr, err := p.parseInExpression(left)
			if err != nil {
				return nil, err
			}
			left = inExpr
		} else {
			operator := p.curToken.Literal
			p.nextToken()

			right, err := p.parsePrimaryExpression()
			if err != nil {
				return nil, err
			}

			expr := GetBinaryExpression() // Use object pool
			expr.Left = left
			expr.Operator = operator
			expr.Right = right
			left = expr
		}
	}

	return left, nil
}

func (p *Parser) parseInExpression(left Expression) (Expression, error) {
	inExpr := &InExpression{
		Expression: left,
		Not:        false,
	}

	// Move past the IN token
	p.nextToken()

	// Expect opening parenthesis
	if !p.curTokenIs(lexer.LPAREN) {
		return nil, fmt.Errorf("expected '(' after IN, got %s", p.curToken.Literal)
	}

	p.nextToken()

	// Check if this is a subquery (starts with SELECT)
	if p.curTokenIs(lexer.SELECT) {
		// Parse subquery
		subquery, err := p.parseSelectStatement()
		if err != nil {
			return nil, fmt.Errorf("failed to parse subquery in IN clause: %v", err)
		}

		// Wrap in SubqueryExpression
		subqueryExpr := &SubqueryExpression{
			Query: subquery,
		}
		inExpr.Values = []Expression{subqueryExpr}
	} else {
		// Parse list of values
		values := make([]Expression, 0)

		// Parse first value
		if !p.curTokenIs(lexer.RPAREN) {
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			values = append(values, expr)

			// Parse additional values
			for p.curTokenIs(lexer.COMMA) {
				p.nextToken()
				expr, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				values = append(values, expr)
			}
		}

		inExpr.Values = values
	}

	// Expect closing parenthesis
	if !p.curTokenIs(lexer.RPAREN) {
		return nil, fmt.Errorf("expected ')' after IN values, got %s", p.curToken.Literal)
	}

	p.nextToken()

	return inExpr, nil
}

func (p *Parser) parsePrimaryExpression() (Expression, error) {
	switch p.curToken.Type {
	case lexer.IDENT:
		return p.parseIdentifierExpression()
	case lexer.NUMBER:
		return p.parseNumberLiteral()
	case lexer.STRING:
		return p.parseStringLiteral()
	case lexer.NULL:
		// Handle NULL literal
		expr := &Literal{Value: nil}
		p.nextToken()
		return expr, nil
	case lexer.ASTERISK:
		expr := &StarExpression{}
		p.nextToken()
		return expr, nil
	case lexer.LPAREN:
		return p.parseGroupedExpression()
	case lexer.CASE:
		return p.parseCaseExpression()
	case lexer.NOT:
		// Check if it's NOT EXISTS
		p.nextToken()
		if p.curTokenIs(lexer.EXISTS) {
			return p.parseExistsExpression(true)
		}
		return nil, fmt.Errorf("unexpected token after NOT: %s", p.curToken.Literal)
	case lexer.EXISTS:
		return p.parseExistsExpression(false)
	default:
		return nil, fmt.Errorf("unexpected token in expression: %s", p.curToken.Literal)
	}
}

func (p *Parser) parseIdentifierExpression() (Expression, error) {
	firstIdent := p.curToken.Literal
	p.nextToken()

	// Check if it's a qualified column (table.column)
	if p.curTokenIs(lexer.DOT) {
		p.nextToken()
		if !p.curTokenIs(lexer.IDENT) && !p.curTokenIs(lexer.ASTERISK) {
			return nil, fmt.Errorf("expected column name after dot, got %s", p.curToken.Literal)
		}

		if p.curTokenIs(lexer.ASTERISK) {
			expr := &StarExpression{Table: firstIdent}
			p.nextToken()
			return expr, nil
		}

		expr := GetColumnReference() // Use object pool
		expr.Table = firstIdent
		expr.Column = p.curToken.Literal
		p.nextToken()
		return expr, nil
	}

	// Check if it's a function call
	if p.curTokenIs(lexer.LPAREN) {
		return p.parseFunctionCall(firstIdent)
	}

	// It's a simple column reference
	expr := GetColumnReference() // Use object pool
	expr.Column = firstIdent
	return expr, nil
}

func (p *Parser) parseFunctionCall(name string) (Expression, error) {
	if !p.curTokenIs(lexer.LPAREN) {
		return nil, fmt.Errorf("expected '(' for function call, got %s", p.curToken.Literal)
	}

	p.nextToken()

	var arguments []Expression

	if !p.curTokenIs(lexer.RPAREN) {
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)

		for p.curTokenIs(lexer.COMMA) {
			p.nextToken()
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
		}
	}

	if !p.curTokenIs(lexer.RPAREN) {
		return nil, fmt.Errorf("expected ')' to close function call, got %s", p.curToken.Literal)
	}

	p.nextToken() // consume the closing paren

	funcCall := &FunctionCall{
		Name:      name,
		Arguments: arguments,
	}

	// Check if this is a window function (followed by OVER)
	if p.curTokenIs(lexer.OVER) {
		return p.parseWindowFunction(funcCall)
	}

	return funcCall, nil
}

func (p *Parser) parseNumberLiteral() (Expression, error) {
	literal := &Literal{}

	if strings.Contains(p.curToken.Literal, ".") {
		value, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse %q as float", p.curToken.Literal)
		}
		literal.Value = value
	} else {
		value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse %q as integer", p.curToken.Literal)
		}
		literal.Value = value
	}

	p.nextToken()
	return literal, nil
}

func (p *Parser) parseStringLiteral() (Expression, error) {
	literal := &Literal{Value: p.curToken.Literal}
	p.nextToken()
	return literal, nil
}

func (p *Parser) parseGroupedExpression() (Expression, error) {
	p.nextToken()

	// Check if this is a subquery (starts with SELECT)
	if p.curTokenIs(lexer.SELECT) {
		subquery, err := p.parseSelectStatement()
		if err != nil {
			return nil, fmt.Errorf("failed to parse subquery: %v", err)
		}

		if !p.curTokenIs(lexer.RPAREN) {
			return nil, fmt.Errorf("expected ')' to close subquery, got %s", p.curToken.Literal)
		}
		p.nextToken()

		return &SubqueryExpression{Query: subquery}, nil
	}

	exp, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.curTokenIs(lexer.RPAREN) {
		return nil, fmt.Errorf("expected ')' to close grouped expression, got %s", p.curToken.Literal)
	}
	p.nextToken()

	return exp, nil
}

func (p *Parser) parseExistsExpression(not bool) (Expression, error) {
	// Current token is EXISTS
	p.nextToken()

	// Expect opening parenthesis
	if !p.curTokenIs(lexer.LPAREN) {
		return nil, fmt.Errorf("expected '(' after EXISTS, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Expect SELECT
	if !p.curTokenIs(lexer.SELECT) {
		return nil, fmt.Errorf("expected SELECT after EXISTS (, got %s", p.curToken.Literal)
	}

	// Parse the subquery
	subquery, err := p.parseSelectStatement()
	if err != nil {
		return nil, fmt.Errorf("failed to parse subquery in EXISTS: %v", err)
	}

	// Expect closing parenthesis
	if !p.curTokenIs(lexer.RPAREN) {
		return nil, fmt.Errorf("expected ')' to close EXISTS subquery, got %s", p.curToken.Literal)
	}
	p.nextToken()

	return &ExistsExpression{
		Subquery: subquery,
		Not:      not,
	}, nil
}

func (p *Parser) isInfixOperator(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.ASSIGN, lexer.EQ, lexer.NOT_EQ, lexer.LT, lexer.GT, lexer.LTE, lexer.GTE,
		lexer.AND, lexer.OR, lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH,
		lexer.LIKE, lexer.IN, lexer.IS:
		return true
	default:
		return false
	}
}

// Stub implementations for other statement types
func (p *Parser) parseInsertStatement() (*InsertStatement, error) {
	stmt := &InsertStatement{}

	if !p.curTokenIs(lexer.INSERT) {
		return nil, fmt.Errorf("expected INSERT, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Expect INTO keyword
	if !p.curTokenIs(lexer.INTO) {
		return nil, fmt.Errorf("expected INTO after INSERT, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse table name
	table, err := p.parseTableReference()
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name: %v", err)
	}
	stmt.Table = *table

	// Optional: column list (col1, col2, ...)
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()

		// Check if it's a column list or VALUES
		if p.curTokenIs(lexer.IDENT) {
			// Parse column list
			for {
				if !p.curTokenIs(lexer.IDENT) {
					return nil, fmt.Errorf("expected column name, got %s", p.curToken.Literal)
				}
				stmt.Columns = append(stmt.Columns, p.curToken.Literal)
				p.nextToken()

				if !p.curTokenIs(lexer.COMMA) {
					break
				}
				p.nextToken() // consume comma
			}

			if !p.curTokenIs(lexer.RPAREN) {
				return nil, fmt.Errorf("expected ')' after column list, got %s", p.curToken.Literal)
			}
			p.nextToken()
		}
	}

	// VALUES or SELECT
	if p.curTokenIs(lexer.VALUES) {
		p.nextToken()

		// Parse VALUES rows: (val1, val2, ...), (val3, val4, ...)
		for {
			if !p.curTokenIs(lexer.LPAREN) {
				return nil, fmt.Errorf("expected '(' for VALUES row, got %s", p.curToken.Literal)
			}
			p.nextToken()

			var row []Expression
			for {
				expr, err := p.parseExpression()
				if err != nil {
					return nil, fmt.Errorf("failed to parse value: %v", err)
				}
				row = append(row, expr)

				if !p.curTokenIs(lexer.COMMA) {
					break
				}
				p.nextToken() // consume comma
			}

			if !p.curTokenIs(lexer.RPAREN) {
				return nil, fmt.Errorf("expected ')' after VALUES row, got %s", p.curToken.Literal)
			}
			p.nextToken()

			stmt.Values = append(stmt.Values, row)

			// Check for more rows
			if !p.curTokenIs(lexer.COMMA) {
				break
			}
			p.nextToken() // consume comma between rows
		}
	} else if p.curTokenIs(lexer.SELECT) {
		// INSERT INTO ... SELECT
		selectStmt, err := p.parseSelectStatement()
		if err != nil {
			return nil, fmt.Errorf("failed to parse SELECT in INSERT: %v", err)
		}
		stmt.Select = selectStmt
	} else {
		return nil, fmt.Errorf("expected VALUES or SELECT after table name, got %s", p.curToken.Literal)
	}

	return stmt, nil
}

func (p *Parser) parseUpdateStatement() (*UpdateStatement, error) {
	stmt := &UpdateStatement{}

	if !p.curTokenIs(lexer.UPDATE) {
		return nil, fmt.Errorf("expected UPDATE, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse table name
	table, err := p.parseTableReference()
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name: %v", err)
	}
	stmt.Table = *table

	// Expect SET keyword
	if !p.curTokenIs(lexer.SET) {
		return nil, fmt.Errorf("expected SET after table name, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse SET assignments: col1 = val1, col2 = val2, ...
	for {
		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected column name in SET clause, got %s", p.curToken.Literal)
		}

		assignment := &Assignment{
			Column: p.curToken.Literal,
		}
		p.nextToken()

		// Expect = operator
		if !p.curTokenIs(lexer.ASSIGN) {
			return nil, fmt.Errorf("expected '=' after column name, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Parse value expression
		value, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse value in SET clause: %v", err)
		}
		assignment.Value = value

		stmt.Set = append(stmt.Set, assignment)

		// Check for more assignments
		if !p.curTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma
	}

	// Optional: WHERE clause
	if p.curTokenIs(lexer.WHERE) {
		p.nextToken()
		where, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse WHERE clause: %v", err)
		}
		stmt.Where = where
	}

	// Optional: ORDER BY clause (MySQL/SQLite)
	if p.curTokenIs(lexer.ORDER) {
		orderBy, err := p.parseOrderByClause()
		if err != nil {
			return nil, fmt.Errorf("failed to parse ORDER BY: %v", err)
		}
		stmt.OrderBy = orderBy
	}

	// Optional: LIMIT clause (MySQL/SQLite)
	if p.curTokenIs(lexer.LIMIT) {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, fmt.Errorf("failed to parse LIMIT: %v", err)
		}
		stmt.Limit = limit
	}

	return stmt, nil
}

func (p *Parser) parseDeleteStatement() (*DeleteStatement, error) {
	stmt := &DeleteStatement{}

	if !p.curTokenIs(lexer.DELETE) {
		return nil, fmt.Errorf("expected DELETE, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Expect FROM keyword
	if !p.curTokenIs(lexer.FROM) {
		return nil, fmt.Errorf("expected FROM after DELETE, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse table name
	table, err := p.parseTableReference()
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name: %v", err)
	}
	stmt.From = *table

	// Optional: WHERE clause
	if p.curTokenIs(lexer.WHERE) {
		p.nextToken()
		where, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse WHERE clause: %v", err)
		}
		stmt.Where = where
	}

	// Optional: ORDER BY clause (MySQL/SQLite)
	if p.curTokenIs(lexer.ORDER) {
		orderBy, err := p.parseOrderByClause()
		if err != nil {
			return nil, fmt.Errorf("failed to parse ORDER BY: %v", err)
		}
		stmt.OrderBy = orderBy
	}

	// Optional: LIMIT clause (MySQL/SQLite)
	if p.curTokenIs(lexer.LIMIT) {
		limit, err := p.parseLimitClause()
		if err != nil {
			return nil, fmt.Errorf("failed to parse LIMIT: %v", err)
		}
		stmt.Limit = limit
	}

	return stmt, nil
}
