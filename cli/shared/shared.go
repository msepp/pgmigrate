package shared

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/log"

	"github.com/peterldowns/pgmigrate/internal/schema"
)

type Flags struct {
	LogFormat  *string // see logger.go
	Database   *string // see root.go
	Migrations *string // see root.go
	Configfile *string // see root.go
}
type Config struct {
	Database   string        `yaml:"database"`
	LogFormat  LogFormat     `yaml:"log_format"`
	Migrations string        `yaml:"migrations"`
	Schema     schema.Config `yaml:"schema"`
}

type StateT struct {
	Flags  Flags
	Config Config
}

var State StateT

func (state *StateT) Parse() {
	cf := state.Configfile()
	if !cf.IsSet() {
		return
	}
	file, err := os.OpenFile(cf.Value(), os.O_RDWR, 0o644)
	if err != nil {
		panic(fmt.Errorf("open config: %w", err))
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("read config: %w", err))
	}
	if err := yaml.Unmarshal(contents, &state.Config); err != nil {
		panic(fmt.Errorf("parse config: %w", err))
	}
}

func (state StateT) Configfile() Variable[string] {
	return NewVariable(
		"config-file",
		*state.Flags.Configfile,
		os.Getenv("PGM_CONFIGFILE"),
		checkpath(".pgmigrate.yaml"), // in cwd
		repopath(".pgmigrate.yaml"),  // in repo root
		"",                           // default to missing
	)
}

func (state StateT) Database() Variable[string] {
	return NewVariable(
		"database",
		*state.Flags.Database,
		os.Getenv("PGM_DATABASE"),
		state.Config.Database,
		"", // default to missing
	)
}

func (state StateT) LogFormat() Variable[LogFormat] {
	return NewVariable(
		"log-format",
		LogFormat(*state.Flags.LogFormat),
		LogFormat(os.Getenv("PGM_LOG_FORMAT")),
		state.Config.LogFormat,
		LogFormatText, // default to LogFormatText
	)
}

func (state StateT) Migrations() Variable[string] {
	return NewVariable(
		"migrations",
		*state.Flags.Migrations,
		os.Getenv("PGM_MIGRATIONS"),
		state.Config.Migrations,
		"", // default to missing
	)
}

func (state StateT) Logger() (*log.Logger, LogAdapter) {
	var logger *log.Logger
	format := state.LogFormat().Value()
	switch format {
	case LogFormatText:
		logger = log.NewWithOptions(os.Stdout, log.Options{Formatter: log.TextFormatter})
	case LogFormatJSON:
		logger = log.NewWithOptions(os.Stdout, log.Options{Formatter: log.JSONFormatter})
	default:
		panic(fmt.Errorf("unknown log format: %s", format))
	}
	return logger, LogAdapter{logger}
}

func repopath(p string) string {
	root, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}
	rootConfig := path.Join(strings.TrimSpace(string(root)), p)
	return checkpath(rootConfig)
}

func checkpath(p string) string {
	p, err := filepath.Abs(p)
	if err != nil {
		return ""
	}
	if _, err := os.Stat(p); err != nil {
		return ""
	}
	return p
}