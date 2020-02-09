package ctrls

import (
	"github.com/labstack/echo"
	"github.com/wakuwaku3/account-book.api/src/adapter/web/ctrls/responses"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/application/usecases"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	notificationRules struct {
		useCase usecases.NotificationRules
	}
	NotificationRules interface {
		GetNotificationRules(c echo.Context) error
		GetNotificationRule(c echo.Context) error
		Create(c echo.Context) error
		Update(c echo.Context) error
		Delete(c echo.Context) error
	}
	getNotificationRulesResponse struct {
		NotificationRules []getNotificationRuleResponse `json:"notificationRules"`
	}
	getNotificationRuleResponse struct {
		NotificationRuleID string `json:"id"`
		Metrics            string `json:"metrics"`
		Threshold          int    `json:"threshold"`
	}
	notificationRuleRequest struct {
		Metrics   string `json:"metrics"`
		Threshold int    `json:"threshold"`
	}
	createNotificationRuleResponse struct {
		NotificationRuleID string `json:"id"`
	}
)

// NewNotificationRules is create instance
func NewNotificationRules(useCase usecases.NotificationRules) NotificationRules {
	return &notificationRules{useCase}
}

func (t *notificationRules) GetNotificationRules(c echo.Context) error {
	res := t.useCase.GetNotificationRules()
	return responses.WriteResponse(c, getNotificationRulesResponse{
		NotificationRules: convertNotificationRules(res.NotificationRules),
	})
}
func convertNotificationRules(notificationRules []usecases.GetNotificationRuleResult) []getNotificationRuleResponse {
	x := make([]getNotificationRuleResponse, len(notificationRules))
	for i, notificationRule := range notificationRules {
		x[i] = convertNotificationRule(notificationRule)
	}
	return x
}
func convertNotificationRule(t usecases.GetNotificationRuleResult) getNotificationRuleResponse {
	return getNotificationRuleResponse{
		NotificationRuleID: t.NotificationRuleID,
		Metrics:            t.Metrics,
		Threshold:          t.Threshold,
	}
}
func (t *notificationRules) GetNotificationRule(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, core.NewError(application.RequiredID))
	}
	res, err := t.useCase.GetNotificationRule(&id)
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, convertNotificationRule(*res))
}
func (t *notificationRules) Create(c echo.Context) error {
	request := new(notificationRuleRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	res, err := t.useCase.Create(request.convert())
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, createNotificationRuleResponse{
		NotificationRuleID: res.NotificationRuleID,
	})
}
func (t *notificationRuleRequest) convert() *usecases.NotificationRuleArgs {
	return &usecases.NotificationRuleArgs{
		Metrics:   t.Metrics,
		Threshold: t.Threshold,
	}
}

func (t *notificationRules) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, core.NewError(application.RequiredID))
	}
	request := new(notificationRuleRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	if err := t.useCase.Update(&id, request.convert()); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
func (t *notificationRules) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, core.NewError(application.RequiredID))
	}
	if err := t.useCase.Delete(&id); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
