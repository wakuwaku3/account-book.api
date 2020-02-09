package usecases

import (
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
	"github.com/wakuwaku3/account-book.api/src/enterprise/domains/notifications"
)

type (
	alerts struct {
		query      AlertsQuery
		repository notifications.AlertsRepository
	}
	// Alerts is AlertsUseCases
	Alerts interface {
		GetAlerts() *GetAlertsResult
		GetAlert(id *string) (*GetAlertResult, core.Error)
		Create(args *AlertArgs) (*CreateAlertResult, core.Error)
		Update(id *string, args *AlertArgs) core.Error
		Delete(id *string) core.Error
	}
	// GetAlertsResult は結果です
	GetAlertsResult struct {
		Alerts []GetAlertResult
	}
	// GetAlertResult は結果です
	GetAlertResult struct {
		AlertID   string
		Metrics   string
		Threshold int
	}
	// AlertArgs は引数です
	AlertArgs struct {
		Metrics   string
		Threshold int
	}
	// CreateAlertResult は結果です
	CreateAlertResult struct {
		AlertID string
	}
)

// NewAlerts is create instance
func NewAlerts(
	query AlertsQuery,
	repository notifications.AlertsRepository,
) Alerts {
	return &alerts{query, repository}
}
func (t *alerts) GetAlerts() *GetAlertsResult {
	return t.query.GetAlerts()
}
func (t *alerts) GetAlert(id *string) (*GetAlertResult, core.Error) {
	info, err := t.query.GetAlert(id)
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (t *alerts) Create(args *AlertArgs) (*CreateAlertResult, core.Error) {
	metrics, err := notifications.NewMetrics(args.Metrics)
	if err != nil {
		return nil, err
	}
	threshold := notifications.NewThreshold(args.Threshold)
	res := t.repository.New(metrics, threshold)
	return &CreateAlertResult{
		AlertID: res.GetID(),
	}, nil
}
func (t *alerts) Update(id *string, args *AlertArgs) core.Error {
	alert, err := t.repository.GetByID(id)
	metrics, err := notifications.NewMetrics(args.Metrics)
	if err != nil {
		return err
	}
	alert.SetMetrics(metrics)
	alert.SetThreshold(notifications.NewThreshold(args.Threshold))
	t.repository.Save(alert)
	return nil
}
func (t *alerts) Delete(id *string) core.Error {
	return t.repository.Delete(id)
}
