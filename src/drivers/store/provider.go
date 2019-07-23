package store

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/wakuwaku3/account-book.api/src/application"
	"google.golang.org/api/option"
)

type (
	provider struct {
		env    application.Env
		app    *firebase.App
		client *firestore.Client
	}
	// Provider はstoreへのアクセサです
	Provider interface {
		Initialize() error
		GetClient() *firestore.Client
	}
)

// NewProvider はProviderインスタンスを生成します
func NewProvider(env application.Env) Provider {
	return &provider{env: env}
}
func (provider *provider) Initialize() error {
	ctx := context.Background()
	app, err := provider.createApp(ctx)
	if err != nil {
		return err
	}
	provider.app = app
	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	provider.client = client
	return nil
}

func (provider *provider) createApp(ctx context.Context) (*firebase.App, error) {
	credentialsFilePath := provider.env.GetCredentialsFilePath()
	keepaliveOption := option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                30 * time.Millisecond,
		Timeout:             20 * time.Millisecond,
		PermitWithoutStream: true,
	}))
	if *credentialsFilePath != "" {
		sa := option.WithCredentialsFile(*credentialsFilePath)

		return firebase.NewApp(ctx, nil, sa, keepaliveOption)
	}
	return firebase.NewApp(ctx, nil, keepaliveOption)
}
func (provider *provider) GetClient() *firestore.Client {
	return provider.client
}
