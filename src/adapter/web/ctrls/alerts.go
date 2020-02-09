package ctrls

import (
	"github.com/labstack/echo"
	"github.com/wakuwaku3/account-book.api/src/adapter/web/ctrls/responses"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/application/usecases"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	alerts struct {
		useCase usecases.Alerts
	}
	Alerts interface {
		GetAlerts(c echo.Context) error
		GetAlert(c echo.Context) error
		Create(c echo.Context) error
		Update(c echo.Context) error
		Delete(c echo.Context) error
	}
	getAlertsResponse struct {
		Alerts []getAlertResponse `json:"alerts"`
	}
	getAlertResponse struct {
		AlertID   string `json:"id"`
		Metrics   string `json:"metrics"`
		Threshold int    `json:"threshold"`
	}
	alertRequest struct {
		Metrics   string `json:"metrics"`
		Threshold int    `json:"threshold"`
	}
	createAlertResponse struct {
		AlertID string `json:"id"`
	}
)

// NewAlerts is create instance
func NewAlerts(useCase usecases.Alerts) Alerts {
	return &alerts{useCase}
}

func (t *alerts) GetAlerts(c echo.Context) error {
	res := t.useCase.GetAlerts()
	return responses.WriteResponse(c, getAlertsResponse{
		Alerts: convertAlerts(res.Alerts),
	})
}
func convertAlerts(alerts []usecases.GetAlertResult) []getAlertResponse {
	x := make([]getAlertResponse, len(alerts))
	for i, alert := range alerts {
		x[i] = convertAlert(alert)
	}
	return x
}
func convertAlert(t usecases.GetAlertResult) getAlertResponse {
	return getAlertResponse{
		AlertID:   t.AlertID,
		Metrics:   t.Metrics,
		Threshold: t.Threshold,
	}
}
func (t *alerts) GetAlert(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, core.NewError(application.RequiredID))
	}
	res, err := t.useCase.GetAlert(&id)
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, convertAlert(*res))
}
func (t *alerts) Create(c echo.Context) error {
	request := new(alertRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	res, err := t.useCase.Create(request.convert())
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, createAlertResponse{
		AlertID: res.AlertID,
	})
}
func (t *alertRequest) convert() *usecases.AlertArgs {
	return &usecases.AlertArgs{
		Metrics:   t.Metrics,
		Threshold: t.Threshold,
	}
}

func (t *alerts) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, core.NewError(application.RequiredID))
	}
	request := new(alertRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	if err := t.useCase.Update(&id, request.convert()); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
func (t *alerts) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, core.NewError(application.RequiredID))
	}
	if err := t.useCase.Delete(&id); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
