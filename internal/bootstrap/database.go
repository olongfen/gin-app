package bootstrap

import (
	"context"
	"errors"
	"gin-app/pkg/scontext"
	"gin-app/pkg/serror"
	"gin-app/pkg/sslog"
	"log/slog"

	"gin-app/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	gormgenerics "github.com/olongfen/gorm-generics"
	"gorm.io/gorm"
)

// NewDatabase 新建数据库
func NewDatabase(conf *Conf) (gormgenerics.Database, error) {
	var logger *slog.Logger
	if conf.LogConf.IsProd {
		c := conf.LogConf
		c.LogPath = "log/db.log"
		l := slog.NewJSONHandler(NewLumberjack(c), nil)
		logger = slog.New(l)
	} else {
		logger = slog.Default()
	}
	dataBase, err := gormgenerics.NewDataBase(conf.DBDriver, conf.DBDsn,
		gormgenerics.WithAutoMigrate(conf.DBAutoMigrate),
		gormgenerics.WithAutoMigrateDst([]any{&domain.User{}}),
		gormgenerics.WithLogger(sslog.NewDBLog(logger)),
		gormgenerics.WithOpentracingPlugin(&gormgenerics.OpentracingPlugin{}),
		gormgenerics.WithTranslateError(translateErr),
	)
	if err != nil {
		return nil, err
	}
	return dataBase, nil
}

func translateErr(ctx context.Context, db *gorm.DB) (err error) {
	err = db.Error
	lan := scontext.GetLanguage(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		switch db.Statement.Table {
		case "users":
			return serror.Error(serror.ErrUserRecordNotFound, lan)
		default:
			return serror.Error(serror.ErrRecordNotFound, lan)
		}
	}

	var errV *pgconn.PgError
	ok := errors.As(err, &errV)
	if !ok {
		return
	}
	switch errV.Code {
	case "23505":
		switch errV.TableName {
		case "users":
			return domain.TranslateUserDBErr(errV, lan)
		}
	}

	return err
}
