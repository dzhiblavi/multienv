package subprocess

type Command struct {
	Name  string
	Stdin string
	Args  []string
}

type RemoteCommand struct {
	Command
	Hostname string
}

type RunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}
