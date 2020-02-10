package handler

import "github.com/wakuwaku3/account-book.api/src/enterprise/core"

import "fmt"

type (
	alert struct{ provider core.ClaimsProvider }
	// Alert は ハンドラーです
	Alert interface {
		Notify(arg *NotifyArgs) error
		NotifyDeadLetter(body *string) error
	}
	// NotifyArgs は Notify の Args です
	NotifyArgs struct {
		ID     string `json:"id"`
		UserID string `json:"userID"`
	}
)

// NewAlert はインスタンスを生成します
func NewAlert(provider core.ClaimsProvider) Alert {
	return &alert{provider}
}
func (t *alert) Notify(arg *NotifyArgs) error {
	userID := *t.provider.GetUserID()
	fmt.Println(*arg)
	fmt.Println(userID)
	return nil
}
func (t *alert) NotifyDeadLetter(body *string) error {
	fmt.Println(*body)
	return nil
}
