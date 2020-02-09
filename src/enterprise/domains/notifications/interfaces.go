package notifications

import "github.com/wakuwaku3/account-book.api/src/enterprise/domains/core"

type ( // AlertsRepository は計画のリポジトリです
	AlertsRepository interface {
		Get() *[]Alert
		GetByID(id *string) (Alert, core.Error)
		New(metrics Metrics, threshold Threshold) Alert
		Save(alert Alert)
		Delete(id *string) core.Error
	}
)
