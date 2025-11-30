package parser

import (
	"fmt"

	"github.com/Chahine-tech/sql-parser-go/pkg/lexer"
)

// parseCreateProcedureStatement parses CREATE PROCEDURE statement
// Note: CREATE keyword has already been consumed by parseCreateStatement()
func (p *Parser) parseCreateProcedureStatement() (Statement, error) {
	stmt := &CreateProcedureStatement{
		Parameters: make([]*ProcedureParameter, 0),
		Options:    make(map[string]string),
	}

	// PROCEDURE keyword (already positioned here by parseCreateStatement)
	if !p.curTokenIs(lexer.PROCEDURE) {
		return nil, fmt.Errorf("expected PROCEDURE, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Procedure name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected procedure name, got %s", p.curToken.Literal)
	}
	stmt.Name = p.curToken.Literal
	p.nextToken()

	// Parameters: (param1 type, param2 type, ...)
	if p.curTokenIs(lexer.LPAREN) {
		params, err := p.parseProcedureParameters()
		if err != nil {
			return nil, err
		}
		stmt.Parameters = params
	}

	// Parse procedure options (dialect-specific)
	if err := p.parseProcedureOptions(stmt); err != nil {
		return nil, err
	}

	// Parse procedure body (AS/IS BEGIN ... END)
	body, err := p.parseProcedureBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	return stmt, nil
}

// parseCreateFunctionStatement parses CREATE FUNCTION statement
// Note: CREATE keyword has already been consumed by parseCreateStatement()
func (p *Parser) parseCreateFunctionStatement() (Statement, error) {
	stmt := &CreateFunctionStatement{
		Parameters: make([]*ProcedureParameter, 0),
		Options:    make(map[string]string),
	}

	// FUNCTION keyword (already positioned here by parseCreateStatement)
	if !p.curTokenIs(lexer.FUNCTION) {
		return nil, fmt.Errorf("expected FUNCTION, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Function name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected function name, got %s", p.curToken.Literal)
	}
	stmt.Name = p.curToken.Literal
	p.nextToken()

	// Parameters: (param1 type, param2 type, ...)
	if p.curTokenIs(lexer.LPAREN) {
		params, err := p.parseProcedureParameters()
		if err != nil {
			return nil, err
		}
		stmt.Parameters = params
	}

	// RETURNS clause (required for functions)
	if p.curTokenIs(lexer.RETURNS) {
		p.nextToken()
		returnType, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		stmt.ReturnType = returnType
	} else {
		return nil, fmt.Errorf("expected RETURNS clause for function")
	}

	// Parse function options (dialect-specific)
	if err := p.parseFunctionOptions(stmt); err != nil {
		return nil, err
	}

	// Parse function body (AS/IS BEGIN ... END or expression)
	body, err := p.parseProcedureBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	return stmt, nil
}

// parseProcedureParameters parses procedure/function parameters
func (p *Parser) parseProcedureParameters() ([]*ProcedureParameter, error) {
	params := make([]*ProcedureParameter, 0)

	// Consume opening parenthesis
	p.nextToken()

	// Empty parameter list
	if p.curTokenIs(lexer.RPAREN) {
		p.nextToken()
		return params, nil
	}

	for {
		param := &ProcedureParameter{
			Mode: "IN", // Default mode
		}

		// Check for parameter mode (IN, OUT, INOUT)
		if p.curTokenIs(lexer.IN) {
			param.Mode = "IN"
			p.nextToken()
		} else if p.curTokenIs(lexer.OUT) {
			param.Mode = "OUT"
			p.nextToken()
		} else if p.curTokenIs(lexer.INOUT) {
			param.Mode = "INOUT"
			p.nextToken()
		}

		// Check for VARIADIC (PostgreSQL)
		if p.curTokenIs(lexer.VARIADIC) {
			param.IsVariadic = true
			p.nextToken()
		}

		// Parameter name
		if !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected parameter name, got %s", p.curToken.Literal)
		}
		param.Name = p.curToken.Literal
		p.nextToken()

		// Parameter type
		dataType, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		param.DataType = dataType

		// Optional DEFAULT value
		if p.curTokenIs(lexer.DEFAULT) || p.curTokenIs(lexer.ASSIGN) {
			p.nextToken()
			defaultValue, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			param.Default = defaultValue
		}

		params = append(params, param)

		// Check for more parameters
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken()
			continue
		}

		// End of parameters
		if p.curTokenIs(lexer.RPAREN) {
			p.nextToken()
			break
		}

		return nil, fmt.Errorf("expected comma or closing parenthesis, got %s", p.curToken.Literal)
	}

	return params, nil
}

