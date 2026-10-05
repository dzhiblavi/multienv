package redis

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (c *Credentials) Default() {
}

type Config struct {
	Address         string `yaml:"address"`
	Port            int    `yaml:"port"`
	DB              int    `yaml:"db"`
	CredentialsPath string `yaml:"credentials_path"`
	credentials     Credentials
}

func (c *Config) Default() {
	c.Address = "127.0.0.1"
	c.Port = 12123
}
