package bootstrap

import (
	"gin-app/internal/infra/cache"
	"gin-app/pkg/serror"
	"gin-app/pkg/slog"

	gormgenerics "github.com/olongfen/gorm-generics"
	"go.uber.org/zap"
)

type Application struct {
	Conf     *Conf
	Log      *zap.Logger
	Database gormgenerics.Database
	Rdb      cache.Cache
}

var GlobalLog *zap.Logger

func App(confPath string) (*Application, error) {
	logger := slog.NewProduceLogger()
	GlobalLog = logger
	// 初始化多语言错误
	if err := serror.InitI18n(); err != nil {
		return nil, err
	}
	conf, err := NewConf(confPath)
	if err != nil {
		return nil, err
	}
	database, err := NewDatabase(conf, logger)
	if err != nil {
		return nil, err
	}
	rdb, err := cache.NewRDB(cache.Config{})
	if err != nil {
		return nil, err
	}
	app := &Application{Database: database, Conf: conf, Log: logger, Rdb: rdb}
	return app, nil
}

func (a *Application) Close() error {
	if a == nil {
		return nil
	}
	a.Database.Close()
	return nil
}
