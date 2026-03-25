package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Chahine-tech/sql-parser-go/pkg/monitor"
)

func TestStatSnapshot_String(t *testing.T) {
	snap := monitor.StatSnapshot{
		TotalLines:    100,
		ParsedQueries: 80,
		FailedParses:  5,
		SkippedLines:  15,
		TotalDuration: 120.5,
		SlowQueries:   3,
		SlowThreshold: 1.0,
		SelectCount:   50,
		InsertCount:   15,
		UpdateCount:   10,
		DeleteCount:   5,
		OtherCount:    0,
		StartTime:     time.Now().Add(-5 * time.Minute),
		LastQueryTime: time.Now(),
		Uptime:        5 * time.Minute,
	}
	s := snap.String()
	for _, expected := range []string{"Total Lines", "Parsed", "SELECT", "INSERT", "Uptime"} {
		if !strings.Contains(s, expected) {
			t.Errorf("StatSnapshot.String() missing %q", expected)
		}
	}
}

func TestStatSnapshot_String_ZeroParsed(t *testing.T) {
	// avgDuration branch: ParsedQueries == 0
	snap := monitor.StatSnapshot{}
	s := snap.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestLogProcessor_GetStatistics(t *testing.T) {
	p := monitor.NewLogProcessor("mysql")
	stats := p.GetStatistics()
	if stats == nil {
		t.Fatal("expected non-nil statistics")
	}
}

func TestLogProcessor_ProcessLines(t *testing.T) {
	p := monitor.NewLogProcessor("mysql")

	var processed []*monitor.ProcessedQuery
	p.SetQueryHandler(func(pq *monitor.ProcessedQuery) {
		processed = append(processed, pq)
	})

	lines := []string{
		"SELECT * FROM users WHERE id = 1",
		"INSERT INTO orders (user_id) VALUES (1)",
		"UPDATE products SET price = 10.0 WHERE id = 5",
		"DELETE FROM temp WHERE created < '2020-01-01'",
		"  ", // blank line — should be skipped
		"not sql at all ???!!!",
	}

	for _, line := range lines {
		p.GetStatistics() // just ensure stats are accessible
		_ = line
	}

	// Exercise Start via a simple channel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lineCh := make(chan string, len(lines))
	for _, l := range lines {
		lineCh <- l
	}
	close(lineCh)
	p.Start(ctx, lineCh)

	snap := p.GetStatistics().GetSnapshot()
	if snap.TotalLines == 0 {
		t.Error("expected some lines processed")
	}
}

func TestStatistics_SetSlowThreshold(t *testing.T) {
	s := monitor.NewStatistics()
	s.SetSlowThreshold(2.0)
	snap := s.GetSnapshot()
	if snap.SlowThreshold != 2.0 {
		t.Errorf("expected threshold 2.0, got %f", snap.SlowThreshold)
	}
}
