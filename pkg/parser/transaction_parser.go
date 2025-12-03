package parser

import (
	"fmt"

	"github.com/Chahine-tech/sql-parser-go/pkg/lexer"
)

// parseBeginTransaction parses BEGIN or START TRANSACTION statements
// Syntax:
//   - BEGIN [WORK | TRANSACTION]
//   - START TRANSACTION
//
// Note: If BEGIN is followed by TRY, it's a TRY...CATCH block, not a transaction
func (p *Parser) parseBeginTransaction() (Statement, error) {
	if p.curTokenIs(lexer.BEGIN) {
		// Peek ahead to see if it's BEGIN TRY (SQL Server TRY...CATCH)
		if p.peekTokenIs(lexer.TRY) {
			p.nextToken() // consume BEGIN
			return p.parseTryStatement()
		}

		stmt := &BeginTransactionStatement{}
		stmt.UseStart = false
		p.nextToken()

		// Optional WORK or TRANSACTION keyword
		if p.curTokenIs(lexer.WORK) || p.curTokenIs(lexer.TRANSACTION) {
			p.nextToken()
		}
		return stmt, nil
	} else if p.curTokenIs(lexer.START) {
		stmt := &BeginTransactionStatement{}
		stmt.UseStart = true
		p.nextToken()

		// Expect TRANSACTION keyword
		if !p.curTokenIs(lexer.TRANSACTION) {
			return nil, fmt.Errorf("expected TRANSACTION after START, got %s", p.curToken.Literal)
		}
		p.nextToken()
		return stmt, nil
	}

	return nil, fmt.Errorf("expected BEGIN or START for transaction, got %s", p.curToken.Literal)
}

// parseCommit parses COMMIT statements
// Syntax:
//   - COMMIT [WORK]
func (p *Parser) parseCommit() (Statement, error) {
	stmt := &CommitStatement{}
	p.nextToken() // consume COMMIT

	// Optional WORK keyword
	if p.curTokenIs(lexer.WORK) {
		stmt.Work = true
		p.nextToken()
	}

	return stmt, nil
}

// parseRollback parses ROLLBACK statements
// Syntax:
//   - ROLLBACK [WORK]
//   - ROLLBACK TO SAVEPOINT name
func (p *Parser) parseRollback() (Statement, error) {
	stmt := &RollbackStatement{}
	p.nextToken() // consume ROLLBACK

	// Check for TO SAVEPOINT
	if p.curTokenIs(lexer.IDENT) && p.curToken.Literal == "TO" {
		p.nextToken() // consume TO

		if !p.curTokenIs(lexer.SAVEPOINT) {
			return nil, fmt.Errorf("expected SAVEPOINT after TO, got %s", p.curToken.Literal)
		}
		p.nextToken() // consume SAVEPOINT

		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected savepoint name, got %s", p.curToken.Literal)
		}
		stmt.ToSavepoint = p.curToken.Literal
		p.nextToken()
	} else if p.curTokenIs(lexer.WORK) {
		// Optional WORK keyword
		stmt.Work = true
		p.nextToken()
	}

	return stmt, nil
}

// parseSavepoint parses SAVEPOINT statements
// Syntax:
//   - SAVEPOINT name
func (p *Parser) parseSavepoint() (Statement, error) {
	stmt := &SavepointStatement{}
	p.nextToken() // consume SAVEPOINT

	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected savepoint name, got %s", p.curToken.Literal)
	}

	stmt.Name = p.curToken.Literal
	p.nextToken()

	return stmt, nil
}

// parseReleaseSavepoint parses RELEASE SAVEPOINT statements
// Syntax:
//   - RELEASE SAVEPOINT name
func (p *Parser) parseReleaseSavepoint() (Statement, error) {
	stmt := &ReleaseSavepointStatement{}
	p.nextToken() // consume RELEASE

	if !p.curTokenIs(lexer.SAVEPOINT) {
		return nil, fmt.Errorf("expected SAVEPOINT after RELEASE, got %s", p.curToken.Literal)
	}
	p.nextToken() // consume SAVEPOINT

	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected savepoint name, got %s", p.curToken.Literal)
	}

	stmt.Name = p.curToken.Literal
	p.nextToken()

	return stmt, nil
}
