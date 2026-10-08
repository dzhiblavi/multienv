package model

type Window struct {
	Id             string `json:"id"`
	Index          int    `json:"index"`
	Name           string `json:"name"`
	Active         bool   `json:"active"`
	CurrentCommand string `json:"current_command"`
	CurrentPath    string `json:"current_path"`
}
