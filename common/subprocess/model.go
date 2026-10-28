package subprocess

type Command struct {
	Name  string   `json:"command"`
	Stdin string   `json:"stdin,omitempty"`
	Args  []string `json:"args,omitempty"`
}

type RemoteCommand struct {
	Command
	Hostname string
	Compress bool
}

type RunResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}
