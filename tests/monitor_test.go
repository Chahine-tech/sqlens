package tests

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Chahine-tech/sql-parser-go/pkg/monitor"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func TestLogWatcher(t *testing.T) {
	// Create a temporary log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// Write initial content
	err := os.WriteFile(logFile, []byte("SELECT * FROM users\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Create watcher
	watcher := monitor.NewLogWatcher(logFile)
	lines := make(chan string, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start watching
	err = watcher.Start(ctx, lines)
	if err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Append a new line
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open log file for writing: %v", err)
	}
	_, err = f.WriteString("INSERT INTO orders (id, total) VALUES (1, 100)\n")
	f.Close()
	if err != nil {
		t.Fatalf("Failed to write to log file: %v", err)
	}

	// Wait for the new line
	select {
	case line := <-lines:
		if line != "INSERT INTO orders (id, total) VALUES (1, 100)" {
			t.Errorf("Expected INSERT query, got: %s", line)
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for new line")
	}

	cancel()
	watcher.Stop()
}

func TestLogWatcherWithTail(t *testing.T) {
	// Create a temporary log file with multiple lines
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test_tail.log")

	content := `SELECT * FROM users WHERE id = 1
SELECT * FROM orders WHERE id = 2
SELECT * FROM products WHERE id = 3
SELECT * FROM categories WHERE id = 4
SELECT * FROM inventory WHERE id = 5
`
	err := os.WriteFile(logFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Create watcher
	watcher := monitor.NewLogWatcher(logFile)
	lines := make(chan string, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start watching with tail of 3 lines
	err = watcher.StartWithTail(ctx, lines, 3)
	if err != nil {
		t.Fatalf("Failed to start watcher with tail: %v", err)
	}

	// Collect tailed lines
	tailedLines := []string{}
	timeout := time.After(1 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case line := <-lines:
			tailedLines = append(tailedLines, line)
		case <-timeout:
			t.Fatalf("Timeout waiting for tailed lines, got %d lines", len(tailedLines))
		}
	}

	// Should get the last 3 lines
	if len(tailedLines) != 3 {
		t.Errorf("Expected 3 tailed lines, got %d", len(tailedLines))
	}

	// The last line should contain "inventory"
	found := false
	for _, line := range tailedLines {
		if contains(line, "inventory") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find 'inventory' in tailed lines, got: %v", tailedLines)
	}

	cancel()
	watcher.Stop()
}

func TestLogProcessor(t *testing.T) {
	processor := monitor.NewLogProcessor("mysql")

	// Track processed queries (atomic to avoid race conditions)
	var processedCount int64
	processor.SetQueryHandler(func(pq *monitor.ProcessedQuery) {
		atomic.AddInt64(&processedCount, 1)
		if pq.Query == "" {
			t.Error("Processed query should not be empty")
		}
		if pq.Statement == nil {
			t.Logf("Query failed to parse: %s", pq.Query)
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lines := make(chan string, 10)

	// Start processor
	go processor.Start(ctx, lines)

	// Send test queries
	testQueries := []string{
		"SELECT * FROM users WHERE id = 1",
		"INSERT INTO orders (id, total) VALUES (1, 100)",
		"UPDATE products SET price = 99.99 WHERE id = 1",
		"DELETE FROM sessions WHERE expires_at < NOW()",
		"",                // Empty line should be skipped
		"This is not SQL", // Should be skipped
		"SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id",
	}

	for _, query := range testQueries {
		lines <- query
	}

	// Give processor time to process
	time.Sleep(500 * time.Millisecond)

	// Should have processed the queries (some may fail to parse due to incomplete SQL)
	count := atomic.LoadInt64(&processedCount)
	if count < 4 {
		t.Errorf("Expected at least 4 processed queries, got %d", count)
	}

	// Check statistics - total lines should include all processed queries
	stats := processor.GetStatistics().GetSnapshot()
	if stats.TotalLines < 4 {
		t.Errorf("Expected at least 4 total lines processed, got %d", stats.TotalLines)
	}

	cancel()
}

func TestStatistics(t *testing.T) {
	stats := monitor.NewStatistics()

	// Set slow threshold
	stats.SetSlowThreshold(1.0)

	// Record some queries
	for i := 0; i < 5; i++ {
		pq := &monitor.ProcessedQuery{
			Timestamp: time.Now(),
			Query:     "SELECT * FROM test",
			Duration:  float64(i) * 0.5, // 0, 0.5, 1.0, 1.5, 2.0
		}
		stats.RecordQuery(pq)
	}

	snapshot := stats.GetSnapshot()

	if snapshot.TotalLines != 5 {
		t.Errorf("Expected 5 total lines, got %d", snapshot.TotalLines)
	}

	// Queries with duration >= 1.0: 1.0, 1.5, 2.0 = 3 slow queries
	if snapshot.SlowQueries != 3 {
		t.Errorf("Expected 3 slow queries (>= 1.0s), got %d", snapshot.SlowQueries)
	}
}

func TestAlertManager(t *testing.T) {
	alertMgr := monitor.NewAlertManager()

	// Add slow query rule
	alertMgr.AddRule(&monitor.SlowQueryRule{Threshold: 1.0})

	// Track triggered alerts
	triggeredAlerts := []*monitor.Alert{}
	alertMgr.AddHandler(func(alert *monitor.Alert) {
		triggeredAlerts = append(triggeredAlerts, alert)
	})

	// Test with slow query
	slowQuery := &monitor.ProcessedQuery{
		Timestamp: time.Now(),
		Query:     "SELECT * FROM large_table",
		Duration:  2.5,
	}
	alertMgr.Check(slowQuery)

	if len(triggeredAlerts) != 1 {
		t.Errorf("Expected 1 alert for slow query, got %d", len(triggeredAlerts))
	}

	if triggeredAlerts[0].Type != "SLOW_QUERY" {
		t.Errorf("Expected SLOW_QUERY alert type, got %s", triggeredAlerts[0].Type)
	}

	// Test with fast query
	fastQuery := &monitor.ProcessedQuery{
		Timestamp: time.Now(),
		Query:     "SELECT * FROM users WHERE id = 1",
		Duration:  0.05,
	}
	alertMgr.Check(fastQuery)

	// Should still be 1 alert (fast query doesn't trigger)
	if len(triggeredAlerts) != 1 {
		t.Errorf("Expected still 1 alert, got %d", len(triggeredAlerts))
	}
}

func TestSlowQueryRule(t *testing.T) {
	rule := &monitor.SlowQueryRule{Threshold: 1.0}

	// Test query below threshold
	fastQuery := &monitor.ProcessedQuery{
		Duration: 0.5,
	}
	alert := rule.Check(fastQuery)
	if alert != nil {
		t.Error("Expected no alert for fast query")
	}

	// Test query at threshold
	thresholdQuery := &monitor.ProcessedQuery{
		Duration: 1.0,
	}
	alert = rule.Check(thresholdQuery)
	if alert == nil {
		t.Error("Expected alert for query at threshold")
	}
	if alert.Level != monitor.AlertWarning {
		t.Errorf("Expected WARNING level, got %v", alert.Level)
	}

	// Test very slow query
	verySlowQuery := &monitor.ProcessedQuery{
		Duration: 10.0, // 10x threshold
	}
	alert = rule.Check(verySlowQuery)
	if alert == nil {
		t.Error("Expected alert for very slow query")
	}
	if alert.Level != monitor.AlertCritical {
		t.Errorf("Expected CRITICAL level for 10x threshold, got %v", alert.Level)
	}
}

func TestParseErrorRule(t *testing.T) {
	rule := &monitor.ParseErrorRule{}

	// Test query with parse error (Statement is nil)
	errorQuery := &monitor.ProcessedQuery{
		Query:     "INVALID SQL SYNTAX!!!",
		Statement: nil,
	}
	alert := rule.Check(errorQuery)
	if alert == nil {
		t.Error("Expected alert for parse error")
	}
	if alert.Type != "PARSE_ERROR" {
		t.Errorf("Expected PARSE_ERROR type, got %s", alert.Type)
	}

	// Test valid query (Statement is not nil - we'll use a real SelectStatement)
	validQuery := &monitor.ProcessedQuery{
		Query:     "SELECT * FROM users",
		Statement: &parser.SelectStatement{}, // Real statement type
	}
	alert = rule.Check(validQuery)
	if alert != nil {
		t.Error("Expected no alert for valid query")
	}
}

func TestQueryExtraction(t *testing.T) {
	processor := monitor.NewLogProcessor("mysql")

	testCases := []struct {
		name     string
		logLine  string
		expected bool // whether we expect to extract a query
	}{
		{
			name:     "Direct SQL",
			logLine:  "SELECT * FROM users",
			expected: true,
		},
		{
			name:     "MySQL log format",
			logLine:  "2024-01-15T10:30:45.123Z  123  Query  SELECT * FROM orders",
			expected: true,
		},
		{
			name:     "INSERT statement",
			logLine:  "INSERT INTO products (name, price) VALUES ('Test', 99.99)",
			expected: true,
		},
		{
			name:     "UPDATE statement",
			logLine:  "UPDATE users SET status = 'active' WHERE id = 1",
			expected: true,
		},
		{
			name:     "DELETE statement",
			logLine:  "DELETE FROM logs WHERE created_at < '2024-01-01'",
			expected: true,
		},
		{
			name:     "Non-SQL line",
			logLine:  "This is just a regular log message",
			expected: false,
		},
		{
			name:     "Empty line",
			logLine:  "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Process the line (use atomic int32 for thread-safe access)
			var queryExtracted int32
			processor.SetQueryHandler(func(pq *monitor.ProcessedQuery) {
				if pq.Query != "" {
					atomic.StoreInt32(&queryExtracted, 1)
				}
			})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			lines := make(chan string, 1)
			go processor.Start(ctx, lines)

			lines <- tc.logLine
			time.Sleep(100 * time.Millisecond)

			extracted := atomic.LoadInt32(&queryExtracted) == 1
			if extracted != tc.expected {
				t.Errorf("Query extraction mismatch: expected %v, got %v for line: %s",
					tc.expected, extracted, tc.logLine)
			}

			cancel()
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
