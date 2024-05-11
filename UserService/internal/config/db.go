package config

import (
	"UserSVC/internal/entities"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMODE  string
}

var config = Config{}

func Connect() (*gorm.DB, error) {
	config.Read()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMODE,
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	DB = conn

	err = conn.AutoMigrate(
		&entities.User{},
		&entities.Product{},
		&entities.Purchase{},
		&entities.Token{},
	)

	if err != nil {
		log.Fatal(err)
	}

	return conn, err
}

func (c *Config) Read() {
	c.Host = os.Getenv("DB_HOST")
	c.User = os.Getenv("DB_USER")
	c.Password = os.Getenv("DB_PASSWORD")
	c.DBName = os.Getenv("DB_NAME")
	c.Port = os.Getenv("DB_PORT")
	c.SSLMODE = os.Getenv("DB_SSLMODE")
}
