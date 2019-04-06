package queries

import (
	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"

	"github.com/wakuwaku3/account-book.api/src/application/usecases"
)

type plans struct {
	repos domains.PlansRepository
}

// NewPlans はインスタンスを生成します
func NewPlans(
	repos domains.PlansRepository,
) usecases.PlansQuery {
	return &plans{
		repos,
	}
}
func (t *plans) GetPlans() (*usecases.GetPlansResult, error) {
	records, err := t.repos.Get()
	if err != nil {
		return nil, err
	}
	plans := make([]usecases.GetPlanResult, len(*records))
	for i, record := range *records {
		r := &record
		plans[i] = *convertPlan(r)
	}
	return &usecases.GetPlansResult{Plans: plans}, nil
}
func (t *plans) GetPlan(id *string) (
	*usecases.GetPlanResult,
	error,
) {
	model, err := t.repos.GetByID(id)
	if err != nil {
		return nil, err
	}
	return convertPlan(model), nil
}
func convertPlan(t *models.Plan) *usecases.GetPlanResult {
	return &usecases.GetPlanResult{
		PlanID:     t.PlanID,
		PlanName:   t.PlanName,
		IsIncome:   t.IsIncome,
		PlanAmount: t.PlanAmount,
		Interval:   t.Interval,
		Start:      t.Start,
		End:        t.End,
	}
}
