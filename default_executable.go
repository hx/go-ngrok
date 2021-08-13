package ngrok

import "os/exec"

// DefaultExecutable is an Executable roughly equivalent to running `ngrok` in the working directory.
var DefaultExecutable = Executable(exec.Command("ngrok").Path)

// IsAvailable is true when the default executable exists and has the necessary permissions to be executed.
func IsAvailable() bool { return DefaultExecutable.IsAvailable() }

// NewProcess creates a new Process with the given arguments, using the default executable.
func NewProcess(args ...string) *Process { return DefaultExecutable.NewProcess(args...) }

// NewTunnel creates a new generic tunnel from the given request.
func NewTunnel(request TunnelRequest) *Tunnel { return DefaultExecutable.NewTunnel(request) }

// NewHTTPTunnel creates a new HTTP tunnel from the given request.
func NewHTTPTunnel(request *HTTPTunnelRequest) *HTTPTunnel {
	return DefaultExecutable.NewHTTPTunnel(request)
}
