package config

type Redis struct {
	Host string `yaml:"host"`
	Port uint16   `yaml:"port"`
	Ssl  string `yaml:"ssl"`
}

type MailTp struct {
	Host string `yaml:"host"`
	Port uint16   `yaml:"port"`
	Ssl  bool `yaml:"ssl"`
}

type Mail struct {
	User     string `yaml:"user"`
	Password string `yaml:"pwd"`
	Smtp     MailTp `yaml:"smtp"`
	Imap     MailTp `yaml:"imap"`
}

type Config struct {
	Host         string  `yaml:"host"`
	Port         uint16    `yaml:"port"`
	Name         string  `yaml:"name"`
	Mail         Mail    `yaml:"mail"`
	RedisCluster []Redis `yaml:"redis"`
}
