package services

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/application"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
)

type (
	plans struct {
		repos domains.PlansRepository
		clock cmn.Clock
	}
	// Plans is PlansService
	Plans interface {
		Create(args *PlanArgs) (*CreatePlanResult, error)
		Update(id *string, args *PlanArgs) error
		Remove(id *string) error
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
func NewPlans(repos domains.PlansRepository, clock cmn.Clock) Plans {
	return &plans{repos, clock}
}
func (t *plans) Create(args *PlanArgs) (*CreatePlanResult, error) {
	id, err := t.repos.Create(args.convert(t.clock.Now()))
	if err != nil {
		return nil, err
	}
	return &CreatePlanResult{PlanID: *id}, nil
}
func (t *PlanArgs) convert(now time.Time) *models.Plan {
	return &models.Plan{
		PlanName:   t.PlanName,
		IsIncome:   t.IsIncome,
		PlanAmount: t.PlanAmount,
		Interval:   t.Interval,
		Start:      t.Start,
		End:        t.End,
	}
}
func (t *plans) Update(id *string, args *PlanArgs) error {
	model, err := t.repos.GetByID(id)
	if err != nil {
		return err
	}
	if model.IsDeleted {
		return application.NewClientError(application.IsDeleted)
	}

	model.PlanName = args.PlanName
	model.IsIncome = args.IsIncome
	model.PlanAmount = args.PlanAmount
	model.Interval = args.Interval
	model.Start = args.Start
	model.End = args.End

	if err := t.repos.Update(id, model); err != nil {
		return err
	}
	return nil
}
func (t *plans) Remove(id *string) error {
	model, err := t.repos.GetByID(id)
	if err != nil {
		return err
	}
	if model.IsDeleted {
		return application.NewClientError(application.IsDeleted)
	}
	model.IsDeleted = true
	if err := t.repos.Update(id, model); err != nil {
		return err
	}
	return nil
}
