package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var (
	LocalConfig ConfigInfo
)

const (
	kAppName       = "APP_NAME"
	kConfigServer  = "CONFIG_SERVER"
	kConfigProfile = "CONFIG_PROFILE"
	kConfigType    = "CONFIG_TYPE"
)

type ConfigInfo struct {
	Greetings struct {
		String string `json:"string"`
	} `json:"greetings"`
}

func GetConfig(configServerHost string, serverName string, profile string) ConfigInfo {
	var config ConfigInfo

	viper.AutomaticEnv()

	viper.SetDefault(kAppName, serverName)
	viper.SetDefault(kConfigServer, configServerHost)
	viper.SetDefault(kConfigProfile, profile)
	viper.SetDefault(kConfigType, "yaml")

	if err := loadRemoteConfig(); err != nil {
		log.Fatal("Fail to load config", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Fail to parse config", err)
	}

	return config
}

func loadRemoteConfig() (err error) {
	confAddr := fmt.Sprintf("%v/%v-%v.%v",
		viper.Get(kConfigServer), viper.Get(kAppName), viper.Get(kConfigProfile), viper.Get(kConfigType))
	resp, err := http.Get(confAddr)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	viper.SetConfigType(viper.GetString(kConfigType))
	if err = viper.ReadConfig(resp.Body); err != nil {
		return
	}
	log.Println("Load config from: ", confAddr)
	return
}