package host

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"multienv/common/host/model"
	"multienv/common/host/operations"
	"multienv/common/log"
	"multienv/common/subprocess"
	"multienv/server/internal/dispatcher"
)

const maxParallelInvocations = 8
const defaultInvokeTimeout = 30 * time.Second

func (d Dispatcher) Invoke(ctx context.Context, req operations.InvokeRequest) (*operations.InvokeResponse, error) {
	if err := validateInvokeRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %w", dispatcher.ErrInvalidParams, err)
	}

	timeout, err := parseInvokeTimeout(req.Timeout)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dispatcher.ErrInvalidParams, err)
	}

	hostnames, err := d.resolveHostnames(req.Hostnames)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dispatcher.ErrInvalidParams, err)
	}

	response := &operations.InvokeResponse{Results: make([]model.Invocation, len(hostnames))}
	slots := make(chan struct{}, maxParallelInvocations)
	var wg sync.WaitGroup
	for i, hostname := range hostnames {
		wg.Go(func() {
			select {
			case slots <- struct{}{}:
				defer func() { <-slots }()
			case <-ctx.Done():
				response.Results[i] = model.Invocation{
					Hostname: hostname, Error: ctx.Err().Error(), RunResult: subprocess.RunResult{ExitCode: -1},
				}
				return
			}

			response.Results[i] = invokeHost(ctx, req, hostname, timeout)
		})
	}

	wg.Wait()
	return response, nil
}

func validateInvokeRequest(req operations.InvokeRequest) error {
	if len(req.Hostnames) == 0 || strings.TrimSpace(req.Name) == "" || strings.ContainsRune(req.Name, 0) {
		return fmt.Errorf("hostnames and a non-empty command without NUL bytes are required")
	}

	for _, hostname := range req.Hostnames {
		if strings.TrimSpace(hostname) == "" || strings.ContainsAny(hostname, "\x00\t\r\n") {
			return fmt.Errorf("invalid hostname %q", hostname)
		}
	}

	for _, arg := range req.Args {
		if strings.ContainsRune(arg, 0) {
			return fmt.Errorf("command arguments must not contain NUL bytes")
		}
	}

	return nil
}

func parseInvokeTimeout(value string) (time.Duration, error) {
	if value == "" {
		return defaultInvokeTimeout, nil
	}

	timeout, err := time.ParseDuration(value)
	if err != nil || timeout <= 0 {
		return 0, fmt.Errorf("timeout must be a positive duration")
	}

	return timeout, nil
}

func invokeHost(
	ctx context.Context,
	req operations.InvokeRequest,
	hostname string,
	timeout time.Duration,
) model.Invocation {
	ctx = log.WithCustomLogData(ctx, "host", hostname)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result, err := subprocess.RunRemote(ctx, subprocess.RemoteCommand{
		Hostname: hostname,
		Command:  req.Command,
		Compress: req.Compress,
	})
	response := model.Invocation{Hostname: hostname, RunResult: result}
	if ctx.Err() != nil {
		err = ctx.Err()
	}

	if err != nil {
		response.Error = err.Error()
		response.ExitCode = -1
	}

	return response
}
