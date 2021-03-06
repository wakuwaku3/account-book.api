package services

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
	accountbook "github.com/wakuwaku3/account-book.api/src/enterprise/domains/accountBook"
	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
)

type (
	actual struct {
		dashboardRepos     application.DashboardRepository
		clock              core.Clock
		assetsChangedEvent accountbook.AssetsChangedEvent
	}
	// Actual is ActualService
	Actual interface {
		Enter(args *ActualArgs) error
	}
	// ActualArgs は引数です
	ActualArgs struct {
		models.ActualKey
		ActualAmount  int
		PlanName      string
		PlanAmount    int
		IsIncome      bool
		PlanCreatedAt time.Time
	}
	// CreateActualResult は結果です
	CreateActualResult struct {
		ActualID string
	}
)

// NewActual is create instance
func NewActual(
	dashboardRepos application.DashboardRepository,
	clock core.Clock,
	assetsChangedEvent accountbook.AssetsChangedEvent,
) Actual {
	return &actual{dashboardRepos, clock, assetsChangedEvent}
}
func (t *actual) Enter(args *ActualArgs) error {
	if args.DashboardID == nil {
		dashboardID, err := t.dashboardRepos.Create(args.SelectedMonth)
		if err != nil {
			return err
		}
		args.DashboardID = dashboardID
	}

	if args.ActualID == nil {
		id, err := t.dashboardRepos.ExistsActual(args.DashboardID, &args.PlanID)
		if err != nil {
			return err
		}
		if id == nil {
			_, err = t.dashboardRepos.CreateActual(args.convert())
			if err != nil {
				return err
			}
			t.assetsChangedEvent.Trigger()
			return nil
		}
	}

	model, err := t.dashboardRepos.GetActual(args.DashboardID, args.ActualID)
	if err != nil {
		return err
	}

	model.IsIncome = args.IsIncome
	model.ActualAmount = args.ActualAmount
	model.PlanAmount = args.PlanAmount
	model.PlanID = args.PlanID
	model.PlanName = args.PlanName

	if err := t.dashboardRepos.UpdateActual(args.DashboardID, args.ActualID, model); err != nil {
		return err
	}
	t.assetsChangedEvent.Trigger()
	return nil
}
func (t *ActualArgs) convert() (*string, *models.Actual) {
	return t.DashboardID, &models.Actual{
		ActualAmount:  t.ActualAmount,
		IsIncome:      t.IsIncome,
		PlanAmount:    t.PlanAmount,
		PlanID:        t.PlanID,
		PlanName:      t.PlanName,
		PlanCreatedAt: t.PlanCreatedAt,
	}
}
