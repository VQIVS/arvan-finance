package config

type Config struct {
	DB     DBConfig     `json:"db"`
	Server ServerConfig `json:"server"`
	Rabbit RabbitConfig `json:"rabbit"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"database"`
	Schema   string `json:"schema"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type ServerConfig struct {
	HttpPort uint   `json:"httpPort"`
	Secret   string `json:"secret"`
}

type RedisConfig struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

type RabbitConfig struct {
	URL    string        `json:"url"`
	Queues []QueueConfig `json:"queues"`
}

type QueueConfig struct {
	Name    string `json:"name"`
	Durable bool   `json:"durable"`
}
