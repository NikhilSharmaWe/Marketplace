package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Shop represents a shop entity.
type Shop struct {
	ID                  string    `bson:"_id,omitempty"`
	Name                string    `bson:"name"`
	Location            string    `bson:"location"`
	OperationHours      string    `bson:"operation_hours"`
	ServiceableProducts []Product `bson:"products"`
}

// Product represents a product entity in the catalog.
type Product struct {
	ID          string  `bson:"_id,omitempty"`
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	Price       float64 `bson:"price"`
}

// Inventory represents the inventory of a shop for a specific product.
type Inventory struct {
	ID        string `bson:"_id,omitempty"`
	ShopID    string `bson:"shop_id"`
	ProductID string `bson:"product_id"`
	Quantity  int    `bson:"quantity"`
}

// ServiceableProduct represents a product that is serviceable by a shop.
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

// User represents a user entity.
type User struct {
	ID       string `bson:"_id,omitempty"`
	Name     string `bson:"name"`
	Location string `bson:"location"`
}

// Neighbour represents a neighbor user for community finder.
type Neighbour struct {
	ID          string  `bson:"_id,omitempty"`
	UserID      string  `bson:"user_id"`
	NeighbourID string  `bson:"neighbour_id"`
	Distance    float64 `bson:"distance"`
}

func CreateMongoCollection(ctx context.Context, name string) *mongo.Collection {
	// err := godotenv.Load("vars.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// url := os.Getenv("127.0.0.1:27017")
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// collection := client.Database(dbName).Collection(collectionName)
	collection := client.Database("market").Collection(name)

	return collection
}
