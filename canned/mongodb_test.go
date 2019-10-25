package canned

import (
	"context"
	"testing"

	testcontainers "github.com/testcontainers/testcontainers-go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertDocument(t *testing.T) {
	ctx := context.Background()

	c, err := NewMongoDBContainer(ctx, MongoDBContainerRequest{})
	if err != nil {
		t.Fatal(err.Error())
	}

	defer c.Container.Terminate(ctx)

	mongoClient, err := c.GetClient(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	collection := mongoClient.Database("testdatabase").Collection("persons")

	result, err := collection.InsertOne(ctx, bson.D{primitive.E{Key: "name", Value: "John Doe"}})
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.InsertedID == nil {
		t.Fatal("Insert failed")
	}

	mongoClient.Disconnect(ctx)

	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestInsertDocumentWithMongoDBContainerRequestParameters(t *testing.T) {
	ctx := context.Background()

	testDbName := "testdb"

	c, err := NewMongoDBContainer(ctx, MongoDBContainerRequest{
		Database: testDbName,
		User:     "top",
		Password: "secret",
	})

	defer c.Container.Terminate(ctx)

	mongoClient, err := c.GetClient(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	collection := mongoClient.Database(testDbName).Collection("persons")

	result, err := collection.InsertOne(ctx, bson.D{primitive.E{Key: "name", Value: "John Doe"}})
	if err != nil {
		t.Fatal(err.Error())
	}

	if result.InsertedID == nil {
		t.Fatal("Insert failed")
	}

	mongoClient.Disconnect(ctx)

	if err != nil {
		t.Fatal(err.Error())
	}
}

func ExampleMongoDBContainerRequest() {
	// Optional
	containerRequest := testcontainers.ContainerRequest{
		Image: "docker.io/library/mongo:4.2.0",
	}

	genericContainerRequest := testcontainers.GenericContainerRequest{
		Started:          true,
		ContainerRequest: containerRequest,
	}

	// Database, User, and Password are optional,
	// the driver will use default ones if not provided
	mongoContainerRequest := MongoDBContainerRequest{
		User:                    "anyuser",
		Password:                "yoursecurepassword",
		Database:                "mycustomdatabase",
		GenericContainerRequest: genericContainerRequest,
	}

	mongoContainerRequest.Validate()
}

func ExampleNewMongoDBContainer() {
	ctx := context.Background()

	c, _ := NewMongoDBContainer(ctx, MongoDBContainerRequest{
		GenericContainerRequest: testcontainers.GenericContainerRequest{
			Started: true,
		},
	})

	defer c.Container.Terminate(ctx)
}

func ExampleMongoDBContainer_GetClient() {
	ctx := context.Background()

	c, _ := NewMongoDBContainer(ctx, MongoDBContainerRequest{
		GenericContainerRequest: testcontainers.GenericContainerRequest{
			Started: true,
		},
	})

	mongoClient, _ := c.GetClient(ctx)

	mongoClient.Ping(ctx, nil)
}
