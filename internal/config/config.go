package config

import (
	"time"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xyconfig"
)

var config *xyconfig.Config

func init() {
	config = xyconfig.GetConfig("xyauth")

	xycond.AssertNil(config.ReadFile("configs/10-default.ini", true))

	d := config.GetDefault("general.config_watch", time.Minute).MustDuration()
	config.SetWatchInterval(d)

	if files, ok := config.Get("general.additions"); ok {
		for _, f := range files.MustArray() {
			if val, ok := f.AsString(); ok && val != "" {
				xycond.AssertNil(config.Read(val))
			}
		}
	}

	// Load environment variables.
	xycond.AssertNil(config.Read("env"))

	if config.GetDefault("general.environment", "dev").MustString() == "dev" {
		xycond.AssertNil(config.Read(".env"))
	}
}

func Add(f string) error {
	return config.Read(f)
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
