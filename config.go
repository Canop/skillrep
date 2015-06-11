package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type DbConfig struct {
	Name     string
	User     string
	Password string
}

type Config struct {
	DB     DbConfig
	ApiKey string
}

var config Config

func (dbc DbConfig) queryString() string {
	return fmt.Sprintf("user=%s dbname=%s password=%s", dbc.User, dbc.Name, dbc.Password)
}

func ReadConfig() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatal("Error while reading config: ", err)
	}
}
