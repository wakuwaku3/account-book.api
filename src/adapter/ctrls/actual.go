package ctrls

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/usecases"

	"github.com/wakuwaku3/account-book.api/src/adapter/ctrls/responses"

	"github.com/labstack/echo"
)

type (
	actual struct {
		useCase usecases.Actual
	}
	// Actual is ActualController
	Actual interface {
		Get(c echo.Context) error
		Put(c echo.Context) error
	}
	getActualRequest struct {
		PlanID        string     `json:"planId"`
		ActualID      *string    `json:"actualId"`
		DashboardID   *string    `json:"dashboardId"`
		SelectedMonth *time.Time `json:"selectedMonth"`
	}
	getActualResponse struct {
		PlanName     string `json:"planName"`
		PlanAmount   int    `json:"planAmount"`
		ActualAmount *int   `json:"actualAmount"`
	}
	putActualRequest struct {
		getActualRequest
		ActualAmount int `json:"actualAmount"`
	}
)

// NewActual is create instance
func NewActual(useCase usecases.Actual) Actual {
	return &actual{useCase}
}

func (t *actual) Get(c echo.Context) error {
	request := new(getActualRequest)
	request.PlanID = c.QueryParam("planId")
	if aid := c.QueryParam("actualId"); aid != "" {
		request.ActualID = &aid
	}
	if did := c.QueryParam("dashboardId"); did != "" {
		request.DashboardID = &did
	}
	if month := c.QueryParam("month"); month != "" {
		s, err := time.Parse("2006-01-02", month)
		if err != nil {
			return err
		}
		request.SelectedMonth = &s
	}
	res, err := t.useCase.Get(&usecases.GetActualArgs{ActualKey: request.convert()})
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, convertActual(*res))
}
func (t *getActualRequest) convert() domains.ActualKey {
	return domains.ActualKey{
		ActualID:      t.ActualID,
		PlanID:        t.PlanID,
		DashboardID:   t.DashboardID,
		SelectedMonth: t.SelectedMonth,
	}
}
func convertActual(t usecases.GetActualResult) getActualResponse {
	return getActualResponse{
		ActualAmount: t.ActualAmount,
		PlanAmount:   t.PlanAmount,
		PlanName:     t.PlanName,
	}
}

func (t *actual) Put(c echo.Context) error {
	request := new(putActualRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	if err := t.useCase.Enter(request.convert()); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
func (t *putActualRequest) convert() *usecases.EnterActualArgs {
	return &usecases.EnterActualArgs{
		ActualKey:    t.getActualRequest.convert(),
		ActualAmount: t.ActualAmount,
	}
}
