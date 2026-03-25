package tests

import (
	"strings"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/internal/performance"
)

func TestGetMemoryStats(t *testing.T) {
	stats := performance.GetMemoryStats()
	if stats.SysBytes == 0 {
		t.Error("expected non-zero sys bytes")
	}
}

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input    uint64
		contains string
	}{
		{0, "B"},
		{512, "B"},
		{1023, "B"},
		{1024, "KB"},
		{1024 * 1024, "MB"},
		{1024 * 1024 * 1024, "GB"},
	}
	for _, tc := range cases {
		result := performance.FormatBytes(tc.input)
		if !strings.Contains(result, tc.contains) {
			t.Errorf("FormatBytes(%d) = %q, expected to contain %q", tc.input, result, tc.contains)
		}
	}
}

func TestNewPerformanceMonitor(t *testing.T) {
	pm := performance.NewPerformanceMonitor()
	if pm == nil {
		t.Fatal("expected non-nil monitor")
	}
}

func TestGetMetrics(t *testing.T) {
	pm := performance.NewPerformanceMonitor()
	metrics := pm.GetMetrics()

	keys := []string{"duration_ms", "memory_used_bytes", "memory_used_human", "memory_delta_bytes", "gc_count_delta", "gc_pause_total_ms"}
	for _, k := range keys {
		if _, ok := metrics[k]; !ok {
			t.Errorf("missing key %q in metrics", k)
		}
	}
	if metrics["duration_ms"].(int64) < 0 {
		t.Error("duration_ms should be non-negative")
	}
}

func TestForceGC(t *testing.T) {
	stats := performance.ForceGC()
	if stats.SysBytes == 0 {
		t.Error("expected non-zero sys bytes after GC")
	}
}
