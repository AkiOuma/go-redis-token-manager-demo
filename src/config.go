package src

import "fmt"

type Config struct {
	address string
	port    string
}

func NewConfig(address, port string) *Config {
	return &Config{
		address: address,
		port:    port,
	}
}

func (c *Config) Url() string {
	return fmt.Sprintf("%s:%s", c.address, c.port)
}
