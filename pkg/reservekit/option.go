package reservekit

// Option represents a client option
type Option func(*Client)

// WithHost sets the API host
func WithHost(host string) Option {
	return func(c *Client) {
		c.host = host
	}
}

// WithVersion sets the API version
func WithVersion(version string) Option {
	return func(c *Client) {
		c.version = version
	}
}
