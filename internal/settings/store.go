package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"pass-generator/internal/password"
	"path/filepath"
)

type jsonConfig struct {
	Length  int    `json:"length"`
	Exclude string `json:"exclude"`
}
type Store struct {
	Path string
}

func New(path string) Store {
	return Store{Path: path}
}

func (store Store) Save(config password.Config) error {
	fileConfig := jsonConfig{
		Length:  config.Length,
		Exclude: config.Exclude,
	}

	data, err := json.MarshalIndent(fileConfig, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(store.Path, data, 0600)
	return err
}

func UserPath() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(baseDir, "PassGenerator")
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return "", err
	}
	result := filepath.Join(dir, "settings.json")
	return result, nil
}

func (store Store) Load() (password.Config, error) {
	data, err := os.ReadFile(store.Path)
	if err != nil {
		return password.Config{}, err
	}
	var fileConfig jsonConfig

	err = json.Unmarshal(data, &fileConfig)
	if err != nil {
		return password.Config{}, err
	}
	if fileConfig.Length <= 0 {
		return password.Config{}, fmt.Errorf("length must be greater than 0")
	}

	config := password.Config{
		Length:  fileConfig.Length,
		Exclude: fileConfig.Exclude,
	}

	return config, nil
}
