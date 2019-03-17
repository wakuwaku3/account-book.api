package ctrls

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/usecases"

	"github.com/wakuwaku3/account-book.api/src/ctrls/responses"

	"github.com/labstack/echo"
)

type (
	dashboard struct {
		useCase usecases.Dashboard
	}
	// Dashboard is DashboardController
	Dashboard interface {
		GetDashboard(c echo.Context) error
	}
	getDashboardResponse struct {
		SelectedMonth time.Time                   `json:"selectedMonth"`
		Summary       getDashboardSummaryResponse `json:"summary"`
		Plans         []getDashboardPlanResponse  `json:"plans"`
	}
	getDashboardSummaryResponse struct {
		Income          int  `json:"income"`
		Expense         int  `json:"expense"`
		PreviousBalance *int `json:"previousBalance,omitempty"`
	}
	getDashboardPlanResponse struct {
		PlanID       string  `json:"id"`
		PlanName     string  `json:"name"`
		IsIncome     bool    `json:"isIncome"`
		ActualAmount *int    `json:"actualAmount,omitempty"`
		PlanAmount   int     `json:"planAmount"`
		ActualID     *string `json:"actualId,omitempty"`
	}
)

// NewDashboard is create instance
func NewDashboard(useCase usecases.Dashboard) Dashboard {
	return &dashboard{useCase}
}

func (t *dashboard) GetDashboard(c echo.Context) error {
	var err error
	var selectedMonth *time.Time
	month := c.QueryParam("month")
	if month != "" {
		s, err := time.Parse("2006-01-02", month)
		if err != nil {
			return err
		}
		selectedMonth = &s
	}
	res, err := t.useCase.GetDashboard(&usecases.GetDashboardArgs{
		SelectedMonth: selectedMonth,
	})
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, convertDashboard(res))
}

func convertDashboard(t *usecases.GetDashboardResult) getDashboardResponse {
	plans := make([]getDashboardPlanResponse, len(t.Plans))
	for i, plan := range t.Plans {
		plans[i] = getDashboardPlanResponse{
			ActualAmount: plan.ActualAmount,
			ActualID:     plan.ActualID,
			IsIncome:     plan.IsIncome,
			PlanAmount:   plan.PlanAmount,
			PlanID:       plan.PlanID,
			PlanName:     plan.PlanName,
		}
	}
	return getDashboardResponse{
		SelectedMonth: t.SelectedMonth,
		Plans:         plans,
		Summary: getDashboardSummaryResponse{
			Expense:         t.Expense,
			Income:          t.Income,
			PreviousBalance: t.PreviousBalance,
		},
	}
}