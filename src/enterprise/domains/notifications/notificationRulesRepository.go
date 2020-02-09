package notifications

import "github.com/wakuwaku3/account-book.api/src/enterprise/core"

type (
	// NotificationRulesRepository は計画のリポジトリです
	NotificationRulesRepository interface {
		Get() *[]NotificationRule
		GetByID(id NotificationRuleID) (NotificationRule, core.Error)
		New(metrics Metrics, threshold Threshold) NotificationRule
		Save(notificationRule NotificationRule)
		Delete(id NotificationRuleID) core.Error
	}
)
