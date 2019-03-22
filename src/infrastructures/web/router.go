package web

import (
	"net/http"

	"github.com/labstack/echo/middleware"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/auth"

	"github.com/labstack/echo"
	ctrl "github.com/wakuwaku3/account-book.api/src/ctrls"
)

func (web *web) setRoute() *web {
	web.echo.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello")
	})
	web.echo.GET("/_ah/warmup", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	// accounts
	// sign-in
	web.echo.POST("/accounts/sign-in", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Accounts) error {
			return controller.SignIn(c)
		})
	})
	// refresh
	web.echo.POST("/accounts/refresh", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Accounts) error {
			return controller.Refresh(c)
		})
	})
	// password-reset-requesting
	web.echo.PUT("/accounts/password-reset-requesting", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Accounts) error {
			return controller.PasswordResetRequesting(c)
		})
	})
	// reset-password
	web.echo.GET("/accounts/reset-password", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Accounts) error {
			return controller.GetResetPasswordModel(c)
		})
	})
	web.echo.POST("/accounts/reset-password", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Accounts) error {
			return controller.ResetPassword(c)
		})
	})

	// auth aria
	jwtSecret := web.env.GetJwtSecret()
	auth := web.echo.Group("", middleware.JWT(*jwtSecret), auth.Authenticate())

	// transactions
	// GET
	auth.GET("/transactions", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Transactions) error {
			return controller.GetTransactions(c)
		})
	})
	// GET
	auth.GET("/transactions/:id", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Transactions) error {
			return controller.GetTransaction(c)
		})
	})
	// POST
	auth.POST("/transactions", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Transactions) error {
			return controller.Create(c)
		})
	})
	// PUT
	auth.PUT("/transactions/:id", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Transactions) error {
			return controller.Update(c)
		})
	})
	// DELETE
	auth.DELETE("/transactions/:id", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Transactions) error {
			return controller.Delete(c)
		})
	})

	// plans
	// GET
	auth.GET("/plans", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Plans) error {
			return controller.GetPlans(c)
		})
	})
	// GET
	auth.GET("/plans/:id", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Plans) error {
			return controller.GetPlan(c)
		})
	})
	// POST
	auth.POST("/plans", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Plans) error {
			return controller.Create(c)
		})
	})
	// PUT
	auth.PUT("/plans/:id", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Plans) error {
			return controller.Update(c)
		})
	})
	// DELETE
	auth.DELETE("/plans/:id", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Plans) error {
			return controller.Delete(c)
		})
	})

	// Dashboard
	// GET
	auth.GET("/dashboard", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Dashboard) error {
			return controller.GetDashboard(c)
		})
	})
	// Approve
	auth.POST("/dashboard/:id", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Dashboard) error {
			return controller.Approve(c)
		})
	})
	// CancelApprove
	auth.DELETE("/dashboard/:id", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Dashboard) error {
			return controller.CancelApprove(c)
		})
	})

	// Actual
	// GET
	auth.GET("/actual", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Actual) error {
			return controller.Get(c)
		})
	})
	// PUT
	auth.PUT("/actual", func(c echo.Context) error {
		container := GetContainer(c)
		return container.Invoke(func(controller ctrl.Actual) error {
			return controller.Put(c)
		})
	})

	return web
}
