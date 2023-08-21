package main

import (
	"context"
	"flag"

	"go.mongodb.org/mongo-driver/mongo"
)

type application struct {
	ctx                          context.Context
	userCollection               *mongo.Collection
	shopCollection               *mongo.Collection
	productCollection            *mongo.Collection
	inventoryCollection          *mongo.Collection
	serviceableProductCollection *mongo.Collection
	neighbourCollection          *mongo.Collection
}

func main() {

	var (
		ctx                          = context.Background()
		userCollection               = CreateMongoCollection(ctx, "user")
		shopCollection               = CreateMongoCollection(ctx, "shop")
		productCollection            = CreateMongoCollection(ctx, "product")
		inventoryCollection          = CreateMongoCollection(ctx, "inventory")
		serviceableProductCollection = CreateMongoCollection(ctx, "serviceableProduct")
		neighbourCollection          = CreateMongoCollection(ctx, "neighbour")
	)

	app := application{
		ctx:                          ctx,
		userCollection:               userCollection,
		shopCollection:               shopCollection,
		productCollection:            productCollection,
		inventoryCollection:          inventoryCollection,
		serviceableProductCollection: serviceableProductCollection,
		neighbourCollection:          neighbourCollection,
	}

	grpcAddr := flag.String("grpc", ":4000", "listen address of the grpc transport")
	flag.Parse()

	makeGRPCServerAndRun(*grpcAddr, app)
}
