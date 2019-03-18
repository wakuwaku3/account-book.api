package repos

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
	"google.golang.org/api/iterator"
)

type dashboard struct {
	provider       store.Provider
	claimsProvider domains.ClaimsProvider
	clock          cmn.Clock
}

// NewDashboard はインスタンスを生成します
func NewDashboard(
	provider store.Provider,
	claimsProvider domains.ClaimsProvider,
	clock cmn.Clock,
) domains.DashboardRepository {
	return &dashboard{provider, claimsProvider, clock}
}

func (t *dashboard) dashboardsRef(client *firestore.Client) *firestore.CollectionRef {
	userID := t.claimsProvider.GetUserID()
	return client.Collection("users").Doc(*userID).Collection("dashboards")
}

func (t *dashboard) GetLatestClosedDashboard() (*models.Dashboard, error) {
	client := t.provider.GetClient()
	ctx := context.Background()

	iter := t.dashboardsRef(client).Where("state", "==", "closed").OrderBy("date", firestore.Desc).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var model models.Dashboard
	if err := doc.DataTo(&model); err != nil {
		iter.Stop()
		return nil, err
	}
	model.DashboardID = doc.Ref.ID
	iter.Stop()
	actual, err := t.getActual(ctx, client, model.DashboardID)
	if err != nil {
		return nil, err
	}
	model.Actual = *actual
	return &model, nil
}
func (t *dashboard) getActual(
	ctx context.Context,
	client *firestore.Client,
	documentID string,
) (*[]models.Actual, error) {
	iter := t.dashboardsRef(client).Doc(documentID).Collection("actual").Documents(ctx)
	slice := make([]models.Actual, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var model models.Actual
		if err := doc.DataTo(&model); err != nil {
			return nil, err
		}
		model.ActualID = doc.Ref.ID
		slice = append(slice, model)
	}
	return &slice, nil
}
func (t *dashboard) GetOldestOpenDashboard() (*models.Dashboard, error) {
	client := t.provider.GetClient()
	ctx := context.Background()

	iter := t.dashboardsRef(client).Where("state", "==", "open").OrderBy("date", firestore.Asc).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var model models.Dashboard
	if err := doc.DataTo(&model); err != nil {
		iter.Stop()
		return nil, err
	}
	model.DashboardID = doc.Ref.ID
	iter.Stop()
	actual, err := t.getActual(ctx, client, model.DashboardID)
	if err != nil {
		return nil, err
	}
	model.Actual = *actual
	return &model, nil
}
func (t *dashboard) GetByMonth(month *time.Time) (*models.Dashboard, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	start := t.clock.GetMonthStartDay(month)

	iter := t.dashboardsRef(client).Where("date", "==", start).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var model models.Dashboard
	if err := doc.DataTo(&model); err != nil {
		return nil, err
	}
	model.DashboardID = doc.Ref.ID
	actual, err := t.getActual(ctx, client, model.DashboardID)
	if err != nil {
		return nil, err
	}
	model.Actual = *actual
	return &model, nil
}
