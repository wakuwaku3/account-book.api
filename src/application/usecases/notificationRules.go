package usecases

import (
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
	"github.com/wakuwaku3/account-book.api/src/enterprise/domains/notifications"
)

type (
	notificationRules struct {
		query      NotificationRulesQuery
		repository notifications.NotificationRulesRepository
	}
	// NotificationRules is NotificationRulesUseCases
	NotificationRules interface {
		GetNotificationRules() *GetNotificationRulesResult
		GetNotificationRule(id *string) (*GetNotificationRuleResult, core.Error)
		Create(args *NotificationRuleArgs) (*CreateNotificationRuleResult, core.Error)
		Update(id *string, args *NotificationRuleArgs) core.Error
		Delete(id *string) core.Error
	}
	// GetNotificationRulesResult は結果です
	GetNotificationRulesResult struct {
		NotificationRules []GetNotificationRuleResult
	}
	// GetNotificationRuleResult は結果です
	GetNotificationRuleResult struct {
		NotificationRuleID string
		Metrics            string
		Threshold          int
	}
	// NotificationRuleArgs は引数です
	NotificationRuleArgs struct {
		Metrics   string
		Threshold int
	}
	// CreateNotificationRuleResult は結果です
	CreateNotificationRuleResult struct {
		NotificationRuleID string
	}
)

// NewNotificationRules is create instance
func NewNotificationRules(
	query NotificationRulesQuery,
	repository notifications.NotificationRulesRepository,
) NotificationRules {
	return &notificationRules{query, repository}
}
func (t *notificationRules) GetNotificationRules() *GetNotificationRulesResult {
	return t.query.GetNotificationRules()
}
func (t *notificationRules) GetNotificationRule(id *string) (*GetNotificationRuleResult, core.Error) {
	info, err := t.query.GetNotificationRule(id)
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (t *notificationRules) Create(args *NotificationRuleArgs) (*CreateNotificationRuleResult, core.Error) {
	metrics, err := notifications.NewMetrics(args.Metrics)
	if err != nil {
		return nil, err
	}
	threshold := notifications.NewThreshold(args.Threshold)
	res := t.repository.New(metrics, threshold)
	return &CreateNotificationRuleResult{
		NotificationRuleID: *res.GetID(),
	}, nil
}
func (t *notificationRules) Update(id *string, args *NotificationRuleArgs) core.Error {
	notificationRule, err := t.repository.GetByID(notifications.NotificationRuleID(id))
	metrics, err := notifications.NewMetrics(args.Metrics)
	if err != nil {
		return err
	}
	notificationRule.SetMetrics(metrics)
	notificationRule.SetThreshold(notifications.NewThreshold(args.Threshold))
	t.repository.Save(notificationRule)
	return nil
}
func (t *notificationRules) Delete(id *string) core.Error {
	return t.repository.Delete(notifications.NotificationRuleID(id))
}
