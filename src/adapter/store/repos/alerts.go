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
	notificationRules struct {
		provider       store.Provider
		clock          core.Clock
		claimsProvider application.ClaimsProvider
	}
	notificationRuleEntity struct {
		Metrics   string `firestore:"metrics"`
		Threshold int    `firestore:"threshold"`
	}
)

// NewNotificationRules はインスタンスを生成します
func NewNotificationRules(
	provider store.Provider,
	clock core.Clock,
	claimsProvider application.ClaimsProvider,
) notifications.NotificationRulesRepository {
	return &notificationRules{provider, clock, claimsProvider}
}
func (t *notificationRules) notificationRulesRef(client *firestore.Client) *firestore.CollectionRef {
	userID := t.claimsProvider.GetUserID()
	return client.Collection("users").Doc(*userID).Collection("notificationRules")
}
func (t *notificationRules) GetByID(id notifications.NotificationRuleID) (notifications.NotificationRule, core.Error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := t.notificationRulesRef(client).Doc(*id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, core.NewError(core.NotFound)
		}
		panic(err)
	}
	var entity notificationRuleEntity
	if err := doc.DataTo(&entity); err != nil {
		panic(err)
	}
	return entity.newNotificationRule(id), nil
}
func (t notificationRuleEntity) newNotificationRule(id notifications.NotificationRuleID) notifications.NotificationRule {
	metrics, err := notifications.NewMetrics(t.Metrics)
	if err != nil {
		panic(err)
	}
	return notifications.NewNotificationRule(id, metrics, notifications.NewThreshold(t.Threshold))
}
func (t *notificationRules) Get() *[]notifications.NotificationRule {
	client := t.provider.GetClient()
	ctx := context.Background()

	notificationRules := make([]notifications.NotificationRule, 0)
	iter := t.notificationRulesRef(client).Where("isDeleted", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		var entity notificationRuleEntity
		if err := doc.DataTo(&entity); err != nil {
			panic(err)
		}
		notificationRules = append(notificationRules, entity.newNotificationRule(notifications.NotificationRuleID(&doc.Ref.ID)))
	}
	return &notificationRules
}
func (t *notificationRules) New(metrics notifications.Metrics, threshold notifications.Threshold) notifications.NotificationRule {
	client := t.provider.GetClient()
	ctx := context.Background()
	entity := newNotificationRuleEntity(metrics, threshold)
	ref, _, err := t.notificationRulesRef(client).Add(ctx, entity)
	if err != nil {
		panic(err)
	}
	return notifications.NewNotificationRule(notifications.NotificationRuleID(&ref.ID), metrics, threshold)
}
func (t *notificationRules) Save(notificationRule notifications.NotificationRule) {
	client := t.provider.GetClient()
	ctx := context.Background()
	entity := newNotificationRuleEntity(notificationRule.GetMetrics(), notificationRule.GetThreshold())
	_, err := t.notificationRulesRef(client).Doc(*notificationRule.GetID()).Set(ctx, entity)
	if err != nil {
		panic(err)
	}
}
func newNotificationRuleEntity(metrics notifications.Metrics, threshold notifications.Threshold) *notificationRuleEntity {
	return &notificationRuleEntity{
		Metrics:   metrics.Get(),
		Threshold: threshold.Get(),
	}
}
func (t *notificationRules) Delete(id notifications.NotificationRuleID) core.Error {
	client := t.provider.GetClient()
	ctx := context.Background()
	_, err := t.notificationRulesRef(client).Doc(*id).Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return core.NewError(core.NotFound)
		}
		panic(err)
	}
	return nil
}
