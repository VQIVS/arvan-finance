package config

type Config struct {
	Server   Server   `yaml:"server"`
	DB       DB       `yaml:"database"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RabbitMQ struct {
	URI    string  `yaml:"uri"`
	Queues []Queue `yaml:"queues"`
}

type Queue struct {
	Name     string `yaml:"name"`
	Exchange string `yaml:"exchange"`
	Routing  string `yaml:"routing"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
}
