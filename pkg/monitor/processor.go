package monitor

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Chahine-tech/sql-parser-go/pkg/analyzer"
	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// LogProcessor processes log lines in real-time
type LogProcessor struct {
	dialectName  string
	dialect      dialect.Dialect
	queryHandler func(*ProcessedQuery)
	stats        *Statistics
	mu           sync.RWMutex
}

// ProcessedQuery represents a parsed query from a log
type ProcessedQuery struct {
	Timestamp    time.Time
	Query        string
	Duration     float64
	RowsAffected int64
	Database     string
	User         string

	// Parsed information
	Statement parser.Statement
	Analysis  *analyzer.QueryAnalysis

	// Log metadata
	LogFormat string
	Severity  string
}

// NewLogProcessor creates a new log processor
func NewLogProcessor(dialectName string) *LogProcessor {
	return &LogProcessor{
		dialectName: dialectName,
		dialect:     dialect.GetDialect(dialectName),
		stats:       NewStatistics(),
	}
}

// SetQueryHandler sets the callback for processed queries
func (p *LogProcessor) SetQueryHandler(handler func(*ProcessedQuery)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.queryHandler = handler
}

// Start begins processing log lines from the channel
func (p *LogProcessor) Start(ctx context.Context, lines <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-lines:
			if !ok {
				return // Channel closed
			}
			p.processLine(line)
		}
	}
}

// processLine processes a single log line
func (p *LogProcessor) processLine(line string) {
	// Skip empty lines
	if strings.TrimSpace(line) == "" {
		return
	}

	// Simple log parsing - extract SQL query from line
	// This is a basic implementation - can be enhanced with proper log format parsers
	query := p.extractQueryFromLine(line)
	if query == "" {
		p.stats.IncrementSkipped()
		return
	}

	// Create processed query
	pq := &ProcessedQuery{
		Timestamp:    time.Now(),
		Query:        query,
		Duration:     0, // Would be extracted from log format
		RowsAffected: 0,
		Database:     "",
		User:         "",
		LogFormat:    "generic",
		Severity:     "INFO",
	}

	// Parse the SQL query
	ctx := context.Background()
	sqlParser := parser.NewWithDialect(ctx, query, p.dialect)
	stmt, err := sqlParser.ParseStatement()
	if err != nil {
		// Failed to parse, but still record it
		p.stats.IncrementParseFailed()
		pq.Statement = nil
	} else {
		pq.Statement = stmt
		p.stats.IncrementParsed()
	}

	// Analyze the query if parsing succeeded
	if stmt != nil && err == nil {
		a := analyzer.NewWithDialect(p.dialect)
		analysis := a.Analyze(stmt)
		pq.Analysis = &analysis
	}

	// Update statistics
	p.stats.RecordQuery(pq)

	// Call handler if set
	p.mu.RLock()
	handler := p.queryHandler
	p.mu.RUnlock()

	if handler != nil {
		handler(pq)
	}
}

// GetStatistics returns current processing statistics
func (p *LogProcessor) GetStatistics() *Statistics {
	return p.stats
}

// Statistics tracks processing statistics
type Statistics struct {
	mu sync.RWMutex

	TotalLines    int64
	ParsedQueries int64
	FailedParses  int64
	SkippedLines  int64

	// Query performance
	TotalDuration float64
	SlowQueries   int64
	SlowThreshold float64

	// By query type
	SelectCount int64
	InsertCount int64
	UpdateCount int64
	DeleteCount int64
	OtherCount  int64

	// Timing
	StartTime     time.Time
	LastQueryTime time.Time
}

// NewStatistics creates a new statistics tracker
func NewStatistics() *Statistics {
	return &Statistics{
		StartTime:     time.Now(),
		SlowThreshold: 1.0, // Default: 1 second
	}
}

// SetSlowThreshold sets the threshold for slow queries in seconds
func (s *Statistics) SetSlowThreshold(threshold float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SlowThreshold = threshold
}

// IncrementParsed increments the parsed queries counter
func (s *Statistics) IncrementParsed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ParsedQueries++
}

// IncrementParseFailed increments the failed parse counter
func (s *Statistics) IncrementParseFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.FailedParses++
}

// IncrementSkipped increments the skipped lines counter
func (s *Statistics) IncrementSkipped() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SkippedLines++
}

// RecordQuery records statistics for a processed query
func (s *Statistics) RecordQuery(pq *ProcessedQuery) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TotalLines++
	s.TotalDuration += pq.Duration
	s.LastQueryTime = pq.Timestamp

	// Check if slow
	if pq.Duration >= s.SlowThreshold {
		s.SlowQueries++
	}

	// Count by type
	if pq.Statement != nil {
		switch pq.Statement.(type) {
		case *parser.SelectStatement:
			s.SelectCount++
		case *parser.InsertStatement:
			s.InsertCount++
		case *parser.UpdateStatement:
			s.UpdateCount++
		case *parser.DeleteStatement:
			s.DeleteCount++
		default:
			s.OtherCount++
		}
	}
}

// GetSnapshot returns a snapshot of current statistics
func (s *Statistics) GetSnapshot() StatSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return StatSnapshot{
		TotalLines:    s.TotalLines,
		ParsedQueries: s.ParsedQueries,
		FailedParses:  s.FailedParses,
		SkippedLines:  s.SkippedLines,
		TotalDuration: s.TotalDuration,
		SlowQueries:   s.SlowQueries,
		SlowThreshold: s.SlowThreshold,
		SelectCount:   s.SelectCount,
		InsertCount:   s.InsertCount,
		UpdateCount:   s.UpdateCount,
		DeleteCount:   s.DeleteCount,
		OtherCount:    s.OtherCount,
		StartTime:     s.StartTime,
		LastQueryTime: s.LastQueryTime,
		Uptime:        time.Since(s.StartTime),
	}
}

// StatSnapshot is a point-in-time snapshot of statistics
type StatSnapshot struct {
	TotalLines    int64
	ParsedQueries int64
	FailedParses  int64
	SkippedLines  int64
	TotalDuration float64
	SlowQueries   int64
	SlowThreshold float64
	SelectCount   int64
	InsertCount   int64
	UpdateCount   int64
	DeleteCount   int64
	OtherCount    int64
	StartTime     time.Time
	LastQueryTime time.Time
	Uptime        time.Duration
}

// String returns a formatted string of the statistics
func (s StatSnapshot) String() string {
	avgDuration := 0.0
	if s.ParsedQueries > 0 {
		avgDuration = s.TotalDuration / float64(s.ParsedQueries)
	}

	return fmt.Sprintf(`Statistics:
  Total Lines:     %d
  Parsed Queries:  %d
  Failed Parses:   %d
  Skipped Lines:   %d

  Query Types:
    SELECT:        %d
    INSERT:        %d
    UPDATE:        %d
    DELETE:        %d
    OTHER:         %d

  Performance:
    Total Duration: %.2fs
    Avg Duration:   %.4fs
    Slow Queries:   %d (threshold: %.2fs)

  Timing:
    Uptime:         %s
    Last Query:     %s`,
		s.TotalLines,
		s.ParsedQueries,
		s.FailedParses,
		s.SkippedLines,
		s.SelectCount,
		s.InsertCount,
		s.UpdateCount,
		s.DeleteCount,
		s.OtherCount,
		s.TotalDuration,
		avgDuration,
		s.SlowQueries,
		s.SlowThreshold,
		s.Uptime.Round(time.Second),
		s.LastQueryTime.Format("2006-01-02 15:04:05"),
	)
}
