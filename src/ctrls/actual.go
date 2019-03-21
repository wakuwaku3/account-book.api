package ctrls

import (
	"github.com/wakuwaku3/account-book.api/src/domains/apperrors"
	"github.com/wakuwaku3/account-book.api/src/usecases"

	"github.com/wakuwaku3/account-book.api/src/ctrls/responses"

	"github.com/labstack/echo"
)

type (
	actual struct {
		useCase usecases.Actual
	}
	// Actual is ActualController
	Actual interface {
		Get(c echo.Context) error
		Create(c echo.Context) error
		Update(c echo.Context) error
	}
	getActualResponse struct {
		ActualID     string `json:"actualId"`
		ActualAmount int    `json:"actualAmount"`
		PlanID       string `json:"planId"`
		PlanName     string `json:"planName"`
		PlanAmount   int    `json:"planAmount"`
	}
	actualRequest struct {
		ActualAmount int    `json:"actualAmount"`
		PlanID       string `json:"planId"`
		DashboardID  string `json:"dashboardId"`
	}
	createActualResponse struct {
		ActualID string `json:"id"`
	}
)

// NewActual is create instance
func NewActual(useCase usecases.Actual) Actual {
	return &actual{useCase}
}

func convertActual(t usecases.GetActualResult) getActualResponse {
	return getActualResponse{
		ActualAmount: t.ActualAmount,
		ActualID:     t.ActualID,
		PlanAmount:   t.PlanAmount,
		PlanID:       t.PlanID,
		PlanName:     t.PlanName,
	}
}
func (t *actual) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, apperrors.NewClientError(apperrors.RequiredID))
	}
	dashboardID := c.Param("dashboardID")
	if dashboardID == "" {
		return responses.WriteErrorResponse(c, apperrors.NewClientError(apperrors.RequiredID))
	}
	res, err := t.useCase.Get(&dashboardID, &id)
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, convertActual(*res))
}
func (t *actual) Create(c echo.Context) error {
	request := new(actualRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	res, err := t.useCase.Create(request.convert())
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, createActualResponse{
		ActualID: res.ActualID,
	})
}
func (t *actualRequest) convert() *usecases.ActualArgs {
	return &usecases.ActualArgs{
		ActualAmount: t.ActualAmount,
		DashboardID:  t.DashboardID,
		PlanID:       t.PlanID,
	}
}

func (t *actual) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, apperrors.NewClientError(apperrors.RequiredID))
	}
	request := new(actualRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	if err := t.useCase.Update(&id, request.convert()); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
