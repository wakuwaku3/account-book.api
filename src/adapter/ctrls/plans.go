package ctrls

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains/apperrors"
	"github.com/wakuwaku3/account-book.api/src/application/usecases"

	"github.com/wakuwaku3/account-book.api/src/adapter/ctrls/responses"

	"github.com/labstack/echo"
)

type (
	plans struct {
		useCase usecases.Plans
	}
	// Plans is PlansController
	Plans interface {
		GetPlans(c echo.Context) error
		GetPlan(c echo.Context) error
		Create(c echo.Context) error
		Update(c echo.Context) error
		Delete(c echo.Context) error
	}
	getPlansResponse struct {
		Plans []getPlanResponse `json:"plans"`
	}
	getPlanResponse struct {
		PlanID     string     `json:"id"`
		PlanName   string     `json:"name"`
		IsIncome   bool       `json:"isIncome"`
		PlanAmount int        `json:"amount"`
		Interval   int        `json:"interval"`
		Start      *time.Time `json:"start"`
		End        *time.Time `json:"end"`
	}
	planRequest struct {
		PlanName   string     `json:"name"`
		IsIncome   bool       `json:"isIncome"`
		PlanAmount int        `json:"amount"`
		Interval   int        `json:"interval"`
		Start      *time.Time `json:"start"`
		End        *time.Time `json:"end"`
	}
	createPlanResponse struct {
		PlanID string `json:"id"`
	}
)

// NewPlans is create instance
func NewPlans(useCase usecases.Plans) Plans {
	return &plans{useCase}
}

func (t *plans) GetPlans(c echo.Context) error {
	res, err := t.useCase.GetPlans()
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, getPlansResponse{
		Plans: convertPlans(res.Plans),
	})
}
func convertPlans(plans []usecases.GetPlanResult) []getPlanResponse {
	x := make([]getPlanResponse, len(plans))
	for i, plan := range plans {
		x[i] = convertPlan(plan)
	}
	return x
}
func convertPlan(t usecases.GetPlanResult) getPlanResponse {
	return getPlanResponse{
		PlanID:     t.PlanID,
		PlanName:   t.PlanName,
		IsIncome:   t.IsIncome,
		PlanAmount: t.PlanAmount,
		Interval:   t.Interval,
		Start:      t.Start,
		End:        t.End,
	}
}
func (t *plans) GetPlan(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, apperrors.NewClientError(apperrors.RequiredID))
	}
	res, err := t.useCase.GetPlan(&id)
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, convertPlan(*res))
}
func (t *plans) Create(c echo.Context) error {
	request := new(planRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	res, err := t.useCase.Create(request.convert())
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, createPlanResponse{
		PlanID: res.PlanID,
	})
}
func (t *planRequest) convert() *usecases.PlanArgs {
	return &usecases.PlanArgs{
		PlanName:   t.PlanName,
		IsIncome:   t.IsIncome,
		PlanAmount: t.PlanAmount,
		Interval:   t.Interval,
		Start:      t.Start,
		End:        t.End,
	}
}

func (t *plans) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, apperrors.NewClientError(apperrors.RequiredID))
	}
	request := new(planRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	if err := t.useCase.Update(&id, request.convert()); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
func (t *plans) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responses.WriteErrorResponse(c, apperrors.NewClientError(apperrors.RequiredID))
	}
	if err := t.useCase.Remove(&id); err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
