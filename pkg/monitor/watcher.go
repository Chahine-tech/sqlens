package monitor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// LogWatcher watches a log file for new entries and streams them
type LogWatcher struct {
	filePath     string
	pollInterval time.Duration
	file         *os.File
	reader       *bufio.Reader
	position     int64
}

// NewLogWatcher creates a new log file watcher
func NewLogWatcher(filePath string) *LogWatcher {
	return &LogWatcher{
		filePath:     filePath,
		pollInterval: 500 * time.Millisecond, // Default poll every 500ms
		position:     0,
	}
}

// SetPollInterval sets how often to check for new lines
func (w *LogWatcher) SetPollInterval(interval time.Duration) {
	w.pollInterval = interval
}

// Start begins watching the log file and sends new lines to the channel
func (w *LogWatcher) Start(ctx context.Context, lines chan<- string) error {
	// Open the file
	file, err := os.Open(w.filePath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	w.file = file

	// Seek to end of file to only read new lines
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat log file: %w", err)
	}
	w.position = info.Size()
	_, err = file.Seek(w.position, io.SeekStart)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to seek log file: %w", err)
	}

	w.reader = bufio.NewReader(file)

	// Start watching in a goroutine
	go w.watch(ctx, lines)

	return nil
}

// watch is the main loop that polls for new lines
func (w *LogWatcher) watch(ctx context.Context, lines chan<- string) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()
	defer w.file.Close()
	defer close(lines)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Check if file has new content
			if err := w.checkForNewLines(lines); err != nil {
				// Log error but continue watching
				fmt.Fprintf(os.Stderr, "Error reading log file: %v\n", err)
			}
		}
	}
}

// checkForNewLines reads any new lines added to the file
func (w *LogWatcher) checkForNewLines(lines chan<- string) error {
	// Check current file size
	info, err := w.file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	currentSize := info.Size()

	// If file was truncated (rotated), reopen from beginning
	if currentSize < w.position {
		w.file.Close()
		file, err := os.Open(w.filePath)
		if err != nil {
			return fmt.Errorf("failed to reopen log file: %w", err)
		}
		w.file = file
		w.reader = bufio.NewReader(file)
		w.position = 0
	}

	// Read new lines
	for {
		line, err := w.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// No more lines, update position and return
				w.position, _ = w.file.Seek(0, io.SeekCurrent)
				return nil
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		// Send line to channel (remove trailing newline)
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}

		select {
		case lines <- line:
			// Line sent successfully
		default:
			// Channel full, skip this line (shouldn't happen with unbuffered channel)
		}
	}
}

// Stop stops watching the log file
func (w *LogWatcher) Stop() {
	if w.file != nil {
		w.file.Close()
	}
}

// TailMode starts watching from the last N lines instead of end of file
func (w *LogWatcher) StartWithTail(ctx context.Context, lines chan<- string, tailLines int) error {
	// Open the file
	file, err := os.Open(w.filePath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	w.file = file

	// Read last N lines first
	lastLines, err := w.readLastNLines(file, tailLines)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to read tail: %w", err)
	}

	// Send tail lines to channel
	for _, line := range lastLines {
		select {
		case lines <- line:
		case <-ctx.Done():
			file.Close()
			return ctx.Err()
		}
	}

	// Get current position
	w.position, err = file.Seek(0, io.SeekCurrent)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to get position: %w", err)
	}

	w.reader = bufio.NewReader(file)

	// Start watching in a goroutine
	go w.watch(ctx, lines)

	return nil
}

// readLastNLines reads the last N lines from the file
func (w *LogWatcher) readLastNLines(file *os.File, n int) ([]string, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := info.Size()
	if fileSize == 0 {
		return []string{}, nil
	}

	// Simple implementation: read entire file and get last N lines
	// For very large files, this could be optimized
	content := make([]byte, fileSize)
	_, err = file.ReadAt(content, 0)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// Split into lines
	allLines := []string{}
	start := 0
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			line := string(content[start:i])
			if len(line) > 0 {
				allLines = append(allLines, line)
			}
			start = i + 1
		}
	}

	// Add last line if it doesn't end with newline
	if start < len(content) {
		line := string(content[start:])
		if len(line) > 0 {
			allLines = append(allLines, line)
		}
	}

	// Return last N lines
	if len(allLines) <= n {
		return allLines, nil
	}
	return allLines[len(allLines)-n:], nil
}
