package ngrok

import (
	"bytes"
	"encoding/json"
)

// messageStream is a delegate for Process that implements io.Writer to receive data from the process's STDOUT.
type messageStream struct {
	buffer     bytes.Buffer
	dispatcher *Dispatcher
}

func (m *messageStream) Write(p []byte) (n int, err error) {
	parts := bytes.Split(p, []byte{'\n'})
	for i, part := range parts {
		m.buffer.Write(part)
		if i < len(parts)-1 {
			if err = m.flush(); err != nil {
				return
			}
		}
	}
	return len(p), nil
}

func (m *messageStream) flush() (err error) {
	msg := new(LogMessage)
	if err = json.Unmarshal(m.buffer.Bytes(), msg); err != nil {
		return
	}
	msg.RawJSON = m.buffer.String()
	m.buffer.Reset()
	m.dispatcher.dispatch(msg)
	return
}
