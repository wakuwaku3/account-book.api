package usecases

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains/services"
)

type (
	dashboard struct {
		query   DashboardQuery
		service services.Dashboard
	}
	// Dashboard is DashboardUseCases
	Dashboard interface {
		GetDashboard(args *GetDashboardArgs) (*GetDashboardResult, error)
	}
	// GetDashboardArgs は引数です
	GetDashboardArgs struct {
		SelectedMonth *time.Time
	}
	// GetDashboardResult は結果です
	GetDashboardResult struct {
		SelectedMonth   time.Time
		Income          int
		Expense         int
		PreviousBalance *int
		Plans           []PlanResult
	}
	// PlanResult は結果です
	PlanResult struct {
		PlanID       string
		PlanName     string
		IsIncome     bool
		ActualAmount *int
		PlanAmount   int
		ActualID     *string
	}
)

// NewDashboard is create instance
func NewDashboard(
	query DashboardQuery,
	service services.Dashboard,
) Dashboard {
	return &dashboard{
		query,
		service,
	}
}
func (t *dashboard) GetDashboard(args *GetDashboardArgs) (*GetDashboardResult, error) {
	info, err := t.query.GetSummary(args)
	if err != nil {
		return nil, err
	}
	return info, nil
}
