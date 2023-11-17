package app

import "time"

type Config struct {
	Endpoint struct {
		Addr                  string
		HTTPReadHeaderTimeout time.Duration
		Hosts                 []string
	}
}
