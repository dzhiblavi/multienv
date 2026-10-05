package operations

type DeleteSessionRequest struct {
	Hostname  string `json:"hostname"`
	SessionId string `json:"session_id"`
}

type DeleteSessionResponse struct{}
