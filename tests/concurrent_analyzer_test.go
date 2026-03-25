package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/analyzer"
	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func parseStmt(t *testing.T, sql string) parser.Statement {
	t.Helper()
	ctx := context.Background()
	p := parser.NewWithDialect(ctx, sql, dialect.GetDialect("mysql"))
	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("parse error for %q: %v", sql, err)
	}
	return stmt
}

func TestConcurrentAnalyzer_Basic(t *testing.T) {
	ca := analyzer.NewConcurrentAnalyzer(2)
	jobs := []analyzer.AnalysisJob{
		{ID: "1", Stmt: parseStmt(t, "SELECT id FROM users")},
		{ID: "2", Stmt: parseStmt(t, "SELECT name FROM products")},
		{ID: "3", Stmt: parseStmt(t, "SELECT * FROM orders WHERE id = 1")},
	}
	ctx := context.Background()
	results := ca.AnalyzeConcurrently(ctx, jobs)
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Error != nil {
			t.Errorf("unexpected error for job %s: %v", r.ID, r.Error)
		}
	}
}

func TestConcurrentAnalyzer_DefaultWorkerCount(t *testing.T) {
	// workerCount <= 0 should default to runtime.NumCPU()
	ca := analyzer.NewConcurrentAnalyzer(0)
	jobs := []analyzer.AnalysisJob{
		{ID: "a", Stmt: parseStmt(t, "SELECT 1")},
	}
	results := ca.AnalyzeConcurrently(context.Background(), jobs)
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestConcurrentAnalyzer_EmptyJobs(t *testing.T) {
	ca := analyzer.NewConcurrentAnalyzer(2)
	results := ca.AnalyzeConcurrently(context.Background(), nil)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestConcurrentAnalyzer_CacheHit(t *testing.T) {
	ca := analyzer.NewConcurrentAnalyzer(1)
	stmt := parseStmt(t, "SELECT id FROM users")
	jobs := []analyzer.AnalysisJob{
		{ID: "cached", Stmt: stmt},
		{ID: "cached", Stmt: stmt}, // same ID — second should hit cache
	}
	ctx := context.Background()
	results := ca.AnalyzeConcurrently(ctx, jobs)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestConcurrentAnalyzer_CancelContext(t *testing.T) {
	ca := analyzer.NewConcurrentAnalyzer(1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	jobs := []analyzer.AnalysisJob{
		{ID: "1", Stmt: parseStmt(t, "SELECT id FROM users")},
	}
	// Should not block or panic
	_ = ca.AnalyzeConcurrently(ctx, jobs)
}

func TestConcurrentAnalyzer_GetCacheStats(t *testing.T) {
	ca := analyzer.NewConcurrentAnalyzer(2)
	jobs := []analyzer.AnalysisJob{
		{ID: "x", Stmt: parseStmt(t, "SELECT 1")},
	}
	ca.AnalyzeConcurrently(context.Background(), jobs)

	stats := ca.GetCacheStats()
	if _, ok := stats["cache_size"]; !ok {
		t.Error("missing cache_size in stats")
	}
	if _, ok := stats["worker_count"]; !ok {
		t.Error("missing worker_count in stats")
	}
}

func TestConcurrentAnalyzer_ClearCache(t *testing.T) {
	ca := analyzer.NewConcurrentAnalyzer(1)
	jobs := []analyzer.AnalysisJob{
		{ID: "y", Stmt: parseStmt(t, "SELECT 1")},
	}
	ca.AnalyzeConcurrently(context.Background(), jobs)
	ca.ClearCache()

	stats := ca.GetCacheStats()
	if stats["cache_size"].(int) != 0 {
		t.Error("expected empty cache after ClearCache")
	}
}

func TestAnalyzer_AnalyzeWithCache(t *testing.T) {
	a := analyzer.New()
	stmt := parseStmt(t, "SELECT id, name FROM users WHERE id = 1")

	r1 := a.AnalyzeWithCache(stmt, "key1")
	r2 := a.AnalyzeWithCache(stmt, "key1") // cache hit
	if r1.QueryType != r2.QueryType {
		t.Error("cache hit should return same result")
	}

	// Empty cache key — should still work
	r3 := a.AnalyzeWithCache(stmt, "")
	if r3.QueryType != "SELECT" {
		t.Errorf("expected SELECT, got %s", r3.QueryType)
	}
}

func TestAnalyzer_SuggestOptimizations_SelectStar(t *testing.T) {
	ctx := context.Background()
	p := parser.NewWithDialect(ctx, "SELECT * FROM users", dialect.GetDialect("mysql"))
	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatal(err)
	}
	sel := stmt.(*parser.SelectStatement)

	a := analyzer.New()
	a.Analyze(stmt)
	suggestions := a.SuggestOptimizations(sel)

	found := false
	for _, s := range suggestions {
		if s.Type == "SELECT_ALL" {
			found = true
		}
	}
	if !found {
		t.Error("expected SELECT_ALL suggestion")
	}
}

func TestAnalyzer_SuggestOptimizations_CartesianProduct(t *testing.T) {
	ctx := context.Background()
	p := parser.NewWithDialect(ctx, "SELECT * FROM users, orders", dialect.GetDialect("mysql"))
	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatal(err)
	}
	sel := stmt.(*parser.SelectStatement)

	a := analyzer.New()
	a.Analyze(stmt)
	suggestions := a.SuggestOptimizations(sel)

	found := false
	for _, s := range suggestions {
		if s.Type == "CARTESIAN_PRODUCT" {
			found = true
		}
	}
	if !found {
		t.Error("expected CARTESIAN_PRODUCT suggestion")
	}
}

func TestAnalyzer_DDL_Statements(t *testing.T) {
	cases := []struct {
		sql      string
		dialect  string
		wantType string
	}{
		{"CREATE TABLE t (id INT)", "mysql", "CREATE TABLE"},
		{"DROP TABLE t", "mysql", "DROP TABLE"},
		{"ALTER TABLE t ADD COLUMN name VARCHAR(100)", "mysql", "ALTER TABLE"},
		{"CREATE INDEX idx ON t (id)", "mysql", "CREATE INDEX"},
	}
	for _, tc := range cases {
		ctx := context.Background()
		p := parser.NewWithDialect(ctx, tc.sql, dialect.GetDialect(tc.dialect))
		stmt, err := p.ParseStatement()
		if err != nil {
			t.Errorf("parse error for %q: %v", tc.sql, err)
			continue
		}
		a := analyzer.New()
		result := a.Analyze(stmt)
		if result.QueryType != tc.wantType {
			t.Errorf("sql=%q: QueryType=%q, want %q", tc.sql, result.QueryType, tc.wantType)
		}
	}
}

func TestAnalyzer_DML_Statements(t *testing.T) {
	cases := []struct {
		sql      string
		wantType string
	}{
		{"INSERT INTO users (id, name) VALUES (1, 'Alice')", "INSERT"},
		{"UPDATE users SET name = 'Bob' WHERE id = 1", "UPDATE"},
		{"DELETE FROM users WHERE id = 1", "DELETE"},
	}
	for _, tc := range cases {
		ctx := context.Background()
		p := parser.NewWithDialect(ctx, tc.sql, dialect.GetDialect("mysql"))
		stmt, err := p.ParseStatement()
		if err != nil {
			t.Errorf("parse error for %q: %v", tc.sql, err)
			continue
		}
		a := analyzer.New()
		result := a.Analyze(stmt)
		if result.QueryType != tc.wantType {
			t.Errorf("sql=%q: QueryType=%q, want %q", tc.sql, result.QueryType, tc.wantType)
		}
		if len(result.Tables) == 0 {
			t.Errorf("sql=%q: expected at least one table", tc.sql)
		}
	}
}

func TestAnalyzer_Nil_Statement(t *testing.T) {
	a := analyzer.New()
	result := a.Analyze(nil)
	if result.QueryType != "" {
		t.Errorf("expected empty query type for nil stmt, got %q", result.QueryType)
	}
}
