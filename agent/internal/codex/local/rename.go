package local

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

func RenameSession(ctx context.Context, id, newName string) error {
	newName = strings.TrimSpace(newName)
	if id == "" || newName == "" || strings.IndexFunc(newName, unicode.IsControl) >= 0 {
		return fmt.Errorf("session_id and a non-empty, single-line new_name are required")
	}

	path, err := codexDatabasePath(ctx, "state")
	if err != nil {
		return err
	}

	db, err := openCodexDatabaseFile(path, true)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := renameSessionInDatabase(ctx, db, id, newName); err != nil {
		return fmt.Errorf("rename Codex session: %w", err)
	}

	if err := appendSessionName(filepath.Dir(path), id, newName); err != nil {
		return fmt.Errorf("session renamed, but updating the Codex name index failed: %w", err)
	}

	return nil
}

func renameSessionInDatabase(ctx context.Context, db *sql.DB, id, newName string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var mode string
	if err := tx.QueryRowContext(ctx, "SELECT history_mode FROM threads WHERE id = ?", id).Scan(&mode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("session %s not found", id)
		}

		return err
	}

	var query string
	switch mode {
	case "paginated":
		query = "UPDATE threads SET name = ? WHERE id = ?"
	case "legacy":
		query = "UPDATE threads SET title = ? WHERE id = ?"
	default:
		return fmt.Errorf("unsupported Codex history mode %q", mode)
	}

	if _, err := tx.ExecContext(ctx, query, newName, id); err != nil {
		return err
	}

	return tx.Commit()
}
