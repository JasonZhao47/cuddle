package startup

import (
	"github.com/jasonzhao47/cuddle/internal/logger"
	"go.uber.org/zap"
)

func InitLog() logger.Logger {
	// TODO: replace with no-op logger
	return logger.NewLogger(zap.L())
}
