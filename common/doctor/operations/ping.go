package operations

type PingRequest struct{}

type PingResponse struct {
	Service  string `json:"service"`
	Hostname string `json:"hostname"`
}
