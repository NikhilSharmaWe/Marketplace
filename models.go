package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Shop struct {
	ID                    string     `bson:"_id,omitempty"`
	Name                  string     `bson:"name"`
	Location              string     `bson:"location"`
	OperationHours        string     `bson:"operation_hours"`
	ServiceableProductsId []string   `bson:"products"`
	Coordinates           [2]float64 `bson:"coordinates"`
}

type Product struct {
	ID          string  `bson:"_id,omitempty"`
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	Price       float64 `bson:"price"`
}

type Inventory struct {
	ID        string `bson:"_id,omitempty"`
	ShopID    string `bson:"shop_id"`
	ProductID string `bson:"product_id"`
	Quantity  int    `bson:"quantity"`
}

type ServiceableProduct struct {
	ID        string `bson:"_id,omitempty"`
	ProductID string `bson:"product_id"`
	ShopID    string `bson:"shop_id"`
}

type ShopsServiceableProducts struct {
	ID       string    `bson:"_id,omitempty"`
	ShopID   string    `bson:"shop_id"`
	Products []Product `bson:"products"`
}

type User struct {
	ID          string     `bson:"_id,omitempty"`
	Name        string     `bson:"name"`
	Location    string     `bson:"location"`
	Coordinates [2]float64 `bson:"coordinates"`
}

func CreateMongoCollection(ctx context.Context, name string) *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("market").Collection(name)

	return collection
}
