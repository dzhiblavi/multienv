package backend

import (
	"context"
	"errors"
	"os"
	"slices"

	agentclient "multienv/agent/client/agent"
	agenthttp "multienv/agent/client/agent/http"
	codexhttp "multienv/agent/codex/client/http"
	codexservice "multienv/agent/codex/service"
	agentconfig "multienv/agent/config"
	doctorservice "multienv/agent/doctor/service"
	filesystemservice "multienv/agent/filesystem/service"
	tmuxhttp "multienv/agent/tmux/client/http"
	tmuxservice "multienv/agent/tmux/service"
	agentops "multienv/common/agent/operations"
	doctormodel "multienv/common/doctor/model"
	"multienv/common/util"
	"multienv/tui/internal/config"
)

type Services struct {
	AgentClient agentclient.Client
	Diagnose    func(context.Context) doctormodel.Report
	Tmux        *tmuxservice.Agent
	Codex       *codexservice.Agent
	Filesystem  *filesystemservice.Agent
}

func New(ctx context.Context) (Services, error) {
	var cfg agentconfig.Config
	if err := util.ReadYaml(ctx, config.AgentPath, &cfg); err != nil {
		return Services{}, err
	}

	forwarding := agenthttp.New(cfg.Client)
	tmux, err := tmuxservice.New(tmuxhttp.New(cfg.Client), forwarding)
	if err != nil {
		return Services{}, err
	}

	codex, err := codexservice.New(codexhttp.New(cfg.Client), forwarding)
	if err != nil {
		return Services{}, err
	}

	filesystem, err := filesystemservice.New(forwarding)
	if err != nil {
		return Services{}, err
	}

	return Services{AgentClient: forwarding, Tmux: tmux, Codex: codex, Filesystem: filesystem,
		Diagnose: func(ctx context.Context) doctormodel.Report {
			return doctorservice.Global(ctx, config.AgentPath)
		},
	}, nil
}

func (s Services) Refresh(ctx context.Context, domain string, hosts []string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	hosts = slices.Clone(hosts)
	slices.Sort(hosts)
	var failures []error
	for _, host := range slices.Compact(hosts) {
		if host == hostname {
			if domain == "tmux" {
				err = s.Tmux.SyncLocal(ctx)
			} else {
				err = s.Codex.SyncLocal(ctx)
			}
		} else {
			_, err = s.AgentClient.Forward(ctx, agentops.ForwardRequest{
				Hostname: host, Domain: domain, Command: "sync-local",
			})
		}

		failures = append(failures, err)
	}

	return errors.Join(failures...)
}
