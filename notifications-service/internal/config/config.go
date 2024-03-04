package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Port        int    `json:"port"`
	RabbitmqUrl string `json:"rabbitmqUrl"`
}

var Config Configuration

func Load() error {
	file, err := os.Open("../configs/config.json")
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		return err
	}

	return nil
}
