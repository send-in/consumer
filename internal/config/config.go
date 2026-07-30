package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	godotenv.Load(".env")
	return &Config{
		Server: ServerConfig{
			Port: ":" + GetEnv("PORT", "8001"),
			Passkey: GetEnv("PASSKEY", ""),
		},
		RabbitMQ: RabbitMQConfig{
			Username: GetEnv("RABBITMQ_USERNAME", "guest"),
			Password: GetEnv("RABBITMQ_PASSWORD", "feetlover"),
			Host: GetEnv("RABBITMQ_HOST", "localhost"),
			Port: GetEnv("RABBITMQ_PORT", "5672"),

			Queue: GetEnv("RABBITMQ_QUEUE", "jobs"),
			Type: GetEnv("RABBITMQ_TYPE", "message-send"),
			DeadQueue: GetEnv("RABBITMQ_DEAD_QUEUE", "dead-jobs"),

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