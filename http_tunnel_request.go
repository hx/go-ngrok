package ngrok

import "strconv"

// HTTPTunnelRequest represents the input (arguments) for an HTTPTunnel.
type HTTPTunnelRequest struct {
	Port      int
	Subdomain string
	Hostname  string
	HTTPS     bool
}

func (r *HTTPTunnelRequest) Args(version *Version) (args []string) {
	args = []string{"http"}

	if r.Subdomain != "" {
		args = append(args, "-subdomain="+r.Subdomain)
	}

	host := strconv.Itoa(r.Port)

	if host == "0" {
		if r.HTTPS {
			host = "443"
		} else {
			host = "80"
		}
	}

	if r.Hostname != "" {
		host = r.Hostname + ":" + host
	} else if r.HTTPS {
		host = "localhost:" + host
	}

	if r.HTTPS {
		host = "https://" + host
	}

	args = append(args, host)

	if version.Major >= 3 {
		args = append(args, "--scheme", "http,https")
	}

	return args
}
