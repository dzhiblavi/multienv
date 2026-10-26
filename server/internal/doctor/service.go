package doctor

import (
	"context"
	"os"
	"sort"
	"sync"
	"time"

	checks "multienv/common/doctor"
	"multienv/common/doctor/model"
)

type StorageChecker interface {
	Check(context.Context) error
}

type HostRunner interface {
	Diagnose(context.Context, string) (model.Report, error)
}

type Service struct {
	hosts    []string
	storage  StorageChecker
	runner   HostRunner
	hostname string
}

func New(hosts []string, storage StorageChecker, runner HostRunner) (*Service, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	unique := map[string]bool{hostname: true}

	for _, host := range hosts {
		unique[host] = true
	}

	ordered := make([]string, 0, len(unique))
	for host := range unique {
		ordered = append(ordered, host)
	}
	sort.Strings(ordered)

	return &Service{hosts: ordered, storage: storage, runner: runner, hostname: hostname}, nil
}

func (s *Service) Diagnose(ctx context.Context) model.Report {
	ctx, cancel := context.WithTimeout(ctx, 55*time.Second)
	defer cancel()
	report := model.Report{Hosts: []model.HostReport{s.checkStorage(ctx)}}
	hosts := make([][]model.HostReport, len(s.hosts))
	slots := make(chan struct{}, 8)
	var wg sync.WaitGroup
	for i, hostname := range s.hosts {
		wg.Go(func() {
			select {
			case slots <- struct{}{}:
				defer func() { <-slots }()
			case <-ctx.Done():
				hosts[i] = []model.HostReport{failedHost(hostname, s.source(hostname), ctx.Err())}
				return
			}

			hosts[i] = s.checkHost(ctx, hostname)
		})
	}

	wg.Wait()

	for _, host := range hosts {
		report.Hosts = append(report.Hosts, host...)
	}

	return report
}

func (s *Service) checkStorage(ctx context.Context) model.HostReport {
	check := checks.Measure(ctx, "redis", "storage", func(ctx context.Context) (string, error) {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		return "Redis JSON commands and scripting available", s.storage.Check(ctx)
	})

	return model.HostReport{Hostname: s.hostname, Source: "server", Checks: []model.Check{check}}
}

func (s *Service) checkHost(ctx context.Context, hostname string) []model.HostReport {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var report model.Report
	check := checks.Measure(ctx, "agent.local", "connection", func(ctx context.Context) (string, error) {
		var err error
		report, err = s.runner.Diagnose(ctx, hostname)
		return "agent local report received", err
	})

	if check.Status == model.Fail {
		return []model.HostReport{{Hostname: hostname, Source: s.source(hostname), Checks: []model.Check{check}}}
	}

	for i := range report.Hosts {
		report.Hosts[i].Source = s.source(hostname)
		report.Hosts[i].Checks = append([]model.Check{check}, report.Hosts[i].Checks...)
	}

	return report.Hosts
}

func failedHost(hostname, source string, err error) model.HostReport {
	return model.HostReport{Hostname: hostname, Source: source, Checks: []model.Check{{
		ID: "agent.local", Group: "connection", Status: model.Fail, Message: err.Error(),
	}}}
}

func (s *Service) source(hostname string) string {
	if hostname == s.hostname {
		return "local on master"
	}

	return "ssh from " + s.hostname
}
