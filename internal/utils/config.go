package utils

import (
	"log"

	"github.com/xybor-x/xyconfig"
)

var config *xyconfig.Config

func initConfig() {
	config = xyconfig.GetConfig("xyauth")

	if err := config.ReadFile("configs/default.ini", true); err != nil {
		log.Panic(err)
	}

	d := config.GetDefault("general.env_watch_cycle", 0).MustDuration()
	if err := config.LoadEnv(d); err != nil {
		log.Panic(err)
	}

	if config.GetDefault("general.environment", "dev").MustString() == "dev" {
		if err := config.ReadFile(".env", true); err != nil {
			log.Panic(err)
		}
	}
}

func AddConfig(f string) error {
	return GetConfig().ReadFile(f, true)
}

func GetConfig() *xyconfig.Config {
	return config
}
