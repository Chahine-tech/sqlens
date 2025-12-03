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
			// DATA is not a reserved keyword, so check if it's an IDENT with literal "DATA"
			if p.curTokenIs(lexer.IDENT) && p.curToken.Literal == "DATA" {
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

	// Parse declarations (DECLARE variables, cursors, handlers)
	for p.curTokenIs(lexer.DECLARE) {
		p.nextToken()

		// Check if it's a handler declaration (MySQL)
		// DECLARE CONTINUE|EXIT|UNDO HANDLER FOR condition statement
		if p.curToken.Literal == "CONTINUE" || p.curToken.Literal == "EXIT" || p.curToken.Literal == "UNDO" {
			handler, err := p.parseHandlerDeclaration()
			if err != nil {
				return nil, err
			}
			// For now, store handlers as statements
			body.Statements = append(body.Statements, handler)

			// Consume semicolon if present
			if p.curTokenIs(lexer.SEMICOLON) {
				p.nextToken()
			}
			continue
		}

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

	// Parse statements until END or EXCEPTION
	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.EXCEPTION) && !p.curTokenIs(lexer.EOF) {
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

	// Check for EXCEPTION block (PostgreSQL/Oracle)
	if p.curTokenIs(lexer.EXCEPTION) {
		exceptionBlock, err := p.parseExceptionBlock()
		if err != nil {
			return nil, err
		}
		body.ExceptionBlock = exceptionBlock
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

	case lexer.DEALLOCATE:
		// DEALLOCATE cursor
		return p.parseDeallocateStatement()

	case lexer.EXIT:
		// EXIT loop
		return p.parseExitStatement()

	case lexer.CONTINUE, lexer.ITERATE:
		// CONTINUE loop
		return p.parseContinueStatement()

	case lexer.REPEAT:
		// REPEAT...UNTIL loop (MySQL)
		return p.parseRepeatStatement()

	case lexer.TRY:
		// TRY...CATCH (SQL Server)
		return p.parseTryStatement()

	case lexer.RAISE:
		// RAISE (PostgreSQL/Oracle)
		return p.parseRaiseStatement()

	case lexer.THROW:
		// THROW (SQL Server)
		return p.parseThrowStatement()

	case lexer.SIGNAL:
		// SIGNAL (MySQL)
		return p.parseSignalStatement()

	case lexer.BEGIN, lexer.START:
		// BEGIN TRY...CATCH (SQL Server) - nested OR BEGIN TRANSACTION
		if p.curTokenIs(lexer.BEGIN) && p.peekTokenIs(lexer.TRY) {
			p.nextToken() // consume BEGIN
			return p.parseTryStatement()
		}
		// BEGIN/START TRANSACTION
		return p.parseBeginTransaction()

	case lexer.COMMIT:
		// COMMIT transaction
		return p.parseCommit()

	case lexer.ROLLBACK:
		// ROLLBACK transaction
		return p.parseRollback()

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
	stmt := &IfStatement{}

	// Consume IF
	p.nextToken()

	// Parse condition
	condition, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to parse IF condition: %w", err)
	}
	stmt.Condition = condition

	// Expect THEN
	if !p.curTokenIs(lexer.THEN) {
		return nil, fmt.Errorf("expected THEN after IF condition, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse THEN block
	for !p.curTokenIs(lexer.ELSEIF) && !p.curTokenIs(lexer.ELSIF) && !p.curTokenIs(lexer.ELSE) && !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.ENDIF) && !p.curTokenIs(lexer.EOF) {
		blockStmt, err := p.parseProcedureStatement()
		if err != nil {
			return nil, err
		}
		if blockStmt != nil {
			stmt.ThenBlock = append(stmt.ThenBlock, blockStmt)
		}

		// Consume semicolon if present
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	// Parse ELSEIF/ELSIF blocks
	for p.curTokenIs(lexer.ELSEIF) || p.curTokenIs(lexer.ELSIF) {
		p.nextToken() // Consume ELSEIF/ELSIF

		elseIfBlock := &ElseIfBlock{}

		// Parse condition
		condition, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse ELSEIF condition: %w", err)
		}
		elseIfBlock.Condition = condition

		// Expect THEN
		if !p.curTokenIs(lexer.THEN) {
			return nil, fmt.Errorf("expected THEN after ELSEIF condition, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Parse ELSEIF block
		for !p.curTokenIs(lexer.ELSEIF) && !p.curTokenIs(lexer.ELSIF) && !p.curTokenIs(lexer.ELSE) && !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.ENDIF) && !p.curTokenIs(lexer.EOF) {
			blockStmt, err := p.parseProcedureStatement()
			if err != nil {
				return nil, err
			}
			if blockStmt != nil {
				elseIfBlock.Block = append(elseIfBlock.Block, blockStmt)
			}

			// Consume semicolon if present
			if p.curTokenIs(lexer.SEMICOLON) {
				p.nextToken()
			}
		}

		stmt.ElseIfList = append(stmt.ElseIfList, elseIfBlock)
	}

	// Parse ELSE block (optional)
	if p.curTokenIs(lexer.ELSE) {
		p.nextToken() // Consume ELSE

		for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.ENDIF) && !p.curTokenIs(lexer.EOF) {
			blockStmt, err := p.parseProcedureStatement()
			if err != nil {
				return nil, err
			}
			if blockStmt != nil {
				stmt.ElseBlock = append(stmt.ElseBlock, blockStmt)
			}

			// Consume semicolon if present
			if p.curTokenIs(lexer.SEMICOLON) {
				p.nextToken()
			}
		}
	}

	// Expect END IF or ENDIF
	if p.curTokenIs(lexer.END) {
		p.nextToken()
		// Optional IF keyword after END
		if p.curTokenIs(lexer.IF) {
			p.nextToken()
		}
	} else if p.curTokenIs(lexer.ENDIF) {
		p.nextToken()
	} else {
		return nil, fmt.Errorf("expected END IF, got %s", p.curToken.Literal)
	}

	return stmt, nil
}

func (p *Parser) parseWhileStatement() (Statement, error) {
	stmt := &WhileStatement{}

	// Consume WHILE
	p.nextToken()

	// Parse condition
	condition, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to parse WHILE condition: %w", err)
	}
	stmt.Condition = condition

	// Expect DO
	if !p.curTokenIs(lexer.DO) {
		return nil, fmt.Errorf("expected DO after WHILE condition, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse loop body
	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.ENDWHILE) && !p.curTokenIs(lexer.EOF) {
		blockStmt, err := p.parseProcedureStatement()
		if err != nil {
			return nil, err
		}
		if blockStmt != nil {
			stmt.Block = append(stmt.Block, blockStmt)
		}

		// Consume semicolon if present
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	// Expect END WHILE or ENDWHILE
	if p.curTokenIs(lexer.END) {
		p.nextToken()
		// Optional WHILE keyword after END
		if p.curTokenIs(lexer.WHILE) {
			p.nextToken()
		}
	} else if p.curTokenIs(lexer.ENDWHILE) {
		p.nextToken()
	} else {
		return nil, fmt.Errorf("expected END WHILE, got %s", p.curToken.Literal)
	}

	return stmt, nil
}

func (p *Parser) parseLoopStatement() (Statement, error) {
	stmt := &LoopStatement{}

	// Check for optional label before LOOP
	// Example: <<my_loop>> LOOP ... END LOOP;
	// For now, we skip label parsing

	// Consume LOOP
	p.nextToken()

	// Parse loop body
	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.ENDLOOP) && !p.curTokenIs(lexer.EOF) {
		blockStmt, err := p.parseProcedureStatement()
		if err != nil {
			return nil, err
		}
		if blockStmt != nil {
			stmt.Block = append(stmt.Block, blockStmt)
		}

		// Consume semicolon if present
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	// Expect END LOOP or ENDLOOP
	if p.curTokenIs(lexer.END) {
		p.nextToken()
		// Optional LOOP keyword after END
		if p.curTokenIs(lexer.LOOP) {
			p.nextToken()
		}
	} else if p.curTokenIs(lexer.ENDLOOP) {
		p.nextToken()
	} else {
		return nil, fmt.Errorf("expected END LOOP, got %s", p.curToken.Literal)
	}

	return stmt, nil
}

func (p *Parser) parseForStatement() (Statement, error) {
	stmt := &ForStatement{}

	// Consume FOR
	p.nextToken()

	// Parse loop variable
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected loop variable name, got %s", p.curToken.Literal)
	}
	stmt.Variable = p.curToken.Literal
	p.nextToken()

	// Expect IN
	if !p.curTokenIs(lexer.IN) {
		return nil, fmt.Errorf("expected IN after loop variable, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Check for REVERSE (PostgreSQL)
	if p.curTokenIs(lexer.REVERSE) {
		stmt.IsReverse = true
		p.nextToken()
	}

	// Parse start value
	start, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to parse FOR start value: %w", err)
	}
	stmt.Start = start

	// Expect .. (range operator) - may be tokenized as two DOT tokens
	if p.curTokenIs(lexer.DOT) {
		p.nextToken()
		if !p.curTokenIs(lexer.DOT) {
			return nil, fmt.Errorf("expected .. for range, got single .")
		}
		p.nextToken()
	} else if p.curToken.Literal == ".." {
		p.nextToken()
	} else {
		return nil, fmt.Errorf("expected .. for range, got %s", p.curToken.Literal)
	}

	// Parse end value
	end, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to parse FOR end value: %w", err)
	}
	stmt.End = end

	// Optional BY step (some dialects)
	if p.curTokenIs(lexer.BY) {
		p.nextToken()
		step, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse FOR step value: %w", err)
		}
		stmt.Step = step
	}

	// Expect LOOP
	if !p.curTokenIs(lexer.LOOP) {
		return nil, fmt.Errorf("expected LOOP after FOR range, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse loop body
	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.ENDFOR) && !p.curTokenIs(lexer.EOF) {
		blockStmt, err := p.parseProcedureStatement()
		if err != nil {
			return nil, err
		}
		if blockStmt != nil {
			stmt.Block = append(stmt.Block, blockStmt)
		}

		// Consume semicolon if present
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	// Expect END LOOP or ENDFOR
	if p.curTokenIs(lexer.END) {
		p.nextToken()
		// Optional LOOP/FOR keyword after END
		if p.curTokenIs(lexer.LOOP) || p.curTokenIs(lexer.FOR) {
			p.nextToken()
		}
	} else if p.curTokenIs(lexer.ENDFOR) {
		p.nextToken()
	} else {
		return nil, fmt.Errorf("expected END LOOP, got %s", p.curToken.Literal)
	}

	return stmt, nil
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

	// Check for fetch direction
	switch p.curToken.Type {
	case lexer.NEXT:
		stmt.Direction = "NEXT"
		p.nextToken()
	case lexer.PRIOR:
		stmt.Direction = "PRIOR"
		p.nextToken()
	case lexer.FIRST:
		stmt.Direction = "FIRST"
		p.nextToken()
	case lexer.LAST:
		stmt.Direction = "LAST"
		p.nextToken()
	case lexer.ABSOLUTE:
		stmt.Direction = "ABSOLUTE"
		p.nextToken()
		// Expect number
		if p.curToken.Type == lexer.NUMBER {
			// Parse count (simple conversion, real implementation might need more robust parsing)
			stmt.Count = 1 // Placeholder, should parse the actual number
			p.nextToken()
		}
	case lexer.RELATIVE:
		stmt.Direction = "RELATIVE"
		p.nextToken()
		// Expect number
		if p.curToken.Type == lexer.NUMBER {
			stmt.Count = 1 // Placeholder
			p.nextToken()
		}
	}

	// Optional FROM keyword (PostgreSQL)
	if p.curTokenIs(lexer.FROM) || p.curTokenIs(lexer.IN) {
		p.nextToken()
	}

	// Parse cursor name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected cursor name, got %s", p.curToken.Literal)
	}
	stmt.CursorName = p.curToken.Literal
	p.nextToken()

	// Optional INTO clause
	if p.curTokenIs(lexer.INTO) {
		p.nextToken()
		for {
			if p.curToken.Type != lexer.IDENT {
				return nil, fmt.Errorf("expected variable name, got %s", p.curToken.Literal)
			}
			stmt.Variables = append(stmt.Variables, p.curToken.Literal)
			p.nextToken()

			if !p.curTokenIs(lexer.COMMA) {
				break
			}
			p.nextToken()
		}
	}

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

func (p *Parser) parseDeallocateStatement() (Statement, error) {
	stmt := &DeallocateStatement{}
	p.nextToken() // Consume DEALLOCATE

	// Optional PREPARE keyword (some dialects use DEALLOCATE PREPARE)
	if p.curToken.Type == lexer.IDENT && p.curToken.Literal == "PREPARE" {
		p.nextToken()
	}

	// Parse cursor/statement name
	if !p.curTokenIs(lexer.IDENT) {
		return nil, fmt.Errorf("expected cursor name, got %s", p.curToken.Literal)
	}
	stmt.CursorName = p.curToken.Literal
	p.nextToken()

	return stmt, nil
}

func (p *Parser) parseExitStatement() (Statement, error) {
	stmt := &ExitStatement{}
	p.nextToken() // Consume EXIT

	// Optional label
	if p.curTokenIs(lexer.IDENT) {
		stmt.Label = p.curToken.Literal
		p.nextToken()
	}

	// Optional WHEN condition
	if p.curTokenIs(lexer.WHEN) {
		p.nextToken()
		condition, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse EXIT WHEN condition: %w", err)
		}
		stmt.Condition = condition
	}

	return stmt, nil
}

func (p *Parser) parseContinueStatement() (Statement, error) {
	stmt := &ContinueStatement{}
	p.nextToken() // Consume CONTINUE/ITERATE

	// Optional label
	if p.curTokenIs(lexer.IDENT) {
		stmt.Label = p.curToken.Literal
		p.nextToken()
	}

	// Optional WHEN condition
	if p.curTokenIs(lexer.WHEN) {
		p.nextToken()
		condition, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse CONTINUE WHEN condition: %w", err)
		}
		stmt.Condition = condition
	}

	return stmt, nil
}

func (p *Parser) parseRepeatStatement() (Statement, error) {
	stmt := &RepeatStatement{}

	// Consume REPEAT
	p.nextToken()

	// Parse loop body
	for !p.curTokenIs(lexer.UNTIL) && !p.curTokenIs(lexer.EOF) {
		blockStmt, err := p.parseProcedureStatement()
		if err != nil {
			return nil, err
		}
		if blockStmt != nil {
			stmt.Body = append(stmt.Body, blockStmt)
		}

		// Consume semicolon if present
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	// Expect UNTIL
	if !p.curTokenIs(lexer.UNTIL) {
		return nil, fmt.Errorf("expected UNTIL after REPEAT body, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse UNTIL condition
	condition, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to parse REPEAT UNTIL condition: %w", err)
	}
	stmt.Condition = condition

	return stmt, nil
}

// ====================================================================
// Exception Handling Parsers
// ====================================================================

// parseTryStatement parses TRY...CATCH (SQL Server)
func (p *Parser) parseTryStatement() (Statement, error) {
	stmt := &TryStatement{}

	// Expect BEGIN TRY
	if !p.curTokenIs(lexer.TRY) {
		return nil, fmt.Errorf("expected TRY, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse TRY block statements
	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.EOF) {
		blockStmt, err := p.parseProcedureStatement()
		if err != nil {
			return nil, err
		}
		if blockStmt != nil {
			stmt.TryBlock = append(stmt.TryBlock, blockStmt)
		}

		// Consume semicolon if present
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	// Expect END TRY
	if !p.curTokenIs(lexer.END) {
		return nil, fmt.Errorf("expected END TRY, got %s", p.curToken.Literal)
	}
	p.nextToken()

	if !p.curTokenIs(lexer.TRY) {
		return nil, fmt.Errorf("expected TRY after END, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Expect BEGIN CATCH
	if !p.curTokenIs(lexer.BEGIN) {
		return nil, fmt.Errorf("expected BEGIN CATCH, got %s", p.curToken.Literal)
	}
	p.nextToken()

	if !p.curTokenIs(lexer.CATCH) {
		return nil, fmt.Errorf("expected CATCH, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse CATCH block
	stmt.CatchBlock = &CatchBlock{}
	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.EOF) {
		blockStmt, err := p.parseProcedureStatement()
		if err != nil {
			return nil, err
		}
		if blockStmt != nil {
			stmt.CatchBlock.Body = append(stmt.CatchBlock.Body, blockStmt)
		}

		// Consume semicolon if present
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}
	}

	// Expect END CATCH
	if !p.curTokenIs(lexer.END) {
		return nil, fmt.Errorf("expected END CATCH, got %s", p.curToken.Literal)
	}
	p.nextToken()

	if !p.curTokenIs(lexer.CATCH) {
		return nil, fmt.Errorf("expected CATCH after END, got %s", p.curToken.Literal)
	}
	p.nextToken()

	return stmt, nil
}

// parseRaiseStatement parses RAISE (PostgreSQL/Oracle)
func (p *Parser) parseRaiseStatement() (Statement, error) {
	stmt := &RaiseStatement{}

	// Consume RAISE
	p.nextToken()

	// Optional level (EXCEPTION, NOTICE, WARNING, etc.)
	if p.curTokenIs(lexer.EXCEPTION) || p.curToken.Literal == "NOTICE" ||
		p.curToken.Literal == "WARNING" || p.curToken.Literal == "INFO" {
		stmt.Level = p.curToken.Literal
		p.nextToken()
	}

	// Parse message (if present)
	if !p.curTokenIs(lexer.SEMICOLON) && !p.curTokenIs(lexer.EOF) {
		message, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse RAISE message: %w", err)
		}
		stmt.Message = message
	}

	return stmt, nil
}

// parseThrowStatement parses THROW (SQL Server)
func (p *Parser) parseThrowStatement() (Statement, error) {
	stmt := &ThrowStatement{}

	// Consume THROW
	p.nextToken()

	// Check if it's a re-throw (no parameters)
	if p.curTokenIs(lexer.SEMICOLON) || p.curTokenIs(lexer.EOF) {
		return stmt, nil
	}

	// Parse error number
	if !p.curTokenIs(lexer.NUMBER) {
		return nil, fmt.Errorf("expected error number after THROW, got %s", p.curToken.Literal)
	}
	// Convert to int (simplified - should handle errors)
	stmt.ErrorNumber = 50000 // Placeholder
	p.nextToken()

	// Expect comma
	if !p.curTokenIs(lexer.COMMA) {
		return nil, fmt.Errorf("expected comma after error number, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse message
	message, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to parse THROW message: %w", err)
	}
	stmt.Message = message

	// Expect comma
	if !p.curTokenIs(lexer.COMMA) {
		return nil, fmt.Errorf("expected comma after message, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse state
	if !p.curTokenIs(lexer.NUMBER) {
		return nil, fmt.Errorf("expected state number, got %s", p.curToken.Literal)
	}
	stmt.State = 1 // Placeholder
	p.nextToken()

	return stmt, nil
}

// parseSignalStatement parses SIGNAL (MySQL)
func (p *Parser) parseSignalStatement() (Statement, error) {
	stmt := &SignalStatement{
		Properties: make(map[string]string),
	}

	// Consume SIGNAL
	p.nextToken()

	// Expect SQLSTATE
	if !p.curTokenIs(lexer.SQLSTATE) {
		return nil, fmt.Errorf("expected SQLSTATE after SIGNAL, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse SQLSTATE value (should be a string like '45000')
	if !p.curTokenIs(lexer.STRING) {
		return nil, fmt.Errorf("expected SQLSTATE value, got %s", p.curToken.Literal)
	}
	stmt.SqlState = p.curToken.Literal
	p.nextToken()

	// Optional SET clause
	if p.curTokenIs(lexer.SET) {
		p.nextToken()

		// Parse properties (MESSAGE_TEXT = 'value', etc.)
		for {
			// Property name
			if !p.curTokenIs(lexer.IDENT) && !p.curTokenIs(lexer.MESSAGE_TEXT) {
				break
			}
			propName := p.curToken.Literal
			p.nextToken()

			// Expect =
			if !p.curTokenIs(lexer.ASSIGN) {
				return nil, fmt.Errorf("expected = after property name, got %s", p.curToken.Literal)
			}
			p.nextToken()

			// Property value (can be string or number)
			if !p.curTokenIs(lexer.STRING) && !p.curTokenIs(lexer.NUMBER) {
				return nil, fmt.Errorf("expected property value, got %s", p.curToken.Literal)
			}
			propValue := p.curToken.Literal
			p.nextToken()

			stmt.Properties[propName] = propValue

			// Check for comma (more properties)
			if p.curTokenIs(lexer.COMMA) {
				p.nextToken()
			} else {
				break
			}
		}
	}

	return stmt, nil
}

// parseExceptionBlock parses EXCEPTION...WHEN block (PostgreSQL/Oracle)
// EXCEPTION
//
//	WHEN exception_name THEN
//	    statements
//	WHEN OTHERS THEN
//	    statements
func (p *Parser) parseExceptionBlock() (*ExceptionBlock, error) {
	block := &ExceptionBlock{
		WhenClauses: make([]*WhenExceptionClause, 0),
	}

	// Consume EXCEPTION keyword
	if !p.curTokenIs(lexer.EXCEPTION) {
		return nil, fmt.Errorf("expected EXCEPTION, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Parse WHEN clauses
	for p.curTokenIs(lexer.WHEN) {
		p.nextToken()

		whenClause := &WhenExceptionClause{
			Body: make([]Statement, 0),
		}

		// Exception name (can be OTHERS, SQLEXCEPTION, or specific exception like division_by_zero)
		if !p.curTokenIs(lexer.OTHERS) && !p.curTokenIs(lexer.SQLEXCEPTION) && !p.curTokenIs(lexer.IDENT) {
			return nil, fmt.Errorf("expected exception name after WHEN, got %s", p.curToken.Literal)
		}
		whenClause.ExceptionName = p.curToken.Literal
		p.nextToken()

		// Expect THEN
		if !p.curTokenIs(lexer.THEN) {
			return nil, fmt.Errorf("expected THEN after exception name, got %s", p.curToken.Literal)
		}
		p.nextToken()

		// Parse handler statements until next WHEN or END
		for !p.curTokenIs(lexer.WHEN) && !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.EOF) {
			stmt, err := p.parseProcedureStatement()
			if err != nil {
				return nil, err
			}
			if stmt != nil {
				whenClause.Body = append(whenClause.Body, stmt)
			}

			// Consume semicolon if present
			if p.curTokenIs(lexer.SEMICOLON) {
				p.nextToken()
			}
		}

		block.WhenClauses = append(block.WhenClauses, whenClause)
	}

	return block, nil
}

// parseHandlerDeclaration parses DECLARE HANDLER (MySQL)
// DECLARE handler_type HANDLER FOR condition_value statement
// handler_type: CONTINUE | EXIT | UNDO
// condition_value: SQLEXCEPTION | SQLWARNING | NOT FOUND | SQLSTATE 'value' | error_code
func (p *Parser) parseHandlerDeclaration() (*HandlerDeclaration, error) {
	stmt := &HandlerDeclaration{
		Body: make([]Statement, 0),
	}

	// Handler type (CONTINUE, EXIT, UNDO)
	if p.curToken.Literal != "CONTINUE" && p.curToken.Literal != "EXIT" && p.curToken.Literal != "UNDO" {
		return nil, fmt.Errorf("expected handler type (CONTINUE, EXIT, UNDO), got %s", p.curToken.Literal)
	}
	stmt.HandlerType = p.curToken.Literal
	p.nextToken()

	// Expect HANDLER keyword
	if !p.curTokenIs(lexer.HANDLER) {
		return nil, fmt.Errorf("expected HANDLER, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Expect FOR
	if !p.curTokenIs(lexer.FOR) {
		return nil, fmt.Errorf("expected FOR, got %s", p.curToken.Literal)
	}
	p.nextToken()

	// Condition value
	if p.curTokenIs(lexer.SQLEXCEPTION) {
		stmt.Condition = "SQLEXCEPTION"
		p.nextToken()
	} else if p.curTokenIs(lexer.SQLWARNING) {
		stmt.Condition = "SQLWARNING"
		p.nextToken()
	} else if p.curTokenIs(lexer.NOT) {
		// NOT FOUND
		p.nextToken()
		if !p.curTokenIs(lexer.FOUND) {
			return nil, fmt.Errorf("expected FOUND after NOT, got %s", p.curToken.Literal)
		}
		stmt.Condition = "NOT FOUND"
		p.nextToken()
	} else if p.curTokenIs(lexer.SQLSTATE) {
		// SQLSTATE 'value'
		p.nextToken()
		if !p.curTokenIs(lexer.STRING) {
			return nil, fmt.Errorf("expected SQLSTATE value, got %s", p.curToken.Literal)
		}
		stmt.Condition = "SQLSTATE " + p.curToken.Literal
		p.nextToken()
	} else if p.curTokenIs(lexer.NUMBER) {
		// MySQL error code (e.g., 1062 for duplicate key)
		stmt.Condition = p.curToken.Literal
		p.nextToken()
	} else {
		return nil, fmt.Errorf("expected condition value (SQLEXCEPTION, SQLWARNING, NOT FOUND, SQLSTATE, or error code), got %s", p.curToken.Literal)
	}

	// Handler body - can be a single statement or BEGIN...END block
	if p.curTokenIs(lexer.BEGIN) {
		p.nextToken()

		// Parse statements until END
		for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.EOF) {
			bodyStmt, err := p.parseProcedureStatement()
			if err != nil {
				return nil, err
			}
			if bodyStmt != nil {
				stmt.Body = append(stmt.Body, bodyStmt)
			}

			// Consume semicolon if present
			if p.curTokenIs(lexer.SEMICOLON) {
				p.nextToken()
			}
		}

		// Consume END
		if !p.curTokenIs(lexer.END) {
			return nil, fmt.Errorf("expected END for handler body, got %s", p.curToken.Literal)
		}
		p.nextToken()
	} else {
		// Single statement
		bodyStmt, err := p.parseProcedureStatement()
		if err != nil {
			return nil, err
		}
		if bodyStmt != nil {
			stmt.Body = append(stmt.Body, bodyStmt)
		}
	}

	return stmt, nil
}
