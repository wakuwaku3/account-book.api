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
func (t *actual) Get(args *usecases.GetActualArgs) (
	*usecases.GetActualResult,
	error,
) {
	plan, err := t.plansRepos.GetByID(&args.PlanID)
	if err != nil {
		return nil, err
	}

	result := &usecases.GetActualResult{
		PlanAmount: plan.PlanAmount,
		PlanName:   plan.PlanName,
	}
	if args.DashboardID != nil && args.ActualID != nil {
		actual, err := t.dashboardRepos.GetActual(args.DashboardID, args.ActualID)
		if err != nil {
			return nil, err
		}
		result.ActualAmount = &actual.ActualAmount
	}
	return result, nil
}
func convertActual(t *models.Actual, p *models.Plan) *usecases.GetActualResult {
	return &usecases.GetActualResult{
		ActualAmount: &t.ActualAmount,
		PlanAmount:   p.PlanAmount,
		PlanName:     p.PlanName,
	}
}
func (t *actual) GetActualInfo(key *domains.ActualKey) (*usecases.ActualInfo, error) {
	plan, err := t.plansRepos.GetByID(&key.PlanID)
	if err != nil {
		return nil, err
	}
	return &usecases.ActualInfo{
		PlanID:        plan.PlanID,
		PlanName:      plan.PlanName,
		PlanAmount:    plan.PlanAmount,
		IsIncome:      plan.IsIncome,
		PlanCreatedAt: plan.CreatedAt,
	}, nil
}