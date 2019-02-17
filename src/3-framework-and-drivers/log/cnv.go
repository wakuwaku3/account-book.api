package log

import (
	"github.com/labstack/gommon/log"
	entitieslog "github.com/wakuwaku3/account-book.api/src/0-enterprise-business-rules/entities/log"
)

// CnvLvl はログレベルを変換します
func CnvLvl(lvl entitieslog.Lvl) log.Lvl {
	switch lvl {
	case entitieslog.Fatal:
	case entitieslog.Error:
		return log.ERROR
	case entitieslog.Warn:
		return log.WARN
	case entitieslog.Info:
		return log.INFO
	}
	return log.DEBUG
}
