package parser

import (
	"fmt"
	"strconv"

	"github.com/Chahine-tech/sql-parser-go/pkg/lexer"
)

// parseCreateStatement handles CREATE TABLE, CREATE INDEX, etc.
func (p *Parser) parseCreateStatement() (Statement, error) {
	if !p.curTokenIs(lexer.CREATE) {
		return nil, fmt.Errorf("expected CREATE, got %s", p.curToken.Literal)
	}
	p.nextToken()

	switch p.curToken.Type {
	case lexer.OR:
		// CREATE OR REPLACE PROCEDURE/FUNCTION
		p.nextToken()
		if !p.curTokenIs(lexer.REPLACE) {
			return nil, fmt.Errorf("expected REPLACE after OR, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Now check what comes after REPLACE
		if p.curTokenIs(lexer.PROCEDURE) {
			stmt, err := p.parseCreateProcedureStatement()
			if err == nil && stmt != nil {
				stmt.(*CreateProcedureStatement).OrReplace = true
			}
			return stmt, err
		} else if p.curTokenIs(lexer.FUNCTION) {
			stmt, err := p.parseCreateFunctionStatement()
			if err == nil && stmt != nil {
				stmt.(*CreateFunctionStatement).OrReplace = true
			}
			return stmt, err
		} else if p.curTokenIs(lexer.VIEW) || p.curTokenIs(lexer.MATERIALIZED) {
			stmt, err := p.parseCreateViewStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				stmt.OrReplace = true
			}
			return stmt, nil
		} else if p.curTokenIs(lexer.TRIGGER) {
			stmt, err := p.parseCreateTriggerStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				stmt.OrReplace = true
			}
			return stmt, nil
		}
		return nil, fmt.Errorf("expected PROCEDURE, FUNCTION, VIEW, or TRIGGER after CREATE OR REPLACE, got %s", p.curToken.Literal)

	case lexer.TABLE:
		return p.parseCreateTableStatement()
	case lexer.INDEX, lexer.UNIQUE:
		return p.parseCreateIndexStatement()
	case lexer.VIEW, lexer.MATERIALIZED:
		return p.parseCreateViewStatement()
	case lexer.PROCEDURE:
		return p.parseCreateProcedureStatement()
	case lexer.FUNCTION:
		return p.parseCreateFunctionStatement()
	case lexer.TRIGGER:
		return p.parseCreateTriggerStatement()
	default:
		return nil, fmt.Errorf("unsupported CREATE statement: CREATE %s", p.curToken.Literal)
	}
}

