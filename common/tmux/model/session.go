package model

type Session struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	WindowCount int    `json:"window_count"`
}

type SessionsById map[string]Session

type SessionsByHost map[string]SessionsById
