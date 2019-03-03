package web

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tampopos/dijct"
	entitieslog "github.com/wakuwaku3/account-book.api/src/domains/entities/log"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/di"
)

type web struct {
	echo      *echo.Echo
	container dijct.Container
}

// Web はWebサーバーのインターフェイスです
type Web interface {
	Start()
}

// NewWeb は Web を生成します
func NewWeb() (Web, error) {
	echo := echo.New()
	container, err := di.CreateContainer()
	if err != nil {
		return nil, err
	}
	web := &web{echo, container}
	web.setLogger(entitieslog.Info)
	web.echo.Use(middleware.Recover())
	web.setRoute()
	return web, nil

}

func (web *web) Start() {
	web.echo.Logger.Fatal(web.echo.Start(":8080"))
}