// parseCreateTableStatement parses CREATE TABLE statements
func (p *Parser) parseCreateTableStatement() (*CreateTableStatement, error) {
	stmt := &CreateTableStatement{}

	// Expect TABLE
	if !p.curTokenIs(lexer.TABLE) {
		return nil, fmt.Errorf("expected TABLE, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Check for IF NOT EXISTS
	if p.curTokenIs(lexer.IF) {
		p.nextToken()
		if !p.curTokenIs(lexer.NOT) {
			return nil, fmt.Errorf("expected NOT after IF, got %s", p.curToken.Literal)
		}
		p.nextToken()
		if !p.curTokenIs(lexer.EXISTS) {
			return nil, fmt.Errorf("expected EXISTS after IF NOT, got %s", p.curToken.Literal)
		}
		stmt.IfNotExists = true
		p.nextToken()
	}

	// Parse table name
	table, err := p.parseTableReference()
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name: %w", err)
	}
	stmt.Table = *table

	// Expect opening parenthesis
	if !p.curTokenIs(lexer.LPAREN) {
		return nil, fmt.Errorf("expected '(' after table name, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse column definitions and constraints
	for !p.curTokenIs(lexer.RPAREN) && !p.curTokenIs(lexer.EOF) {
		// Check if this is a table constraint
		if p.curTokenIs(lexer.PRIMARY) || p.curTokenIs(lexer.FOREIGN) ||
			p.curTokenIs(lexer.UNIQUE) || p.curTokenIs(lexer.CONSTRAINT) {
			constraint, err := p.parseTableConstraint()
			if err != nil {
				return nil, err
			}
			stmt.Constraints = append(stmt.Constraints, constraint)
		} else {
			// Parse column definition
			column, err := p.parseColumnDefinition()
			if err != nil {
				return nil, err
			}
			stmt.Columns = append(stmt.Columns, column)
		}

		// Check for comma
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken()
		} else if !p.curTokenIs(lexer.RPAREN) {
			return nil, fmt.Errorf("expected ',' or ')', got %s", p.curToken.Literal)
		}
	}

	// Expect closing parenthesis
	if !p.curTokenIs(lexer.RPAREN) {
		return nil, fmt.Errorf("expected ')' to close CREATE TABLE, got %s", p.curToken.Literal)
	}
	p.nextToken()

	return stmt, nil
}

// parseColumnDefinition parses a column definition
func (p *Parser) parseColumnDefinition() (*ColumnDefinition, error) {
	col := &ColumnDefinition{}

	// Column name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected column name, got %s", p.curToken.Literal)
	}
	col.Name = p.curToken.Literal
	p.nextToken()

	// Data type
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected data type, got %s", p.curToken.Literal)
	}
	col.DataType = p.curToken.Literal
	p.nextToken()

	// Check for length/precision: VARCHAR(255), DECIMAL(10,2)
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()
		if !p.curTokenIs(lexer.NUMBER) {
			return nil, fmt.Errorf("expected number for type length, got %s", p.curToken.Literal)
		}
		length, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			return nil, fmt.Errorf("invalid length: %w", err)
		}
		col.Length = length
		p.nextToken()

		// Check for scale (DECIMAL(10,2))
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken()
			if !p.curTokenIs(lexer.NUMBER) {
				return nil, fmt.Errorf("expected number for type scale, got %s", p.curToken.Literal)
			}
			scale, err := strconv.Atoi(p.curToken.Literal)
			if err != nil {
				return nil, fmt.Errorf("invalid scale: %w", err)
			}
			col.Precision = col.Length
			col.Scale = scale
			p.nextToken()
		}

		if !p.curTokenIs(lexer.RPAREN) {
			return nil, fmt.Errorf("expected ')' after type parameters, got %s", p.curToken.Literal)
		}
		p.nextToken()
	}

	// Parse column constraints
	for {
		switch p.curToken.Type {
		case lexer.NOT:
			p.nextToken()
			if !p.curTokenIs(lexer.NULL) {
				return nil, fmt.Errorf("expected NULL after NOT, got %s", p.curToken.Literal)
			}
			col.NotNull = true
			p.nextToken()

		case lexer.NULL:
			// NULL is allowed but we don't need to set anything
			p.nextToken()

		case lexer.PRIMARY:
			p.nextToken()
			if !p.curTokenIs(lexer.KEY) {
				return nil, fmt.Errorf("expected KEY after PRIMARY, got %s", p.curToken.Literal)
			}
			col.PrimaryKey = true
			p.nextToken()

		case lexer.UNIQUE:
			col.Unique = true
			p.nextToken()

		case lexer.AUTO_INCREMENT, lexer.AUTOINCREMENT, lexer.IDENTITY:
			col.AutoIncrement = true
			p.nextToken()

		case lexer.DEFAULT:
			p.nextToken()
			// Parse default value expression
			defaultExpr, err := p.parseExpression()
			if err != nil {
				return nil, fmt.Errorf("failed to parse DEFAULT value: %w", err)
			}
			col.Default = defaultExpr

		case lexer.REFERENCES:
			// Inline foreign key
			fkRef, err := p.parseForeignKeyReference()
			if err != nil {
				return nil, err
			}
			col.References = fkRef

		default:
			// No more column constraints
			return col, nil
		}
	}
}

