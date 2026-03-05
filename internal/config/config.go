package config

import "time"

type Config struct {
	HTTP     http     `koanf:"http"`
	Database Database `koanf:"database"`
	Events   Events   `koanf:"events"`
	Delivery Delivery `koanf:"delivery"`
}

type http struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`

	OpenAPI openAPIConfig `koanf:"openapi"`
}

type openAPIConfig struct {
	Enabled    bool   `koanf:"enabled"`
	PublicHost string `koanf:"public_host"`
	PublicPath string `koanf:"public_path"`
}

type Database struct {
	URL string `koanf:"url"`

	ConnMaxIdleTime time.Duration `koanf:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
	MaxOpenConns    int           `koanf:"max_open_conns"`
	MaxIdleConns    int           `koanf:"max_idle_conns"`
}

type Events struct {
	Timeout time.Duration `koanf:"timeout"`
}

type Delivery struct {
	Timeout   time.Duration `koanf:"timeout"`
	UserAgent string        `koanf:"user_agent"`
}

func Default() Config {
	//nolint:gosec // default values
	return Config{
		HTTP: http{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
			OpenAPI: openAPIConfig{
				Enabled:    false,
				PublicHost: "",
				PublicPath: "",
			},
		},
		Database: Database{
			URL:             "mysql://hook-relay:hook-relay@127.0.0.1:3306/hook-relay?charset=utf8mb4&parseTime=true",
			ConnMaxIdleTime: 0,
			ConnMaxLifetime: 0,
			MaxOpenConns:    0,
			MaxIdleConns:    0,
		},
		Events: Events{
			Timeout: time.Second,
		},
		Delivery: Delivery{
			Timeout:   time.Second,
			UserAgent: "hook-relay",
		},
	}
}
