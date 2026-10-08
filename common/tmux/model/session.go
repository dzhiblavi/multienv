package model

import "time"

type Session struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	WindowCount int       `json:"window_count"`
	SyncedAt    time.Time `json:"synced_at"`
}

type SessionsById map[string]Session

type SessionsByHost map[string]SessionsById
