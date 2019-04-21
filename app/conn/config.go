package conn

import (
	"encoding/json"
	"io/ioutil"

	"github.com/indrenicloud/tricloud-agent/app/logg"
)

type Config struct {
	UUID   string
	ApiKey string
}

// GetConfig gives config file if it exists
func GetConfig() *Config {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		logg.Log("Could not read config.")
		return nil
	}
	c := &Config{}

	err = json.Unmarshal(data, c)
	if err == nil {
		return c
	}

	return nil
}

// SaveConfig saves config file if it exists
func SaveConfig(c *Config) {
	rawc, err := json.Marshal(c)
	if err != nil {
		logg.Log("Could not save config.")
	}
	ioutil.WriteFile("config.json", rawc, 0644)
}
