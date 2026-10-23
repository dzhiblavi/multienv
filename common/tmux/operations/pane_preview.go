package operations

type PanePreviewRequest struct {
	Hostname string `json:"hostname"`
	Target   string `json:"target"`
}

type PanePreviewResponse struct {
	Content string `json:"content"`
}
