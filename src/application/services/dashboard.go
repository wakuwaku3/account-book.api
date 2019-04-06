package services

import (
	"errors"
	"sort"
	"time"

	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
)

type (
	dashboard struct {
		repos             application.DashboardRepository
		transactionsRepos application.TransactionsRepository
		plansRepos        application.PlansRepository
		clock             cmn.Clock
	}
	// Dashboard is DashboardService
	Dashboard interface {
		Approve(id *string) error
		CancelApprove(id *string) error
	}
)

// NewDashboard is create instance
func NewDashboard(
	repos application.DashboardRepository,
	transactionsRepos application.TransactionsRepository,
	plansRepos application.PlansRepository,
	clock cmn.Clock,
) Dashboard {
	return &dashboard{
		repos,
		transactionsRepos,
		plansRepos,
		clock,
	}
}
func (t *dashboard) Approve(id *string) error {
	current, err := t.repos.GetByID(id)
	if err != nil {
		return err
	}
	if current == nil {
		return errors.New("dashboard is not found")
	}
	if current.State == "closed" {
		return errors.New("this dashboard is already closed")
	}

	previous, err := t.repos.GetLatestClosedDashboard()
	if err != nil {
		return err
	}
	if previous != nil && previous.State != "closed" {
		return errors.New("previous dashboard is not closed")
	}
	chError := make(chan error)

	chTrn := t.getTransactionsWorker(&current.Date, chError)
	chPln := t.getPlansWorker(&current.Date, chError)

	var plans []models.Plan
	select {
	case err := <-chError:
		return err
	case p := <-chPln:
		plans = *p
		break
	}

	// 集計
	income := 0
	expense := 0
	actualMap := make(map[string]models.Actual)
	for _, actual := range current.Actual {
		actualMap[actual.PlanID] = actual
		if actual.IsIncome {
			income += actual.ActualAmount
		} else {
			expense += actual.ActualAmount
		}
	}

	for _, plan := range plans {
		_, ok := actualMap[plan.PlanID]
		if !ok {
			return errors.New("has not input plan")
		}
	}

	var trn []models.Transaction
	select {
	case err := <-chError:
		return err
	case t := <-chTrn:
		trn = *t
		break
	}

	// 取引から集計
	dMap := make(map[string]*models.Daily)
	for _, transaction := range trn {
		key := transaction.Date.Format("2006-01-02")
		val, ok := dMap[key]
		if !ok {
			val = &models.Daily{
				Date:    t.clock.GetDay(&transaction.Date),
				Expense: 0,
				Income:  0,
			}
			dMap[key] = val
		}
		if transaction.Category == 5 {
			val.Income += transaction.Amount
		} else {
			val.Expense += transaction.Amount
		}
	}
	dSlice := make([]models.Daily, 0)
	for _, d := range dMap {
		dSlice = append(dSlice, *d)
	}
	sort.SliceStable(dSlice, func(i, j int) bool { return dSlice[i].Date.Before(dSlice[j].Date) })

	current.Income = &income
	current.Expense = &expense
	currentBalance := income - expense
	current.CurrentBalance = &currentBalance
	if previous != nil {
		current.PreviousBalance = previous.Balance
		current.PreviousDashboardID = &previous.DashboardID
	}
	balance := *current.PreviousBalance + currentBalance
	current.Balance = &balance
	current.State = "closed"
	current.Daily = dSlice
	return t.repos.Approve(current)
}
func (t *dashboard) getTransactionsWorker(selectedMonth *time.Time, chError chan error) <-chan *[]models.Transaction {
	ch := make(chan *[]models.Transaction)
	go func() {
		transactions, err := t.transactionsRepos.GetByMonth(selectedMonth)
		if err != nil {
			chError <- err
			return
		}
		ch <- transactions
	}()
	return ch
}
func (t *dashboard) getPlansWorker(selectedMonth *time.Time, chError chan error) <-chan *[]models.Plan {
	ch := make(chan *[]models.Plan)

	go func() {
		plans, err := t.plansRepos.GetByMonth(selectedMonth)
		if err != nil {
			chError <- err
			return
		}
		ch <- plans
	}()
	return ch
}
func (t *dashboard) CancelApprove(id *string) error {
	current, err := t.repos.GetByID(id)
	if err != nil {
		return err
	}
	if current == nil {
		return errors.New("dashboard is not found")
	}
	if current.State != "closed" {
		return errors.New("this dashboard is not closed")
	}
	if err := t.repos.ExistsClosedNext(id); err != nil {
		return err
	}
	current.Balance = nil
	current.CurrentBalance = nil
	current.Expense = nil
	current.Income = nil
	current.PreviousBalance = nil
	current.PreviousDashboardID = nil
	current.State = "open"
	return t.repos.CancelApprove(current)
}
