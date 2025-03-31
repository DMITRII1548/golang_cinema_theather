package config

import "sync"

var DatabaseConfig *databaseConfig
var onceDB sync.Once

type databaseConfig struct {
	Host string
	User string
	Pass string
	Name string
}

func InitDatabaseConfig() {
	onceDB.Do(func() {
		DatabaseConfig = &databaseConfig{
			Host: getEnv("DB_HOST", "localhost"),
			User: getEnv("DB_USER", "root"),
			Pass: getEnv("DB_PASS", "secret"),
			Name: getEnv("DB_NAME", "mydb"),
		}
	})
}