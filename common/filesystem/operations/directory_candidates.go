package operations

type DirectoryCandidatesRequest struct {
	Hostname string `json:"hostname"`
}

type DirectoryCandidatesResponse struct {
	Directories []string `json:"directories"`
}
