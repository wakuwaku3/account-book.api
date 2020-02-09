package notifications

import (
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	metrics struct {
		value string
	}
	Metrics interface {
		Get() string
		Set(string) core.Error
		Valid() core.Error
	}
)

const (
	ExpenseBase = "ExpenseBase"
	BalanceBase = "BalanceBase"
	// NotSupportedMetricsFormat :メトリクスとしてサポートしていない形式です
	NotSupportedMetricsFormat core.ErrorCode = "notifications-00001"
)

func NewMetrics(value string) (Metrics, core.Error) {
	ins := &metrics{value}
	if err := ins.Valid(); err != nil {
		return nil, err
	}
	return ins, nil
}
func (t *metrics) Get() string { return t.value }
func (t *metrics) Set(value string) core.Error {
	t.value = value
	return t.Valid()
}
func (t *metrics) Equal(o Metrics) bool { return t.value == o.Get() }
func (t *metrics) Valid() core.Error {
	if t.value == ExpenseBase || t.value == BalanceBase {
		return nil
	}
	return core.NewError(NotSupportedMetricsFormat)
}
