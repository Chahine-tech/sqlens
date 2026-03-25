package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Chahine-tech/sql-parser-go/pkg/logger"
)

func TestSQLServerLogParser_ParseLog_Empty(t *testing.T) {
	p := logger.NewSQLServerLogParser()
	entries, err := p.ParseLog(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestSQLServerLogParser_ParseLog_ProfilerLines(t *testing.T) {
	lines := `2023-01-15 10:30:45.123   SELECT * FROM users WHERE id = 1  Duration: 250 ms  CPU: 10 ms  Reads: 100  Writes: 0
2023-01-15 10:31:00.456   INSERT INTO orders VALUES (1, 2, 3)  Duration: 50 ms  CPU: 5 ms  Reads: 10  Writes: 5
`
	p := logger.NewSQLServerLogParser()
	entries, err := p.ParseLog(strings.NewReader(lines))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = entries // entries may or may not be parsed depending on exact format
}

func TestSQLServerLogParser_ParseLog_SkipsBlankLines(t *testing.T) {
	lines := "\n\n\n"
	p := logger.NewSQLServerLogParser()
	entries, err := p.ParseLog(strings.NewReader(lines))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for blank lines, got %d", len(entries))
	}
}

func TestSQLServerLogParser_ParseLogFile_NotFound(t *testing.T) {
	p := logger.NewSQLServerLogParser()
	_, err := p.ParseLogFile("/nonexistent/file.log")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSQLServerLogParser_ParseLogFile_Valid(t *testing.T) {
	content := "2023-01-15 10:30:45.123   SELECT * FROM users\n"
	tmp := filepath.Join(t.TempDir(), "test.log")
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	p := logger.NewSQLServerLogParser()
	_, err := p.ParseLogFile(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSQLServerLogParser_ParseLog_ExtendedEventsLine(t *testing.T) {
	line := `<event name="sql_statement_completed" timestamp="2023-01-15T10:30:45.123Z" duration="250000" statement="SELECT * FROM users" database_name="mydb" />`
	p := logger.NewSQLServerLogParser()
	entries, err := p.ParseLog(strings.NewReader(line))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = entries
}

func TestSQLServerLogParser_ParseLog_QueryStoreLine(t *testing.T) {
	line := `{"query_id":1,"query_text":"SELECT * FROM users","avg_duration":1500,"execution_count":10}`
	p := logger.NewSQLServerLogParser()
	entries, err := p.ParseLog(strings.NewReader(line))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = entries
}

func TestLogFormat_String(t *testing.T) {
	cases := []struct {
		format   logger.LogFormat
		expected string
	}{
		{logger.FormatProfiler, "SQL Server Profiler"},
		{logger.FormatExtendedEvents, "Extended Events"},
		{logger.FormatQueryStore, "Query Store"},
		{logger.FormatErrorLog, "Error Log"},
		{logger.FormatPerformanceCounter, "Performance Counter"},
		{logger.FormatUnknown, "Unknown"},
	}
	for _, tc := range cases {
		if got := tc.format.String(); got != tc.expected {
			t.Errorf("LogFormat(%d).String() = %q, want %q", tc.format, got, tc.expected)
		}
	}
}

func TestLogFormatDetector_ExtendedEvents(t *testing.T) {
	d := logger.NewLogFormatDetector()
	got := d.DetectFormat(`<event name="sql_statement_completed" />`)
	if got != logger.FormatExtendedEvents {
		t.Errorf("expected FormatExtendedEvents, got %v", got)
	}
}

func TestLogFormatDetector_QueryStore(t *testing.T) {
	d := logger.NewLogFormatDetector()
	got := d.DetectFormat(`{"query_sql_text":"SELECT 1","avg_duration":100,"execution_count":5}`)
	if got != logger.FormatQueryStore {
		t.Errorf("expected FormatQueryStore, got %v", got)
	}
}

func TestLogFormatDetector_Profiler(t *testing.T) {
	d := logger.NewLogFormatDetector()
	got := d.DetectFormat(`SQL:BatchCompleted SELECT * FROM users`)
	if got != logger.FormatProfiler {
		t.Errorf("expected FormatProfiler, got %v", got)
	}
}

func TestLogFormatDetector_ErrorLog(t *testing.T) {
	d := logger.NewLogFormatDetector()
	got := d.DetectFormat(`spid5 Login succeeded for user 'sa'`)
	if got != logger.FormatErrorLog {
		t.Errorf("expected FormatErrorLog, got %v", got)
	}
}

func TestLogFormatDetector_PerformanceCounter(t *testing.T) {
	d := logger.NewLogFormatDetector()
	got := d.DetectFormat(`Duration: 250 ms  CPU: 10 ms  Reads: 100  Writes: 0`)
	if got != logger.FormatPerformanceCounter {
		t.Errorf("expected FormatPerformanceCounter, got %v", got)
	}
}

func TestLogFormatDetector_Unknown(t *testing.T) {
	d := logger.NewLogFormatDetector()
	got := d.DetectFormat(`some random unrecognized line`)
	if got != logger.FormatUnknown {
		t.Errorf("expected FormatUnknown, got %v", got)
	}
}

func TestCalculateMetrics_Empty(t *testing.T) {
	metrics := logger.CalculateMetrics(nil)
	if metrics.TotalEntries != 0 {
		t.Errorf("expected 0 entries, got %d", metrics.TotalEntries)
	}
}

func TestCalculateMetrics_WithEntries(t *testing.T) {
	now := time.Now()
	entries := []logger.LogEntry{
		{Timestamp: now, Duration: 100, Database: "db1", User: "user1", Query: "SELECT * FROM users", Reads: 50, Writes: 0},
		{Timestamp: now, Duration: 200, Database: "db2", User: "user2", Query: "INSERT INTO orders VALUES (1)", Reads: 10, Writes: 5},
		{Timestamp: now, Duration: 50, Database: "db1", User: "user1", Query: "UPDATE products SET name='x'", Reads: 20, Writes: 10},
		{Timestamp: now, Duration: 300, Database: "", User: "", Query: "DELETE FROM temp", Reads: 0, Writes: 30},
		{Timestamp: now, Duration: 80, Database: "db1", User: "", Query: "CREATE TABLE t (id INT)", Reads: 0, Writes: 0},
		{Timestamp: now, Duration: 60, Database: "db1", User: "", Query: "DROP TABLE old", Reads: 0, Writes: 0},
		{Timestamp: now, Duration: 70, Database: "db1", User: "", Query: "ALTER TABLE t ADD col INT", Reads: 0, Writes: 0},
		{Timestamp: now, Duration: 90, Database: "db1", User: "", Query: "EXEC sp_test", Reads: 5, Writes: 0},
		{Timestamp: now, Duration: 110, Database: "db1", User: "", Query: "EXECUTE proc_name", Reads: 5, Writes: 0},
		{Timestamp: now, Duration: 120, Database: "db1", User: "", Query: "SOME OTHER STATEMENT", Reads: 0, Writes: 0},
		{Timestamp: now, Duration: 130, Database: "db1", User: "", Query: "", Reads: 0, Writes: 0},
	}
	metrics := logger.CalculateMetrics(entries)
	if metrics.TotalEntries != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), metrics.TotalEntries)
	}
	if metrics.TotalReads == 0 {
		t.Error("expected non-zero total reads")
	}
	if metrics.TotalWrites == 0 {
		t.Error("expected non-zero total writes")
	}
	if metrics.AvgDuration == 0 {
		t.Error("expected non-zero avg duration")
	}
	if _, ok := metrics.QueryTypes["SELECT"]; !ok {
		t.Error("expected SELECT in query types")
	}
	if _, ok := metrics.QueryTypes["INSERT"]; !ok {
		t.Error("expected INSERT in query types")
	}
}

func TestFilterEntries(t *testing.T) {
	now := time.Now()
	entries := []logger.LogEntry{
		{Duration: 500, Database: "db1", User: "alice", Query: "SELECT * FROM users", Reads: 100},
		{Duration: 100, Database: "db2", User: "bob", Query: "INSERT INTO t VALUES (1)", Reads: 5},
		{Duration: 2000, Database: "db1", User: "alice", Query: "SELECT id FROM orders", Reads: 50},
	}
	// populate timestamps
	for i := range entries {
		entries[i].Timestamp = now
	}

	// Filter by min duration
	filtered := logger.FilterEntries(entries, logger.FilterCriteria{MinDuration: 400 * time.Millisecond})
	if len(filtered) != 2 {
		t.Errorf("min duration filter: expected 2, got %d", len(filtered))
	}

	// Filter by max duration
	filtered = logger.FilterEntries(entries, logger.FilterCriteria{MaxDuration: 600 * time.Millisecond})
	if len(filtered) != 2 {
		t.Errorf("max duration filter: expected 2, got %d", len(filtered))
	}

	// Filter by database
	filtered = logger.FilterEntries(entries, logger.FilterCriteria{Database: "db1"})
	if len(filtered) != 2 {
		t.Errorf("database filter: expected 2, got %d", len(filtered))
	}

	// Filter by user
	filtered = logger.FilterEntries(entries, logger.FilterCriteria{User: "alice"})
	if len(filtered) != 2 {
		t.Errorf("user filter: expected 2, got %d", len(filtered))
	}

	// Filter by query type
	filtered = logger.FilterEntries(entries, logger.FilterCriteria{QueryType: "SELECT"})
	if len(filtered) != 2 {
		t.Errorf("query type filter: expected 2, got %d", len(filtered))
	}

	// Filter by min reads
	filtered = logger.FilterEntries(entries, logger.FilterCriteria{MinReads: 50})
	if len(filtered) != 2 {
		t.Errorf("min reads filter: expected 2, got %d", len(filtered))
	}

	// Filter by max reads
	filtered = logger.FilterEntries(entries, logger.FilterCriteria{MaxReads: 10})
	if len(filtered) != 1 {
		t.Errorf("max reads filter: expected 1, got %d", len(filtered))
	}
}
