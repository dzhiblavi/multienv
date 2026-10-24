package filesystem

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"multienv/common/subprocess"
)

func DirectoryCandidates(ctx context.Context) ([]string, error) {
	program, err := exec.LookPath("zoxide")
	if errors.Is(err, exec.ErrNotFound) {
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find zoxide: %w", err)
	}

	result, err := subprocess.RunSuccess(ctx, subprocess.Command{
		Name: program,
		Args: []string{"query", "--list"},
	})
	if err != nil {
		return nil, fmt.Errorf("query directory candidates: %w", err)
	}

	directories := []string{}
	for directory := range strings.SplitSeq(strings.TrimSuffix(result.Stdout, "\n"), "\n") {
		if directory != "" {
			directories = append(directories, directory)
		}
	}

	return directories, nil
}
