package repos

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"

	"github.com/wakuwaku3/account-book.api/src/domains/models"
	"google.golang.org/api/iterator"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
)

type (
	plans struct {
		provider       store.Provider
		clock          cmn.Clock
		claimsProvider domains.ClaimsProvider
	}
)

// NewPlans はインスタンスを生成します
func NewPlans(
	provider store.Provider,
	clock cmn.Clock,
	claimsProvider domains.ClaimsProvider,
) domains.PlansRepository {
	return &plans{provider, clock, claimsProvider}
}
func (t *plans) plansRef(client *firestore.Client) *firestore.CollectionRef {
	userID := t.claimsProvider.GetUserID()
	return client.Collection("users").Doc(*userID).Collection("plans")
}
func (t *plans) GetByID(id *string) (*models.Plan, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := t.plansRef(client).Doc(*id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var plan models.Plan
	if err := doc.DataTo(&plan); err != nil {
		return nil, err
	}
	plan.PlanID = *id
	return &plan, nil
}
func (t *plans) Get() (*[]models.Plan, error) {
	client := t.provider.GetClient()
	ctx := context.Background()

	plans := make([]models.Plan, 0)
	iter := t.plansRef(client).Where("is-deleted", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var plan models.Plan
		if err := doc.DataTo(&plan); err != nil {
			return nil, err
		}
		plan.PlanID = doc.Ref.ID
		plans = append(plans, plan)
	}
	return &plans, nil
}
func (t *plans) Create(model *models.Plan) (*string, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	ref, _, err := t.plansRef(client).Add(ctx, model)
	if err != nil {
		return nil, err
	}
	return &ref.ID, nil
}
func (t *plans) Update(id *string, model *models.Plan) error {
	client := t.provider.GetClient()
	ctx := context.Background()
	_, err := t.plansRef(client).Doc(*id).Set(ctx, model)
	if err != nil {
		return err
	}
	return nil
}
func (t *plans) Delete(id *string) error {
	client := t.provider.GetClient()
	ctx := context.Background()
	_, err := t.plansRef(client).Doc(*id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
