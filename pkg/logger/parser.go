package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Duration  int64     `json:"duration_ms"`
	Database  string    `json:"database"`
	User      string    `json:"user"`
	Query     string    `json:"query"`
	Reads     int64     `json:"logical_reads"`
	Writes    int64     `json:"writes"`
	CPU       int64     `json:"cpu_ms"`
	SPID      int       `json:"spid"`
}

// Package-level compiled regexes — allocated once, reused for every log line.
var (
	reProfiler     = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})\s+(.+)`)
	reErrorLog     = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{2})\s+(\w+)\s+(.+)`)
	rePerfCounter  = regexp.MustCompile(`Duration:\s*(\d+)\s*ms.*CPU:\s*(\d+)\s*ms.*Reads:\s*(\d+).*Writes:\s*(\d+)`)
	reXETimestamp  = regexp.MustCompile(`timestamp="([^"]+)"`)
	reXEDuration   = regexp.MustCompile(`duration="(\d+)"`)
	reXEStatement  = regexp.MustCompile(`statement="([^"]+)"`)
	reXEDatabase   = regexp.MustCompile(`database_name="([^"]+)"`)
)

type SQLServerLogParser struct{}

func NewSQLServerLogParser() *SQLServerLogParser {
	return &SQLServerLogParser{}
}

func (p *SQLServerLogParser) ParseLog(reader io.Reader) ([]LogEntry, error) {
	var entries []LogEntry
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Try to detect the log format and parse accordingly
		if entry := p.parseProfilerLine(line); entry != nil {
			entries = append(entries, *entry)
		} else if entry := p.parseExtendedEventsLine(line); entry != nil {
			entries = append(entries, *entry)
		} else if entry := p.parseQueryStoreLine(line); entry != nil {
			entries = append(entries, *entry)
		} else if entry := p.parseErrorLogLine(line); entry != nil {
			entries = append(entries, *entry)
		}
		// For continuation lines, we'd need more sophisticated logic
		// This is a simplified version
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	return entries, nil
}

func (p *SQLServerLogParser) parseProfilerLine(line string) *LogEntry {
	// Example: 2024-01-01 10:30:45.123 SQL:BatchCompleted SELECT * FROM Users
	matches := reProfiler.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05.000", matches[1])
	if err != nil {
		return nil
	}

	content := matches[2]

	entry := &LogEntry{
		Timestamp: timestamp,
	}

	// Extract performance metrics if present
	if perfMatch := rePerfCounter.FindStringSubmatch(content); len(perfMatch) >= 5 {
		if duration, err := strconv.ParseInt(perfMatch[1], 10, 64); err == nil {
			entry.Duration = duration
		}
		if cpu, err := strconv.ParseInt(perfMatch[2], 10, 64); err == nil {
			entry.CPU = cpu
		}
		if reads, err := strconv.ParseInt(perfMatch[3], 10, 64); err == nil {
			entry.Reads = reads
		}
		if writes, err := strconv.ParseInt(perfMatch[4], 10, 64); err == nil {
			entry.Writes = writes
		}
	}

	// Extract SQL query
	if sqlIndex := strings.Index(content, "SELECT"); sqlIndex != -1 {
		entry.Query = strings.TrimSpace(content[sqlIndex:])
	} else if sqlIndex := strings.Index(content, "INSERT"); sqlIndex != -1 {
		entry.Query = strings.TrimSpace(content[sqlIndex:])
	} else if sqlIndex := strings.Index(content, "UPDATE"); sqlIndex != -1 {
		entry.Query = strings.TrimSpace(content[sqlIndex:])
	} else if sqlIndex := strings.Index(content, "DELETE"); sqlIndex != -1 {
		entry.Query = strings.TrimSpace(content[sqlIndex:])
	}

	return entry
}

