package usecases

import (
	"github.com/wakuwaku3/account-book.api/src/domains/services"
)

type (
	actual struct {
		query   ActualQuery
		service services.Actual
	}
	// Actual is ActualUseCases
	Actual interface {
		Get(dashboardID *string, id *string) (*GetActualResult, error)
		Create(args *ActualArgs) (*CreateActualResult, error)
		Update(id *string, args *ActualArgs) error
	}
	// GetActualResult は結果です
	GetActualResult struct {
		ActualID     string
		ActualAmount int
		PlanID       string
		PlanName     string
		PlanAmount   int
	}
	// ActualInfo は実績登録のための情報です
	ActualInfo struct {
		PlanID     string
		PlanName   string
		PlanAmount int
		IsIncome   bool
	}
	// ActualArgs は引数です
	ActualArgs struct {
		ActualAmount int
		PlanID       string
		DashboardID  string
	}
	// CreateActualResult は結果です
	CreateActualResult struct {
		ActualID string
	}
)

// NewActual is create instance
func NewActual(
	query ActualQuery,
	service services.Actual,
) Actual {
	return &actual{
		query,
		service,
	}
}
func (t *actual) Get(dashboardID *string, id *string) (*GetActualResult, error) {
	return t.query.Get(dashboardID, id)
}
func (t *actual) Create(args *ActualArgs) (*CreateActualResult, error) {
	info, err := t.query.GetActualInfo(&args.PlanID)
	if err != nil {
		return nil, err
	}
	res, err := t.service.Create(args.convert(info))
	if err != nil {
		return nil, err
	}
	return &CreateActualResult{
		ActualID: res.ActualID,
	}, nil
}
func (t *ActualArgs) convert(info *ActualInfo) *services.ActualArgs {
	return &services.ActualArgs{
		ActualAmount: t.ActualAmount,
		DashboardID:  t.DashboardID,
		IsIncome:     info.IsIncome,
		PlanAmount:   info.PlanAmount,
		PlanID:       info.PlanID,
		PlanName:     info.PlanName,
	}
}
func (t *actual) Update(id *string, args *ActualArgs) error {
	info, err := t.query.GetActualInfo(&args.PlanID)
	if err != nil {
		return err
	}
	return t.service.Update(id, args.convert(info))
}
