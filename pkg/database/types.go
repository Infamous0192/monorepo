package database

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	Username string `yaml:"username" env:"USERNAME" env-default:"username"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"password"`
	Database string `yaml:"database" env:"DATABASE" env-default:"database"`
}

type Config struct {
	Redis   DatabaseConfig `yaml:"redis" env-prefix:"REDIS_"`
	MongoDB DatabaseConfig `yaml:"mongodb" env-prefix:"MONGODB_"`
}
