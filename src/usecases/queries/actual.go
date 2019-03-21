package queries

import (
	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"

	"github.com/wakuwaku3/account-book.api/src/usecases"
)

type actual struct {
	dashboardRepos domains.DashboardRepository
	plansRepos     domains.PlansRepository
}

// NewActual はインスタンスを生成します
func NewActual(
	dashboardRepos domains.DashboardRepository,
	plansRepos domains.PlansRepository,
) usecases.ActualQuery {
	return &actual{
		dashboardRepos,
		plansRepos,
	}
}
func (t *actual) Get(dashboardID *string, id *string) (
	*usecases.GetActualResult,
	error,
) {
	model, err := t.dashboardRepos.GetActual(dashboardID, id)
	if err != nil {
		return nil, err
	}
	plan, err := t.plansRepos.GetByID(&model.PlanID)
	if err != nil {
		return nil, err
	}
	return convertActual(model, plan), nil
}
func convertActual(t *models.Actual, p *models.Plan) *usecases.GetActualResult {
	return &usecases.GetActualResult{
		ActualAmount: t.ActualAmount,
		ActualID:     t.ActualID,
		PlanAmount:   p.PlanAmount,
		PlanID:       p.PlanID,
		PlanName:     p.PlanName,
	}
}
func (t *actual) GetActualInfo(planID *string) (*usecases.ActualInfo, error) {
	plan, err := t.plansRepos.GetByID(planID)
	if err != nil {
		return nil, err
	}
	return &usecases.ActualInfo{
		PlanID:     plan.PlanID,
		PlanName:   plan.PlanName,
		PlanAmount: plan.PlanAmount,
		IsIncome:   plan.IsIncome,
	}, nil
}
