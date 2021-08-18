package ngrok

import (
	"encoding/json"
	"net/url"
	"regexp"
)

var schemaPattern = regexp.MustCompile(`^\w+://`)

// URL is an extension of url.URL that implements json.Unmarshaler.
type URL struct{ *url.URL }

func (u URL) MarshalJSON() ([]byte, error) { return json.Marshal(u.String()) }

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
		u.URL = parsed
	}
	return err
}
