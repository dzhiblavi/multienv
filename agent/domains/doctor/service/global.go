package service

import (
	"context"

	checks "multienv/common/doctor"
	"multienv/common/doctor/model"
	"multienv/common/doctor/operations"
)

func Global(ctx context.Context, configPath string) model.Report {
	report, client := diagnoseLocal(ctx, configPath)
	localChecks := report.Hosts[0].Checks
	if client == nil || localChecks[len(localChecks)-1].Status != model.OK {
		return report
	}

	var response *operations.GlobalResponse
	check := checks.Measure(ctx, "master.doctor", "connection", func(ctx context.Context) (string, error) {
		var err error
		response, err = client.Global(ctx, operations.GlobalRequest{})
		return "global report received", err
	})

	report.Hosts[0].Checks = append(report.Hosts[0].Checks, check)
	if check.Status == model.OK {
		report.Hosts = append(report.Hosts, response.Report.Hosts...)
	}

	return report
}
