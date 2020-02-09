package repos

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
	"github.com/wakuwaku3/account-book.api/src/enterprise/domains/notifications"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wakuwaku3/account-book.api/src/adapter/store"
	"github.com/wakuwaku3/account-book.api/src/application"
)

type (
	alerts struct {
		provider       store.Provider
		clock          core.Clock
		claimsProvider application.ClaimsProvider
	}
	alertEntity struct {
		Metrics   string `firestore:"metrics"`
		Threshold int    `firestore:"threshold"`
	}
)

// NewAlerts はインスタンスを生成します
func NewAlerts(
	provider store.Provider,
	clock core.Clock,
	claimsProvider application.ClaimsProvider,
) notifications.AlertsRepository {
	return &alerts{provider, clock, claimsProvider}
}
func (t *alerts) alertsRef(client *firestore.Client) *firestore.CollectionRef {
	userID := t.claimsProvider.GetUserID()
	return client.Collection("users").Doc(*userID).Collection("alerts")
}
func (t *alerts) GetByID(id notifications.AlertID) (notifications.Alert, core.Error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := t.alertsRef(client).Doc(*id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, core.NewError(core.NotFound)
		}
		panic(err)
	}
	var entity alertEntity
	if err := doc.DataTo(&entity); err != nil {
		panic(err)
	}
	return entity.newAlert(id), nil
}
func (t alertEntity) newAlert(id notifications.AlertID) notifications.Alert {
	metrics, err := notifications.NewMetrics(t.Metrics)
	if err != nil {
		panic(err)
	}
	return notifications.NewAlert(id, metrics, notifications.NewThreshold(t.Threshold))
}
func (t *alerts) Get() *[]notifications.Alert {
	client := t.provider.GetClient()
	ctx := context.Background()

	alerts := make([]notifications.Alert, 0)
	iter := t.alertsRef(client).Where("isDeleted", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		var entity alertEntity
		if err := doc.DataTo(&entity); err != nil {
			panic(err)
		}
		alerts = append(alerts, entity.newAlert(notifications.AlertID(&doc.Ref.ID)))
	}
	return &alerts
}
func (t *alerts) New(metrics notifications.Metrics, threshold notifications.Threshold) notifications.Alert {
	client := t.provider.GetClient()
	ctx := context.Background()
	entity := newAlertEntity(metrics, threshold)
	ref, _, err := t.alertsRef(client).Add(ctx, entity)
	if err != nil {
		panic(err)
	}
	return notifications.NewAlert(notifications.AlertID(&ref.ID), metrics, threshold)
}
func (t *alerts) Save(alert notifications.Alert) {
	client := t.provider.GetClient()
	ctx := context.Background()
	entity := newAlertEntity(alert.GetMetrics(), alert.GetThreshold())
	_, err := t.alertsRef(client).Doc(*alert.GetID()).Set(ctx, entity)
	if err != nil {
		panic(err)
	}
}
func newAlertEntity(metrics notifications.Metrics, threshold notifications.Threshold) *alertEntity {
	return &alertEntity{
		Metrics:   metrics.Get(),
		Threshold: threshold.Get(),
	}
}
func (t *alerts) Delete(id notifications.AlertID) core.Error {
	client := t.provider.GetClient()
	ctx := context.Background()
	_, err := t.alertsRef(client).Doc(*id).Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return core.NewError(core.NotFound)
		}
		panic(err)
	}
	return nil
}
