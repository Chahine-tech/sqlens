package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Chahine-tech/sql-parser-go/internal/config"
	"github.com/Chahine-tech/sql-parser-go/internal/performance"
	"github.com/Chahine-tech/sql-parser-go/pkg/analyzer"
	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/logger"
	"github.com/Chahine-tech/sql-parser-go/pkg/monitor"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

const banner = `
 ███████╗ ██████╗ ██╗     ███████╗███╗   ██╗███████╗
 ██╔════╝██╔═══██╗██║     ██╔════╝████╗  ██║██╔════╝
 ███████╗██║   ██║██║     █████╗  ██╔██╗ ██║███████╗
 ╚════██║██║▄▄ ██║██║     ██╔══╝  ██║╚██╗██║╚════██║
 ███████║╚██████╔╝███████╗███████╗██║ ╚████║███████║
 ╚══════╝ ╚══▀▀═╝ ╚══════╝╚══════╝╚═╝  ╚═══╝╚══════╝

 Welcome to SQLENS — your multi-dialect SQL analysis tool! 🚀
 Supported: MySQL • PostgreSQL • SQL Server • SQLite • Oracle
`

func main() {
	var (
		queryFile     = flag.String("query", "", "File containing the SQL query")
		queryText     = flag.String("sql", "", "SQL query string")
		logFile       = flag.String("log", "", "SQL Server log file")
		outputFormat  = flag.String("output", "json", "Output format (json, table)")
		verbose       = flag.Bool("verbose", false, "Verbose mode")
		configFile    = flag.String("config", "", "Configuration file path")
		dialectFlag   = flag.String("dialect", "", "SQL dialect (mysql, postgresql, sqlserver, sqlite, oracle)")
		showHelp      = flag.Bool("help", false, "Show help")
		watchMode     = flag.Bool("watch", false, "Watch log file for real-time monitoring")
		tailLines     = flag.Int("tail", 10, "Number of lines to tail when starting watch mode")
		slowThreshold = flag.Float64("slow", 1.0, "Slow query threshold in seconds")
	)
	flag.Parse()

	if *showHelp {
		fmt.Print(banner)
		showUsage()
		return
	}

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Printf("Warning: Could not load config: %v\n", err)
		cfg = config.DefaultConfig()
	}

	if *outputFormat != "json" {
		cfg.Output.Format = *outputFormat
	}

	// Override dialect from command line if provided
	if *dialectFlag != "" {
		cfg.Parser.Dialect = *dialectFlag
	}

	if *queryFile != "" {
		if err := analyzeQueryFile(*queryFile, cfg, *verbose); err != nil {
			fmt.Printf("Error analyzing query file: %v\n", err)
			os.Exit(1)
		}
	} else if *queryText != "" {
		if err := analyzeQueryString(*queryText, cfg, *verbose); err != nil {
			fmt.Printf("Error analyzing query: %v\n", err)
			os.Exit(1)
		}
	} else if *logFile != "" {
		if *watchMode {
			if err := watchLogFile(*logFile, cfg, *verbose, *tailLines, *slowThreshold); err != nil {
				fmt.Printf("Error watching log file: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := parseLogFile(*logFile, cfg, *verbose); err != nil {
				fmt.Printf("Error parsing log file: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		showUsage()
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println("SQL Parser Go - Multi-Dialect SQL Query Analysis Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  sqlparser -query file.sql          Analyze SQL query from file")
	fmt.Println("  sqlparser -sql \"SELECT * FROM...\"   Analyze SQL query from string")
	fmt.Println("  sqlparser -log logfile.log          Parse SQL Server log file")
	fmt.Println("  sqlparser -log logfile.log -watch   Watch log file in real-time")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -output FORMAT    Output format: json, table (default: json)")
	fmt.Println("  -dialect DIALECT  SQL dialect: mysql, postgresql, sqlserver, sqlite, oracle (default: sqlserver)")
	fmt.Println("  -verbose          Enable verbose output")
	fmt.Println("  -config FILE      Configuration file path")
	fmt.Println("  -watch            Enable real-time log monitoring (use with -log)")
	fmt.Println("  -tail N           Number of lines to tail when starting watch (default: 10)")
	fmt.Println("  -slow SECONDS     Slow query threshold in seconds (default: 1.0)")
	fmt.Println("  -help             Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  sqlparser -query complex_query.sql -output json -dialect mysql")
	fmt.Println("  sqlparser -sql \"SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id\" -dialect postgresql")
	fmt.Println("  sqlparser -log sqlserver.log -output table -verbose")
	fmt.Println("  sqlparser -log sqlserver.log -watch -tail 20 -slow 2.0 -dialect mysql")
}

func analyzeQueryFile(filename string, cfg *config.Config, verbose bool) error {
	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	return analyzeQueryString(string(content), cfg, verbose)
}

func analyzeQueryString(sql string, cfg *config.Config, verbose bool) error {
	monitor := performance.NewPerformanceMonitor()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if verbose {
		fmt.Print(banner)
		fmt.Printf("Analyzing SQL query...\n")
		fmt.Printf("Query: %s\n", sql)
		fmt.Printf("Dialect: %s\n\n", cfg.Parser.Dialect)
	}

	// Get the dialect
	d := dialect.GetDialect(cfg.Parser.Dialect)

	// Create parser with dialect
	p := parser.NewWithDialect(ctx, sql, d)
	stmt, err := p.ParseStatement()
	if err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}

	if verbose {
		fmt.Printf("Parsed statement type: %s\n", stmt.Type())

		if metrics := p.GetParseMetrics(); metrics != nil {
			fmt.Printf("Parser metrics: %v\n", metrics)
		}
	}

	// Create analyzer with dialect for enhanced optimization suggestions
	a := analyzer.NewWithDialect(d)
	analysis := a.Analyze(stmt)

	var suggestions []analyzer.OptimizationSuggestion
	var enhancedSuggestions []analyzer.EnhancedOptimizationSuggestion

	if cfg.Analyzer.EnableOptimizations {
		// Get enhanced optimization suggestions
		enhancedSuggestions = a.GetEnhancedOptimizations(stmt)

		// For backward compatibility, also get legacy suggestions
		if selectStmt, ok := stmt.(*parser.SelectStatement); ok {
			suggestions = a.SuggestOptimizations(selectStmt)
		}

		// Add enhanced suggestions to analysis for output
		analysis.EnhancedSuggestions = enhancedSuggestions
		analysis.Suggestions = suggestions
	}

	if verbose {
		metrics := monitor.GetMetrics()
		fmt.Printf("Performance metrics: %v\n", metrics)

		if len(enhancedSuggestions) > 0 {
			fmt.Printf("\nOptimization suggestions found: %d\n", len(enhancedSuggestions))
			for _, suggestion := range enhancedSuggestions {
				fmt.Printf("- [%s] %s: %s\n", suggestion.Severity, suggestion.Type, suggestion.Description)
			}
		}
	}

	return outputAnalysis(analysis, suggestions, cfg)
}

func watchLogFile(filename string, cfg *config.Config, verbose bool, tailLines int, slowThreshold float64) error {
	if verbose {
		fmt.Printf("🔍 Starting real-time log monitoring: %s\n", filename)
		fmt.Printf("Dialect: %s\n", cfg.Parser.Dialect)
		fmt.Printf("Slow query threshold: %.2fs\n", slowThreshold)
		fmt.Printf("Tailing last %d lines...\n\n", tailLines)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create watcher
	watcher := monitor.NewLogWatcher(filename)
	lines := make(chan string, 100) // Buffered channel for log lines

	// Start watching (with tail)
	var err error
	if tailLines > 0 {
		err = watcher.StartWithTail(ctx, lines, tailLines)
	} else {
		err = watcher.Start(ctx, lines)
	}
	if err != nil {
		return fmt.Errorf("failed to start watcher: %w", err)
	}

	// Create alert manager
	alertMgr := monitor.NewAlertManager()

	// Add alert rules
	alertMgr.AddRule(&monitor.SlowQueryRule{Threshold: slowThreshold})
	alertMgr.AddRule(&monitor.ParseErrorRule{})
	alertMgr.AddRule(&monitor.OptimizationRule{MinSeverity: "medium"})
	alertMgr.AddRule(&monitor.FullTableScanRule{})

	// Add console alert handler
	alertMgr.AddHandler(monitor.ConsoleAlertHandler)

	// Create processor
	processor := monitor.NewLogProcessor(cfg.Parser.Dialect)
	processor.SetQueryHandler(func(pq *monitor.ProcessedQuery) {
		// Check alerts first
		alertMgr.Check(pq)

		// Print query information
		fmt.Printf("[%s] Duration: %.3fs | Database: %s | User: %s\n",
			pq.Timestamp.Format("15:04:05"),
			pq.Duration,
			pq.Database,
			pq.User)

		if len(pq.Query) > 100 {
			fmt.Printf("  Query: %s...\n", pq.Query[:97])
		} else {
			fmt.Printf("  Query: %s\n", pq.Query)
		}

		// Show analysis if available
		if pq.Analysis != nil {
			if len(pq.Analysis.Tables) > 0 {
				tables := make([]string, len(pq.Analysis.Tables))
				for i, t := range pq.Analysis.Tables {
					tables[i] = t.Name
				}
				fmt.Printf("  Tables: %s\n", strings.Join(tables, ", "))
			}

			// Show optimizations if any
			if len(pq.Analysis.EnhancedSuggestions) > 0 {
				fmt.Printf("  ⚠️  %d optimization suggestions\n", len(pq.Analysis.EnhancedSuggestions))
				for _, opt := range pq.Analysis.EnhancedSuggestions {
					fmt.Printf("    - [%s] %s\n", opt.Severity, opt.Description)
				}
			}
		}

		fmt.Println()
	})

	// Start processor
	go processor.Start(ctx, lines)

	// Print statistics periodically
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)

	fmt.Println("📊 Real-time monitoring started. Press Ctrl+C to stop.")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Main loop
	for {
		select {
		case <-ticker.C:
			// Print statistics
			stats := processor.GetStatistics().GetSnapshot()
			fmt.Println()
			fmt.Println(strings.Repeat("=", 80))
			fmt.Println(stats.String())

			// Print alert counts
			alertCounts := alertMgr.GetAlertCounts()
			if len(alertCounts) > 0 {
				fmt.Println("\nAlerts:")
				for level, count := range alertCounts {
					fmt.Printf("  %s: %d\n", level.String(), count)
				}
			}
			fmt.Println(strings.Repeat("=", 80))
			fmt.Println()

		case <-sigChan:
			fmt.Println("\n\nStopping monitoring...")
			cancel()
			watcher.Stop()

			// Print final statistics
			stats := processor.GetStatistics().GetSnapshot()
			fmt.Println()
			fmt.Println("Final Statistics:")
			fmt.Println(stats.String())
			return nil
		}
	}
}

func parseLogFile(filename string, cfg *config.Config, verbose bool) error {
	if verbose {
		fmt.Printf("Parsing log file: %s\n", filename)
	}

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	// Parse the log
	logParser := logger.NewSQLServerLogParser()
	entries, err := logParser.ParseLog(file)
	if err != nil {
		return fmt.Errorf("failed to parse log: %v", err)
	}

	if verbose {
		fmt.Printf("Parsed %d log entries\n", len(entries))
	}

	// Apply filters
	filteredEntries := logger.FilterEntries(entries, logger.FilterCriteria{
		MinDuration: 0, // Will use config values in real implementation
		// Add more filter criteria based on config
	})

	if verbose {
		fmt.Printf("After filtering: %d entries\n", len(filteredEntries))
	}

	// Calculate metrics
	metrics := logger.CalculateMetrics(filteredEntries)

	// Output results
	return outputLogAnalysis(filteredEntries, metrics, cfg)
}

func outputAnalysis(analysis analyzer.QueryAnalysis, suggestions []analyzer.OptimizationSuggestion, cfg *config.Config) error {
	switch cfg.Output.Format {
	case "json":
		return outputJSON(map[string]interface{}{
			"analysis":    analysis,
			"suggestions": suggestions,
		}, cfg.Output.PrettyJSON)
	case "table":
		return outputTable(analysis, suggestions)
	case "csv":
		return outputCSV(analysis, suggestions)
	default:
		return fmt.Errorf("unsupported output format: %s", cfg.Output.Format)
	}
}

func outputLogAnalysis(entries []logger.LogEntry, metrics logger.LogMetrics, cfg *config.Config) error {
	switch cfg.Output.Format {
	case "json":
		return outputJSON(map[string]interface{}{
			"entries": entries,
			"metrics": metrics,
		}, cfg.Output.PrettyJSON)
	case "table":
		return outputLogTable(entries, metrics)
	default:
		return fmt.Errorf("unsupported output format: %s", cfg.Output.Format)
	}
}

func outputJSON(data interface{}, pretty bool) error {
	var output []byte
	var err error

	if pretty {
		output, err = json.MarshalIndent(data, "", "  ")
	} else {
		output, err = json.Marshal(data)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	fmt.Println(string(output))
	return nil
}

func outputTable(analysis analyzer.QueryAnalysis, suggestions []analyzer.OptimizationSuggestion) error {
	fmt.Println("=== SQL Query Analysis ===")
	fmt.Printf("Query Type: %s\n", analysis.QueryType)
	fmt.Printf("Complexity: %d\n\n", analysis.Complexity)

	// Tables
	if len(analysis.Tables) > 0 {
		fmt.Println("Tables:")
		fmt.Printf("%-20s %-10s %-10s %s\n", "Name", "Schema", "Alias", "Usage")
		fmt.Println(strings.Repeat("-", 60))
		for _, table := range analysis.Tables {
			fmt.Printf("%-20s %-10s %-10s %s\n", table.Name, table.Schema, table.Alias, table.Usage)
		}
		fmt.Println()
	}

	// Columns
	if len(analysis.Columns) > 0 {
		fmt.Println("Columns:")
		fmt.Printf("%-20s %-10s %s\n", "Name", "Table", "Usage")
		fmt.Println(strings.Repeat("-", 40))
		for _, column := range analysis.Columns {
			fmt.Printf("%-20s %-10s %s\n", column.Name, column.Table, column.Usage)
		}
		fmt.Println()
	}

	// Joins
	if len(analysis.Joins) > 0 {
		fmt.Println("Joins:")
		fmt.Printf("%-10s %-15s %-15s %s\n", "Type", "Left Table", "Right Table", "Condition")
		fmt.Println(strings.Repeat("-", 60))
		for _, join := range analysis.Joins {
			fmt.Printf("%-10s %-15s %-15s %s\n", join.Type, join.LeftTable, join.RightTable, join.Condition)
		}
		fmt.Println()
	}

	// Enhanced Optimization Suggestions (new)
	if len(analysis.EnhancedSuggestions) > 0 {
		fmt.Println("🚀 Advanced Optimization Suggestions:")
		fmt.Printf("%-5s %-12s %-8s %-25s %-8s %s\n", "#", "Category", "Impact", "Type", "Severity", "Description")
		fmt.Println(strings.Repeat("-", 100))

		for i, suggestion := range analysis.EnhancedSuggestions {
			icon := getSeverityIcon(suggestion.Severity)
			fmt.Printf("%-5d %-12s %-8s %-25s %s%-8s %s\n",
				i+1,
				suggestion.Category,
				suggestion.Impact,
				suggestion.Type,
				icon,
				suggestion.Severity,
				suggestion.Description)

			if suggestion.FixSuggestion != "" {
				fmt.Printf("      💡 Fix: %s\n", suggestion.FixSuggestion)
			}

			if suggestion.Dialect != "" {
				fmt.Printf("      🔧 Dialect: %s\n", suggestion.Dialect)
			}
		}
		fmt.Println()
	}

	// Legacy Optimization Suggestions (for backward compatibility)
	if len(suggestions) > 0 && len(analysis.EnhancedSuggestions) == 0 {
		fmt.Println("Optimization Suggestions:")
		for i, suggestion := range suggestions {
			fmt.Printf("%d. [%s] %s\n", i+1, suggestion.Severity, suggestion.Description)
		}
		fmt.Println()
	}

	return nil
}

// getSeverityIcon returns an appropriate icon for the severity level
func getSeverityIcon(severity string) string {
	switch severity {
	case "CRITICAL":
		return "🔥"
	case "ERROR":
		return "❌"
	case "WARNING":
		return "⚠️ "
	case "INFO":
		return "ℹ️ "
	default:
		return "   "
	}
}

func outputLogTable(entries []logger.LogEntry, metrics logger.LogMetrics) error {
	fmt.Println("=== SQL Server Log Analysis ===")
	fmt.Printf("Total Entries: %d\n", metrics.TotalEntries)
	fmt.Printf("Average Duration: %.2f ms\n", metrics.AvgDuration)
	fmt.Printf("Max Duration: %d ms\n", metrics.MaxDuration)
	fmt.Printf("Min Duration: %d ms\n\n", metrics.MinDuration)

	// Query Types Distribution
	if len(metrics.QueryTypes) > 0 {
		fmt.Println("Query Types:")
		for queryType, count := range metrics.QueryTypes {
			fmt.Printf("  %s: %d\n", queryType, count)
		}
		fmt.Println()
	}

	// Show first few entries
	if len(entries) > 0 {
		fmt.Println("Recent Entries:")
		fmt.Printf("%-20s %-8s %-8s %-15s %s\n", "Timestamp", "Duration", "Reads", "Database", "Query")
		fmt.Println(strings.Repeat("-", 80))

		limit := 10
		if len(entries) < limit {
			limit = len(entries)
		}

		for i := 0; i < limit; i++ {
			entry := entries[i]
			query := entry.Query
			if len(query) > 40 {
				query = query[:37] + "..."
			}
			fmt.Printf("%-20s %-8d %-8d %-15s %s\n",
				entry.Timestamp.Format("2006-01-02 15:04:05"),
				entry.Duration,
				entry.Reads,
				entry.Database,
				query)
		}
	}

	return nil
}

func outputCSV(analysis analyzer.QueryAnalysis, suggestions []analyzer.OptimizationSuggestion) error {
	fmt.Println("Query Type,Complexity,Table Count,Column Count,Join Count,Suggestions")
	fmt.Printf("%s,%d,%d,%d,%d,%d\n",
		analysis.QueryType,
		analysis.Complexity,
		len(analysis.Tables),
		len(analysis.Columns),
		len(analysis.Joins),
		len(suggestions))

	// Tables CSV
	if len(analysis.Tables) > 0 {
		fmt.Println("\nTables:")
		fmt.Println("Name,Schema,Alias,Usage")
		for _, table := range analysis.Tables {
			fmt.Printf("%s,%s,%s,%s\n", table.Name, table.Schema, table.Alias, table.Usage)
		}
	}

	// Columns CSV
	if len(analysis.Columns) > 0 {
		fmt.Println("\nColumns:")
		fmt.Println("Name,Table,Usage")
		for _, column := range analysis.Columns {
			fmt.Printf("%s,%s,%s\n", column.Name, column.Table, column.Usage)
		}
	}

	// Suggestions CSV
	if len(suggestions) > 0 {
		fmt.Println("\nSuggestions:")
		fmt.Println("Type,Description,Severity")
		for _, suggestion := range suggestions {
			fmt.Printf("%s,\"%s\",%s\n", suggestion.Type, suggestion.Description, suggestion.Severity)
		}
	}

	return nil
}
