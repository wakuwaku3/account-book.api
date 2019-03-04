package web

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/di"
)

type web struct {
	echo *echo.Echo
	env  domains.Env
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
	web := &web{echo, nil}
	container.Invoke(func(env domains.Env) {
		web.env = env
	})
	web.echo.Use(DI(container))
	web.echo.Logger.SetLevel(log.INFO)
	web.echo.Use(middleware.Logger())
	web.echo.Use(middleware.Recover())
	web.setRoute()
	return web, nil
}

func (web *web) Start() {
	web.echo.Logger.Fatal(web.echo.Start(":8080"))
}
