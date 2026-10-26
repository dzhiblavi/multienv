package host

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"multienv/common/host/model"
)

func New(hostnames []string) (Dispatcher, error) {
	master, err := os.Hostname()
	if err != nil {
		return Dispatcher{}, err
	}

	if slices.Contains(hostnames, model.All) {
		return Dispatcher{}, fmt.Errorf("@all cannot be used in the configured host list")
	}

	d := Dispatcher{master: master}
	d.hosts, err = d.resolveHostnames(hostnames)
	return d, err
}

func (d Dispatcher) Hostnames() []string {
	return slices.Clone(d.hosts)
}

func (d Dispatcher) resolveHostnames(requested []string) ([]string, error) {
	var resolved []string
	seen := map[string]bool{}
	for _, hostname := range requested {
		hosts := []string{hostname}
		if hostname == model.All {
			hosts = d.hosts
		}

		for _, host := range hosts {
			if host == model.Master {
				host = d.master
			} else if strings.HasPrefix(host, "@") {
				return nil, fmt.Errorf("unknown host alias %q", host)
			}

			if strings.TrimSpace(host) == "" || strings.ContainsAny(host, "\x00\t\r\n") {
				return nil, fmt.Errorf("invalid hostname %q", host)
			}

			if !seen[host] {
				seen[host] = true
				resolved = append(resolved, host)
			}
		}
	}

	return resolved, nil
}
