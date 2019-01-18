package cfg

import (
	"github.com/zzoe/zoe.github.io/server/define"
	"go.uber.org/zap"
	"sync"
)

func Clear() {
	ProcessMust(define.EventResourceClear)
}

func Regist(e define.Event, key string, fn func() error) {
	em, ok := eventMap[e]
	if ok {
		em.Store(key, fn)
	}
}

func Process(e define.Event) map[string]error {
	var lock sync.Mutex
	errs := make(map[string]error, 0)

	em, ok := eventMap[e]
	if ok {
		var wg sync.WaitGroup
		em.Range(func(key, value interface{}) bool {
			wg.Add(1)

			go func(k string) {
				defer wg.Done()

				fn, _ := em.Load(k)
				err := fn.(func() error)()

				lock.Lock()
				errs[k] = err
				lock.Unlock()

			}(key.(string))

			return true
		})
		wg.Wait()

		for key, err := range errs {
			if err != nil {
				Log.Error("Process event fail", zap.String("key", key), zap.Error(err))
			}
		}
	}
	return errs
}

func ProcessMust(e define.Event) {
	errs := Process(e)
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
