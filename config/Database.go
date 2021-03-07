package config

import (
	"fmt"
	"keyholders/helpers"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
		Host:     helpers.DotEnvVariable("MYSQL_HOST"),
		Port:     3306,
		User:     helpers.DotEnvVariable("MYSQL_USER"),
		Password: helpers.DotEnvVariable("MYSQL_PASSWORD"),
		DBName:   helpers.DotEnvVariable("MYSQL_DBNAME"),
	}
	return &dbConfig
}

func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}
