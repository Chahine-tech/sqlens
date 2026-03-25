package dialect

import "strings"

// SQLiteDialect implements SQLite-specific features
type SQLiteDialect struct{}

func (d *SQLiteDialect) Name() string {
	return "SQLite"
}

func (d *SQLiteDialect) QuoteIdentifier(identifier string) string {
	return `"` + identifier + `"`
}

func (d *SQLiteDialect) SupportsFeature(feature Feature) bool {
	switch feature {
	case FeatureCTE:
		return true
	case FeatureWindowFunctions:
		return true // SQLite 3.25+
	case FeatureJSONSupport:
		return true // SQLite 3.38+
	case FeatureRecursiveCTE:
		return true
	case FeatureFullTextSearch:
		return true // FTS extension
	case FeatureUpsert:
		return true // INSERT ... ON CONFLICT
	default:
		return false
	}
}

func (d *SQLiteDialect) GetKeywords() []string {
	sqlite := []string{
		"AUTOINCREMENT", "PRAGMA", "VACUUM", "ANALYZE", "ATTACH", "DETACH",
		"CONFLICT", "ABORT", "FAIL", "IGNORE", "REPLACE", "ROLLBACK",
		"TEMP", "TEMPORARY", "GLOB", "REGEXP", "MATCH", "ISNULL", "NOTNULL",
		"WITHOUT", "ROWID", "STRICT", "GENERATED", "ALWAYS", "STORED", "VIRTUAL",
	}
	return append(CommonKeywords, sqlite...)
}

func (d *SQLiteDialect) GetDataTypes() []string {
	sqlite := []string{
		"INTEGER", "REAL", "TEXT", "BLOB", "NUMERIC", "NONE",
		// SQLite is flexible with data types, but these are the storage classes
	}
	return append(CommonDataTypes, sqlite...)
}

var sqliteReserved = func() map[string]bool {
	m := make(map[string]bool)
	for _, k := range (&SQLiteDialect{}).GetKeywords() {
		m[k] = true
	}
	return m
}()

func (d *SQLiteDialect) IsReservedWord(word string) bool {
	return sqliteReserved[strings.ToUpper(word)]
}

func (d *SQLiteDialect) GetLimitSyntax() LimitSyntax {
	return LimitSyntaxStandard
}
