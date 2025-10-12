package config

type ServerConfig struct {
	Port string
	Passkey string
}

type RabbitMQConfig struct {
	Username string
	Password string
	Host string
	Port string

	Queue string
	DeadQueue string
	DeadLenght int
	Type string

}

type Config struct {
	Server   ServerConfig
	RabbitMQ RabbitMQConfig
}