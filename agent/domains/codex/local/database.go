package local

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
	path, err := codexDatabasePath(ctx, prefix)
	if err != nil {
		return nil, err
	}

	return openCodexDatabaseFile(path, false)
}

func codexHome() (string, error) {
	if home := os.Getenv("CODEX_HOME"); home != "" {
		return home, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	return filepath.Join(home, ".codex"), nil
}

func codexDatabasePath(ctx context.Context, prefix string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	home, err := codexHome()
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(home)
	if err != nil {
		return "", fmt.Errorf("read Codex directory %q: %w", home, err)
	}

	latestVersion := -1
	var latestPath string
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return "", err
		}

		version := codexDatabaseVersion(entry.Name(), prefix)
		if !entry.IsDir() && version > latestVersion {
			latestVersion = version
			latestPath = filepath.Join(home, entry.Name())
		}
	}

	if latestPath == "" {
		return "", fmt.Errorf("find %s_N.sqlite in %q: %w", prefix, home, os.ErrNotExist)
	}

	return filepath.Abs(latestPath)
}

func codexDatabaseVersion(name, prefix string) int {
	if !strings.HasPrefix(name, prefix+"_") || !strings.HasSuffix(name, ".sqlite") {
		return -1
	}

	version := strings.TrimSuffix(strings.TrimPrefix(name, prefix+"_"), ".sqlite")
	if version == "" || strings.IndexFunc(version, func(r rune) bool { return r < '0' || r > '9' }) >= 0 {
		return -1
	}

	number, err := strconv.Atoi(version)
	if err != nil {
		return -1
	}

	return number
}

func openCodexDatabaseFile(path string, writable bool) (*sql.DB, error) {
	uri := url.URL{Scheme: "file", Path: path}
	query := uri.Query()
	query.Set("mode", "ro")
	if writable {
		query.Set("mode", "rw")
		query.Set("_txlock", "immediate")
		query.Set("_pragma", "busy_timeout(3000)")
	}

	uri.RawQuery = query.Encode()
	db, err := sql.Open("sqlite", uri.String())
	if err != nil {
		return nil, fmt.Errorf("open Codex database %q: %w", path, err)
	}

	if writable {
		db.SetMaxOpenConns(1)
	}

	return db, nil
}
