// appconfig contains logic that defines configuration for a given mod instance
// appconfig setings are user facing settings
package appconfig

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"path"
)

type AppConfig struct {
	GameDir string `json:"gameDir"`
	ModDir  string `json:"modDir"`
}

func writeConfigFile(confFile string, config AppConfig) error {
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(confFile, b, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func ensureConfig(confPath string) (AppConfig, error) {
	defaultConfig := AppConfig{GameDir: "unknown", ModDir: "unknown"}

	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		slog.Debug("write config")
		writeConfigFile(confPath, defaultConfig)

		return defaultConfig, nil
	} else {
		cfg, err := readConfFile(confPath)
		if err != nil {
			return AppConfig{}, err
		}

		// validate ron dir
		if err := validateRonDir(cfg.GameDir); err != nil {
			return AppConfig{}, err
		}
		//
		// validate mod dir
		if err := validateModDir(cfg.ModDir); err != nil {
			return AppConfig{}, err
		}

		return cfg, nil
	}
}

func readConfFile(cfgPath string) (AppConfig, error) {
	file, err := os.ReadFile(cfgPath)
	if err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}

func Setup() error {
	confPath := "./test/"

	// Detect and setup config
	// Create instance dir structure
	// generate initial state file
	configFilePath := path.Join(confPath, "ron-mm.conf.json")

	if err := os.MkdirAll(confPath, 0700); err != nil {
		return err
	}

	_, err := ensureConfig(configFilePath)
	if err != nil {
		return err
	}
	// todo resync stuff?

	return nil
}
