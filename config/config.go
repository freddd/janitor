package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Tracker Tracker `yaml:"tracker"`
}

type Tracker struct {
	Keywords  []string `yaml:"keywords"`
	FileNames []string `yaml:"fileNames"`
	RepoPath  string   `yaml:"repoPath"`
	WhiteList []string `yaml:"whitelist"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Validate(config Config) error {
	return nil
}
