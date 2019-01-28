package httpsrv

import (
	"context"
	"github.com/spf13/viper"
	"github.com/zzoe/zoe.github.io/server/cfg"
	"github.com/zzoe/zoe.github.io/server/define"
	"github.com/zzoe/zoe.github.io/server/util"
	"go.uber.org/zap"
	"net/http"
)

var (
	id      int
	srv     *http.Server
	srvQuit chan struct{}
	log     = cfg.Log
)

func Start(end chan struct{}) {
	srvQuit = end
	cfg.Regist(define.EventCfgChange, "srvRestart", Restart)
	start(id)
}

func Restart() (err error) {
	addr := viper.GetString("http.addr")
	if addr == "" || (srv != nil && addr == srv.Addr) {
		return
	}

	if err = stop(); err != nil {
		log.Error("stop()", zap.Error(err))
	}

	go start(id)
	return
}

func Stop() {
	defer quit()
	util.Warn(stop())
}

func start(srvID int) (err error) {
	srv = &http.Server{
		Addr:    viper.GetString("http.addr"),
		Handler: router(),
	}

	err = srv.ListenAndServeTLS("certificate/server.crt", "certificate/server.key")
	if err != nil && srvID == id {
		log.Error("srv.ListenAndServe()", zap.Any("srv", srv), zap.Error(err))
		quit()
	}

	return
}

func stop() error {
	shutCtx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("http.shuttimeout"))
	defer cancel()

	id++
	return srv.Shutdown(shutCtx)
}

func quit() {
	if srvQuit != nil {
		close(srvQuit)
	}
}
