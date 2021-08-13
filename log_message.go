package ngrok

import (
	"encoding/json"
	"net/url"
	"regexp"
	"time"
)

// URL is an extension of url.URL that implements json.Unmarshaler.
type URL url.URL

var schemaPattern = regexp.MustCompile(`^\w+://`)

func (u *URL) UnmarshalJSON(bytes []byte) error {
	str := ""
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	if !schemaPattern.MatchString(str) {
		str = "tcp://" + str
	}
	parsed, err := url.Parse(str)
	if err == nil {
		*u = URL(*parsed)
	}
	return err
}

// Raw returns a pointer to a raw url.URL.
func (u *URL) Raw() *url.URL {
	if u == nil {
		return nil
	}
	raw := url.URL(*u)
	return &raw
}

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
