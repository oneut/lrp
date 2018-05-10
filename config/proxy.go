package config

var defaultProxyScheme string = "http"
var defaultProxyHost string = ":9000"

type Proxy struct {
	Scheme      string
	Host        string
	StaticPath  string `yaml:"staticPath"`
	BrowserOpen string `yaml:"browserOpen"`
}

func (p *Proxy) GetScheme() string {
	if p.Scheme == "" {
		return defaultProxyScheme
	}

	return p.Scheme
}

func (p *Proxy) GetHost() string {
	if p.Host == "" {
		return defaultProxyHost
	}

	return p.Host
}

func (p *Proxy) IsBrowserOpen() bool {
	if p.BrowserOpen == "none" {
		return false
	}

	return true
}
