package usecases

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/application/services"
)

type (
	dashboard struct {
		query   DashboardQuery
		service services.Dashboard
	}
	// Dashboard is DashboardUseCases
	Dashboard interface {
		GetDashboard(args *GetDashboardArgs) (*GetDashboardResult, error)
		Approve(id *string) error
		CancelApprove(id *string) error
		AdjustBalance(args *AdjustBalanceArgs) error
	}
	// GetDashboardArgs は引数です
	GetDashboardArgs struct {
		SelectedMonth *time.Time
	}
	// GetDashboardResult は結果です
	GetDashboardResult struct {
		DashboardID      string
		SelectedMonth    time.Time
		Income           int
		Expense          int
		Balance          *int
		PreviousBalance  *int
		Plans            []PlanResult
		Daily            []DailyResult
		State            string
		CanApprove       bool
		CanCancelApprove bool
	}
	// PlanResult は結果です
	PlanResult struct {
		PlanID       string
		PlanName     string
		IsIncome     bool
		ActualAmount *int
		PlanAmount   int
		ActualID     *string
		CreatedAt    time.Time
	}
	// DailyResult は結果です
	DailyResult struct {
		Date    time.Time
		Income  int
		Expense int
		Balance int
	}
	// AdjustBalanceArgs は引数です
	AdjustBalanceArgs struct {
		DashboardID string
		Balance     int
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
func (t *dashboard) Approve(id *string) error {
	return t.service.Approve(id)
}
func (t *dashboard) CancelApprove(id *string) error {
	return t.service.CancelApprove(id)
}
func (t *dashboard) AdjustBalance(args *AdjustBalanceArgs) error {
	return t.service.AdjustBalance(&services.AdjustBalanceArgs{
		DashboardID: args.DashboardID,
		Balance:     args.Balance,
	})
}
