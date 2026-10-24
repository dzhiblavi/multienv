package operations

type ForwardRequest struct {
	Hostname string `json:"hostname"`
	Domain   string `json:"domain"`
	Command  string `json:"command"`
	Input    any    `json:"input"`
}

type ForwardResponse struct {
	Output string `json:"output"`
}
