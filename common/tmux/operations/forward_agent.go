package operations

type ForwardAgentRequest struct {
	Hostname string `json:"hostname"`
	Command  string `json:"command"`
	Input    any    `json:"input"`
}

type ForwardAgentResponse struct {
	Output string `json:"output"`
}
