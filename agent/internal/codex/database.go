package codex

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// openCodexDatabase selects the highest numeric database version, as codex-ss
// does. Read-only mode prevents accidentally creating or modifying a database.
func openCodexDatabase(ctx context.Context, prefix string) (*sql.DB, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	home := os.Getenv("CODEX_HOME")
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("get home directory: %w", err)
		}
		home = filepath.Join(userHome, ".codex")
	}

	entries, err := os.ReadDir(home)
	if err != nil {
		return nil, fmt.Errorf("read Codex directory %q: %w", home, err)
	}

	latestVersion := -1
	var latestPath string
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, prefix+"_") || !strings.HasSuffix(name, ".sqlite") {
			continue
		}
		version := strings.TrimSuffix(strings.TrimPrefix(name, prefix+"_"), ".sqlite")
		if version == "" || strings.IndexFunc(version, func(r rune) bool { return r < '0' || r > '9' }) >= 0 {
			continue
		}
		number, err := strconv.Atoi(version)
		if err != nil || number <= latestVersion {
			continue
		}
		latestVersion = number
		latestPath = filepath.Join(home, name)
	}
	if latestPath == "" {
		return nil, fmt.Errorf("find %s_N.sqlite in %q: %w", prefix, home, os.ErrNotExist)
	}

	path, err := filepath.Abs(latestPath)
	if err != nil {
		return nil, fmt.Errorf("resolve database path: %w", err)
	}
	uri := url.URL{Scheme: "file", Path: path}
	query := uri.Query()
	query.Set("mode", "ro")
	uri.RawQuery = query.Encode()

	db, err := sql.Open("sqlite", uri.String())
	if err != nil {
		return nil, fmt.Errorf("open Codex database %q: %w", path, err)
	}
	return db, nil
}
