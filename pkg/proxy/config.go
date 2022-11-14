package proxy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type App struct {
	Name    string
	Ports   []int
	Targets []string
}

type Config struct {
	Apps []App
}

func LoadConfigFromFile(path string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	if !json.Valid(data) {
		return config, errors.New("invalid json file")
	}

	json.Unmarshal(data, &config)

	return config, nil

}
