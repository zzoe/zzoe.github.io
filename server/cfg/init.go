package cfg

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/zzoe/zoe.github.io/server/define"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	Log     *zap.Logger
	BaseDir = filepath.Dir(os.Args[0])

	debug    = flag.Bool("debug", false, "debug or not")
	eventMap = map[define.Event]*sync.Map{
		define.EventInit:          new(sync.Map),
		define.EventCfgChange:     new(sync.Map),
		define.EventResourceClear: new(sync.Map),
	}
)

func init() {
	flag.Parse()
	initLog()
	initConfig()
	Log.Debug("cfg init")
}

func initLog() {
	var enc zapcore.Encoder
	var enab zapcore.LevelEnabler
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/assistant.log",
		MaxSize:    8,  // megabytes
		MaxAge:     32, // days
		MaxBackups: 32,
		LocalTime:  true,
		Compress:   true,
	})

	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02_15:04:05.000"))
	}
	if *debug {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeTime = timeEncoder
		enc = zapcore.NewConsoleEncoder(encoderCfg)
		enab = zapcore.DebugLevel
		w = zapcore.NewMultiWriteSyncer(w, os.Stdout)
	} else {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = timeEncoder
		enc = zapcore.NewJSONEncoder(encoderCfg)
		enab = zapcore.InfoLevel
	}
	Log = zap.New(zapcore.NewCore(enc, w, enab))
}

func initConfig() {
	var err error
	cfgFile := filepath.Join(BaseDir, "config.toml")
	if _, err = os.Stat(cfgFile); err != nil {
		Log.Panic("配置文件状态异常", zap.Error(err), zap.String("cfgFile", cfgFile))
	}

	viper.SetConfigFile(cfgFile)
	if err = viper.ReadInConfig(); err != nil {
		Log.Panic("读取配置文件失败", zap.Error(err))
	}

	//Watching and re-reading config files
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		Log.Info("Config file changed:", zap.Any("e.Op", e.Op))
		errs := Process(define.EventCfgChange)
		for i := range errs {
			Log.Error("配置文件修改回调报错", zap.String("key", i), zap.Error(errs[i]))
		}
	})
}
