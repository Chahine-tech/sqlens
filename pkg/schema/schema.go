package schema

import (
	"fmt"
	"strings"
)

// DataType represents a SQL data type
type DataType struct {
	Name      string // INT, VARCHAR, DECIMAL, etc.
	Length    int    // For VARCHAR(255), CHAR(10)
	Precision int    // For DECIMAL(10,2)
	Scale     int    // For DECIMAL(10,2)
	Nullable  bool   // Can be NULL
}

// String returns the string representation of the data type
func (dt *DataType) String() string {
	if dt.Length > 0 {
		return fmt.Sprintf("%s(%d)", dt.Name, dt.Length)
	}
	if dt.Precision > 0 {
		if dt.Scale > 0 {
			return fmt.Sprintf("%s(%d,%d)", dt.Name, dt.Precision, dt.Scale)
		}
		return fmt.Sprintf("%s(%d)", dt.Name, dt.Precision)
	}
	return dt.Name
}

// IsCompatibleWith checks if this data type is compatible with another
func (dt *DataType) IsCompatibleWith(other *DataType) bool {
	// Exact match
	if dt.Name == other.Name {
		return true
	}

	// Numeric types are compatible with each other
	numericTypes := map[string]bool{
		"INT": true, "INTEGER": true, "BIGINT": true, "SMALLINT": true, "TINYINT": true,
		"DECIMAL": true, "NUMERIC": true, "FLOAT": true, "DOUBLE": true, "REAL": true,
	}
	if numericTypes[dt.Name] && numericTypes[other.Name] {
		return true
	}

	// String types are compatible with each other
	stringTypes := map[string]bool{
		"VARCHAR": true, "CHAR": true, "TEXT": true, "NVARCHAR": true, "NCHAR": true,
	}
	if stringTypes[dt.Name] && stringTypes[other.Name] {
		return true
	}

	// Date/time types
	dateTypes := map[string]bool{
		"DATE": true, "TIME": true, "DATETIME": true, "TIMESTAMP": true,
	}
	if dateTypes[dt.Name] && dateTypes[other.Name] {
		return true
	}

	return false
}

// Column represents a database column
type Column struct {
	Name         string
	DataType     *DataType
	IsPrimaryKey bool
	IsUnique     bool
	IsForeignKey bool
	ForeignKey   *ForeignKeyRef // Reference to another table
	DefaultValue interface{}
}

// ForeignKeyRef represents a foreign key reference
type ForeignKeyRef struct {
	Table  string
	Column string
}

// Table represents a database table
type Table struct {
	Name    string
	Schema  string             // Schema/database name (optional)
	Columns map[string]*Column // Column name -> Column
	Indexes map[string]*Index  // Index name -> Index
}

// NewTable creates a new table
func NewTable(name string) *Table {
	return &Table{
		Name:    name,
		Columns: make(map[string]*Column),
		Indexes: make(map[string]*Index),
	}
}

// AddColumn adds a column to the table
func (t *Table) AddColumn(col *Column) {
	t.Columns[strings.ToLower(col.Name)] = col
}

// GetColumn retrieves a column by name (case-insensitive)
func (t *Table) GetColumn(name string) (*Column, bool) {
	col, ok := t.Columns[strings.ToLower(name)]
	return col, ok
}

// HasColumn checks if a column exists (case-insensitive)
func (t *Table) HasColumn(name string) bool {
	_, ok := t.Columns[strings.ToLower(name)]
	return ok
}

// AddIndex adds an index to the table
func (t *Table) AddIndex(idx *Index) {
	t.Indexes[strings.ToLower(idx.Name)] = idx
}

// GetIndex retrieves an index by name (case-insensitive)
func (t *Table) GetIndex(name string) (*Index, bool) {
	idx, ok := t.Indexes[strings.ToLower(name)]
	return idx, ok
}

// Index represents a database index
type Index struct {
	Name     string
	Table    string
	Columns  []string
	IsUnique bool
}

// Schema represents a database schema (collection of tables)
type Schema struct {
	Name   string
	Tables map[string]*Table // Table name -> Table
}

// NewSchema creates a new schema
func NewSchema(name string) *Schema {
	return &Schema{
		Name:   name,
		Tables: make(map[string]*Table),
	}
}

// AddTable adds a table to the schema
func (s *Schema) AddTable(table *Table) {
	s.Tables[strings.ToLower(table.Name)] = table
}

// GetTable retrieves a table by name (case-insensitive)
func (s *Schema) GetTable(name string) (*Table, bool) {
	table, ok := s.Tables[strings.ToLower(name)]
	return table, ok
}

// HasTable checks if a table exists (case-insensitive)
func (s *Schema) HasTable(name string) bool {
	_, ok := s.Tables[strings.ToLower(name)]
	return ok
}

// GetColumn retrieves a column from a table (case-insensitive)
func (s *Schema) GetColumn(tableName, columnName string) (*Column, error) {
	table, ok := s.GetTable(tableName)
	if !ok {
		return nil, fmt.Errorf("table '%s' not found in schema", tableName)
	}

	column, ok := table.GetColumn(columnName)
	if !ok {
		return nil, fmt.Errorf("column '%s' not found in table '%s'", columnName, tableName)
	}

	return column, nil
}

// Validate performs basic schema validation
func (s *Schema) Validate() error {
	// Check for foreign key references
	for _, table := range s.Tables {
		for _, col := range table.Columns {
			if col.IsForeignKey && col.ForeignKey != nil {
				// Check if referenced table exists
				refTable, ok := s.GetTable(col.ForeignKey.Table)
				if !ok {
					return fmt.Errorf("foreign key in table '%s' column '%s' references non-existent table '%s'",
						table.Name, col.Name, col.ForeignKey.Table)
				}

				// Check if referenced column exists
				if !refTable.HasColumn(col.ForeignKey.Column) {
					return fmt.Errorf("foreign key in table '%s' column '%s' references non-existent column '%s.%s'",
						table.Name, col.Name, col.ForeignKey.Table, col.ForeignKey.Column)
				}
			}
		}
	}

	return nil
}
