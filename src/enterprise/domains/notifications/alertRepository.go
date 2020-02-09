package notifications

import "github.com/wakuwaku3/account-book.api/src/enterprise/core"

type ( // AlertsRepository は計画のリポジトリです
	AlertsRepository interface {
		Get() *[]Alert
		GetByID(id AlertID) (Alert, core.Error)
		New(metrics Metrics, threshold Threshold) Alert
		Save(alert Alert)
		Delete(id AlertID) core.Error
	}
)
