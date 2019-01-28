package cfg

import (
	"github.com/zzoe/zoe.github.io/server/define"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Clear() {
	MustProcess(define.EventResourceClear)
}

func Regist(e define.Event, key string, fn func() error) {
	em, ok := eventMap[e]
	if ok {
		em[key] = fn
	}
}

func Process(e define.Event) {
	em, ok := eventMap[e]
	if ok {
		for k, fn := range em {
			go func(key string, fn func() error) {
				if err := fn(); err != nil {
					Log.Error("Process event fail", zap.Any("event", e), zap.String("key", key), zap.Error(err))
				}
			}(k, fn)
		}
	}
}

func MustProcess(e define.Event) {
	em, ok := eventMap[e]
	if ok{
		var g errgroup.Group
		for _,fn := range em{
			g.Go(fn)
		}

		if err := g.Wait(); err != nil{
			panic(err)
		}
	}
}
