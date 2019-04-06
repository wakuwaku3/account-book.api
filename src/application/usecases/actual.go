package usecases

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/application/services"
)

type (
	actual struct {
		query   ActualQuery
		service services.Actual
	}
	// Actual is ActualUseCases
	Actual interface {
		Get(args *GetActualArgs) (*GetActualResult, error)
		Enter(args *EnterActualArgs) error
	}
	// ActualInfo は実績登録のための情報です
	ActualInfo struct {
		PlanID        string
		PlanName      string
		PlanAmount    int
		IsIncome      bool
		PlanCreatedAt time.Time
	}
	// GetActualArgs は引数です
	GetActualArgs struct {
		application.ActualKey
	}
	// GetActualResult は結果です
	GetActualResult struct {
		PlanName     string
		PlanAmount   int
		ActualAmount *int
	}
	// EnterActualArgs は引数です
	EnterActualArgs struct {
		application.ActualKey
		ActualAmount int
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
func (t *actual) Get(args *GetActualArgs) (*GetActualResult, error) {
	return t.query.Get(args)
}
func (t *actual) Enter(args *EnterActualArgs) error {
	info, err := t.query.GetActualInfo(&args.ActualKey)
	if err != nil {
		return err
	}
	return t.service.Enter(args.convert(info))
}
func (t *EnterActualArgs) convert(info *ActualInfo) *services.ActualArgs {
	return &services.ActualArgs{
		ActualAmount:  t.ActualAmount,
		IsIncome:      info.IsIncome,
		PlanAmount:    info.PlanAmount,
		PlanName:      info.PlanName,
		PlanCreatedAt: info.PlanCreatedAt,
		ActualKey: application.ActualKey{
			DashboardID:   t.DashboardID,
			PlanID:        t.PlanID,
			ActualID:      t.ActualID,
			SelectedMonth: t.SelectedMonth,
		},
	}
}
