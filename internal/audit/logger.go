package audit

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Logger struct {
	path string
	mu   sync.Mutex
}

type Event struct {
	Timestamp time.Time         `json:"timestamp"`
	Action    string            `json:"action"`
	Actor     string            `json:"actor"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

func New(path string) *Logger {
	return &Logger{path: path}
}

func (l *Logger) Append(event Event) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	file, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = file.Write(append(payload, '\n'))
	return err
}
