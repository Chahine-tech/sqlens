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

// schemaData is a shared intermediate representation for JSON and YAML loading.
// Both tag sets are present so a single struct works for both unmarshal calls.
type schemaData struct {
	Name   string      `json:"name"   yaml:"name"`
	Tables []tableData `json:"tables" yaml:"tables"`
}

type tableData struct {
	Name    string       `json:"name"           yaml:"name"`
	Schema  string       `json:"schema,omitempty" yaml:"schema,omitempty"`
	Columns []columnData `json:"columns"        yaml:"columns"`
	Indexes []indexData  `json:"indexes,omitempty" yaml:"indexes,omitempty"`
}

type columnData struct {
	Name         string      `json:"name"                yaml:"name"`
	Type         string      `json:"type"                yaml:"type"`
	Length       int         `json:"length,omitempty"    yaml:"length,omitempty"`
	Precision    int         `json:"precision,omitempty" yaml:"precision,omitempty"`
	Scale        int         `json:"scale,omitempty"     yaml:"scale,omitempty"`
	Nullable     bool        `json:"nullable,omitempty"  yaml:"nullable,omitempty"`
	PrimaryKey   bool        `json:"primary_key,omitempty" yaml:"primary_key,omitempty"`
	Unique       bool        `json:"unique,omitempty"    yaml:"unique,omitempty"`
	ForeignKey   bool        `json:"foreign_key,omitempty" yaml:"foreign_key,omitempty"`
	FKTable      string      `json:"fk_table,omitempty"  yaml:"fk_table,omitempty"`
	FKColumn     string      `json:"fk_column,omitempty" yaml:"fk_column,omitempty"`
	DefaultValue interface{} `json:"default,omitempty"   yaml:"default,omitempty"`
}

type indexData struct {
	Name     string   `json:"name"           yaml:"name"`
	Columns  []string `json:"columns"        yaml:"columns"`
	IsUnique bool     `json:"unique,omitempty" yaml:"unique,omitempty"`
}

// buildSchema constructs a *Schema from the shared intermediate representation.
func buildSchema(data schemaData) *Schema {
	s := NewSchema(data.Name)
	for _, td := range data.Tables {
		table := NewTable(td.Name)
		table.Schema = td.Schema

		for _, cd := range td.Columns {
			col := &Column{
				Name:         cd.Name,
				IsPrimaryKey: cd.PrimaryKey,
				IsUnique:     cd.Unique,
				IsForeignKey: cd.ForeignKey,
				DefaultValue: cd.DefaultValue,
				DataType: &DataType{
					Name:      strings.ToUpper(cd.Type),
					Length:    cd.Length,
					Precision: cd.Precision,
					Scale:     cd.Scale,
					Nullable:  cd.Nullable,
				},
			}
			if cd.ForeignKey && cd.FKTable != "" {
				col.ForeignKey = &ForeignKeyRef{
					Table:  cd.FKTable,
					Column: cd.FKColumn,
				}
			}
			table.AddColumn(col)
		}

		for _, id := range td.Indexes {
			table.AddIndex(&Index{
				Name:     id.Name,
				Table:    td.Name,
				Columns:  id.Columns,
				IsUnique: id.IsUnique,
			})
		}

		s.AddTable(table)
	}
	return s
}

// LoadFromJSON loads a schema from JSON
func (sl *SchemaLoader) LoadFromJSON(data []byte) (*Schema, error) {
	var sd schemaData
	if err := json.Unmarshal(data, &sd); err != nil {
		return nil, fmt.Errorf("failed to parse JSON schema: %w", err)
	}
	return buildSchema(sd), nil
}

// LoadFromYAML loads a schema from YAML
func (sl *SchemaLoader) LoadFromYAML(data []byte) (*Schema, error) {
	var sd schemaData
	if err := yaml.Unmarshal(data, &sd); err != nil {
		return nil, fmt.Errorf("failed to parse YAML schema: %w", err)
	}
	return buildSchema(sd), nil
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

	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".json"):
		return sl.LoadFromJSON(data)
	case strings.HasSuffix(lower, ".yaml"), strings.HasSuffix(lower, ".yml"):
		return sl.LoadFromYAML(data)
	}

	// Try JSON first, then YAML
	if s, err := sl.LoadFromJSON(data); err == nil {
		return s, nil
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
