package config

var defaultProxyHost string = ":9000"

type Proxy struct {
	Host       string
	StaticPath string `yaml:"static_path"`
}

func (p *Proxy) GetHost() string {
	if p.Host == "" {
		return defaultProxyHost
	}

	return p.Host
}