// parseDataType parses a data type definition (VARCHAR(255), INT, DECIMAL(10,2), etc.)
func (p *Parser) parseDataType() (*DataTypeDefinition, error) {
	dataType := &DataTypeDefinition{}

	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected data type, got %s", p.curToken.Literal)
	}

	dataType.Name = p.curToken.Literal
	p.nextToken()

	// Check for size/precision: VARCHAR(255), DECIMAL(10,2)
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()

		// First number (length or precision)
		if !p.curTokenIs(lexer.NUMBER) {
			return nil, fmt.Errorf("expected number for data type size/precision")
		}
		// Parse as int
		var firstNum int
		fmt.Sscanf(p.curToken.Literal, "%d", &firstNum)
		p.nextToken()

		// Check for second number (scale for DECIMAL)
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken()
			if !p.curTokenIs(lexer.NUMBER) {
				return nil, fmt.Errorf("expected number for data type scale")
			}
			var secondNum int
			fmt.Sscanf(p.curToken.Literal, "%d", &secondNum)
			dataType.Precision = firstNum
			dataType.Scale = secondNum
			p.nextToken()
		} else {
			// Only one number = length
			dataType.Length = firstNum
		}

		if !p.curTokenIs(lexer.RPAREN) {
			return nil, fmt.Errorf("expected closing parenthesis after data type size")
		}
		p.nextToken()
	}

	// Check for array type (PostgreSQL): INT[], VARCHAR[]
	// Note: This is a simplified check, real parsing might differ by dialect
	// For now, we'll skip array parsing as it's complex

	return dataType, nil
}

// parseProcedureOptions parses procedure-specific options
func (p *Parser) parseProcedureOptions(stmt *CreateProcedureStatement) error {
	// LANGUAGE (PostgreSQL)
	if p.curTokenIs(lexer.LANGUAGE) {
		p.nextToken()
		if p.curTokenIs(lexer.IDENT) || p.curTokenIs(lexer.SQL) || p.curTokenIs(lexer.PLPGSQL) {
			stmt.Language = p.curToken.Literal
			p.nextToken()
		}
	}

	// SECURITY DEFINER/INVOKER (PostgreSQL)
	if p.curTokenIs(lexer.SECURITY) {
		p.nextToken()
		if p.curTokenIs(lexer.DEFINER) {
			stmt.SecurityDefiner = true
			p.nextToken()
		} else if p.curTokenIs(lexer.INVOKER) {
			stmt.SecurityDefiner = false
			p.nextToken()
		}
	}

	// Additional options can be added here based on dialect
	return nil
}

// parseFunctionOptions parses function-specific options
func (p *Parser) parseFunctionOptions(stmt *CreateFunctionStatement) error {
	// LANGUAGE (PostgreSQL)
	if p.curTokenIs(lexer.LANGUAGE) {
		p.nextToken()
		if p.curTokenIs(lexer.IDENT) || p.curTokenIs(lexer.SQL) || p.curTokenIs(lexer.PLPGSQL) {
			stmt.Language = p.curToken.Literal
			p.nextToken()
		}
	}

	// DETERMINISTIC (MySQL)
	if p.curTokenIs(lexer.DETERMINISTIC) {
		stmt.Deterministic = true
		p.nextToken()
	}

	// SECURITY DEFINER/INVOKER (PostgreSQL)
	if p.curTokenIs(lexer.SECURITY) {
		p.nextToken()
		if p.curTokenIs(lexer.DEFINER) {
			stmt.SecurityDefiner = true
			p.nextToken()
		} else if p.curTokenIs(lexer.INVOKER) {
			stmt.SecurityDefiner = false
			p.nextToken()
		}
	}

	// SQL DATA ACCESS (MySQL): CONTAINS SQL, READS SQL DATA, MODIFIES SQL DATA, NO SQL
	if p.curTokenIs(lexer.CONTAINS) || p.curTokenIs(lexer.READS) || p.curTokenIs(lexer.MODIFIES) || p.curTokenIs(lexer.NO) {
		access := p.curToken.Literal
		p.nextToken()
		if p.curTokenIs(lexer.SQL) {
			access += " SQL"
			p.nextToken()
			if p.curTokenIs(lexer.DATA) {
				access += " DATA"
				p.nextToken()
			}
		}
		stmt.Options["sql_data_access"] = access
	}

	return nil
}

