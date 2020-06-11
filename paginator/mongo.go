package paginator

import (
	"context"
	"time"

	log "github.com/victron/simpleLogger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClient struct {
	client *mongo.Client
	ctx    *context.Context
}

func (mclient *mongoClient) Connect(url string) error {
	var err error
	(*mclient).client, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Error.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = (*mclient).client.Connect(ctx)
	if err != nil {
		log.Error.Fatal(err)
	}
	(*mclient).ctx = &ctx
	return nil
}

func (mclient *mongoClient) Close() error {
	if err := (*mclient).client.Disconnect(*(*mclient).ctx); err != nil {
		log.Error.Fatal(err)
	}
	return nil
}
