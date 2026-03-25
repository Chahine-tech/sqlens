package tests

import (
	"testing"
	"time"

	"github.com/Chahine-tech/sql-parser-go/pkg/analyzer"
	"github.com/Chahine-tech/sql-parser-go/pkg/monitor"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func makePQ(query string, duration float64, stmt parser.Statement, analysis *analyzer.QueryAnalysis) *monitor.ProcessedQuery {
	return &monitor.ProcessedQuery{
		Timestamp: time.Now(),
		Query:     query,
		Duration:  duration,
		Statement: stmt,
		Analysis:  analysis,
	}
}

// --- AlertLevel.String ---

func TestAlertLevel_String(t *testing.T) {
	cases := []struct {
		level    monitor.AlertLevel
		expected string
	}{
		{monitor.AlertInfo, "INFO"},
		{monitor.AlertWarning, "WARNING"},
		{monitor.AlertError, "ERROR"},
		{monitor.AlertCritical, "CRITICAL"},
		{monitor.AlertLevel(99), "UNKNOWN"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.expected {
			t.Errorf("AlertLevel(%d).String() = %q, want %q", tc.level, got, tc.expected)
		}
	}
}

// --- AlertManager ---

func TestAlertManager_NoRules(t *testing.T) {
	am := monitor.NewAlertManager()
	pq := makePQ("SELECT 1", 0, nil, nil)
	am.Check(pq) // should not panic
	counts := am.GetAlertCounts()
	if len(counts) != 0 {
		t.Error("expected no alerts with no rules")
	}
}

func TestAlertManager_AddHandlerAndRule(t *testing.T) {
	am := monitor.NewAlertManager()
	am.AddRule(&monitor.SlowQueryRule{Threshold: 1.0})

	var fired []*monitor.Alert
	am.AddHandler(func(a *monitor.Alert) { fired = append(fired, a) })

	// Should not fire — fast query
	am.Check(makePQ("SELECT 1", 0.1, nil, nil))
	if len(fired) != 0 {
		t.Error("expected no alerts for fast query")
	}

	// Should fire — slow query
	am.Check(makePQ("SELECT 1", 2.0, nil, nil))
	if len(fired) != 1 {
		t.Errorf("expected 1 alert, got %d", len(fired))
	}
}

func TestAlertManager_GetAlertCounts(t *testing.T) {
	am := monitor.NewAlertManager()
	am.AddRule(&monitor.SlowQueryRule{Threshold: 0.5})
	am.Check(makePQ("q", 1.0, nil, nil)) // warning level
	counts := am.GetAlertCounts()
	total := int64(0)
	for _, c := range counts {
		total += c
	}
	if total == 0 {
		t.Error("expected at least one alert recorded")
	}
}

// --- SlowQueryRule ---

func TestSlowQueryRule_Name(t *testing.T) {
	r := &monitor.SlowQueryRule{Threshold: 1.0}
	if r.Name() != "SlowQueryRule" {
		t.Errorf("unexpected name: %s", r.Name())
	}
}

func TestSlowQueryRule_BelowThreshold(t *testing.T) {
	r := &monitor.SlowQueryRule{Threshold: 2.0}
	alert := r.Check(makePQ("SELECT 1", 1.0, nil, nil))
	if alert != nil {
		t.Error("expected no alert below threshold")
	}
}

func TestSlowQueryRule_Warning(t *testing.T) {
	r := &monitor.SlowQueryRule{Threshold: 1.0}
	alert := r.Check(makePQ("SELECT 1", 1.5, nil, nil))
	if alert == nil {
		t.Fatal("expected alert")
	}
	if alert.Level != monitor.AlertWarning {
		t.Errorf("expected WARNING, got %s", alert.Level)
	}
	if alert.Type != "SLOW_QUERY" {
		t.Errorf("expected SLOW_QUERY, got %s", alert.Type)
	}
}

func TestSlowQueryRule_Error(t *testing.T) {
	r := &monitor.SlowQueryRule{Threshold: 1.0}
	alert := r.Check(makePQ("SELECT 1", 2.5, nil, nil))
	if alert == nil {
		t.Fatal("expected alert")
	}
	if alert.Level != monitor.AlertError {
		t.Errorf("expected ERROR, got %s", alert.Level)
	}
}

func TestSlowQueryRule_Critical(t *testing.T) {
	r := &monitor.SlowQueryRule{Threshold: 1.0}
	alert := r.Check(makePQ("SELECT 1", 6.0, nil, nil))
	if alert == nil {
		t.Fatal("expected alert")
	}
	if alert.Level != monitor.AlertCritical {
		t.Errorf("expected CRITICAL, got %s", alert.Level)
	}
}

// --- ParseErrorRule ---

func TestParseErrorRule_Name(t *testing.T) {
	r := &monitor.ParseErrorRule{}
	if r.Name() != "ParseErrorRule" {
		t.Errorf("unexpected name: %s", r.Name())
	}
}

func TestParseErrorRule_NoAlert_WhenParsed(t *testing.T) {
	r := &monitor.ParseErrorRule{}
	pq := makePQ("SELECT 1", 0, &parser.SelectStatement{}, nil)
	if alert := r.Check(pq); alert != nil {
		t.Error("expected no alert when statement is parsed")
	}
}

func TestParseErrorRule_Alert_WhenNilStatement(t *testing.T) {
	r := &monitor.ParseErrorRule{}
	pq := makePQ("INVALID SQL !!!", 0, nil, nil)
	alert := r.Check(pq)
	if alert == nil {
		t.Fatal("expected alert for nil statement with non-empty query")
	}
	if alert.Type != "PARSE_ERROR" {
		t.Errorf("expected PARSE_ERROR, got %s", alert.Type)
	}
}

func TestParseErrorRule_NoAlert_EmptyQuery(t *testing.T) {
	r := &monitor.ParseErrorRule{}
	pq := makePQ("", 0, nil, nil)
	if alert := r.Check(pq); alert != nil {
		t.Error("expected no alert for empty query")
	}
}

// --- FullTableScanRule ---

func TestFullTableScanRule_Name(t *testing.T) {
	r := &monitor.FullTableScanRule{}
	if r.Name() != "FullTableScanRule" {
		t.Errorf("unexpected name: %s", r.Name())
	}
}

func TestFullTableScanRule_NoAlert_NilAnalysis(t *testing.T) {
	r := &monitor.FullTableScanRule{}
	pq := makePQ("SELECT 1", 0, nil, nil)
	if alert := r.Check(pq); alert != nil {
		t.Error("expected no alert with nil analysis")
	}
}

func TestFullTableScanRule_Alert_SelectWithoutWhere(t *testing.T) {
	r := &monitor.FullTableScanRule{}
	analysis := &analyzer.QueryAnalysis{
		Tables: []analyzer.TableInfo{{Name: "users"}},
	}
	pq := makePQ("SELECT * FROM users", 0, &parser.SelectStatement{}, analysis)
	alert := r.Check(pq)
	if alert == nil {
		t.Fatal("expected FULL_TABLE_SCAN alert")
	}
	if alert.Type != "FULL_TABLE_SCAN" {
		t.Errorf("expected FULL_TABLE_SCAN, got %s", alert.Type)
	}
}

func TestFullTableScanRule_NoAlert_SelectWithWhere(t *testing.T) {
	r := &monitor.FullTableScanRule{}
	analysis := &analyzer.QueryAnalysis{
		Tables: []analyzer.TableInfo{{Name: "users"}},
	}
	whereExpr := &parser.BinaryExpression{
		Left:     &parser.ColumnReference{Column: "id"},
		Operator: "=",
		Right:    &parser.Literal{Value: 1},
	}
	pq := makePQ("SELECT * FROM users WHERE id=1", 0, &parser.SelectStatement{Where: whereExpr}, analysis)
	alert := r.Check(pq)
	if alert != nil && alert.Type == "FULL_TABLE_SCAN" {
		t.Error("expected no FULL_TABLE_SCAN when WHERE is present")
	}
}

func TestFullTableScanRule_Alert_UpdateWithoutWhere(t *testing.T) {
	r := &monitor.FullTableScanRule{}
	analysis := &analyzer.QueryAnalysis{
		Tables: []analyzer.TableInfo{{Name: "users"}},
	}
	pq := makePQ("UPDATE users SET name='x'", 0, &parser.UpdateStatement{}, analysis)
	alert := r.Check(pq)
	if alert == nil {
		t.Fatal("expected UNSAFE_UPDATE alert")
	}
	if alert.Type != "UNSAFE_UPDATE" {
		t.Errorf("expected UNSAFE_UPDATE, got %s", alert.Type)
	}
	if alert.Level != monitor.AlertError {
		t.Errorf("expected ERROR level, got %s", alert.Level)
	}
}

func TestFullTableScanRule_Alert_DeleteWithoutWhere(t *testing.T) {
	r := &monitor.FullTableScanRule{}
	analysis := &analyzer.QueryAnalysis{
		Tables: []analyzer.TableInfo{{Name: "users"}},
	}
	pq := makePQ("DELETE FROM users", 0, &parser.DeleteStatement{From: parser.TableReference{Name: "users"}}, analysis)
	alert := r.Check(pq)
	if alert == nil {
		t.Fatal("expected UNSAFE_DELETE alert")
	}
	if alert.Type != "UNSAFE_DELETE" {
		t.Errorf("expected UNSAFE_DELETE, got %s", alert.Type)
	}
}

func TestFullTableScanRule_NoAlert_UpdateWithWhere(t *testing.T) {
	r := &monitor.FullTableScanRule{}
	analysis := &analyzer.QueryAnalysis{
		Tables: []analyzer.TableInfo{{Name: "users"}},
	}
	whereExpr := &parser.BinaryExpression{
		Left:     &parser.ColumnReference{Column: "id"},
		Operator: "=",
		Right:    &parser.Literal{Value: 1},
	}
	stmt := &parser.UpdateStatement{Where: whereExpr}
	pq := makePQ("UPDATE users SET name='x' WHERE id=1", 0, stmt, analysis)
	alert := r.Check(pq)
	if alert != nil && alert.Type == "UNSAFE_UPDATE" {
		t.Error("expected no UNSAFE_UPDATE when WHERE present")
	}
}

// --- OptimizationRule ---

func TestOptimizationRule_Name(t *testing.T) {
	r := &monitor.OptimizationRule{MinSeverity: "low"}
	if r.Name() != "OptimizationRule" {
		t.Errorf("unexpected name: %s", r.Name())
	}
}

func TestOptimizationRule_NoAlert_NilAnalysis(t *testing.T) {
	r := &monitor.OptimizationRule{MinSeverity: "low"}
	pq := makePQ("SELECT 1", 0, nil, nil)
	if alert := r.Check(pq); alert != nil {
		t.Error("expected no alert with nil analysis")
	}
}

func TestOptimizationRule_NoAlert_NoSuggestions(t *testing.T) {
	r := &monitor.OptimizationRule{MinSeverity: "low"}
	analysis := &analyzer.QueryAnalysis{}
	pq := makePQ("SELECT 1", 0, nil, analysis)
	if alert := r.Check(pq); alert != nil {
		t.Error("expected no alert with no suggestions")
	}
}

func TestOptimizationRule_Alert_Low(t *testing.T) {
	r := &monitor.OptimizationRule{MinSeverity: "low"}
	analysis := &analyzer.QueryAnalysis{
		EnhancedSuggestions: []analyzer.EnhancedOptimizationSuggestion{
			{Type: "SELECT_ALL", Severity: "LOW", Description: "avoid select star"},
		},
	}
	pq := makePQ("SELECT * FROM t", 0, nil, analysis)
	alert := r.Check(pq)
	if alert == nil {
		t.Fatal("expected alert for low severity with MinSeverity=low")
	}
	if alert.Type != "OPTIMIZATION_OPPORTUNITY" {
		t.Errorf("expected OPTIMIZATION_OPPORTUNITY, got %s", alert.Type)
	}
}

func TestOptimizationRule_NoAlert_BelowThreshold(t *testing.T) {
	r := &monitor.OptimizationRule{MinSeverity: "high"}
	analysis := &analyzer.QueryAnalysis{
		EnhancedSuggestions: []analyzer.EnhancedOptimizationSuggestion{
			{Type: "SELECT_ALL", Severity: "low", Description: "avoid select star"},
		},
	}
	pq := makePQ("SELECT * FROM t", 0, nil, analysis)
	if alert := r.Check(pq); alert != nil {
		t.Error("expected no alert: suggestion severity below threshold")
	}
}

func TestOptimizationRule_Alert_Medium(t *testing.T) {
	r := &monitor.OptimizationRule{MinSeverity: "medium"}
	analysis := &analyzer.QueryAnalysis{
		EnhancedSuggestions: []analyzer.EnhancedOptimizationSuggestion{
			{Type: "COMPLEX_QUERY", Severity: "medium", Description: "complex"},
		},
	}
	pq := makePQ("SELECT 1", 0, nil, analysis)
	if alert := r.Check(pq); alert == nil {
		t.Error("expected alert for medium severity with MinSeverity=medium")
	}
}

// --- ConsoleAlertHandler (just exercise it, output goes to stdout) ---

func TestConsoleAlertHandler_DoesNotPanic(t *testing.T) {
	alert := &monitor.Alert{
		Level:   monitor.AlertWarning,
		Type:    "SLOW_QUERY",
		Message: "Query took 2.5s",
		Query: &monitor.ProcessedQuery{
			Query:    "SELECT * FROM users",
			Duration: 2.5,
		},
	}
	// Should not panic
	monitor.ConsoleAlertHandler(alert)
}

func TestConsoleAlertHandler_NilQuery(t *testing.T) {
	alert := &monitor.Alert{
		Level:   monitor.AlertInfo,
		Type:    "PARSE_ERROR",
		Message: "failed to parse",
		Query:   nil,
	}
	monitor.ConsoleAlertHandler(alert)
}
