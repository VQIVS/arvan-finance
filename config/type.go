package config

import "google.golang.org/genproto/googleapis/type/decimal"

type Config struct {
	DB     DBConfig        `json:"db"`
	Server ServerConfig    `json:"server"`
	Redis  RedisConfig     `json:"redis"`
	SmsFee decimal.Decimal `json:"sms_fee"`
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
