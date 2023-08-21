package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) Create(ctx context.Context, collection *mongo.Collection, value interface{}) (err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"took":      time.Since(begin),
			"err":       err,
		}).Info("create")
	}(time.Now())

	insertResutlt, err := collection.InsertOne(app.ctx, value)
	if err != nil {
		return err
	}

	log.Printf("insertResult: %+v\n", insertResutlt)

	return nil
}

func (app *application) GetByID(id string, collection *mongo.Collection, result interface{}) error {
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) GetByShopAndProductID(shopId, productId string, collection *mongo.Collection, result interface{}) error {
	filter := bson.M{"shop_id": shopId, "product_id": productId}
	err := collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) UpdateInventory(ctx context.Context, shopId, productId string, newQuantity int) error {
	filter := bson.D{
		primitive.E{Key: "shop_id", Value: shopId},
		primitive.E{Key: "product_id", Value: productId},
	}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "quantity", Value: newQuantity},
	}}}

	inventory := Inventory{}

	err := app.inventoryCollection.FindOneAndUpdate(ctx, filter, update).Decode(&inventory)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) UpdateServiceableProductsForShop(ctx context.Context, shopId string, updatedList []string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: shopId}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "products", Value: updatedList},
	}}}

	shop := Shop{}

	err := app.shopCollection.FindOneAndUpdate(ctx, filter, update).Decode(&shop)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) GetAllUser(ctx context.Context) ([]*User, error) {
	var users []*User
	cur, err := app.userCollection.Find(ctx, bson.M{})
	if err != nil {
		return users, err
	}

	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return users, err
		}

		users = append(users, &user)
	}

	if err := cur.Err(); err != nil {
		return users, nil
	}

	cur.Close(ctx)
	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	return users, nil
}

func (app *application) GetAllShops(ctx context.Context) ([]*Shop, error) {
	var shops []*Shop
	cur, err := app.shopCollection.Find(ctx, bson.M{})
	if err != nil {
		return shops, err
	}

	for cur.Next(ctx) {
		var shop Shop
		err := cur.Decode(&shop)
		if err != nil {
			return shops, err
		}

		shops = append(shops, &shop)
	}

	if err := cur.Err(); err != nil {
		return shops, nil
	}

	cur.Close(ctx)
	if len(shops) == 0 {
		return shops, mongo.ErrNoDocuments
	}

	return shops, nil
}

func (app *application) GetShopsWithProduct(ctx context.Context, productId string) ([]*Shop, error) {
	var shops []*Shop
	cur, err := app.shopCollection.Find(ctx, bson.M{})
	if err != nil {
		return shops, err
	}

	for cur.Next(ctx) {
		var shop Shop
		err := cur.Decode(&shop)
		if err != nil {
			return shops, err
		}

		for _, id := range shop.ServiceableProductsId {
			if id == productId {
				shops = append(shops, &shop)
			}
		}
	}

	if err := cur.Err(); err != nil {
		return shops, nil
	}

	cur.Close(ctx)
	if len(shops) == 0 {
		return shops, mongo.ErrNoDocuments
	}

	return shops, nil
}

func calculateDistance(coord1, coord2 [2]float64) float64 {
	const earthRadiusKm = 6371

	lat1 := degToRad(coord1[0])
	lon1 := degToRad(coord1[1])
	lat2 := degToRad(coord2[0])
	lon2 := degToRad(coord2[1])

	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadiusKm * c
	return distance
}

func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
