package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.Parser.Dialect != "sqlserver" {
		t.Errorf("expected default dialect sqlserver, got %s", cfg.Parser.Dialect)
	}
	if cfg.Parser.MaxQuerySize <= 0 {
		t.Error("expected positive max_query_size")
	}
	if !cfg.Analyzer.EnableOptimizations {
		t.Error("expected optimizations enabled by default")
	}
	if cfg.Output.Format != "json" {
		t.Errorf("expected default output format json, got %s", cfg.Output.Format)
	}
}

func TestLoadConfig_EmptyFilename(t *testing.T) {
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoadConfig_JSON(t *testing.T) {
	data := `{
		"parser": {"dialect": "mysql", "max_query_size": 500000, "strict_mode": true},
		"analyzer": {"enable_optimizations": false, "complexity_threshold": 5},
		"logger": {"default_format": "profiler", "max_file_size_mb": 50},
		"output": {"format": "table", "pretty_json": false}
	}`
	tmp := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(tmp, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.LoadConfig(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Parser.Dialect != "mysql" {
		t.Errorf("expected mysql dialect, got %s", cfg.Parser.Dialect)
	}
	if cfg.Parser.MaxQuerySize != 500000 {
		t.Errorf("expected 500000, got %d", cfg.Parser.MaxQuerySize)
	}
	if cfg.Analyzer.EnableOptimizations {
		t.Error("expected optimizations disabled")
	}
}

func TestLoadConfig_YAML(t *testing.T) {
	data := `
parser:
  dialect: postgresql
  max_query_size: 200000
  strict_mode: false
analyzer:
  enable_optimizations: true
  complexity_threshold: 8
logger:
  default_format: profiler
  max_file_size_mb: 20
output:
  format: json
  pretty_json: true
`
	tmp := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(tmp, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.LoadConfig(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Parser.Dialect != "postgresql" {
		t.Errorf("expected postgresql, got %s", cfg.Parser.Dialect)
	}
}

func TestLoadConfig_YML(t *testing.T) {
	data := `
parser:
  dialect: sqlite
  max_query_size: 100000
logger:
  max_file_size_mb: 10
output:
  format: csv
`
	tmp := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(tmp, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.LoadConfig(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Parser.Dialect != "sqlite" {
		t.Errorf("expected sqlite, got %s", cfg.Parser.Dialect)
	}
}

func TestLoadConfig_UnsupportedExtension(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(tmp, []byte("key=value"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := config.LoadConfig(tmp)
	if err == nil {
		t.Error("expected error for unsupported extension")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := config.LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(tmp, []byte("{invalid json"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := config.LoadConfig(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSaveConfig_JSON(t *testing.T) {
	cfg := config.DefaultConfig()
	tmp := filepath.Join(t.TempDir(), "out.json")
	if err := config.SaveConfig(cfg, tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := config.LoadConfig(tmp)
	if err != nil {
		t.Fatalf("unexpected error loading saved config: %v", err)
	}
	if loaded.Parser.Dialect != cfg.Parser.Dialect {
		t.Error("saved and loaded dialect mismatch")
	}
}

func TestSaveConfig_JSON_NotPretty(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.PrettyJSON = false
	tmp := filepath.Join(t.TempDir(), "out.json")
	if err := config.SaveConfig(cfg, tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSaveConfig_YAML(t *testing.T) {
	cfg := config.DefaultConfig()
	tmp := filepath.Join(t.TempDir(), "out.yaml")
	if err := config.SaveConfig(cfg, tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := config.LoadConfig(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Parser.Dialect != cfg.Parser.Dialect {
		t.Error("dialect mismatch after yaml round-trip")
	}
}

func TestSaveConfig_UnsupportedExtension(t *testing.T) {
	cfg := config.DefaultConfig()
	err := config.SaveConfig(cfg, "/tmp/config.toml")
	if err == nil {
		t.Error("expected error for unsupported extension")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := config.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid default config, got: %v", err)
	}
}

func TestValidate_InvalidMaxQuerySize(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Parser.MaxQuerySize = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero max_query_size")
	}
}

func TestValidate_NegativeComplexityThreshold(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Analyzer.ComplexityThreshold = -1
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative complexity_threshold")
	}
}

func TestValidate_InvalidMaxFileSizeMB(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Logger.MaxFileSizeMB = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero max_file_size_mb")
	}
}

func TestValidate_InvalidOutputFormat(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Output.Format = "xml"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid output format")
	}
}

func TestValidate_InvalidDialect(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Parser.Dialect = "snowflake"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid dialect")
	}
}

func TestValidate_AllValidFormats(t *testing.T) {
	for _, format := range []string{"json", "table", "csv"} {
		cfg := config.DefaultConfig()
		cfg.Output.Format = format
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected format %s to be valid: %v", format, err)
		}
	}
}

func TestValidate_AllValidDialects(t *testing.T) {
	for _, dialect := range []string{"sqlserver", "mysql", "postgresql", "oracle", "sqlite"} {
		cfg := config.DefaultConfig()
		cfg.Parser.Dialect = dialect
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected dialect %s to be valid: %v", dialect, err)
		}
	}
}

func TestGetConfigPath(t *testing.T) {
	p := config.GetConfigPath()
	if p == "" {
		t.Error("expected non-empty config path")
	}
}

func TestEnsureConfigDir(t *testing.T) {
	// Just ensure it doesn't error — the dir may already exist
	if err := config.EnsureConfigDir(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
