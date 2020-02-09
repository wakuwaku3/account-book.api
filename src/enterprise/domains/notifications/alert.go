package notifications

type (
	alert struct {
		id        string
		metrics   Metrics
		threshold Threshold
	}
	Alert interface {
		GetID() string
		GetMetrics() Metrics
		SetMetrics(Metrics)
		GetThreshold() Threshold
		SetThreshold(Threshold)
		Equal(alert Alert) bool
	}
)

func NewAlert(
	id string,
	metrics Metrics,
	threshold Threshold,
) Alert {
	return &alert{id, metrics, threshold}
}

func (t *alert) GetID() string                { return t.id }
func (t *alert) GetMetrics() Metrics          { return t.metrics }
func (t *alert) SetMetrics(value Metrics)     { t.metrics = value }
func (t *alert) GetThreshold() Threshold      { return t.threshold }
func (t *alert) SetThreshold(value Threshold) { t.threshold = value }
func (t *alert) Equal(alert Alert) bool {
	id := t.id
	aID := alert.GetID()
	if id == "" || aID == "" {
		return t == alert
	}
	return id == aID
}
