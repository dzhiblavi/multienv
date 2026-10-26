package model

import (
	"time"
)

type State string

const (
	Unknown        State = "Unknown"
	Archived       State = "Archived"
	NotLoaded      State = "Not Loaded"
	Idle           State = "Idle"
	Working        State = "Working"
	ActionRequired State = "ActionRequired"
)

type Session struct {
	Hostname     string    `json:"hostname"`
	Id           string    `json:"id"`
	WorkDir      string    `json:"work_dir"`
	TmuxPaneId   string    `json:"tmux_pane_id"`
	Title        string    `json:"title"`
	FirstMessage string    `json:"first_message"`
	LastMessage  string    `json:"last_message"`
	State        State     `json:"state"`
	SyncedAt     time.Time `json:"synced_at"`
}

type SessionsById map[string]Session

type SessionsByHost map[string]SessionsById
