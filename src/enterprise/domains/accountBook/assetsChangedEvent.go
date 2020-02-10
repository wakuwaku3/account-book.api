package accountbook

import "github.com/wakuwaku3/account-book.api/src/enterprise/core"

type (
	assetsChangedEvent struct {
		publisher      core.Publisher
		guidFactory    core.GuidFactory
		claimsProvider core.ClaimsProvider
	}
	// AssetsChangedEvent は入出金状態に変更があったときのイベントです
	AssetsChangedEvent interface {
		Trigger()
	}
	// assetsChangedEventMessage はイベントに付属するメッセージです
	assetsChangedEventMessage struct {
		ID     string `json:"id"`
		UserID string `json:"userID"`
	}
)

const (
	// EventName です
	EventName core.EventName = "AssetsChanged"
)

// NewAssetsChangedEvent はインスタンスを生成します
func NewAssetsChangedEvent(
	publisher core.Publisher,
	guidFactory core.GuidFactory,
	claimsProvider core.ClaimsProvider,
) AssetsChangedEvent {
	return &assetsChangedEvent{publisher, guidFactory, claimsProvider}
}
func (t *assetsChangedEvent) Trigger() {
	id, err := t.guidFactory.Create()
	if err != nil {
		panic(err)
	}
	message := assetsChangedEventMessage{ID: *id, UserID: *t.claimsProvider.GetUserID()}
	if err := t.publisher.Publish(EventName, message); err != nil {
		panic(err)
	}
}
