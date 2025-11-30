package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// SchemaLoader loads schemas from various formats
type SchemaLoader struct {
	schemas map[string]*Schema // Schema name -> Schema
}

// NewSchemaLoader creates a new schema loader
func NewSchemaLoader() *SchemaLoader {
	return &SchemaLoader{
		schemas: make(map[string]*Schema),
	}
}

// LoadFromJSON loads a schema from JSON
func (sl *SchemaLoader) LoadFromJSON(data []byte) (*Schema, error) {
	var schemaData struct {
		Name   string `json:"name"`
		Tables []struct {
			Name    string `json:"name"`
			Schema  string `json:"schema,omitempty"`
			Columns []struct {
				Name         string      `json:"name"`
				Type         string      `json:"type"`
				Length       int         `json:"length,omitempty"`
				Precision    int         `json:"precision,omitempty"`
				Scale        int         `json:"scale,omitempty"`
				Nullable     bool        `json:"nullable,omitempty"`
				PrimaryKey   bool        `json:"primary_key,omitempty"`
				Unique       bool        `json:"unique,omitempty"`
				ForeignKey   bool        `json:"foreign_key,omitempty"`
				FKTable      string      `json:"fk_table,omitempty"`
				FKColumn     string      `json:"fk_column,omitempty"`
				DefaultValue interface{} `json:"default,omitempty"`
			} `json:"columns"`
			Indexes []struct {
				Name     string   `json:"name"`
				Columns  []string `json:"columns"`
				IsUnique bool     `json:"unique,omitempty"`
			} `json:"indexes,omitempty"`
		} `json:"tables"`
	}

	if err := json.Unmarshal(data, &schemaData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON schema: %w", err)
	}

	schema := NewSchema(schemaData.Name)

	// Parse tables
	for _, tableData := range schemaData.Tables {
		table := NewTable(tableData.Name)
		table.Schema = tableData.Schema

		// Parse columns
		for _, colData := range tableData.Columns {
			col := &Column{
				Name:         colData.Name,
				IsPrimaryKey: colData.PrimaryKey,
				IsUnique:     colData.Unique,
				IsForeignKey: colData.ForeignKey,
				DefaultValue: colData.DefaultValue,
				DataType: &DataType{
					Name:      strings.ToUpper(colData.Type),
					Length:    colData.Length,
					Precision: colData.Precision,
					Scale:     colData.Scale,
					Nullable:  colData.Nullable,
				},
			}

			if colData.ForeignKey && colData.FKTable != "" {
				col.ForeignKey = &ForeignKeyRef{
					Table:  colData.FKTable,
					Column: colData.FKColumn,
				}
			}

			table.AddColumn(col)
		}

		// Parse indexes
		for _, idxData := range tableData.Indexes {
			idx := &Index{
				Name:     idxData.Name,
				Table:    tableData.Name,
				Columns:  idxData.Columns,
				IsUnique: idxData.IsUnique,
			}
			table.AddIndex(idx)
		}

		schema.AddTable(table)
	}

	return schema, nil
}

// LoadFromYAML loads a schema from YAML
func (sl *SchemaLoader) LoadFromYAML(data []byte) (*Schema, error) {
	var schemaData struct {
		Name   string `yaml:"name"`
		Tables []struct {
			Name    string `yaml:"name"`
			Schema  string `yaml:"schema,omitempty"`
			Columns []struct {
				Name         string      `yaml:"name"`
				Type         string      `yaml:"type"`
				Length       int         `yaml:"length,omitempty"`
				Precision    int         `yaml:"precision,omitempty"`
				Scale        int         `yaml:"scale,omitempty"`
				Nullable     bool        `yaml:"nullable,omitempty"`
				PrimaryKey   bool        `yaml:"primary_key,omitempty"`
				Unique       bool        `yaml:"unique,omitempty"`
				ForeignKey   bool        `yaml:"foreign_key,omitempty"`
				FKTable      string      `yaml:"fk_table,omitempty"`
				FKColumn     string      `yaml:"fk_column,omitempty"`
				DefaultValue interface{} `yaml:"default,omitempty"`
			} `yaml:"columns"`
			Indexes []struct {
				Name     string   `yaml:"name"`
				Columns  []string `yaml:"columns"`
				IsUnique bool     `yaml:"unique,omitempty"`
			} `yaml:"indexes,omitempty"`
		} `yaml:"tables"`
	}

	if err := yaml.Unmarshal(data, &schemaData); err != nil {
		return nil, fmt.Errorf("failed to parse YAML schema: %w", err)
	}

	schema := NewSchema(schemaData.Name)

	// Parse tables
	for _, tableData := range schemaData.Tables {
		table := NewTable(tableData.Name)
		table.Schema = tableData.Schema

		// Parse columns
		for _, colData := range tableData.Columns {
			col := &Column{
				Name:         colData.Name,
				IsPrimaryKey: colData.PrimaryKey,
				IsUnique:     colData.Unique,
				IsForeignKey: colData.ForeignKey,
				DefaultValue: colData.DefaultValue,
				DataType: &DataType{
					Name:      strings.ToUpper(colData.Type),
					Length:    colData.Length,
					Precision: colData.Precision,
					Scale:     colData.Scale,
					Nullable:  colData.Nullable,
				},
			}

			if colData.ForeignKey && colData.FKTable != "" {
				col.ForeignKey = &ForeignKeyRef{
					Table:  colData.FKTable,
					Column: colData.FKColumn,
				}
			}

			table.AddColumn(col)
		}

		// Parse indexes
		for _, idxData := range tableData.Indexes {
			idx := &Index{
				Name:     idxData.Name,
				Table:    tableData.Name,
				Columns:  idxData.Columns,
				IsUnique: idxData.IsUnique,
			}
			table.AddIndex(idx)
		}

		schema.AddTable(table)
	}

	return schema, nil
}

// LoadFromFile loads a schema from a file (auto-detects JSON/YAML)
func (sl *SchemaLoader) LoadFromFile(filename string) (*Schema, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open schema file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Auto-detect format based on file extension
	if strings.HasSuffix(strings.ToLower(filename), ".json") {
		return sl.LoadFromJSON(data)
	} else if strings.HasSuffix(strings.ToLower(filename), ".yaml") || strings.HasSuffix(strings.ToLower(filename), ".yml") {
		return sl.LoadFromYAML(data)
	}

	// Try JSON first, then YAML
	schema, err := sl.LoadFromJSON(data)
	if err == nil {
		return schema, nil
	}

	return sl.LoadFromYAML(data)
}

// AddSchema adds a schema to the loader's cache
func (sl *SchemaLoader) AddSchema(schema *Schema) {
	sl.schemas[strings.ToLower(schema.Name)] = schema
}

// GetSchema retrieves a schema by name (case-insensitive)
func (sl *SchemaLoader) GetSchema(name string) (*Schema, bool) {
	schema, ok := sl.schemas[strings.ToLower(name)]
	return schema, ok
}

// HasSchema checks if a schema exists (case-insensitive)
func (sl *SchemaLoader) HasSchema(name string) bool {
	_, ok := sl.schemas[strings.ToLower(name)]
	return ok
}
