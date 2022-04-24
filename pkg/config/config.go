package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Config struct {
	ListenAddr string `json:"listen_addr"`
	RootDir    string `json:"root_dir"`
}

var defaultConfig = Config{
	ListenAddr: ":8080",
	RootDir:    "/var/pics",
}

func NewConfig(filename string) Config {
	var config Config
	fh, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultConfig
		}
		panic(err)
	}
	bytes, err := io.ReadAll(fh)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}
	return config
}
