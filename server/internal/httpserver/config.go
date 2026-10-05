package httpserver

type Config struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func (c *Config) Default() {
	c.Address = "127.0.0.1"
	c.Port = 12123
}
