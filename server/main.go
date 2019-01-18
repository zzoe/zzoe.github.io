package main

import (
	"github.com/spf13/viper"
	"github.com/zzoe/zoe.github.io/server/cfg"
	"github.com/zzoe/zoe.github.io/server/httpsrv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	log = cfg.Log
)

func main() {
	log.Debug("main begin")
	defer func() {
		if err := log.Sync(); err != nil {
			panic(err)
		}
	}()

	end := make(chan struct{})
	go httpsrv.Start(end)

	sysQuit := make(chan os.Signal)
	signal.Notify(sysQuit, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case <-end:
		log.Info("End of program")
	case <-sysQuit:
		go httpsrv.Stop()
		select {
		case <-end:
			log.Info("Program is closed")
		case <-time.After(viper.GetDuration("server.shuttimeout")):
			log.Info("Program forced shutdown")
		}
	}
}
