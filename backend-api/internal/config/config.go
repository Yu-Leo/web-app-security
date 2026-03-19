package config

type Env string

const (
	EnvProd  = "prod"
	EnvStage = "stage"
	EnvLocal = "local"
	EnvTests = "tests"
)

// Config - service config.
type Config struct {
	Env Env

	HTTPServer ServerConfig
	GRPCServer ServerConfig

	PostgresMaster string `required:"true" split_words:"true"`
}

type ServerConfig struct {
	Host string `required:"true"`
	Port int    `required:"true"`
}
