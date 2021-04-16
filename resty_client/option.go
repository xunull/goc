package resty_client

type option struct {
	TokenLower     bool
	TokenValue     string
	AddTokenHeader bool
	HostUrl        string
	Proxy          string
}

type RestyOption func(o *option)

func getDefaultOption(ops ...RestyOption) *option {
	d := &option{

	}
	for _, o := range ops {
		o(d)
	}
	return d
}

func WithProxy(proxy string) RestyOption {
	return func(o *option) {
		o.Proxy = proxy
	}
}

func WithLowerTokenHeader(token string) RestyOption {
	return func(o *option) {
		o.TokenLower = true
		o.TokenValue = token
		o.AddTokenHeader = true
	}
}

func WithTokenHeader(token string) RestyOption {
	return func(o *option) {
		o.TokenLower = false
		o.TokenValue = token
		o.AddTokenHeader = true
	}
}

func WithHostUrl(url string) RestyOption {
	return func(o *option) {
		o.HostUrl = url
	}
}
