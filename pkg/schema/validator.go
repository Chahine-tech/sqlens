package schema

import (
	"fmt"

	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// ValidationError represents a schema validation error
type ValidationError struct {
	Type    string // TABLE_NOT_FOUND, COLUMN_NOT_FOUND, TYPE_MISMATCH, etc.
	Message string
	Table   string
	Column  string
}

// Error implements the error interface
func (ve *ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", ve.Type, ve.Message)
}

// Validator validates SQL statements against a schema
type Validator struct {
	schema *Schema
}

// NewValidator creates a new validator
func NewValidator(schema *Schema) *Validator {
	return &Validator{
		schema: schema,
	}
}

// ValidateStatement validates a SQL statement against the schema
func (v *Validator) ValidateStatement(stmt parser.Statement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	// Validate that schema is not nil
	if v.schema == nil {
		errors = append(errors, &ValidationError{
			Type:    "SCHEMA_NOT_LOADED",
			Message: "Schema is not loaded",
		})
		return errors
	}

	switch s := stmt.(type) {
	case *parser.SelectStatement:
		errors = append(errors, v.validateSelectStatement(s)...)
	case *parser.InsertStatement:
		errors = append(errors, v.validateInsertStatement(s)...)
	case *parser.UpdateStatement:
		errors = append(errors, v.validateUpdateStatement(s)...)
	case *parser.DeleteStatement:
		errors = append(errors, v.validateDeleteStatement(s)...)
	case *parser.WithStatement:
		errors = append(errors, v.validateWithStatement(s)...)
	}

	return errors
}

// validateSelectStatement validates a SELECT statement
func (v *Validator) validateSelectStatement(stmt *parser.SelectStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	// Validate FROM clause
	if stmt.From != nil {
		for _, tableRef := range stmt.From.Tables {
			if !v.validateTableReference(&tableRef) {
				errors = append(errors, &ValidationError{
					Type:    "TABLE_NOT_FOUND",
					Message: fmt.Sprintf("Table '%s' not found in schema", tableRef.Name),
					Table:   tableRef.Name,
				})
			}
		}
	}

	// Validate JOINs
	for _, join := range stmt.Joins {
		if !v.validateTableReference(&join.Table) {
			errors = append(errors, &ValidationError{
				Type:    "TABLE_NOT_FOUND",
				Message: fmt.Sprintf("Table '%s' in JOIN not found in schema", join.Table.Name),
				Table:   join.Table.Name,
			})
		}
	}

	// Validate columns in SELECT list
	for _, col := range stmt.Columns {
		errors = append(errors, v.validateExpression(col, stmt)...)
	}

	// Validate WHERE clause
	if stmt.Where != nil {
		errors = append(errors, v.validateExpression(stmt.Where, stmt)...)
	}

	// Validate GROUP BY
	for _, expr := range stmt.GroupBy {
		errors = append(errors, v.validateExpression(expr, stmt)...)
	}

	// Validate HAVING
	if stmt.Having != nil {
		errors = append(errors, v.validateExpression(stmt.Having, stmt)...)
	}

	return errors
}

// validateInsertStatement validates an INSERT statement
func (v *Validator) validateInsertStatement(stmt *parser.InsertStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	// Validate table
	if !v.validateTableReference(&stmt.Table) {
		errors = append(errors, &ValidationError{
			Type:    "TABLE_NOT_FOUND",
			Message: fmt.Sprintf("Table '%s' not found in schema", stmt.Table.Name),
			Table:   stmt.Table.Name,
		})
		return errors // Can't continue without valid table
	}

	table, _ := v.schema.GetTable(stmt.Table.Name)

	// Validate columns
	for _, colName := range stmt.Columns {
		if !table.HasColumn(colName) {
			errors = append(errors, &ValidationError{
				Type:    "COLUMN_NOT_FOUND",
				Message: fmt.Sprintf("Column '%s' not found in table '%s'", colName, stmt.Table.Name),
				Table:   stmt.Table.Name,
				Column:  colName,
			})
		}
	}

	return errors
}

// validateUpdateStatement validates an UPDATE statement
func (v *Validator) validateUpdateStatement(stmt *parser.UpdateStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	// Validate table
	if !v.validateTableReference(&stmt.Table) {
		errors = append(errors, &ValidationError{
			Type:    "TABLE_NOT_FOUND",
			Message: fmt.Sprintf("Table '%s' not found in schema", stmt.Table.Name),
			Table:   stmt.Table.Name,
		})
		return errors
	}

	table, _ := v.schema.GetTable(stmt.Table.Name)

	// Validate SET columns
	for _, assignment := range stmt.Set {
		if !table.HasColumn(assignment.Column) {
			errors = append(errors, &ValidationError{
				Type:    "COLUMN_NOT_FOUND",
				Message: fmt.Sprintf("Column '%s' not found in table '%s'", assignment.Column, stmt.Table.Name),
				Table:   stmt.Table.Name,
				Column:  assignment.Column,
			})
		}
	}

	// Validate WHERE clause
	if stmt.Where != nil {
		errors = append(errors, v.validateExpressionForTable(stmt.Where, &stmt.Table)...)
	}

	return errors
}

// validateDeleteStatement validates a DELETE statement
func (v *Validator) validateDeleteStatement(stmt *parser.DeleteStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	// Validate table
	if !v.validateTableReference(&stmt.From) {
		errors = append(errors, &ValidationError{
			Type:    "TABLE_NOT_FOUND",
			Message: fmt.Sprintf("Table '%s' not found in schema", stmt.From.Name),
			Table:   stmt.From.Name,
		})
		return errors
	}

	// Validate WHERE clause
	if stmt.Where != nil {
		errors = append(errors, v.validateExpressionForTable(stmt.Where, &stmt.From)...)
	}

	return errors
}

// validateWithStatement validates a WITH (CTE) statement
func (v *Validator) validateWithStatement(stmt *parser.WithStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	// Validate main query
	if selectStmt, ok := stmt.Query.(*parser.SelectStatement); ok {
		errors = append(errors, v.validateSelectStatement(selectStmt)...)
	}

	return errors
}

// validateTableReference validates a table reference
func (v *Validator) validateTableReference(tableRef *parser.TableReference) bool {
	// Skip validation for derived tables (subqueries)
	if tableRef.Subquery != nil {
		return true
	}

	// Skip validation for empty table names
	if tableRef.Name == "" {
		return true
	}

	// Check if table exists in schema
	return v.schema.HasTable(tableRef.Name)
}

// validateExpression validates an expression against the schema
func (v *Validator) validateExpression(expr parser.Expression, stmt *parser.SelectStatement) []*ValidationError {
	errors := make([]*ValidationError, 0)

	switch e := expr.(type) {
	case *parser.ColumnReference:
		// Validate column reference
		if e.Table != "" {
			// Fully qualified column (table.column)
			if !v.schema.HasTable(e.Table) {
				errors = append(errors, &ValidationError{
					Type:    "TABLE_NOT_FOUND",
					Message: fmt.Sprintf("Table '%s' not found", e.Table),
					Table:   e.Table,
				})
			} else {
				table, _ := v.schema.GetTable(e.Table)
				if !table.HasColumn(e.Column) {
					errors = append(errors, &ValidationError{
						Type:    "COLUMN_NOT_FOUND",
						Message: fmt.Sprintf("Column '%s' not found in table '%s'", e.Column, e.Table),
						Table:   e.Table,
						Column:  e.Column,
					})
				}
			}
		} else {
			// Unqualified column - check in all tables in FROM clause
			found := false
			if stmt.From != nil {
				for _, tableRef := range stmt.From.Tables {
					if table, ok := v.schema.GetTable(tableRef.Name); ok {
						if table.HasColumn(e.Column) {
							found = true
							break
						}
					}
				}
			}
			if !found {
				errors = append(errors, &ValidationError{
					Type:    "COLUMN_NOT_FOUND",
					Message: fmt.Sprintf("Column '%s' not found in any table", e.Column),
					Column:  e.Column,
				})
			}
		}

	case *parser.BinaryExpression:
		// Recursively validate left and right sides
		errors = append(errors, v.validateExpression(e.Left, stmt)...)
		errors = append(errors, v.validateExpression(e.Right, stmt)...)

	case *parser.FunctionCall:
		// Validate function arguments
		for _, arg := range e.Arguments {
			errors = append(errors, v.validateExpression(arg, stmt)...)
		}

	case *parser.SubqueryExpression:
		// Validate subquery
		errors = append(errors, v.validateSelectStatement(e.Query)...)
	}

	return errors
}

// validateExpressionForTable validates an expression for a specific table
func (v *Validator) validateExpressionForTable(expr parser.Expression, tableRef *parser.TableReference) []*ValidationError {
	errors := make([]*ValidationError, 0)

	switch e := expr.(type) {
	case *parser.ColumnReference:
		table, ok := v.schema.GetTable(tableRef.Name)
		if !ok {
			return errors // Table not found error already reported
		}

		if !table.HasColumn(e.Column) {
			errors = append(errors, &ValidationError{
				Type:    "COLUMN_NOT_FOUND",
				Message: fmt.Sprintf("Column '%s' not found in table '%s'", e.Column, tableRef.Name),
				Table:   tableRef.Name,
				Column:  e.Column,
			})
		}

	case *parser.BinaryExpression:
		errors = append(errors, v.validateExpressionForTable(e.Left, tableRef)...)
		errors = append(errors, v.validateExpressionForTable(e.Right, tableRef)...)
	}

	return errors
}
