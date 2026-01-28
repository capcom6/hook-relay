package config

import "time"

type Config struct {
	Server   Server   `koanf:"server"`
	Database Database `koanf:"database"`
	Events   Events   `koanf:"events"`
}

type Server struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`
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

func Default() Config {
	return Config{
		Server: Server{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
		},
		Database: Database{
			URL:             "mysql://hook-relay:hook-relay@127.0.0.1:3306/hook-relay?charset=utf8mb4,utf8&parseTime=true",
			ConnMaxIdleTime: 0,
			ConnMaxLifetime: 0,
			MaxOpenConns:    0,
			MaxIdleConns:    0,
		},
		Events: Events{
			Timeout: time.Second,
		},
	}
}
