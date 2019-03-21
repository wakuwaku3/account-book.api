package services

import (
	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
)

type (
	actual struct {
		dashboardRepos domains.DashboardRepository
		clock          cmn.Clock
	}
	// Actual is ActualService
	Actual interface {
		Create(args *ActualArgs) (*CreateActualResult, error)
		Update(id *string, args *ActualArgs) error
	}
	// ActualArgs は引数です
	ActualArgs struct {
		DashboardID  string
		ActualAmount int
		PlanID       string
		PlanName     string
		PlanAmount   int
		IsIncome     bool
	}
	// CreateActualResult は結果です
	CreateActualResult struct {
		ActualID string
	}
)

// NewActual is create instance
func NewActual(dashboardRepos domains.DashboardRepository, clock cmn.Clock) Actual {
	return &actual{dashboardRepos, clock}
}
func (t *actual) Create(args *ActualArgs) (*CreateActualResult, error) {
	id, err := t.dashboardRepos.ExistsActual(&args.DashboardID, &args.PlanID)
	if err != nil {
		return nil, err
	}
	if id != nil {
		err := t.Update(id, args)
		if err != nil {
			return nil, err
		}
		return &CreateActualResult{ActualID: *id}, nil
	}
	id, err = t.dashboardRepos.CreateActual(args.convert())
	if err != nil {
		return nil, err
	}
	return &CreateActualResult{ActualID: *id}, nil
}
func (t *ActualArgs) convert() (*string, *models.Actual) {
	return &t.DashboardID, &models.Actual{
		ActualAmount: t.ActualAmount,
		IsIncome:     t.IsIncome,
		PlanAmount:   t.PlanAmount,
		PlanID:       t.PlanID,
		PlanName:     t.PlanName,
	}
}
func (t *actual) Update(id *string, args *ActualArgs) error {
	model, err := t.dashboardRepos.GetActual(&args.DashboardID, id)
	if err != nil {
		return err
	}

	model.IsIncome = args.IsIncome
	model.ActualAmount = args.ActualAmount
	model.PlanAmount = args.PlanAmount
	model.PlanID = args.PlanID
	model.PlanName = args.PlanName

	if err := t.dashboardRepos.UpdateActual(&args.DashboardID, id, model); err != nil {
		return err
	}
	return nil
}
