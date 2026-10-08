package model

import "time"

type Session struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	WindowCount    int       `json:"window_count"`
	Windows        []Window  `json:"windows"`
	Attached       int       `json:"attached"`
	LastAttached   int64     `json:"last_attached"`
	CurrentCommand string    `json:"current_command"`
	CurrentPath    string    `json:"current_path"`
	SyncedAt       time.Time `json:"synced_at"`
}

type SessionsById map[string]Session

type SessionsByHost map[string]SessionsById
