package utils

import (
	"fmt"
	"os"

	"github.com/xybor-x/xyconfig"
	"github.com/xybor-x/xylog"
)

var logger *xylog.Logger

func GetLogger() *xylog.Logger {
	return logger
}

func initLogger() {
	emitter := xylog.NewStreamEmitter(os.Stdout)
	handler := xylog.GetHandler("xybor.auth")
	handler.AddMacro("time", "asctime")
	handler.AddMacro("level", "levelname")
	handler.AddEmitter(emitter)

	logger = xylog.GetLogger("xybor.auth")
	logger.SetLevel(GetConfig().GetDefault("general.loglevel", xylog.INFO).MustInt())
	logger.AddHandler(handler)

	config.AddHook("general.loglevel", func(e xyconfig.Event) {
		logger.SetLevel(e.New.MustInt())
		fmt.Println("Set log level to ", e.New.MustInt())
	})
}
