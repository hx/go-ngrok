package ngrok

import (
	"context"
	"errors"
	"net/url"
	"time"
)

const HTTPTunnelURLsTimeout = time.Second * 10

// HTTPTunnel represents a process running a pair of HTTP/HTTPS tunnels.
type HTTPTunnel struct {
	*Tunnel
	InsecureURL *url.URL
	SecureURL   *url.URL
}

// NewHTTPTunnel creates a new HTTP tunnel from the given request.
func (e Executable) NewHTTPTunnel(request *HTTPTunnelRequest) *HTTPTunnel {
	return &HTTPTunnel{Tunnel: e.NewTunnel(request)}
}

// Start opens the tunnel, and returns when it is open and both SecureURL and InsecureURL are set. If
// HTTPTunnelURLsTimeout elapses before they are set, an error will be returned.
func (t *HTTPTunnel) Start() (err error) {
	ready := make(chan error, 1) // Buffer allows for early return if tunnel doesn't start
	go func() { ready <- t.waitForURLs() }()
	if err = t.Tunnel.Start(); err != nil {
		return
	}
	err = <-ready

	// A canceled context means the process exited.
	if errors.Is(err, context.Canceled) {
		err = t.Wait()
	}
	return
}

func (t *HTTPTunnel) waitForURLs() (err error) {
	timeout, cancel := context.WithTimeout(context.Background(), HTTPTunnelURLsTimeout)
	err = t.WaitOneContext(timeout, func(message *LogMessage) bool {
		if message.Object == "tunnels" && message.URL != nil {
			switch message.URL.Scheme {
			case "http":
				t.InsecureURL = message.URL.Raw()
			case "https":
				t.SecureURL = message.URL.Raw()
			}
		}
		return t.InsecureURL != nil && t.SecureURL != nil
	})
	cancel()
	return
}
