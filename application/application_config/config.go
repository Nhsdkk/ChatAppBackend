package application_config

type Environment = string

const (
	Development Environment = "Development"
	Production              = "Production"
)

type ApplicationConfig struct {
	Url         string      `env:"URL"`
	Environment Environment `env:"ENVIRONMENT"`
}
