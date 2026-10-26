package service

import (
	"context"
	"fmt"
	"os"

	"multienv/agent/config"
	"multienv/agent/domains/doctor/client"
	transport "multienv/agent/domains/doctor/client/http"
	checks "multienv/common/doctor"
	"multienv/common/doctor/model"
	"multienv/common/doctor/operations"
	"multienv/common/util"
)

func Local(ctx context.Context, configPath string) model.Report {
	report, _ := diagnoseLocal(ctx, configPath)
	return report
}

func diagnoseLocal(ctx context.Context, configPath string) (model.Report, client.Client) {
	hostname, err := os.Hostname()
	host := model.HostReport{Hostname: hostname, Source: "local", Checks: []model.Check{}}
	if err != nil {
		host.Checks = append(host.Checks, model.Check{
			ID: "hostname", Group: "agent", Status: model.Fail, Message: err.Error(),
		})
	}

	cfg, configuration := checkConfiguration(ctx, configPath)
	host.Checks = append(host.Checks, configuration)
	host.Checks = append(host.Checks, checkDependencies(ctx)...)
	if configuration.Status != model.OK {
		host.Checks = append(host.Checks, model.Check{
			ID: "master.ping", Group: "connection", Status: model.Skipped,
			Message: "master address unavailable because agent configuration is invalid",
		})

		return model.Report{Hosts: []model.HostReport{host}}, nil
	}

	master := transport.New(cfg.Client)
	host.Checks = append(host.Checks, checkMaster(ctx, master, cfg.Client.URL))

	return model.Report{Hosts: []model.HostReport{host}}, master
}

func checkConfiguration(ctx context.Context, path string) (config.Config, model.Check) {
	var cfg config.Config
	check := checks.Measure(ctx, "config", "agent", func(ctx context.Context) (string, error) {
		if err := util.ReadYaml(ctx, path, &cfg); err != nil {
			return "", err
		}

		if cfg.Client.Timeout <= 0 {
			return "", fmt.Errorf("client.timeout must be positive")
		}

		return path, nil
	})

	return cfg, check
}

func checkMaster(ctx context.Context, master client.Client, address string) model.Check {
	return checks.Measure(ctx, "master.ping", "connection", func(ctx context.Context) (string, error) {
		response, err := master.Ping(ctx, operations.PingRequest{})
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s (%s)", address, response.Hostname), nil
	})
}
