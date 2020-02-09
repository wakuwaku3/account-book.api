package core

type (
	// EventName はイベント名です
	EventName string
	// JSONMessage はイベントに付与するJsonです
	JSONMessage interface{}
	// Publisher はイベントを発行します
	Publisher interface {
		Publish(name EventName, jsonMessage JSONMessage) error
	}
)
