package delivery

import "time"

type Config struct {
	Timeout   time.Duration
	UserAgent string
}
