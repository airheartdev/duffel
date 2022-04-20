package duffel

import "net/http"

// WithAPIToken sets the API host to the default Duffel production host.
func WithDefaultAPI() Option {
	return func(c *Options) {
		c.Host = "https://api.duffel.com"
	}
}

// WithHost allows you to specify the Duffel API host to use for making requests.
func WithHost(host string) Option {
	return func(c *Options) {
		c.Host = host
	}
}

// WithVersion allows you to specify "Duffel-Version" header for the API version that you are targeting.
func WithAPIVersion(version string) Option {
	return func(c *Options) {
		c.Version = version
	}
}

// WithUserAgent allows you to specify a custom user agent string to use for making requests.
func WithUserAgent(ua string) Option {
	return func(c *Options) {
		c.UserAgent = ua
	}
}

// WithHTTPClient allows you to specify a custom http.Client to use for making requests.
// This is useful if you want to use a custom transport or proxy.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Options) {
		c.HttpDoer = client
	}
}

// WithDebug enables debug logging of requests and responses.
// DO NOT USE IN PRODUCTION.
func WithDebug(client *http.Client) Option {
	return func(c *Options) {
		c.Debug = true
	}
}
