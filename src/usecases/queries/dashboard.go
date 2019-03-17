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
	selectedMonth := args.SelectedMonth

	// 当月のダッシュボード取得
	currentDashboard, err := t.GetCurrentDashboard(selectedMonth)
	if err != nil {
		return nil, err
	}

	if currentDashboard != nil && currentDashboard.State == "closed" {
		// 締め処理済みの場合、集計済みなのでそのまま返す
		plans := make([]usecases.PlanResult, len(currentDashboard.Actual))
		for i, actual := range currentDashboard.Actual {
			plans[i] = usecases.PlanResult{
				ActualID:     &actual.ActualID,
				IsIncome:     actual.IsIncome,
				ActualAmount: &actual.ActualAmount,
				PlanAmount:   actual.PlanAmount,
				PlanID:       actual.PlanID,
				PlanName:     actual.PlanName,
			}
		}
		return &usecases.GetDashboardResult{
			Expense:         *currentDashboard.Expense,
			Income:          *currentDashboard.Income,
			PreviousBalance: currentDashboard.PreviousBalance,
			SelectedMonth:   *selectedMonth,
			Plans:           plans,
		}, nil
	}

	if selectedMonth == nil {
		if currentDashboard != nil {
			selectedMonth = &currentDashboard.Date
		}
	}

	// 前月のダッシュボード取得
	previousDashboard, err := t.GetPreviousDashboard(selectedMonth)
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

	// 当月の取引取得
	transactions, err := t.transactionsRepos.GetByMonth(selectedMonth)
	if err != nil {
		return nil, err
	}

	// 当月の計画取得
	plans, err := t.plansRepos.GetByMonth(selectedMonth)
	if err != nil {
		return nil, err
	}

	// 集計処理
	result := new(usecases.GetDashboardResult)
	result.SelectedMonth = *selectedMonth
	if previousDashboard != nil {
		// 前月ダッシュボードが取得できた場合のみ設定
		result.PreviousBalance = previousDashboard.Balance
	}

	// 計画と実績をマージする
	pMap := map[string]usecases.PlanResult{}
	for _, plan := range *plans {
		pMap[plan.PlanID] = usecases.PlanResult{
			IsIncome:   plan.IsIncome,
			PlanAmount: plan.PlanAmount,
			PlanID:     plan.PlanID,
			PlanName:   plan.PlanName,
		}
	}
	if currentDashboard != nil {
		// 実績がある場合優先して上書きする
		for _, actual := range currentDashboard.Actual {
			pMap[actual.PlanID] = usecases.PlanResult{
				ActualID:     &actual.ActualID,
				IsIncome:     actual.IsIncome,
				ActualAmount: &actual.ActualAmount,
				PlanAmount:   actual.PlanAmount,
				PlanID:       actual.PlanID,
				PlanName:     actual.PlanName,
			}
		}
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
	// 計画から集計(ついでに戻り値として型を整形する)
	ps := make([]usecases.PlanResult, len(pMap))
	index := 0
	for _, pm := range pMap {
		if pm.ActualID == nil {
			if pm.IsIncome {
				income += pm.PlanAmount
			} else {
				expense += pm.PlanAmount
			}
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

	// 集計した値を設定
	result.Income = income
	result.Expense = expense
	result.Plans = ps
	return result, nil
}

func (t *dashboard) GetCurrentDashboard(selectedMonth *time.Time) (*models.Dashboard, error) {
	if selectedMonth != nil {
		return t.repos.GetByMonth(selectedMonth)
	}
	return t.repos.GetOldestOpenDashboard()
}
func (t *dashboard) GetPreviousDashboard(selectedMonth *time.Time) (*models.Dashboard, error) {
	if selectedMonth != nil {
		previousMonth := selectedMonth.AddDate(0, -1, 0)
		return t.repos.GetByMonth(&previousMonth)
	}
	return t.repos.GetLatestClosedDashboard()
}
