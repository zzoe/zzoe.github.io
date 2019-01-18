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
	restarting bool
	srv     *http.Server
	srvQuit chan struct{}
	log     = cfg.Log
)

func Start(q chan struct{}) {
	srvQuit = q
	cfg.Regist(define.EventCfgChange, "srvRestart", Restart)
	util.Warn(start())
}

func Stop() {
	defer quit()
	util.Warn(stop())
}

func Restart() (err error) {
	addr := viper.GetString("http.addr")
	if addr == "" || (srv != nil && addr == srv.Addr) {
		return
	}

	if err = stop(); err != nil{
		log.Error("stop()", zap.Error(err))
		return
	}

	return start()
}

func start() (err error) {
	srv = &http.Server{
		Addr: viper.GetString("http.addr"),
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			if !restarting{
				quit()
				log.Error("srv.ListenAndServe()", zap.Any("srv", srv), zap.Error(err))
			}
		}
	}()

	return
}

func stop() (err error){
	restarting = true
	defer func() {
		restarting = false
	}()

	return srv.Shutdown(context.Background())
}

func quit(){
	if srvQuit != nil{
		close(srvQuit)
	}
}
