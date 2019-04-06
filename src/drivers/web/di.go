package web

import (
	"github.com/labstack/echo"
	"github.com/tampopos/dijct"
)

const diContainerKey = "di-container"

// DI は DIコンテナのミドルウェアです
func DI(container dijct.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			childContainer := container.CreateChildContainer()
			c.Set(diContainerKey, childContainer)
			return next(c)
		}
	}
}

// GetContainer はコンテナを取得します
func GetContainer(c echo.Context) dijct.Container {
	return c.Get(diContainerKey).(dijct.Container)
}
