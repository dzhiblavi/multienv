package util

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heetch/confita/backend/file"
)

type WithDefault interface {
	Default()
}

func ReadYaml[T WithDefault](ctx context.Context, path string, config T) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	path = os.ExpandEnv(path)
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home directory: %w", err)
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	config.Default()

	// Load the file directly: confita.Loader can skip its file backend when
	// the target has no config tags, even if it has YAML tags.
	if err := file.NewBackend(path).Unmarshal(ctx, config); err != nil {
		return fmt.Errorf("load config %q: %w", path, err)
	}

	return nil
}
