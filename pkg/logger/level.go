package pkglogger

import (
	"os"

	"github.com/rs/zerolog"
)

func SetLogLevel() {
	levelStr := os.Getenv("LOG_LEVEL")

	level, err := zerolog.ParseLevel(levelStr)
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
}
