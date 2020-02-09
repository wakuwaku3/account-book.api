package queries

import (
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
	"github.com/wakuwaku3/account-book.api/src/enterprise/domains/notifications"

	"github.com/wakuwaku3/account-book.api/src/application/usecases"
)

type notificationRules struct {
	repos notifications.NotificationRulesRepository
}

// NewNotificationRules はインスタンスを生成します
func NewNotificationRules(
	repos notifications.NotificationRulesRepository,
) usecases.NotificationRulesQuery {
	return &notificationRules{repos}
}
func (t *notificationRules) GetNotificationRules() *usecases.GetNotificationRulesResult {
	records := t.repos.Get()
	notificationRules := make([]usecases.GetNotificationRuleResult, len(*records))
	for i, record := range *records {
		r := &record
		notificationRules[i] = *convertNotificationRule(r)
	}
	return &usecases.GetNotificationRulesResult{NotificationRules: notificationRules}
}
func (t *notificationRules) GetNotificationRule(id *string) (*usecases.GetNotificationRuleResult, core.Error) {
	notificationRule, err := t.repos.GetByID(id)
	if err != nil {
		return nil, err
	}
	return convertNotificationRule(&notificationRule), nil
}
func convertNotificationRule(t *notifications.NotificationRule) *usecases.GetNotificationRuleResult {
	notificationRule := *t
	return &usecases.GetNotificationRuleResult{
		NotificationRuleID: *notificationRule.GetID(),
		Metrics:            notificationRule.GetMetrics().Get(),
		Threshold:          notificationRule.GetThreshold().Get(),
	}
}
