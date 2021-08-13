package ngrok

// TunnelRequest is the base interface for requests to make Tunnel instances.
type TunnelRequest interface {
	Args() []string
}
