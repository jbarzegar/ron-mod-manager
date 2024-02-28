package config

import (
	"encoding/json"
	"os"

	"github.com/jbarzegar/ron-mod-manager/types"
)

var conf types.MMConfig

var ConfPath string

var DBPath string

func ReadConfFile(p string) types.MMConfig {
	file, err := os.ReadFile(p)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(file, &conf)

	return conf
}

func SetConfig(c types.MMConfig) types.MMConfig {
	conf = c

	return conf
}

func GetConfig() types.MMConfig {

	return conf
}
