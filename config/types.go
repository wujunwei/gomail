package config

type Redis struct {
	Address string `yaml:"address"`
	Name    string `yaml:"name"`
	Network string `yaml:"network"`
}

type Mail struct {
	User     string `yaml:"user"`
	Password string `yaml:"pwd"`
	Smtp     string `yaml:"smtp"`
	Imap     string `yaml:"imap"`
}

type Config struct {
	Host         string  `yaml:"host"`
	Port         uint16  `yaml:"port"`
	Name         string  `yaml:"name"`
	Mail         Mail    `yaml:"mail"`
	RedisCluster []Redis `yaml:"redis"`
}
