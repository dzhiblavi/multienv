package doctor

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	checks "multienv/common/doctor"
	"multienv/common/doctor/model"
)

type dependency struct {
	name     string
	group    string
	args     []string
	features []string
	optional bool
}

func dependencies() []dependency {
	return []dependency{
		{name: "tmux", group: "tmux", args: []string{"-V"}},
		{name: "ssh", group: "connection", args: []string{"-V"}},
		{name: "codex", group: "codex", optional: true, args: []string{"--no-daemon", "--help"},
			features: []string{"resume", "archive", "unarchive", "delete"}},
		{name: "lsof", group: "codex", optional: true, args: []string{"-v"}},
		{name: "zoxide", group: "filesystem", optional: true, args: []string{"query", "--help"},
			features: []string{"--list"}},
	}
}

func checkDependencies(ctx context.Context) []model.Check {
	items := dependencies()
	if runtime.GOOS == "darwin" {
		items = append(items, dependency{name: "osascript", group: "desktop", optional: true},
			dependency{name: "/usr/bin/open", group: "desktop", optional: true})
	}

	results := make([]model.Check, len(items))
	var wg sync.WaitGroup
	for i, item := range items {
		wg.Go(func() { results[i] = checkDependency(ctx, item) })
	}

	wg.Wait()

	return results
}

func checkDependency(ctx context.Context, item dependency) model.Check {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	check := checks.Measure(ctx, "binary."+item.name, item.group, func(ctx context.Context) (string, error) {
		return inspectDependency(ctx, item)
	})

	if check.Status == model.Fail && item.optional {
		check.Status = model.Warn
	}

	return check
}

func inspectDependency(ctx context.Context, item dependency) (string, error) {
	path, err := exec.LookPath(item.name)
	if err != nil {
		return "", err
	}

	if len(item.args) == 0 {
		return path, nil
	}

	output, err := exec.CommandContext(ctx, path, item.args...).CombinedOutput()
	if ctx.Err() != nil {
		return "", fmt.Errorf("%s: %w", path, ctx.Err())
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w: %s", path, err, strings.TrimSpace(string(output)))
	}

	for _, feature := range item.features {
		if !strings.Contains(string(output), feature) {
			return "", fmt.Errorf("%s: required capability %s not found", path, feature)
		}
	}

	detail := strings.TrimSpace(string(output))
	if len(item.features) > 0 {
		detail = "capabilities: " + strings.Join(item.features, ", ")
	} else if len(detail) > 160 {
		detail = strings.SplitN(detail, "\n", 2)[0]
	}

	return path + " — " + detail, nil
}
