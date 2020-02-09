package notifications

type (
	notificationRule struct {
		id        NotificationRuleID
		metrics   Metrics
		threshold Threshold
	}
	// NotificationRuleID は通知ルールの ID です
	NotificationRuleID *string

	// NotificationRule は通知ルールの Entity です
	NotificationRule interface {
		GetID() NotificationRuleID
		GetMetrics() Metrics
		SetMetrics(Metrics)
		GetThreshold() Threshold
		SetThreshold(Threshold)
		Equal(notificationRule NotificationRule) bool
	}
)

// NewNotificationRule は通知ルールを生成します
func NewNotificationRule(
	id NotificationRuleID,
	metrics Metrics,
	threshold Threshold,
) NotificationRule {
	return &notificationRule{id, metrics, threshold}
}

func (t *notificationRule) GetID() NotificationRuleID    { return t.id }
func (t *notificationRule) GetMetrics() Metrics          { return t.metrics }
func (t *notificationRule) SetMetrics(value Metrics)     { t.metrics = value }
func (t *notificationRule) GetThreshold() Threshold      { return t.threshold }
func (t *notificationRule) SetThreshold(value Threshold) { t.threshold = value }
func (t *notificationRule) Equal(notificationRule NotificationRule) bool {
	id := t.id
	aID := notificationRule.GetID()
	if id == nil || *id == "" || aID == nil || *aID == "" {
		return t == notificationRule
	}
	return id == aID
}
