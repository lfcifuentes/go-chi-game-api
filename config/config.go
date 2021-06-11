package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Engine_sql    string
	Port          string
	Username      string
	Password      string
	Database      string
	Cluster       string
	Host          string
	SSL_root_cert string
}

func (c *Config) LoadEnv() {
	godotenv.Load(".env")
	c.Engine_sql = os.Getenv("DB_ENGINE_SQL")
	c.Port = os.Getenv("DB_PORT")
	c.Username = os.Getenv("DB_USERNAME")
	c.Password = os.Getenv("DB_PASSWORD")
	c.Database = os.Getenv("DB_DATABASE")
	c.Cluster = os.Getenv("DB_CLUSTER")
	c.SSL_root_cert = os.Getenv("DB_SSL_ROOT_CERT")
	c.Host = os.Getenv("DB_HOST")
}