// parseTableConstraint parses table-level constraints
func (p *Parser) parseTableConstraint() (*TableConstraint, error) {
	constraint := &TableConstraint{}

	// Optional CONSTRAINT name
	if p.curTokenIs(lexer.CONSTRAINT) {
		p.nextToken()
		if p.curTokenIs(lexer.IDENT) {
			constraint.Name = p.curToken.Literal
			p.nextToken()
		}
	}

	switch p.curToken.Type {
	case lexer.PRIMARY:
		constraint.ConstraintType = "PRIMARY_KEY"
		p.nextToken()
		if !p.curTokenIs(lexer.KEY) {
			return nil, fmt.Errorf("expected KEY after PRIMARY, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Parse column list
		if !p.curTokenIs(lexer.LPAREN) {
			return nil, fmt.Errorf("expected '(' after PRIMARY KEY, got %s", p.curToken.Literal)
		}
		p.nextToken()

		for !p.curTokenIs(lexer.RPAREN) {
			if !p.curTokenIs(lexer.IDENT) {
				return nil, fmt.Errorf("expected column name, got %s", p.curToken.Literal)
			}
			constraint.Columns = append(constraint.Columns, p.curToken.Literal)
			p.nextToken()

			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			}
		}
		p.nextToken() // consume )

	case lexer.FOREIGN:
		constraint.ConstraintType = "FOREIGN_KEY"
		p.nextToken()
		if !p.curTokenIs(lexer.KEY) {
			return nil, fmt.Errorf("expected KEY after FOREIGN, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Parse column list
		if !p.curTokenIs(lexer.LPAREN) {
			return nil, fmt.Errorf("expected '(' after FOREIGN KEY, got %s", p.curToken.Literal)
		}
		p.nextToken()

		for !p.curTokenIs(lexer.RPAREN) {
			if !p.curTokenIs(lexer.IDENT) {
				return nil, fmt.Errorf("expected column name, got %s", p.curToken.Literal)
			}
			constraint.Columns = append(constraint.Columns, p.curToken.Literal)
			p.nextToken()

			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			}
		}
		p.nextToken() // consume )

		// Parse REFERENCES
		fkRef, err := p.parseForeignKeyReference()
		if err != nil {
			return nil, err
		}
		constraint.References = fkRef

	case lexer.UNIQUE:
		constraint.ConstraintType = "UNIQUE"
		p.nextToken()

		// Optional KEY keyword
		if p.curTokenIs(lexer.KEY) {
			p.nextToken()
		}

		// Parse column list
		if !p.curTokenIs(lexer.LPAREN) {
			return nil, fmt.Errorf("expected '(' after UNIQUE, got %s", p.curToken.Literal)
		}
		p.nextToken()

		for !p.curTokenIs(lexer.RPAREN) {
			if !p.curTokenIs(lexer.IDENT) {
				return nil, fmt.Errorf("expected column name, got %s", p.curToken.Literal)
			}
			constraint.Columns = append(constraint.Columns, p.curToken.Literal)
			p.nextToken()

			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			}
		}
		p.nextToken() // consume )

	default:
		return nil, fmt.Errorf("unexpected constraint type: %s", p.curToken.Literal)
	}

	return constraint, nil
}

// parseForeignKeyReference parses REFERENCES clause
func (p *Parser) parseForeignKeyReference() (*ForeignKeyReference, error) {
	if !p.curTokenIs(lexer.REFERENCES) {
		return nil, fmt.Errorf("expected REFERENCES, got %s", p.curToken.Literal)
	}
	p.nextToken()

	fkRef := &ForeignKeyReference{}

	// Table name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected table name after REFERENCES, got %s", p.curToken.Literal)
	}
	fkRef.Table = p.curToken.Literal
	p.nextToken()

	// Column list (optional in some dialects, but we'll require it)
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()
		for !p.curTokenIs(lexer.RPAREN) {
			if !p.curTokenIs(lexer.IDENT) {
				return nil, fmt.Errorf("expected column name, got %s", p.curToken.Literal)
			}
			fkRef.Columns = append(fkRef.Columns, p.curToken.Literal)
			p.nextToken()

			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			}
		}
		p.nextToken() // consume )
	}

	// Optional ON DELETE / ON UPDATE
	for p.curTokenIs(lexer.ON) {
		p.nextToken()
		if p.curTokenIs(lexer.DELETE) {
			p.nextToken()
			fkRef.OnDelete = p.parseReferentialAction()
		} else if p.curTokenIs(lexer.UPDATE) {
			p.nextToken()
			fkRef.OnUpdate = p.parseReferentialAction()
		} else {
			return nil, fmt.Errorf("expected DELETE or UPDATE after ON, got %s", p.curToken.Literal)
		}
	}

	return fkRef, nil
}

// parseReferentialAction parses CASCADE, SET NULL, etc.
func (p *Parser) parseReferentialAction() string {
	action := ""

	// Handle SET NULL, SET DEFAULT
	if p.curTokenIs(lexer.SET) {
		action = "SET"
		p.nextToken()
		if p.curTokenIs(lexer.NULL) {
			action = "SET NULL"
			p.nextToken()
		} else if p.curTokenIs(lexer.DEFAULT) {
			action = "SET DEFAULT"
			p.nextToken()
		}
	} else if p.curTokenIs(lexer.IDENT) {
		action = p.curToken.Literal // CASCADE, RESTRICT, etc.
		p.nextToken()

		// Handle NO ACTION
		if action == "NO" && p.curTokenIs(lexer.IDENT) && p.curToken.Literal == "ACTION" {
			action = "NO ACTION"
			p.nextToken()
		}
	}

	return action
}

