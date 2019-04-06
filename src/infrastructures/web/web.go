package web

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/di"
)

type web struct {
	echo *echo.Echo
	env  application.Env
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
	container.Invoke(func(env application.Env) {
		web.env = env
	})
	web.echo.Logger.SetLevel(log.INFO)
	corsConfig := middleware.DefaultCORSConfig
	corsConfig.AllowOrigins = *web.env.GetAllowOrigins()
	web.echo.Use(middleware.CORSWithConfig(corsConfig))
	web.echo.Use(middleware.Logger())
	web.echo.Use(middleware.Recover())
	web.echo.Use(DI(container))
	web.setRoute()
	return web, nil
}

func (web *web) Start() {
	web.echo.Logger.Fatal(web.echo.Start(":8080"))
}
