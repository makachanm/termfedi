package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func LoadConfig() (Configuration, error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error in getting fs: ", err)
		return Configuration{}, err
	}

	if _, err := os.Stat(filepath.Join(wd, "config.json")); os.IsNotExist(err) {
		excuteable, _ := os.Executable()
		wd = path.Dir(excuteable)
		if _, err := os.Stat(filepath.Join(wd, "config.json")); os.IsNotExist(err) {
			return Configuration{}, errors.New("config file not exist")
		}
	}

	ctx, err := os.ReadFile(filepath.Join(wd, "config.json"))
	if err != nil {
		return Configuration{}, err
	}

	var cfg = Configuration{}
	err = json.Unmarshal(ctx, &cfg)
	if err != nil {
		return Configuration{}, err
	}

	return cfg, nil
}
