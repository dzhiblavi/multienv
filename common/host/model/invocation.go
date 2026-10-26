package model

import (
	"fmt"

	"multienv/common/subprocess"
)

const Master = "@master"
const All = "@all"

type Invocation struct {
	subprocess.RunResult

	Hostname string `json:"hostname"`
	Error    string `json:"error,omitempty"`
}

func (i Invocation) Err() error {
	if i.Error != "" {
		return fmt.Errorf("%s: %s", i.Hostname, i.Error)
	}

	if i.ExitCode != 0 {
		return fmt.Errorf("%s: command exited with code %d: %s", i.Hostname, i.ExitCode, i.Stderr)
	}

	return nil
}
