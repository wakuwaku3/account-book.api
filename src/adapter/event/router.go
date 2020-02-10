package event

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tampopos/dijct"
	handler "github.com/wakuwaku3/account-book.api/src/adapter/event/handlers"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	router struct {
		routingTable map[core.QueueName]Handler
	}
	// Router は Queue に対する ハンドラーの振り分けを設定します
	Router interface {
		Route(queueName core.QueueName, handler Handler) error
		Each(func(queueName core.QueueName, handler Handler) error) error
	}
	// Handler は キューイベントの処理です
	Handler func(dijct.Container, *string) error
)

// NewRouter はインスタンスを生成します
func NewRouter() Router {
	routingTable := map[core.QueueName]Handler{}
	instance := &router{routingTable}
	instance.Route("DeadLetter", func(container dijct.Container, message *string) error {
		return container.Invoke(func(alert handler.Alert) error {
			return alert.NotifyDeadLetter(message)
		})
	})
	instance.Route("NotifyAlert", func(container dijct.Container, message *string) error {
		var args handler.NotifyArgs
		if err := unmarshalJSON(message, &args); err != nil {
			return err
		}
		return container.Invoke(func(alert handler.Alert) error {
			return alert.Notify(&args)
		})
	})
	return instance
}
func unmarshalJSON(message *string, obj interface{}) error {
	if message == nil {
		return errors.New("messageが存在しないためUnmarshalできません")
	}
	bytes := []byte(*message)
	if err := json.Unmarshal(bytes, obj); err != nil {
		return err
	}
	return nil
}
func (t *router) Route(queueName core.QueueName, handler Handler) error {
	if _, ok := t.routingTable[queueName]; ok {
		return fmt.Errorf("指定されたキュー(%s)のルーティングは既に設定されています", queueName)
	}
	t.routingTable[queueName] = handler
	return nil
}
func (t *router) Each(fn func(queueName core.QueueName, handler Handler) error) error {
	for key := range t.routingTable {
		f, _ := t.routingTable[key]
		if err := fn(key, f); err != nil {
			return err
		}
	}
	return nil
}
