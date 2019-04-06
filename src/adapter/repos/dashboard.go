package repos

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
	"github.com/wakuwaku3/account-book.api/src/enterprise/helpers"
	"github.com/wakuwaku3/account-book.api/src/drivers/store"
	"google.golang.org/api/iterator"
)

type dashboard struct {
	provider       store.Provider
	claimsProvider application.ClaimsProvider
	clock          helpers.Clock
}

// NewDashboard はインスタンスを生成します
func NewDashboard(
	provider store.Provider,
	claimsProvider application.ClaimsProvider,
	clock helpers.Clock,
) application.DashboardRepository {
	return &dashboard{provider, claimsProvider, clock}
}

func (t *dashboard) dashboardsRef(client *firestore.Client) *firestore.CollectionRef {
	userID := t.claimsProvider.GetUserID()
	return client.Collection("users").Doc(*userID).Collection("dashboards")
}

func (t *dashboard) actualRef(client *firestore.Client, dashboardID *string) *firestore.CollectionRef {
	return t.dashboardsRef(client).Doc(*dashboardID).Collection("actual")
}

func (t *dashboard) dailyRef(client *firestore.Client, dashboardID *string) *firestore.CollectionRef {
	return t.dashboardsRef(client).Doc(*dashboardID).Collection("daily")
}

func (t *dashboard) GetByID(id *string) (*models.Dashboard, error) {
	client := t.provider.GetClient()
	ctx := context.Background()

	doc, err := t.dashboardsRef(client).Doc(*id).Get(ctx)
	if err != nil {
		if !doc.Exists() {
			return nil, nil
		}
		return nil, err
	}
	var model models.Dashboard
	doc.DataTo(&model)
	model.Date = model.Date.In(t.clock.DefaultLocation())
	model.DashboardID = doc.Ref.ID
	actual, err := t.getActual(ctx, client, model.DashboardID)
	if err != nil {
		return nil, err
	}
	model.Actual = *actual
	return &model, nil
}
func (t *dashboard) ExistsClosedNext(id *string) error {
	client := t.provider.GetClient()
	ctx := context.Background()

	iter := t.dashboardsRef(client).Where("previousDashboardId", "==", *id).Where("state", "==", "closed").Documents(ctx)
	_, err := iter.Next()
	if err == iterator.Done {
		return nil
	}
	if err != nil {
		return err
	}
	iter.Stop()
	return errors.New("next dashboard is already closed")
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
	model.Date = model.Date.In(t.clock.DefaultLocation())
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
	dashboardID string,
) (*[]models.Actual, error) {
	iter := t.actualRef(client, &dashboardID).OrderBy("planCreatedAt", firestore.Asc).Documents(ctx)
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
func (t *dashboard) getDaily(
	ctx context.Context,
	client *firestore.Client,
	dashboardID string,
) (*[]models.Daily, error) {
	iter := t.dailyRef(client, &dashboardID).OrderBy("date", firestore.Asc).Documents(ctx)
	slice := make([]models.Daily, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var model models.Daily
		if err := doc.DataTo(&model); err != nil {
			return nil, err
		}
		model.DailyID = doc.Ref.ID
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
	model.Date = model.Date.In(t.clock.DefaultLocation())
	model.DashboardID = doc.Ref.ID
	iter.Stop()
	actual, err := t.getActual(ctx, client, model.DashboardID)
	if err != nil {
		return nil, err
	}
	model.Actual = *actual
	daily, err := t.getDaily(ctx, client, model.DashboardID)
	if err != nil {
		return nil, err
	}
	model.Daily = *daily
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
	model.Date = model.Date.In(t.clock.DefaultLocation())
	model.DashboardID = doc.Ref.ID
	actual, err := t.getActual(ctx, client, model.DashboardID)
	if err != nil {
		return nil, err
	}
	model.Actual = *actual
	daily, err := t.getDaily(ctx, client, model.DashboardID)
	if err != nil {
		return nil, err
	}
	model.Daily = *daily
	return &model, nil
}
func (t *dashboard) Create(month *time.Time) (*string, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	start := t.clock.GetMonthStartDay(month)

	iter := t.dashboardsRef(client).Where("date", "==", start).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		doc, _, err := t.dashboardsRef(client).Add(ctx, models.Dashboard{
			Date:  start,
			State: "open",
		})
		if err != nil {
			return nil, err
		}
		return &doc.ID, nil
	}
	if err != nil {
		return nil, err
	}
	return &doc.Ref.ID, nil
}
func (t *dashboard) Approve(model *models.Dashboard) error {
	client := t.provider.GetClient()
	ctx := context.Background()
	batch := client.Batch()

	ref := t.dashboardsRef(client).Doc(model.DashboardID)
	batch.Set(ref, model)

	for _, daily := range model.Daily {
		newRef := t.dailyRef(client, &model.DashboardID).NewDoc()
		batch.Create(newRef, daily)
	}

	_, err := batch.Commit(ctx)
	return err
}
func (t *dashboard) CancelApprove(model *models.Dashboard) error {
	client := t.provider.GetClient()
	ctx := context.Background()
	batch := client.Batch()

	ref := t.dashboardsRef(client).Doc(model.DashboardID)
	batch.Set(ref, model)

	iter := t.dailyRef(client, &model.DashboardID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		batch.Delete(doc.Ref)
	}

	_, err := batch.Commit(ctx)
	return err
}
func (t *dashboard) GetActual(dashboardID *string, id *string) (*models.Actual, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := t.actualRef(client, dashboardID).Doc(*id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var model models.Actual
	doc.DataTo(&model)
	model.ActualID = doc.Ref.ID
	return &model, nil
}
func (t *dashboard) ExistsActual(dashboardID *string, planID *string) (*string, error) {
	client := t.provider.GetClient()
	ctx := context.Background()

	iter := t.actualRef(client, dashboardID).Where("planId", "==", *planID).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &doc.Ref.ID, nil
}
func (t *dashboard) CreateActual(dashboardID *string, model *models.Actual) (*string, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	ref, _, err := t.actualRef(client, dashboardID).Add(ctx, model)
	if err != nil {
		return nil, err
	}
	return &ref.ID, nil
}
func (t *dashboard) UpdateActual(dashboardID *string, id *string, model *models.Actual) error {
	client := t.provider.GetClient()
	ctx := context.Background()
	_, err := t.actualRef(client, dashboardID).Doc(*id).Set(ctx, model)
	if err != nil {
		return err
	}
	return nil
}
