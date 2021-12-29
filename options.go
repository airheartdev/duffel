package duffel

func WithDefaultAPI() Option {
	return func(c *Options) {
		c.Host = "https://api.duffel.com"
	}
}

func WithHost(host string) Option {
	return func(c *Options) {
		c.Host = host
	}
}

func WithAPIVersion(version string) Option {
	return func(c *Options) {
		c.Version = version
	}
}

func WithUserAgent(ua string) Option {
	return func(c *Options) {
		c.UserAgent = ua
	}
}