// parseDropStatement handles DROP TABLE, DROP DATABASE, DROP INDEX
func (p *Parser) parseDropStatement() (*DropStatement, error) {
	stmt := &DropStatement{}

	if !p.curTokenIs(lexer.DROP) {
		return nil, fmt.Errorf("expected DROP, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Object type: TABLE, DATABASE, INDEX, VIEW
	switch p.curToken.Type {
	case lexer.TABLE:
		stmt.ObjectType = "TABLE"
	case lexer.DATABASE, lexer.SCHEMA:
		stmt.ObjectType = "DATABASE"
	case lexer.INDEX:
		stmt.ObjectType = "INDEX"
	case lexer.VIEW:
		stmt.ObjectType = "VIEW"
	case lexer.MATERIALIZED:
		// DROP MATERIALIZED VIEW (PostgreSQL)
		stmt.ObjectType = "MATERIALIZED VIEW"
		p.nextToken()
		if !p.curTokenIs(lexer.VIEW) {
			return nil, fmt.Errorf("expected VIEW after MATERIALIZED, got %s", p.curToken.Literal)
		}
	case lexer.TRIGGER:
		stmt.ObjectType = "TRIGGER"
	default:
		return nil, fmt.Errorf("expected TABLE, DATABASE, INDEX, VIEW, or TRIGGER, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Check for IF EXISTS
	if p.curTokenIs(lexer.IF) {
		p.nextToken()
		if !p.curTokenIs(lexer.EXISTS) {
			return nil, fmt.Errorf("expected EXISTS after IF, got %s", p.curToken.Literal)
		}
		stmt.IfExists = true
		p.nextToken()
	}

	// Object name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected object name, got %s", p.curToken.Literal)
	}
	stmt.ObjectName = p.curToken.Literal
	p.nextToken()

	// For DROP INDEX, might have ON table_name
	if stmt.ObjectType == "INDEX" && p.curTokenIs(lexer.ON) {
		p.nextToken()
		// Just consume the table name, we don't store it
		if p.curTokenIs(lexer.IDENT) {
			p.nextToken()
		}
	}

	// Optional CASCADE
	if p.curTokenIs(lexer.IDENT) && p.curToken.Literal == "CASCADE" {
		stmt.Cascade = true
		p.nextToken()
	}

	return stmt, nil
}

// parseAlterStatement handles ALTER TABLE
func (p *Parser) parseAlterStatement() (*AlterTableStatement, error) {
	stmt := &AlterTableStatement{}

	if !p.curTokenIs(lexer.ALTER) {
		return nil, fmt.Errorf("expected ALTER, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Expect TABLE
	if !p.curTokenIs(lexer.TABLE) {
		return nil, fmt.Errorf("expected TABLE after ALTER, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse table name
	table, err := p.parseTableReference()
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name: %w", err)
	}
	stmt.Table = *table

	// Parse ALTER action
	action, err := p.parseAlterAction()
	if err != nil {
		return nil, err
	}
	stmt.Action = action

	return stmt, nil
}

// parseAlterAction parses ADD/DROP/MODIFY/CHANGE
func (p *Parser) parseAlterAction() (*AlterAction, error) {
	action := &AlterAction{}

	switch p.curToken.Type {
	case lexer.ADD:
		action.ActionType = "ADD"
		p.nextToken()

		// Check if adding a constraint or column
		if p.curTokenIs(lexer.CONSTRAINT) || p.curTokenIs(lexer.PRIMARY) ||
			p.curTokenIs(lexer.FOREIGN) || p.curTokenIs(lexer.UNIQUE) {
			constraint, err := p.parseTableConstraint()
			if err != nil {
				return nil, err
			}
			action.Constraint = constraint
		} else {
			// Optional COLUMN keyword
			if p.curTokenIs(lexer.COLUMN) {
				p.nextToken()
			}

			// Parse column definition
			col, err := p.parseColumnDefinition()
			if err != nil {
				return nil, err
			}
			action.Column = col
		}

	case lexer.DROP:
		action.ActionType = "DROP"
		p.nextToken()

		// Check if dropping a CONSTRAINT or COLUMN
		if p.curTokenIs(lexer.CONSTRAINT) {
			p.nextToken()
			// Constraint name
			if !p.curTokenIs(lexer.IDENT) {
				return nil, fmt.Errorf("expected constraint name, got %s", p.curToken.Literal)
			}
			action.ColumnName = p.curToken.Literal // Reuse ColumnName for constraint name
			p.nextToken()
		} else {
			// Optional COLUMN keyword
			if p.curTokenIs(lexer.COLUMN) {
				p.nextToken()
			}

			// Column name
			if !p.curTokenIs(lexer.IDENT) {
				return nil, fmt.Errorf("expected column name, got %s", p.curToken.Literal)
			}
			action.ColumnName = p.curToken.Literal
			p.nextToken()
		}

	case lexer.MODIFY:
		action.ActionType = "MODIFY"
		p.nextToken()

		// Optional COLUMN keyword
		if p.curTokenIs(lexer.COLUMN) {
			p.nextToken()
		}

		// Parse new column definition
		col, err := p.parseColumnDefinition()
		if err != nil {
			return nil, err
		}
		action.Column = col

	case lexer.CHANGE:
		action.ActionType = "CHANGE"
		p.nextToken()

		// Optional COLUMN keyword
		if p.curTokenIs(lexer.COLUMN) {
			p.nextToken()
		}

		// Old column name
		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected old column name, got %s", p.curToken.Literal)
		}
		action.ColumnName = p.curToken.Literal
		p.nextToken()

		// New column definition
		col, err := p.parseColumnDefinition()
		if err != nil {
			return nil, err
		}
		action.NewColumn = col

	default:
		return nil, fmt.Errorf("expected ADD, DROP, MODIFY, or CHANGE, got %s", p.curToken.Literal)
	}

	return action, nil
}

// parseCreateIndexStatement parses CREATE INDEX
func (p *Parser) parseCreateIndexStatement() (*CreateIndexStatement, error) {
	stmt := &CreateIndexStatement{}

	// Check for UNIQUE
	if p.curTokenIs(lexer.UNIQUE) {
		stmt.Unique = true
		p.nextToken()
	}

	// Expect INDEX
	if !p.curTokenIs(lexer.INDEX) {
		return nil, fmt.Errorf("expected INDEX, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Check for IF NOT EXISTS
	if p.curTokenIs(lexer.IF) {
		p.nextToken()
		if !p.curTokenIs(lexer.NOT) {
			return nil, fmt.Errorf("expected NOT after IF, got %s", p.curToken.Literal)
		}
		p.nextToken()
		if !p.curTokenIs(lexer.EXISTS) {
			return nil, fmt.Errorf("expected EXISTS after IF NOT, got %s", p.curToken.Literal)
		}
		stmt.IfNotExists = true
		p.nextToken()
	}

	// Index name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected index name, got %s", p.curToken.Literal)
	}
	stmt.IndexName = p.curToken.Literal
	p.nextToken()

	// Expect ON
	if !p.curTokenIs(lexer.ON) {
		return nil, fmt.Errorf("expected ON after index name, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Table name
	table, err := p.parseTableReference()
	if err != nil {
		return nil, fmt.Errorf("failed to parse table name: %w", err)
	}
	stmt.Table = *table

	// Column list
	if !p.curTokenIs(lexer.LPAREN) {
		return nil, fmt.Errorf("expected '(' for column list, got %s", p.curToken.Literal)
	}
	p.nextToken()

	for !p.curTokenIs(lexer.RPAREN) {
		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected column name, got %s", p.curToken.Literal)
		}
		stmt.Columns = append(stmt.Columns, p.curToken.Literal)
		p.nextToken()

		if p.curTokenIs(lexer.COMMA) {
			p.nextToken()
		}
	}
	p.nextToken() // consume )

	return stmt, nil
}

// parseCreateViewStatement parses CREATE VIEW and CREATE MATERIALIZED VIEW statements
func (p *Parser) parseCreateViewStatement() (*CreateViewStatement, error) {
	stmt := &CreateViewStatement{
		Options: make(map[string]string),
	}

	// Check for MATERIALIZED VIEW (PostgreSQL)
	if p.curTokenIs(lexer.MATERIALIZED) {
		stmt.Materialized = true
		p.nextToken()
	}

	// Expect VIEW
	if !p.curTokenIs(lexer.VIEW) {
		return nil, fmt.Errorf("expected VIEW, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Check for IF NOT EXISTS
	if p.curTokenIs(lexer.IF) {
		p.nextToken()
		if !p.curTokenIs(lexer.NOT) {
			return nil, fmt.Errorf("expected NOT after IF, got %s", p.curToken.Literal)
		}
		p.nextToken()
		if !p.curTokenIs(lexer.EXISTS) {
			return nil, fmt.Errorf("expected EXISTS after NOT, got %s", p.curToken.Literal)
		}
		stmt.IfNotExists = true
		p.nextToken()
	}

	// Parse view name (can have schema)
	// We need a simpler approach than parseTableReference because AS is coming
	viewName := TableReference{}

	// Check for schema.table format
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected view name, got %s", p.curToken.Literal)
	}

	firstPart := p.curToken.Literal
	p.nextToken()

	if p.curTokenIs(lexer.DOT) {
		// schema.view_name
		viewName.Schema = firstPart
		p.nextToken()
		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected view name after schema, got %s", p.curToken.Literal)
		}
		viewName.Name = p.curToken.Literal
		p.nextToken()
	} else {
		// Just view_name
		viewName.Name = firstPart
	}

	stmt.ViewName = viewName

	// Optional: column list (col1, col2, ...)
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()
		for !p.curTokenIs(lexer.RPAREN) {
			if !p.curTokenIs(lexer.IDENT) {
				return nil, fmt.Errorf("expected column name in view column list, got %s", p.curToken.Literal)
			}
			stmt.Columns = append(stmt.Columns, p.curToken.Literal)
			p.nextToken()

			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			}
		}
		if !p.curTokenIs(lexer.RPAREN) {
			return nil, fmt.Errorf("expected ) after column list, got %s", p.curToken.Literal)
		}
		p.nextToken()
	}

	// Expect AS
	if !p.curTokenIs(lexer.AS) {
		return nil, fmt.Errorf("expected AS, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse SELECT statement
	if !p.curTokenIs(lexer.SELECT) {
		return nil, fmt.Errorf("expected SELECT after AS, got %s", p.curToken.Literal)
	}
	selectStmt, err := p.parseSelectStatement()
	if err != nil {
		return nil, fmt.Errorf("failed to parse SELECT in view: %w", err)
	}
	stmt.SelectStmt = selectStmt

	// Check for WITH CHECK OPTION
	if p.curTokenIs(lexer.WITH) {
		p.nextToken()
		if p.curTokenIs(lexer.CHECK) {
			p.nextToken()
			if !p.curTokenIs(lexer.OPTION) {
				return nil, fmt.Errorf("expected OPTION after WITH CHECK, got %s", p.curToken.Literal)
			}
			stmt.WithCheck = true
			p.nextToken()
		} else {
			// If it's not CHECK, move back - might be end of statement
			// For now, we'll just ignore unknown WITH clauses
		}
	}

	return stmt, nil
}

// parseCreateTriggerStatement parses CREATE TRIGGER statements
func (p *Parser) parseCreateTriggerStatement() (*CreateTriggerStatement, error) {
	stmt := &CreateTriggerStatement{
		Options: make(map[string]string),
	}

	// Expect TRIGGER
	if !p.curTokenIs(lexer.TRIGGER) {
		return nil, fmt.Errorf("expected TRIGGER, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Check for IF NOT EXISTS (MySQL)
	if p.curTokenIs(lexer.IF) {
		p.nextToken()
		if !p.curTokenIs(lexer.NOT) {
			return nil, fmt.Errorf("expected NOT after IF, got %s", p.curToken.Literal)
		}
		p.nextToken()
		if !p.curTokenIs(lexer.EXISTS) {
			return nil, fmt.Errorf("expected EXISTS after NOT, got %s", p.curToken.Literal)
		}
		stmt.IfNotExists = true
		p.nextToken()
	}

	// Parse trigger name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected trigger name, got %s", p.curToken.Literal)
	}
	stmt.TriggerName = p.curToken.Literal
	p.nextToken()

	// Parse timing: BEFORE, AFTER, or INSTEAD OF
	if p.curTokenIs(lexer.BEFORE) {
		stmt.Timing = "BEFORE"
		p.nextToken()
	} else if p.curTokenIs(lexer.AFTER) {
		stmt.Timing = "AFTER"
		p.nextToken()
	} else if p.curTokenIs(lexer.INSTEAD) {
		// INSTEAD OF (SQL Server, Oracle)
		stmt.Timing = "INSTEAD OF"
		p.nextToken()
		if !p.curTokenIs(lexer.OF) {
			return nil, fmt.Errorf("expected OF after INSTEAD, got %s", p.curToken.Literal)
		}
		p.nextToken()
	} else {
		return nil, fmt.Errorf("expected BEFORE, AFTER, or INSTEAD, got %s", p.curToken.Literal)
	}

	// Parse events: INSERT, UPDATE, DELETE (can be multiple with OR)
	for {
		if p.curTokenIs(lexer.INSERT) {
			stmt.Events = append(stmt.Events, "INSERT")
			p.nextToken()
		} else if p.curTokenIs(lexer.UPDATE) {
			stmt.Events = append(stmt.Events, "UPDATE")
			p.nextToken()
		} else if p.curTokenIs(lexer.DELETE) {
			stmt.Events = append(stmt.Events, "DELETE")
			p.nextToken()
		} else {
			break
		}

		// Check for OR (multiple events)
		if p.curTokenIs(lexer.OR) {
			p.nextToken()
		} else {
			break
		}
	}

	if len(stmt.Events) == 0 {
		return nil, fmt.Errorf("expected at least one trigger event (INSERT, UPDATE, DELETE)")
	}

	// Expect ON
	if !p.curTokenIs(lexer.ON) {
		return nil, fmt.Errorf("expected ON, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse table name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected table name, got %s", p.curToken.Literal)
	}
	tableName := p.curToken.Literal
	p.nextToken()

	// Check for schema.table
	if p.curTokenIs(lexer.DOT) {
		p.nextToken()
		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected table name after schema, got %s", p.curToken.Literal)
		}
		stmt.TableName = TableReference{
			Schema: tableName,
			Name:   p.curToken.Literal,
		}
		p.nextToken()
	} else {
		stmt.TableName = TableReference{
			Name: tableName,
		}
	}

	// Check for FOR EACH ROW/STATEMENT
	if p.curTokenIs(lexer.FOR) {
		p.nextToken()
		if !p.curTokenIs(lexer.EACH) {
			return nil, fmt.Errorf("expected EACH after FOR, got %s", p.curToken.Literal)
		}
		p.nextToken()
		if p.curTokenIs(lexer.ROW) {
			stmt.ForEachRow = true
			p.nextToken()
		} else if p.curTokenIs(lexer.IDENT) && p.curToken.Literal == "STATEMENT" {
			stmt.ForEachRow = false
			p.nextToken()
		} else {
			return nil, fmt.Errorf("expected ROW or STATEMENT after FOR EACH, got %s", p.curToken.Literal)
		}
	}

	// Check for WHEN condition (PostgreSQL)
	if p.curTokenIs(lexer.WHEN) {
		p.nextToken()
		if !p.curTokenIs(lexer.LPAREN) {
			return nil, fmt.Errorf("expected ( after WHEN, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Parse condition expression
		// For now, we'll skip detailed expression parsing and just consume until )
		// In a full implementation, you'd call p.parseExpression() here
		parenCount := 1
		for parenCount > 0 && !p.curTokenIs(lexer.EOF) {
			if p.curTokenIs(lexer.LPAREN) {
				parenCount++
			} else if p.curTokenIs(lexer.RPAREN) {
				parenCount--
			}
			if parenCount > 0 {
				p.nextToken()
			}
		}

		if !p.curTokenIs(lexer.RPAREN) {
			return nil, fmt.Errorf("expected ) after WHEN condition, got %s", p.curToken.Literal)
		}
		p.nextToken()
	}

	// Parse trigger body
	// For now, we'll handle BEGIN...END blocks simply
	if p.curTokenIs(lexer.BEGIN) {
		// Skip to END for now - full body parsing would be more complex
		body := &ProcedureBody{}
		p.nextToken()

		// Consume tokens until we find END
		for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.EOF) {
			p.nextToken()
		}

		if !p.curTokenIs(lexer.END) {
			return nil, fmt.Errorf("expected END for trigger body, got %s", p.curToken.Literal)
		}
		p.nextToken()

		stmt.Body = body
	}

	return stmt, nil
}
