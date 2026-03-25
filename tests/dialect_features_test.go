package tests

import (
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
)

// Exercise SupportsFeature for all dialects across all features so coverage
// picks up every branch of each switch statement.

var allFeatures = []dialect.Feature{
	dialect.FeatureCTE,
	dialect.FeatureWindowFunctions,
	dialect.FeatureJSONSupport,
	dialect.FeatureArraySupport,
	dialect.FeatureRecursiveCTE,
	dialect.FeaturePartitioning,
	dialect.FeatureFullTextSearch,
	dialect.FeatureXMLSupport,
	dialect.FeatureUpsert,
	dialect.FeatureReturningClause,
}

func TestDialect_SupportsFeature_AllDialects(t *testing.T) {
	dialects := []string{"mysql", "postgresql", "sqlserver", "sqlite", "oracle"}
	for _, name := range dialects {
		d := dialect.GetDialect(name)
		for _, f := range allFeatures {
			// Just exercise — we don't assert specific values, just no panics
			_ = d.SupportsFeature(f)
		}
	}
}

func TestDialect_MySQL_Features(t *testing.T) {
	d := dialect.GetDialect("mysql")
	if !d.SupportsFeature(dialect.FeatureCTE) {
		t.Error("MySQL should support CTE")
	}
	if !d.SupportsFeature(dialect.FeatureWindowFunctions) {
		t.Error("MySQL should support window functions")
	}
	if !d.SupportsFeature(dialect.FeatureJSONSupport) {
		t.Error("MySQL should support JSON")
	}
	if !d.SupportsFeature(dialect.FeatureFullTextSearch) {
		t.Error("MySQL should support full text search")
	}
	if !d.SupportsFeature(dialect.FeatureUpsert) {
		t.Error("MySQL should support upsert")
	}
}

func TestDialect_PostgreSQL_Features(t *testing.T) {
	d := dialect.GetDialect("postgresql")
	if !d.SupportsFeature(dialect.FeatureCTE) {
		t.Error("PostgreSQL should support CTE")
	}
	if !d.SupportsFeature(dialect.FeatureWindowFunctions) {
		t.Error("PostgreSQL should support window functions")
	}
	if !d.SupportsFeature(dialect.FeatureArraySupport) {
		t.Error("PostgreSQL should support arrays")
	}
	if !d.SupportsFeature(dialect.FeatureReturningClause) {
		t.Error("PostgreSQL should support RETURNING")
	}
}

func TestDialect_SQLServer_Features(t *testing.T) {
	d := dialect.GetDialect("sqlserver")
	if !d.SupportsFeature(dialect.FeatureCTE) {
		t.Error("SQL Server should support CTE")
	}
	if !d.SupportsFeature(dialect.FeatureXMLSupport) {
		t.Error("SQL Server should support XML")
	}
}

func TestDialect_SQLite_Features(t *testing.T) {
	d := dialect.GetDialect("sqlite")
	if !d.SupportsFeature(dialect.FeatureCTE) {
		t.Error("SQLite should support CTE")
	}
}

func TestDialect_Oracle_Features(t *testing.T) {
	d := dialect.GetDialect("oracle")
	if !d.SupportsFeature(dialect.FeatureCTE) {
		t.Error("Oracle should support CTE")
	}
	if !d.SupportsFeature(dialect.FeatureXMLSupport) {
		t.Error("Oracle should support XML")
	}
}

func TestDialect_GetLimitSyntax(t *testing.T) {
	cases := []struct {
		name     string
		expected dialect.LimitSyntax
	}{
		{"mysql", dialect.LimitSyntaxStandard},
		{"postgresql", dialect.LimitSyntaxStandard},
		{"sqlserver", dialect.LimitSyntaxSQLServer},
		{"sqlite", dialect.LimitSyntaxStandard},
		{"oracle", dialect.LimitSyntaxOracle},
	}
	for _, tc := range cases {
		d := dialect.GetDialect(tc.name)
		if got := d.GetLimitSyntax(); got != tc.expected {
			t.Errorf("dialect %s: LimitSyntax = %v, want %v", tc.name, got, tc.expected)
		}
	}
}

func TestDialect_GetKeywords_NonEmpty(t *testing.T) {
	for _, name := range []string{"mysql", "postgresql", "sqlserver", "sqlite", "oracle"} {
		d := dialect.GetDialect(name)
		if len(d.GetKeywords()) == 0 {
			t.Errorf("dialect %s should have keywords", name)
		}
	}
}

func TestDialect_GetDataTypes_NonEmpty(t *testing.T) {
	for _, name := range []string{"mysql", "postgresql", "sqlserver", "sqlite", "oracle"} {
		d := dialect.GetDialect(name)
		if len(d.GetDataTypes()) == 0 {
			t.Errorf("dialect %s should have data types", name)
		}
	}
}

func TestDialect_IsReservedWord(t *testing.T) {
	cases := []struct {
		dialect string
		word    string
		want    bool
	}{
		{"mysql", "SELECT", true},
		{"mysql", "notakeyword123", false},
		{"postgresql", "SELECT", true},
		{"postgresql", "notakeyword123", false},
		{"sqlserver", "SELECT", true},
		{"sqlserver", "notakeyword123", false},
		{"sqlite", "SELECT", true},
		{"sqlite", "notakeyword123", false},
		{"oracle", "SELECT", true},
		{"oracle", "notakeyword123", false},
	}
	for _, tc := range cases {
		d := dialect.GetDialect(tc.dialect)
		got := d.IsReservedWord(tc.word)
		if got != tc.want {
			t.Errorf("dialect %s: IsReservedWord(%q) = %v, want %v", tc.dialect, tc.word, got, tc.want)
		}
	}
}

func TestDialect_QuoteIdentifier(t *testing.T) {
	cases := []struct {
		dialect string
		ident   string
		want    string
	}{
		{"mysql", "table_name", "`table_name`"},
		{"postgresql", "table_name", `"table_name"`},
		{"sqlserver", "table_name", "[table_name]"},
		{"sqlite", "table_name", `"table_name"`},
		{"oracle", "table_name", `"table_name"`},
	}
	for _, tc := range cases {
		d := dialect.GetDialect(tc.dialect)
		got := d.QuoteIdentifier(tc.ident)
		if got != tc.want {
			t.Errorf("dialect %s: QuoteIdentifier(%q) = %q, want %q", tc.dialect, tc.ident, got, tc.want)
		}
	}
}

func TestDialect_Name(t *testing.T) {
	cases := []struct{ dialect, want string }{
		{"mysql", "MySQL"},
		{"postgresql", "PostgreSQL"},
		{"sqlserver", "SQL Server"},
		{"sqlite", "SQLite"},
		{"oracle", "Oracle"},
	}
	for _, tc := range cases {
		d := dialect.GetDialect(tc.dialect)
		if got := d.Name(); got != tc.want {
			t.Errorf("dialect %s: Name() = %q, want %q", tc.dialect, got, tc.want)
		}
	}
}
