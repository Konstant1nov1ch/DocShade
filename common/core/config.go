package core

import (
	"os"

	"gitlab.com/docshade/common/log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config interface {
	GetPostgresConfig() PostgresConfig
	GetS3Config() S3Config
	GetRedisConfig() RedisConfig
	GetRabbitMQConfig() RabbitMQConfig
	GetAnonymizerConfig() AnonymizerConfig
	// AddHandler добавить ручку в конфигурацию
	AddHandler(handler Handler) Config
	// GetHandlerList получить список ручек из конфигурации
	GetHandlerList() []Handler
	LoadConfig() error
	// GetLogConfig получить конфигурацию для логгера
	GetLogConfig() log.LoggerConfig
	// GetPort получить порт приложения
	GetPort() string
}

const (
	defaultPort = "8080"
)

type config struct {
	name     string
	services Services
	port     string
	path     string
	handlers []Handler
}

type PostgresConfig struct {
	Host     string `yaml:"postgres_host"`
	DBName   string `yaml:"postgres_db_name"`
	UserName string `yaml:"postgres_user_name"`
	Password string `yaml:"postgres_password"`
	Port     string `yaml:"postgres_port"`
}

type RedisConfig struct {
	Host     string `yaml:"redis_host"`
	Password string `yaml:"redis_password"`
	DB       int    `yaml:"redis_db"`
	Port     string `yaml:"redis_port"`
}

type S3Config struct {
	Endpoint        string `yaml:"s3_endpoint"`
	AccessKeyID     string `yaml:"s3_accessKeyID"`
	SecretAccessKey string `yaml:"s3_secretAccessKey"`
}

type RabbitMQConfig struct {
	URI      string `yaml:"rabbitmq_uri"`
	Username string `yaml:"rabbitmq_username"`
	Password string `yaml:"rabbitmq_password"`
}

type AnonymizerConfig struct {
	URI string `yaml:"anonymizer_url"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Services struct {
	LogConfig        log.LoggerConfig
	PostgresConfig   PostgresConfig   `yaml:"postgres"`
	RedisConfig      RedisConfig      `yaml:"redis"`
	S3Config         S3Config         `yaml:"s3"`
	RabbitMQConfig   RabbitMQConfig   `yaml:"rabbitmq"`
	ServerConfig     ServerConfig     `yaml:"server"`
	AnonymizerConfig AnonymizerConfig `yaml:"py_anonymizer"`
}

func NewConfig(name string) Config {
	return &config{
		name:     name,
		services: Services{},
		path:     "",
	}
}

// GetHandlerList получить список ручек из конфигурации
func (c *config) GetHandlerList() []Handler {
	return c.handlers
}

// AddHandler добавить ручку в конфигурацию
func (c *config) AddHandler(handler Handler) Config {
	c.handlers = append(c.handlers, handler)

	return c
}

func (c *config) GetPostgresConfig() PostgresConfig {
	return c.services.PostgresConfig
}

func (c *config) GetAnonymizerConfig() AnonymizerConfig {
	return c.services.AnonymizerConfig
}

func (c *config) GetRedisConfig() RedisConfig {
	return c.services.RedisConfig
}

func (c *config) LoadConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		c.path = configPath
	}
	if err := cleanenv.ReadConfig(configPath, &c.services); err != nil {
		return err
	}

	c.port = c.services.ServerConfig.Port

	return nil
}

func (c *config) GetPort() string {
	if c.port == "" {
		c.port = defaultPort
	}

	return c.port
}

func (c *config) GetLogConfig() log.LoggerConfig {
	return c.services.LogConfig
}

func (c *config) GetS3Config() S3Config {
	return c.services.S3Config
}

func (c *config) GetRabbitMQConfig() RabbitMQConfig {
	return c.services.RabbitMQConfig
}
