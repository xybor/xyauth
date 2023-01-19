package logger

import (
	"os"

	"github.com/xybor-x/xyconfig"
	"github.com/xybor-x/xylog"
	"github.com/xybor/xyauth/internal/config"
)

var logger *xylog.Logger

// Event creates an EventLogger which logs key-value pairs.
func Event(e string) *xylog.EventLogger {
	return logger.Event(e)
}

func init() {
	emitter := xylog.NewStreamEmitter(os.Stdout)
	handler := xylog.GetHandler("xybor.auth")
	handler.AddMacro("time", "asctime")
	handler.AddMacro("level", "levelname")
	handler.AddEmitter(emitter)

	logger = xylog.GetLogger("xybor.auth")
	logger.SetLevel(config.GetDefault("general.loglevel", xylog.INFO).MustInt())
	logger.AddHandler(handler)

	config.AddHook("general.loglevel", func(e xyconfig.Event) {
		logger.SetLevel(e.New.MustInt())
	})
}
