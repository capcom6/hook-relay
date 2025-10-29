package config

type Config struct {
	Server Server `koanf:"server"`
}

type Server struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`
}

func Default() Config {
	return Config{
		Server: Server{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
		},
	}
}
