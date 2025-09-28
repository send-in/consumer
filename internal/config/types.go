package config

type ServerConfig struct {
	Port string
}

type RabbitMQConfig struct {
	Username string
	Password string
	Host string
	Port string
	Queue string
	Type string
}

type Config struct {
	Server   ServerConfig
	RabbitMQ RabbitMQConfig
}