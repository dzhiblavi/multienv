package local

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

type sessionNameEntry struct {
	ID         string `json:"id"`
	ThreadName string `json:"thread_name"`
	UpdatedAt  string `json:"updated_at"`
}

func appendSessionName(home, id, name string) error {
	entry := sessionNameEntry{ID: id, ThreadName: name, UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(home, "session_index.jsonl"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	written, err := file.Write(data)
	if err == nil && written != len(data) {
		err = io.ErrShortWrite
	}

	return errors.Join(err, file.Close())
}
