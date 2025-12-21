package main

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	URL string `yaml:"url"`
}

var conf Config

func confDirExists() (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	path := filepath.Join(home, ".ekhoes")

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}

func createEkhoesConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(home, ".ekhoes")

	// 0700 → solo l'utente può accedere
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return err
	}

	configPath := filepath.Join(home, ".ekhoes/conf.yml")

	if _, err := os.Stat(configPath); err == nil {
		return errors.New("il file ~/.ekhoes esiste già")
	} else if !os.IsNotExist(err) {
		return err
	}

	cfg := Config{
		URL: "https://websocket.ekhoes.com",
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	// 0600 → solo l'utente può leggere/scrivere
	return os.WriteFile(configPath, data, 0600)
}

func loadEkhoesConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".ekhoes", "conf.yml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		return err
	}

	return nil
}
