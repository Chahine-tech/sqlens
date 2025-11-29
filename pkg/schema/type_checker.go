package schema

import (
	"fmt"

	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// TypeChecker performs type checking on SQL expressions
type TypeChecker struct {
	schema    *Schema
	validator *Validator
}

// NewTypeChecker creates a new type checker
func NewTypeChecker(schema *Schema) *TypeChecker {
	return &TypeChecker{
		schema:    schema,
		validator: NewValidator(schema),
	}
}

// CheckStatement performs type checking on a statement
func (tc *TypeChecker) CheckStatement(stmt parser.Statement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	switch s := stmt.(type) {
	case *parser.SelectStatement:
		errors = append(errors, tc.checkSelectStatement(s)...)
	case *parser.InsertStatement:
		errors = append(errors, tc.checkInsertStatement(s)...)
	case *parser.UpdateStatement:
		errors = append(errors, tc.checkUpdateStatement(s)...)
	}

	return errors
}

// checkSelectStatement checks types in a SELECT statement
func (tc *TypeChecker) checkSelectStatement(stmt *parser.SelectStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	// Check WHERE clause types
	if stmt.Where != nil {
		errors = append(errors, tc.checkBooleanExpression(stmt.Where, stmt)...)
	}

	// Check HAVING clause types
	if stmt.Having != nil {
		errors = append(errors, tc.checkBooleanExpression(stmt.Having, stmt)...)
	}

	// Check JOIN conditions
	for _, join := range stmt.Joins {
		if join.Condition != nil {
			errors = append(errors, tc.checkBooleanExpression(join.Condition, stmt)...)
		}
	}

	return errors
}

// checkInsertStatement checks types in an INSERT statement
func (tc *TypeChecker) checkInsertStatement(stmt *parser.InsertStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	table, ok := tc.schema.GetTable(stmt.Table.Name)
	if !ok {
		return errors // Table validation already done
	}

	// If columns are specified, check type compatibility for each value
	if len(stmt.Columns) > 0 && len(stmt.Values) > 0 {
		for _, valueRow := range stmt.Values {
			if len(valueRow) != len(stmt.Columns) {
				errors = append(errors, &ValidationError{
					Type:    "COLUMN_COUNT_MISMATCH",
					Message: fmt.Sprintf("Column count mismatch: %d columns specified, %d values provided", len(stmt.Columns), len(valueRow)),
					Table:   stmt.Table.Name,
				})
				continue
			}

			for i, colName := range stmt.Columns {
				col, ok := table.GetColumn(colName)
				if !ok {
					continue // Column validation already done
				}

				// Check type compatibility
				valueType := tc.inferExpressionType(valueRow[i], nil)
				if valueType != nil && !col.DataType.IsCompatibleWith(valueType) {
					errors = append(errors, &ValidationError{
						Type:    "TYPE_MISMATCH",
						Message: fmt.Sprintf("Type mismatch for column '%s': expected %s, got %s", colName, col.DataType, valueType),
						Table:   stmt.Table.Name,
						Column:  colName,
					})
				}
			}
		}
	}

	return errors
}

// checkUpdateStatement checks types in an UPDATE statement
func (tc *TypeChecker) checkUpdateStatement(stmt *parser.UpdateStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	table, ok := tc.schema.GetTable(stmt.Table.Name)
	if !ok {
		return errors
	}

	// Check SET clause type compatibility
	for _, assignment := range stmt.Set {
		col, ok := table.GetColumn(assignment.Column)
		if !ok {
			continue
		}

		valueType := tc.inferExpressionType(assignment.Value, nil)
		if valueType != nil && !col.DataType.IsCompatibleWith(valueType) {
			errors = append(errors, &ValidationError{
				Type:    "TYPE_MISMATCH",
				Message: fmt.Sprintf("Type mismatch for column '%s': expected %s, got %s", assignment.Column, col.DataType, valueType),
				Table:   stmt.Table.Name,
				Column:  assignment.Column,
			})
		}
	}

	// Check WHERE clause
	if stmt.Where != nil {
		fakeSelect := &parser.SelectStatement{
			From: &parser.FromClause{
				Tables: []parser.TableReference{stmt.Table},
			},
		}
		errors = append(errors, tc.checkBooleanExpression(stmt.Where, fakeSelect)...)
	}

	return errors
}

// checkBooleanExpression checks if an expression is boolean
func (tc *TypeChecker) checkBooleanExpression(expr parser.Expression, stmt *parser.SelectStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	switch e := expr.(type) {
	case *parser.BinaryExpression:
		// Check if operator is a comparison operator
		comparisonOps := map[string]bool{
			"=": true, "!=": true, "<>": true, "<": true, ">": true, "<=": true, ">=": true,
			"AND": true, "OR": true, "LIKE": true, "IN": true, "IS": true,
		}

		if !comparisonOps[e.Operator] {
			errors = append(errors, &ValidationError{
				Type:    "NON_BOOLEAN_EXPRESSION",
				Message: fmt.Sprintf("Non-boolean operator '%s' used in boolean context", e.Operator),
			})
		}

		// Check type compatibility of operands
		leftType := tc.inferExpressionType(e.Left, stmt)
		rightType := tc.inferExpressionType(e.Right, stmt)

		if leftType != nil && rightType != nil {
			if !leftType.IsCompatibleWith(rightType) {
				errors = append(errors, &ValidationError{
					Type:    "TYPE_MISMATCH",
					Message: fmt.Sprintf("Type mismatch in comparison: %s vs %s", leftType, rightType),
				})
			}
		}

		// Recursively check sub-expressions
		errors = append(errors, tc.checkBooleanExpression(e.Left, stmt)...)
		errors = append(errors, tc.checkBooleanExpression(e.Right, stmt)...)

	case *parser.UnaryExpression:
		if e.Operator == "NOT" {
			errors = append(errors, tc.checkBooleanExpression(e.Operand, stmt)...)
		}

	case *parser.InExpression:
		// IN expressions are boolean
		exprType := tc.inferExpressionType(e.Expression, stmt)
		for _, val := range e.Values {
			valType := tc.inferExpressionType(val, stmt)
			if exprType != nil && valType != nil && !exprType.IsCompatibleWith(valType) {
				errors = append(errors, &ValidationError{
					Type:    "TYPE_MISMATCH",
					Message: fmt.Sprintf("Type mismatch in IN clause: %s vs %s", exprType, valType),
				})
			}
		}

	case *parser.ExistsExpression:
		// EXISTS is always boolean - check subquery
		if selectStmt, ok := e.Subquery.(*parser.SelectStatement); ok {
			errors = append(errors, tc.checkSelectStatement(selectStmt)...)
		}
	}

	return errors
}

// inferExpressionType infers the data type of an expression
func (tc *TypeChecker) inferExpressionType(expr parser.Expression, stmt *parser.SelectStatement) *DataType {
	switch e := expr.(type) {
	case *parser.Literal:
		// Infer type from literal value
		switch v := e.Value.(type) {
		case int, int64, int32:
			return &DataType{Name: "INT"}
		case float64, float32:
			return &DataType{Name: "FLOAT"}
		case string:
			return &DataType{Name: "VARCHAR", Length: len(v)}
		case bool:
			return &DataType{Name: "BOOLEAN"}
		case nil:
			return &DataType{Name: "NULL", Nullable: true}
		}

	case *parser.ColumnReference:
		// Look up column type from schema
		if e.Table != "" {
			if col, err := tc.schema.GetColumn(e.Table, e.Column); err == nil {
				return col.DataType
			}
		} else if stmt != nil && stmt.From != nil {
			// Try to find column in any table
			for _, tableRef := range stmt.From.Tables {
				if col, err := tc.schema.GetColumn(tableRef.Name, e.Column); err == nil {
					return col.DataType
				}
			}
		}

	case *parser.BinaryExpression:
		// Arithmetic operators return numeric types
		leftType := tc.inferExpressionType(e.Left, stmt)
		rightType := tc.inferExpressionType(e.Right, stmt)

		if leftType != nil && rightType != nil {
			// For arithmetic, promote to the wider type
			if leftType.Name == "FLOAT" || rightType.Name == "FLOAT" {
				return &DataType{Name: "FLOAT"}
			}
			if leftType.Name == "DECIMAL" || rightType.Name == "DECIMAL" {
				return &DataType{Name: "DECIMAL", Precision: 18, Scale: 4}
			}
			return leftType
		}

	case *parser.FunctionCall:
		// Infer type based on function name
		return tc.inferFunctionReturnType(e)
	}

	return nil
}

// inferFunctionReturnType infers the return type of a function
func (tc *TypeChecker) inferFunctionReturnType(fn *parser.FunctionCall) *DataType {
	switch fn.Name {
	case "COUNT", "SUM":
		return &DataType{Name: "INT"}
	case "AVG":
		return &DataType{Name: "FLOAT"}
	case "MAX", "MIN":
		// Return type depends on argument type
		if len(fn.Arguments) > 0 {
			return tc.inferExpressionType(fn.Arguments[0], nil)
		}
	case "CONCAT", "SUBSTRING", "UPPER", "LOWER", "TRIM":
		return &DataType{Name: "VARCHAR"}
	case "NOW", "CURRENT_TIMESTAMP":
		return &DataType{Name: "TIMESTAMP"}
	case "CURRENT_DATE":
		return &DataType{Name: "DATE"}
	}

	return nil
}
