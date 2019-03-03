package web

import (
	"github.com/labstack/echo/middleware"
	entitieslog "github.com/wakuwaku3/account-book.api/src/domains/entities/log"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/log"
)

func (web *web) setLogger(lvl entitieslog.Lvl) *web {
	web.echo.Logger.SetLevel(log.CnvLvl(lvl))
	web.echo.Use(middleware.Logger())
	return web
}
