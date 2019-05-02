package conn

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/indrenicloud/tricloud-agent/app/logg"
)

type Config struct {
	UUID   string
	ApiKey string
	Url    string
}

var c *Config

func init() {
	logg.Log("config init")
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		logg.Log("Could not read config:", err)
		os.Exit(1)
	}
	c = &Config{}

	err = json.Unmarshal(data, c)
	if err != nil {
		logg.Log(err)
		os.Exit(1)
	}

	if c.Url == "" {
		c.Url = "localhost:8081"
	}

}

// GetConfig gives config file
func GetConfig() *Config {
	return c
}

// SaveConfig saves config file if it exists
func SaveConfig() {
	rawc, err := json.Marshal(c)
	if err != nil {
		logg.Log("Could not save config:", err)
		os.Exit(1)
	}
	ioutil.WriteFile("config.json", rawc, 0644)
}
