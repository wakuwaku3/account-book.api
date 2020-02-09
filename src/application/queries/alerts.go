package queries

import (
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
	"github.com/wakuwaku3/account-book.api/src/enterprise/domains/notifications"

	"github.com/wakuwaku3/account-book.api/src/application/usecases"
)

type alerts struct {
	repos notifications.AlertsRepository
}

// NewAlerts はインスタンスを生成します
func NewAlerts(
	repos notifications.AlertsRepository,
) usecases.AlertsQuery {
	return &alerts{repos}
}
func (t *alerts) GetAlerts() *usecases.GetAlertsResult {
	records := t.repos.Get()
	alerts := make([]usecases.GetAlertResult, len(*records))
	for i, record := range *records {
		r := &record
		alerts[i] = *convertAlert(r)
	}
	return &usecases.GetAlertsResult{Alerts: alerts}
}
func (t *alerts) GetAlert(id *string) (*usecases.GetAlertResult, core.Error) {
	alert, err := t.repos.GetByID(id)
	if err != nil {
		return nil, err
	}
	return convertAlert(&alert), nil
}
func convertAlert(t *notifications.Alert) *usecases.GetAlertResult {
	alert := *t
	return &usecases.GetAlertResult{
		AlertID:   alert.GetID(),
		Metrics:   alert.GetMetrics().Get(),
		Threshold: alert.GetThreshold().Get(),
	}
}
