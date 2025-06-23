package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Config struct {
	Db_url           string `json:"db_url"`
	Current_Username string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getHomeDir() (string, error) {

	file_path, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Failed to retrieve home directory for user")
	}

	return file_path, nil
}

// Set username to Config struct and write to a json file
func (config Config) SetUser() {

	data, err := json.Marshal(config)
	if err != nil {
		log.Fatal("Failed to encode Config struct to bytes")
	}

	home_path, err := getHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	full_path := home_path + "/" + configFileName

	err = os.WriteFile(full_path, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Read the json configuration file and return Config struct
func Read() Config {

	// Find the home directory
	file_path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to retrieve home directory for user")
	}

	full_path := file_path + "/" + configFileName

	// Read the configuration file called gatorconfig.json
	//fmt.Println(full_path)
	data, err := os.ReadFile(full_path)
	if err != nil {
		log.Fatal("Failed to read config file")
	}

	// Decode the data into a Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Failed to decode json data into Config struct")
	}

	return config

}
