package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type GlobalConfig struct {
	DB     DbConfig
	ApiKey string
	Port   int
}

type DbConfig struct {
	Name     string
	User     string
	Password string
}

var config GlobalConfig

func Config() *GlobalConfig {
	if config.DB.Name == "" {
		file, _ := os.Open("config.json")
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&config)
		if err != nil {
			log.Fatal("Error while reading config: ", err)
		}
	}
	return &config
}

func (dbc DbConfig) queryString() string {
	return fmt.Sprintf("user=%s dbname=%s password=%s", dbc.User, dbc.Name, dbc.Password)
}
