package canned

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	mongoUser       = "user"
	mongoPassword   = "password"
	mongoDatabase   = "database"
	mongoImage      = "mongo"
	mongoDefaultTag = "4.2.0"
	mongoPort       = "27017/tcp"
)

// MongoDBContainerRequest represents the parameters for requesting a running MongoDB container
type MongoDBContainerRequest struct {
	testcontainers.GenericContainerRequest
	User     string
	Password string
	Database string
}

// MongoDBContainer should always be created via NewMongoDBContainer
type MongoDBContainer struct {
	Container testcontainers.Container
	client    *mongo.Client
	req       MongoDBContainerRequest
}

// GetClient returns a *mongo.Client.
func (c *MongoDBContainer) GetClient(ctx context.Context) (*mongo.Client, error) {

	host, err := c.Container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := c.Container.MappedPort(ctx, mongoPort)
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		c.req.User,
		c.req.Password,
		host,
		mappedPort.Int(),
	))

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewMongoDBContainer represents the running container instance.
func NewMongoDBContainer(ctx context.Context, req MongoDBContainerRequest) (*MongoDBContainer, error) {

	provider, err := req.ProviderType.GetProvider()
	if err != nil {
		return nil, err
	}

	req.ExposedPorts = []string{mongoPort}
	req.Env = map[string]string{}
	req.Started = true

	if req.Image == "" && req.FromDockerfile.Context == "" {
		req.Image = fmt.Sprintf("%s:%s", mongoImage, mongoDefaultTag)
	}

	if req.User == "" {
		req.User = mongoUser
	}

	if req.Password == "" {
		req.Password = mongoPassword
	}

	if req.Database == "" {
		req.Database = mongoDatabase
	}

	req.Env["MONGO_INITDB_ROOT_USERNAME"] = req.User
	req.Env["MONGO_INITDB_ROOT_PASSWORD"] = req.Password
	req.Env["MONGO_INITDB_DATABASE"] = req.Database

	req.WaitingFor = wait.ForLog("waiting for connections on port 27017")

	c, err := provider.CreateContainer(ctx, req.ContainerRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create container")
	}

	mongoC := &MongoDBContainer{
		Container: c,
		req:       req,
	}

	if req.Started {
		if err := c.Start(ctx); err != nil {
			return mongoC, errors.Wrap(err, "failed to start container")
		}
	}

	return mongoC, nil

}
