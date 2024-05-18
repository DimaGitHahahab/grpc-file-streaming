package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoURI string `envconfig:"MONGO_URI" default:"mongodb://mongo:27017"`
	DbName   string `envconfig:"DB_NAME" default:"user_files_db"`
	Port     string `envconfig:"SERVER_PORT" default:"50051"`
}

func LoadConfig() *Config {
	log.Println("Loading env variables")

	if err := godotenv.Load(); err != nil {
		log.Fatalln("Failed to load env variables: ", err)
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalln("Failed to process env variables: ", err)
	}

	log.Println("Loaded env variables successfully")
	return &cfg
}