// parseProcedureBody parses the procedure/function body (AS/IS BEGIN ... END)
func (p *Parser) parseProcedureBody() (*ProcedureBody, error) {
	body := &ProcedureBody{
		Statements: make([]Statement, 0),
		Variables:  make([]*VariableDecl, 0),
		Cursors:    make([]*CursorDecl, 0),
	}

	// AS or IS keyword (Oracle, PostgreSQL)
	if p.curTokenIs(lexer.AS) || p.curTokenIs(lexer.IS) {
		p.nextToken()
	}

	// PostgreSQL: Can use $$ delimiter or BEGIN
	// MySQL: BEGIN

	// BEGIN block
	if !p.curTokenIs(lexer.BEGIN) {
		return nil, fmt.Errorf("expected BEGIN for procedure body, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse declarations (DECLARE variables, cursors)
	for p.curTokenIs(lexer.DECLARE) {
		p.nextToken()

		// Check if it's a cursor or variable
		// Save position to peek ahead
		name := p.curToken.Literal
		p.nextToken()

		if p.curTokenIs(lexer.CURSOR) {
			// It's a cursor declaration
			cursor := &CursorDecl{
				Name: name,
			}
			p.nextToken() // Consume CURSOR

			// FOR or IS
			if p.curTokenIs(lexer.FOR) || p.curTokenIs(lexer.IS) {
				p.nextToken()
			}

			// Parse SELECT statement
			if p.curTokenIs(lexer.SELECT) {
				selectStmt, err := p.parseSelectStatement()
				if err != nil {
					return nil, err
				}
				cursor.Query = selectStmt
			}

			body.Cursors = append(body.Cursors, cursor)

			// Consume semicolon if present
			if p.curTokenIs(lexer.SEMICOLON) {
				p.nextToken()
			}
		} else {
			// It's a variable declaration
			variable := &VariableDecl{
				Name: name,
			}

			// Parse data type
			dataType, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			variable.DataType = dataType

			// Optional DEFAULT/= value
			if p.curTokenIs(lexer.DEFAULT) || p.curTokenIs(lexer.ASSIGN) {
				p.nextToken()
				defaultValue, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				variable.Default = defaultValue
			}

			body.Variables = append(body.Variables, variable)

			// Consume semicolon if present
			if p.curTokenIs(lexer.SEMICOLON) {
				p.nextToken()
			}
		}
	}

	// Parse statements until END
	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.EOF) {
		stmt, err := p.parseProcedureStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			body.Statements = append(body.Statements, stmt)
		}

		// Consume semicolon if present
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	// Consume END
	if !p.curTokenIs(lexer.END) {
		return nil, fmt.Errorf("expected END for procedure body, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Optional procedure name after END
	if p.curTokenIs(lexer.IDENT) {
		p.nextToken()
	}

	// Consume final semicolon if present
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return body, nil
}

// parseProcedureStatement parses a single statement within a procedure body
func (p *Parser) parseProcedureStatement() (Statement, error) {
	switch p.curToken.Type {
	case lexer.SELECT, lexer.INSERT, lexer.UPDATE, lexer.DELETE:
		// Regular SQL statements
		return p.ParseStatement()

	case lexer.SET:
		// Assignment: SET var = value
		return p.parseAssignmentStatement()

	case lexer.IF:
		// IF statement
		return p.parseIfStatement()

	case lexer.WHILE:
		// WHILE loop
		return p.parseWhileStatement()

	case lexer.LOOP:
		// LOOP statement
		return p.parseLoopStatement()

	case lexer.FOR:
		// FOR loop
		return p.parseForStatement()

	case lexer.CASE:
		// CASE statement
		return p.parseCaseStatement()

	case lexer.RETURN:
		// RETURN statement
		return p.parseReturnStatement()

	case lexer.OPEN:
		// OPEN cursor
		return p.parseOpenCursorStatement()

	case lexer.FETCH:
		// FETCH cursor
		return p.parseFetchStatement()

	case lexer.CLOSE:
		// CLOSE cursor
		return p.parseCloseStatement()

	case lexer.EXIT:
		// EXIT loop
		return p.parseExitStatement()

	case lexer.CONTINUE, lexer.ITERATE:
		// CONTINUE loop
		return p.parseContinueStatement()

	default:
		return nil, fmt.Errorf("unexpected statement in procedure body: %s", p.curToken.Literal)
	}
}

// parseAssignmentStatement parses SET var = value
func (p *Parser) parseAssignmentStatement() (Statement, error) {
	stmt := &AssignmentStatement{}

	// Consume SET
	p.nextToken()

	// Variable name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected variable name, got %s", p.curToken.Literal)
	}
	stmt.Variable = p.curToken.Literal
	p.nextToken()

	// = or :=
	if !p.curTokenIs(lexer.ASSIGN) {
		return nil, fmt.Errorf("expected =, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Value expression
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	stmt.Value = value

	return stmt, nil
}

// parseReturnStatement parses RETURN expression
func (p *Parser) parseReturnStatement() (Statement, error) {
	stmt := &ReturnStatement{}

	// Consume RETURN
	p.nextToken()

	// Optional return value
	if !p.curTokenIs(lexer.SEMICOLON) && !p.curTokenIs(lexer.END) {
		value, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Value = value
	}

	return stmt, nil
}

// Placeholder implementations for control flow statements
// These will be implemented in detail as needed

func (p *Parser) parseIfStatement() (Statement, error) {
	// TODO: Implement IF...THEN...ELSE parsing
	return nil, fmt.Errorf("IF statement parsing not yet implemented")
}

func (p *Parser) parseWhileStatement() (Statement, error) {
	// TODO: Implement WHILE loop parsing
	return nil, fmt.Errorf("WHILE statement parsing not yet implemented")
}

func (p *Parser) parseLoopStatement() (Statement, error) {
	// TODO: Implement LOOP parsing
	return nil, fmt.Errorf("LOOP statement parsing not yet implemented")
}

func (p *Parser) parseForStatement() (Statement, error) {
	// TODO: Implement FOR loop parsing
	return nil, fmt.Errorf("FOR statement parsing not yet implemented")
}

func (p *Parser) parseCaseStatement() (Statement, error) {
	// TODO: Implement CASE statement parsing
	return nil, fmt.Errorf("CASE statement parsing not yet implemented")
}

func (p *Parser) parseOpenCursorStatement() (Statement, error) {
	stmt := &OpenCursorStatement{}
	p.nextToken() // Consume OPEN
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected cursor name")
	}
	stmt.CursorName = p.curToken.Literal
	p.nextToken()
	return stmt, nil
}

func (p *Parser) parseFetchStatement() (Statement, error) {
	stmt := &FetchStatement{}
	p.nextToken() // Consume FETCH
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected cursor name")
	}
	stmt.CursorName = p.curToken.Literal
	p.nextToken()
	// TODO: Parse INTO variables
	return stmt, nil
}

func (p *Parser) parseCloseStatement() (Statement, error) {
	stmt := &CloseStatement{}
	p.nextToken() // Consume CLOSE
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected cursor name")
	}
	stmt.CursorName = p.curToken.Literal
	p.nextToken()
	return stmt, nil
}

func (p *Parser) parseExitStatement() (Statement, error) {
	stmt := &ExitStatement{}
	p.nextToken() // Consume EXIT
	// TODO: Parse optional label and WHEN condition
	return stmt, nil
}

func (p *Parser) parseContinueStatement() (Statement, error) {
	stmt := &ContinueStatement{}
	p.nextToken() // Consume CONTINUE/ITERATE
	// TODO: Parse optional label and WHEN condition
	return stmt, nil
}
