package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type (
	GitHubConfig struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	}
	DataBaseConfig struct {
		Url string `toml:"url"`
	}
	Config struct {
		GitHubConfig   *GitHubConfig   `toml:"github"`
		DataBaseConfig *DataBaseConfig `toml:"database"`
	}
)

func LoadFromFile(path string) (*Config, error) {
	config := new(Config)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if _, err := toml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
