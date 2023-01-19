package config

import (
	"log"

	"github.com/xybor-x/xyconfig"
)

var config *xyconfig.Config

func init() {
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

func Add(f string) error {
	return config.ReadFile(f, true)
}

func Get(name string) (xyconfig.Value, bool) {
	return config.Get(name)
}

func GetDefault(name string, def any) xyconfig.Value {
	return config.GetDefault(name, def)
}

func MustGet(name string) xyconfig.Value {
	return config.MustGet(name)
}

func ToMap() map[string]any {
	return config.ToMap()
}

func AddHook(key string, f func(xyconfig.Event)) {
	config.AddHook(key, f)
}
