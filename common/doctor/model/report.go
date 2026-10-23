package model

type Status string

const (
	OK      Status = "ok"
	Warn    Status = "warn"
	Fail    Status = "fail"
	Skipped Status = "skipped"
)

type Check struct {
	ID         string `json:"id"`
	Group      string `json:"group"`
	Status     Status `json:"status"`
	Message    string `json:"message"`
	DurationMS int64  `json:"duration_ms"`
}

type HostReport struct {
	Hostname string  `json:"hostname"`
	Source   string  `json:"source"`
	Checks   []Check `json:"checks"`
}

type Report struct {
	Hosts []HostReport `json:"hosts"`
}

func (r Report) Failed() bool {
	for _, host := range r.Hosts {
		for _, check := range host.Checks {
			if check.Status == Fail {
				return true
			}
		}
	}

	return false
}
