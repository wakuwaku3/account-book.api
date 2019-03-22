package usecases

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains/models"
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
		Approve(id *string) error
		CancelApprove(id *string) error
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
		PreviousBalance  *int
		Plans            []PlanResult
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
	// ApproveInfo は承認処理に必要な情報です
	ApproveInfo struct {
		Income              int
		Expense             int
		PreviousBalance     int
		CurrentBalance      int
		Balance             int
		PreviousDashboardID *string
		Daily               *[]*models.Daily
	}
	// CancelApproveInfo は承認処理に必要な情報です
	CancelApproveInfo struct {
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
