package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	godotenv.Load("../../.env")
	return &Config{
		Server: ServerConfig{
			Port: ":" + getEnv("PORT", "8000"),
		},
		RabbitMQ: RabbitMQConfig{
			Username: getEnv("RABBITMQ_USERNAME", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", "guest"),
			Host: getEnv("RABBITMQ_HOST", "localhost"),
			Port: getEnv("RABBITMQ_PORT", "5672"),
			Queue: getEnv("RABBITMQ_QUEUE", "jobs"),
			Type: getEnv("RABBITMQ_TYPE", "message-send"),
		},
	}, nil
}

func (cfg *RabbitMQConfig) GetRabbitMQURL() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
}