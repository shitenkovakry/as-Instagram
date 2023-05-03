package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	timeout = time.Duration(time.Second * 30)
)

func connect(dbHosts []string, username, password string) (*mongo.Client, error) {
	opts := options.Client()

	if username != "" && password != "" {
		opts = opts.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	client, err := mongo.NewClient(opts.
		SetHosts(dbHosts),
	)

	if err != nil {
		return nil, errors.Wrap(err,
			`client, err := mongo.NewClient(opts.ApplyURI("mongodb://" + dbURL).SetConnectTimeout(connectTimeout))`,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "err := client.Connect(ctx)")
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, "err := client.Ping(ctx, nil)")
	}

	return client, nil
}
