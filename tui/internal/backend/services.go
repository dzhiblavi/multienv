package backend

import (
	"context"
	"errors"
	"os"
	"slices"

	agentclient "multienv/agent/client/agent"
	agentconfig "multienv/agent/config"
	codexservice "multienv/agent/domains/codex/service"
	doctorservice "multienv/agent/domains/doctor/service"
	filesystemservice "multienv/agent/domains/filesystem/service"
	hostclient "multienv/agent/domains/host/client"
	hosthttp "multienv/agent/domains/host/client/http"
	storagehttp "multienv/agent/domains/storage/client/http"
	tmuxservice "multienv/agent/domains/tmux/service"
	doctormodel "multienv/common/doctor/model"
	"multienv/common/util"
	"multienv/tui/internal/config"
)

type Services struct {
	HostClient hostclient.Client
	Diagnose   func(context.Context) doctormodel.Report
	Tmux       *tmuxservice.Agent
	Codex      *codexservice.Agent
	Filesystem *filesystemservice.Agent
}

func New(ctx context.Context) (Services, error) {
	var cfg agentconfig.Config
	if err := util.ReadYaml(ctx, config.AgentPath, &cfg); err != nil {
		return Services{}, err
	}

	hosts := hosthttp.New(cfg.Client)
	data := storagehttp.New(cfg.Client)
	tmux, err := tmuxservice.New(data, hosts)
	if err != nil {
		return Services{}, err
	}

	codex, err := codexservice.New(data, hosts)
	if err != nil {
		return Services{}, err
	}

	filesystem, err := filesystemservice.New(hosts)
	if err != nil {
		return Services{}, err
	}

	return Services{HostClient: hosts, Tmux: tmux, Codex: codex, Filesystem: filesystem,
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
			_, err = agentclient.Invoke(ctx, s.HostClient, host, domain, "sync-local", nil)
		}

		failures = append(failures, err)
	}

	return errors.Join(failures...)
}