func (p *SQLServerLogParser) parseExtendedEventsLine(line string) *LogEntry {
	// Basic XML parsing for Extended Events
	if !strings.Contains(line, "sql_statement_completed") {
		return nil
	}

	entry := &LogEntry{}

	if matches := reXETimestamp.FindStringSubmatch(line); len(matches) >= 2 {
		if timestamp, err := time.Parse(time.RFC3339, matches[1]); err == nil {
			entry.Timestamp = timestamp
		}
	}

	if matches := reXEDuration.FindStringSubmatch(line); len(matches) >= 2 {
		if duration, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			entry.Duration = duration / 1000 // microseconds → milliseconds
		}
	}

	if matches := reXEStatement.FindStringSubmatch(line); len(matches) >= 2 {
		entry.Query = matches[1]
	}

	if matches := reXEDatabase.FindStringSubmatch(line); len(matches) >= 2 {
		entry.Database = matches[1]
	}

	return entry
}

func (p *SQLServerLogParser) parseQueryStoreLine(line string) *LogEntry {
	// Parse JSON format from Query Store
	if !strings.Contains(line, "query_sql_text") {
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(line), &data); err != nil {
		return nil
	}

	entry := &LogEntry{}

	if sqlText, ok := data["query_sql_text"].(string); ok {
		entry.Query = sqlText
	}

	if duration, ok := data["avg_duration"].(float64); ok {
		entry.Duration = int64(duration)
	}

	if reads, ok := data["avg_logical_io_reads"].(float64); ok {
		entry.Reads = int64(reads)
	}

	if writes, ok := data["avg_logical_io_writes"].(float64); ok {
		entry.Writes = int64(writes)
	}

	return entry
}

func (p *SQLServerLogParser) parseErrorLogLine(line string) *LogEntry {
	// Parse general SQL Server error log format
	matches := reErrorLog.FindStringSubmatch(line)
	if len(matches) < 4 {
		return nil
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05.00", matches[1])
	if err != nil {
		return nil
	}

	logLevel := matches[2]
	content := matches[3]

	// Only process lines that might contain SQL queries
	if !p.looksLikeSQL(content) {
		return nil
	}

	entry := &LogEntry{
		Timestamp: timestamp,
		Query:     content,
	}

	// Add log level as metadata (could be extended)
	if logLevel == "Error" {
		// Handle error entries differently if needed
	}

	return entry
}

func (p *SQLServerLogParser) looksLikeSQL(line string) bool {
	upperLine := strings.ToUpper(line)
	sqlKeywords := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "ALTER", "EXEC", "EXECUTE"}

	for _, keyword := range sqlKeywords {
		if strings.Contains(upperLine, keyword) {
			return true
		}
	}

	return false
}

// ParseLogFile parses a log file by path.
func (p *SQLServerLogParser) ParseLogFile(filename string) ([]LogEntry, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()
	return p.ParseLog(f)
}

// FilterEntries filters log entries based on criteria
func FilterEntries(entries []LogEntry, criteria FilterCriteria) []LogEntry {
	var filtered []LogEntry

	for _, entry := range entries {
		if matchesCriteria(entry, criteria) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

type FilterCriteria struct {
	MinDuration time.Duration
	MaxDuration time.Duration
	Database    string
	User        string
	QueryType   string // SELECT, INSERT, UPDATE, DELETE
	MinReads    int64
	MaxReads    int64
}

func matchesCriteria(entry LogEntry, criteria FilterCriteria) bool {
	if criteria.MinDuration > 0 && time.Duration(entry.Duration)*time.Millisecond < criteria.MinDuration {
		return false
	}

	if criteria.MaxDuration > 0 && time.Duration(entry.Duration)*time.Millisecond > criteria.MaxDuration {
		return false
	}

	if criteria.Database != "" && !strings.EqualFold(entry.Database, criteria.Database) {
		return false
	}

	if criteria.User != "" && !strings.EqualFold(entry.User, criteria.User) {
		return false
	}

	if criteria.QueryType != "" {
		upperQuery := strings.ToUpper(entry.Query)
		if !strings.HasPrefix(strings.TrimSpace(upperQuery), criteria.QueryType) {
			return false
		}
	}

	if criteria.MinReads > 0 && entry.Reads < criteria.MinReads {
		return false
	}

	if criteria.MaxReads > 0 && entry.Reads > criteria.MaxReads {
		return false
	}

	return true
}
