package queries

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains/models"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"

	"github.com/wakuwaku3/account-book.api/src/usecases"
)

type dashboard struct {
	repos             domains.DashboardRepository
	transactionsRepos domains.TransactionsRepository
	plansRepos        domains.PlansRepository
	clock             cmn.Clock
}

// NewDashboard はインスタンスを生成します
func NewDashboard(
	repos domains.DashboardRepository,
	transactionsRepos domains.TransactionsRepository,
	plansRepos domains.PlansRepository,
	clock cmn.Clock,
) usecases.DashboardQuery {
	return &dashboard{
		repos,
		transactionsRepos,
		plansRepos,
		clock,
	}
}

func (t *dashboard) GetSummary(args *usecases.GetDashboardArgs) (*usecases.GetDashboardResult, error) {
	if args.SelectedMonth != nil {
		return t.getSummaryByMonth(args.SelectedMonth)
	}

	selectedMonth := args.SelectedMonth

	// 当月のダッシュボード取得
	currentDashboard, err := t.repos.GetOldestOpenDashboard()
	if err != nil {
		return nil, err
	}

	if selectedMonth == nil {
		if currentDashboard != nil {
			selectedMonth = &currentDashboard.Date
			return t.getSummaryByMonthWithCurrentDashboard(selectedMonth, currentDashboard)
		}
	}

	// 前月のダッシュボード取得
	previousDashboard, err := t.getPreviousDashboard(selectedMonth)
	if err != nil {
		return nil, err
	}
	if selectedMonth == nil {
		m := t.clock.GetMonthStartDay(nil)
		if previousDashboard != nil {
			m = previousDashboard.Date.AddDate(0, 1, 0)
		}
		selectedMonth = &m
	}

	result, allInputPlans, err := t.getSummaryByMonthWithoutPreviousDashboard(selectedMonth, currentDashboard)
	if err != nil {
		return nil, err
	}

	if previousDashboard != nil {
		// 前月ダッシュボードが取得できた場合のみ設定
		result.PreviousBalance = previousDashboard.Balance
	}
	result.CanApprove = t.canApprove(result, allInputPlans, previousDashboard)
	result.CanCancelApprove = false
	return result, nil
}
func (t *dashboard) canApprove(result *usecases.GetDashboardResult, allInputPlans bool, previousDashboard *models.Dashboard) bool {
	if !allInputPlans {
		return false
	}
	if previousDashboard != nil && previousDashboard.State == "open" {
		return false
	}
	monthStart := t.clock.GetMonthStartDay(nil)
	return result.SelectedMonth.Before(monthStart)
}
func (t *dashboard) getPreviousDashboard(selectedMonth *time.Time) (*models.Dashboard, error) {
	if selectedMonth != nil {
		previousMonth := selectedMonth.AddDate(0, -1, 0)
		return t.repos.GetByMonth(&previousMonth)
	}
	return t.repos.GetLatestClosedDashboard()
}
func (t *dashboard) getSummaryByMonth(selectedMonth *time.Time) (*usecases.GetDashboardResult, error) {
	currentDashboard, err := t.repos.GetByMonth(selectedMonth)
	if err != nil {
		return nil, err
	}
	if currentDashboard != nil && currentDashboard.State == "closed" {
		chError := make(chan error)
		chNext := t.getDashboardByNextMonthWorker(selectedMonth, chError)

		// 締め処理済みの場合、集計済みなのでそのまま返す
		plans := make([]usecases.PlanResult, len(currentDashboard.Actual))
		for i, actual := range currentDashboard.Actual {
			actualID := actual.ActualID
			actualAmount := actual.ActualAmount
			plans[i] = usecases.PlanResult{
				ActualID:     &actualID,
				IsIncome:     actual.IsIncome,
				ActualAmount: &actualAmount,
				PlanAmount:   actual.PlanAmount,
				PlanID:       actual.PlanID,
				PlanName:     actual.PlanName,
				CreatedAt:    actual.PlanCreatedAt,
			}
		}

		result := &usecases.GetDashboardResult{
			DashboardID:     currentDashboard.DashboardID,
			Expense:         *currentDashboard.Expense,
			Income:          *currentDashboard.Income,
			PreviousBalance: currentDashboard.PreviousBalance,
			SelectedMonth:   *selectedMonth,
			Plans:           plans,
			State:           "closed",
		}

		select {
		case err := <-chError:
			return nil, err
		case nextDashboard := <-chNext:
			result.CanApprove = false
			result.CanCancelApprove = nextDashboard == nil || nextDashboard.State == "open"
			break
		}

		return result, nil
	}

	return t.getSummaryByMonthWithCurrentDashboard(selectedMonth, currentDashboard)
}
func (t *dashboard) getSummaryByMonthWithCurrentDashboard(selectedMonth *time.Time, currentDashboard *models.Dashboard) (*usecases.GetDashboardResult, error) {
	chError := make(chan error)
	chPrevious := t.getDashboardByPreviousMonthWorker(selectedMonth, chError)

	result, allInputPlans, err := t.getSummaryByMonthWithoutPreviousDashboard(selectedMonth, currentDashboard)
	if err != nil {
		return nil, err
	}

	select {
	case err := <-chError:
		return nil, err
	case previousDashboard := <-chPrevious:
		if previousDashboard != nil {
			// 前月ダッシュボードが取得できた場合のみ設定
			result.PreviousBalance = previousDashboard.Balance
		}
		result.CanApprove = t.canApprove(result, allInputPlans, previousDashboard)
		result.CanCancelApprove = false
		break
	}

	return result, nil
}
func (t *dashboard) getSummaryByMonthWithoutPreviousDashboard(selectedMonth *time.Time, currentDashboard *models.Dashboard) (*usecases.GetDashboardResult, bool, error) {
	chError := make(chan error)

	chIncome, chExpense := t.getTransactionsSummaryWorker(selectedMonth, chError)
	chPln := t.getPlansWorker(selectedMonth, chError)

	result := new(usecases.GetDashboardResult)
	result.SelectedMonth = *selectedMonth

	var pMap map[string]usecases.PlanResult
	select {
	case err := <-chError:
		return nil, false, err
	case p := <-chPln:
		pMap = *p
		break
	}

	if currentDashboard != nil {
		result.DashboardID = currentDashboard.DashboardID
		// 実績がある場合優先して上書きする
		for _, actual := range currentDashboard.Actual {
			actualID := actual.ActualID
			actualAmount := actual.ActualAmount
			pMap[actual.PlanID] = usecases.PlanResult{
				ActualID:     &actualID,
				IsIncome:     actual.IsIncome,
				ActualAmount: &actualAmount,
				PlanAmount:   actual.PlanAmount,
				PlanID:       actual.PlanID,
				PlanName:     actual.PlanName,
			}
		}
	}

	// 計画から集計(ついでに戻り値として型を整形する)
	income := 0
	expense := 0
	allInput := true
	ps := make([]usecases.PlanResult, len(pMap))
	index := 0
	for _, pm := range pMap {
		if pm.ActualID == nil {
			if pm.IsIncome {
				income += pm.PlanAmount
			} else {
				expense += pm.PlanAmount
			}
			allInput = false
		} else {
			if pm.IsIncome {
				income += *pm.ActualAmount
			} else {
				expense += *pm.ActualAmount
			}
		}
		ps[index] = pm
		index++
	}

	select {
	case err := <-chError:
		return nil, false, err
	case i := <-chIncome:
		income += i
		break
	}

	select {
	case err := <-chError:
		return nil, false, err
	case e := <-chExpense:
		expense += e
		break
	}

	result.Income = income
	result.Expense = expense
	result.Plans = ps
	result.State = "open"
	return result, allInput, nil
}
func (t *dashboard) getDashboardByMonthWorker(selectedMonth *time.Time, chError chan error) <-chan *models.Dashboard {
	ch := make(chan *models.Dashboard)
	go func() {
		d, err := t.repos.GetByMonth(selectedMonth)
		if err != nil {
			chError <- err
			return
		}
		ch <- d
	}()
	return ch
}
func (t *dashboard) getDashboardByPreviousMonthWorker(selectedMonth *time.Time, chError chan error) <-chan *models.Dashboard {
	previousMonth := selectedMonth.AddDate(0, -1, 0)
	return t.getDashboardByMonthWorker(&previousMonth, chError)
}
func (t *dashboard) getDashboardByNextMonthWorker(selectedMonth *time.Time, chError chan error) <-chan *models.Dashboard {
	previousMonth := selectedMonth.AddDate(0, 1, 0)
	return t.getDashboardByMonthWorker(&previousMonth, chError)
}
func (t *dashboard) getTransactionsSummaryWorker(selectedMonth *time.Time, chError chan error) (<-chan int, <-chan int) {
	chIncome := make(chan int)
	chExpense := make(chan int)
	go func() {
		transactions, err := t.transactionsRepos.GetByMonth(selectedMonth)
		if err != nil {
			chError <- err
			return
		}
		// 収入と支出を集計する
		income := 0
		expense := 0
		// 取引から集計
		for _, transaction := range *transactions {
			if transaction.Category == 5 {
				income += transaction.Amount
			} else {
				expense += transaction.Amount
			}
		}
		chIncome <- income
		chExpense <- expense
	}()
	return chIncome, chExpense
}
func (t *dashboard) getPlansWorker(selectedMonth *time.Time, chError chan error) <-chan *map[string]usecases.PlanResult {
	ch := make(chan *map[string]usecases.PlanResult)

	go func() {
		plans, err := t.plansRepos.GetByMonth(selectedMonth)
		if err != nil {
			chError <- err
			return
		}
		pMap := map[string]usecases.PlanResult{}
		for _, plan := range *plans {
			pMap[plan.PlanID] = usecases.PlanResult{
				IsIncome:   plan.IsIncome,
				PlanAmount: plan.PlanAmount,
				PlanID:     plan.PlanID,
				PlanName:   plan.PlanName,
				CreatedAt:  plan.CreatedAt,
			}
		}
		ch <- &pMap
	}()
	return ch
}
