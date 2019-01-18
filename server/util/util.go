package util

import (
	"github.com/zzoe/zoe.github.io/server/cfg"
	"go.uber.org/zap"
)

var (
	log = cfg.Log
)

func Warn(err error) {
	if err != nil {
		log.Warn("WARN", zap.Error(err))
	}
}
