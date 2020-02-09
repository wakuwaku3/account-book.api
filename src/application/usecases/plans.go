package usecases

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/application/services"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	plans struct {
		query   PlansQuery
		service services.Plans
	}
	// Plans is PlansUseCases
	Plans interface {
		GetPlans() (*GetPlansResult, error)
		GetPlan(id *string) (*GetPlanResult, error)
		Create(args *PlanArgs) (*CreatePlanResult, error)
		Update(id *string, args *PlanArgs) error
		Remove(id *string) error
	}
	// GetPlansResult は結果です
	GetPlansResult struct {
		Plans []GetPlanResult
	}
	// GetPlanResult は結果です
	GetPlanResult struct {
		PlanID     string
		PlanName   string
		IsIncome   bool
		PlanAmount int
		Interval   int
		Start      *time.Time
		End        *time.Time
	}
	// PlanArgs は引数です
	PlanArgs struct {
		PlanName   string
		IsIncome   bool
		PlanAmount int
		Interval   int
		Start      *time.Time
		End        *time.Time
	}
	// CreatePlanResult は結果です
	CreatePlanResult struct {
		PlanID string
	}
)

// NewPlans is create instance
func NewPlans(
	query PlansQuery,
	service services.Plans,
) Plans {
	return &plans{
		query,
		service,
	}
}
func (t *plans) GetPlans() (*GetPlansResult, error) {
	info, err := t.query.GetPlans()
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (t *plans) GetPlan(id *string) (*GetPlanResult, error) {
	info, err := t.query.GetPlan(id)
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (t *plans) Create(args *PlanArgs) (*CreatePlanResult, error) {
	if err := args.valid(); err != nil {
		return nil, err
	}
	res, err := t.service.Create(args.convert())
	if err != nil {
		return nil, err
	}
	return &CreatePlanResult{
		PlanID: res.PlanID,
	}, nil
}
func (t *PlanArgs) valid() error {
	err := core.NewError()
	if t.PlanName == "" {
		err.Append(application.RequiredPlanName)
	}
	if t.Interval <= 0 {
		err.Append(application.MoreThanZeroInterval)
	}
	if t.Start != nil && t.End != nil && t.Start.After(*t.End) {
		err.Append(application.InValidDateRange)
	}
	if err.HasError() {
		return err
	}
	return nil
}
func (t *PlanArgs) convert() *services.PlanArgs {
	return &services.PlanArgs{
		PlanName:   t.PlanName,
		IsIncome:   t.IsIncome,
		PlanAmount: t.PlanAmount,
		Interval:   t.Interval,
		Start:      t.Start,
		End:        t.End,
	}
}
func (t *plans) Update(id *string, args *PlanArgs) error {
	if err := args.valid(); err != nil {
		return err
	}
	return t.service.Update(id, args.convert())
}
func (t *plans) Remove(id *string) error {
	return t.service.Remove(id)
}
