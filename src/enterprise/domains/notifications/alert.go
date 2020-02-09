package notifications

type (
	alert struct {
		id        AlertID
		metrics   Metrics
		threshold Threshold
	}
	AlertID *string
	Alert   interface {
		GetID() AlertID
		GetMetrics() Metrics
		SetMetrics(Metrics)
		GetThreshold() Threshold
		SetThreshold(Threshold)
		Equal(alert Alert) bool
	}
)

func NewAlert(
	id AlertID,
	metrics Metrics,
	threshold Threshold,
) Alert {
	return &alert{id, metrics, threshold}
}

func (t *alert) GetID() AlertID               { return t.id }
func (t *alert) GetMetrics() Metrics          { return t.metrics }
func (t *alert) SetMetrics(value Metrics)     { t.metrics = value }
func (t *alert) GetThreshold() Threshold      { return t.threshold }
func (t *alert) SetThreshold(value Threshold) { t.threshold = value }
func (t *alert) Equal(alert Alert) bool {
	id := t.id
	aID := alert.GetID()
	if id == nil || *id == "" || aID == nil || *aID == "" {
		return t == alert
	}
	return id == aID
}
