package ngrok

import (
	"context"
	"time"
)

const WebServiceTimeout = time.Second * 5

// Tunnel is the base type for HTTPTunnel, and, in the future, other types of tunnels.
type Tunnel struct {
	*Process

	webService     *URL
	webServiceDone chan struct{}
}

func (e Executable) NewTunnel(request TunnelRequest) (*Tunnel, error) {
	version, err := e.Version()
	if err != nil {
		return nil, err
	}
	return &Tunnel{
		Process:        e.NewProcess(request.Args(version)...),
		webServiceDone: make(chan struct{}),
	}, nil
}

func (t *Tunnel) Start() error {
	go t.waitForWebService()
	return t.Process.Start()
}

// WebService returns the URL of the process's web service.
//
// If the process has not reported a web service URL, WebService will block for up to WebServiceTimeout, starting from
// when Start was called.
//
// If no web service URL is reported after WebServiceTimeout has elapsed, WebService will return nil.
func (t *Tunnel) WebService() *URL {
	<-t.webServiceDone
	return t.webService
}

func (t *Tunnel) waitForWebService() {
	timeout, cancel := context.WithTimeout(context.Background(), WebServiceTimeout)
	t.WaitOneContext(timeout, func(message *LogMessage) bool {
		if message.Object == "web" && message.Address != nil {
			u := *message.Address
			if u.Scheme == "tcp" {
				u.Scheme = "http"
			}
			t.webService = &u
			return true
		}
		return false
	})
	cancel()
	close(t.webServiceDone)
}
