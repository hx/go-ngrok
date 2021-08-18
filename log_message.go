package ngrok

import "time"

// LogMessage represents a single log message emitted by a Process.
type LogMessage struct {
	Level   string    `json:"lvl"`
	Message string    `json:"msg"`
	Time    time.Time `json:"t"`
	Object  string    `json:"obj"`
	URL     *URL      `json:"url"`
	Name    string    `json:"name"`
	Address *URL      `json:"addr"`
	Error   string    `json:"err"`
	RawJSON string
}
