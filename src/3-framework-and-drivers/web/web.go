package web

import (
	"github.com/labstack/echo"
	entitieslog "github.com/wakuwaku3/account-book.api/src/0-enterprise-business-rules/entities/log"
	"github.com/wakuwaku3/account-book.api/src/3-framework-and-drivers/di"
	"go.uber.org/dig"
)

type web struct {
	echo      *echo.Echo
	container *dig.Container
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
	return web.setLogger(entitieslog.Info).setRoute(), nil
}

func (web *web) Start() {
	web.echo.Logger.Fatal(web.echo.Start(":8080"))
}
