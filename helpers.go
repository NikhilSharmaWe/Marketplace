package main

import (
	"context"
	"log"
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

func (app *application) UpdateServiceableProductsForShop(ctx context.Context, shopId string, updatedList []Product) error {
	filter := bson.D{primitive.E{Key: "_id", Value: shopId}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "products", Value: updatedList},
	}}}

	shop := Shop{}

	err := app.shopCollection.FindOneAndUpdate(ctx, filter, update).Decode(&shop)
	if err != nil {
		return err
	}

	log.Printf("insertResult: %+v\n", err)
	return nil
}

// func getChatRoom(key string) (model.ChatRoom, error) {
// 	chatRoom := model.ChatRoom{}
// 	filter := bson.M{"key": key}

// 	err := chatroomCollection.FindOne(context.Background(), filter).Decode(&chatRoom)
// 	if err != nil {
// 		return chatRoom, err
// 	}

// 	return chatRoom, nil
// }

// func getChats(crKey string) ([]*model.Chat, error) {
// 	var chats []*model.Chat
// 	filter := bson.M{"key": crKey}
// 	cur, err := chatCollection.Find(ctx, filter)
// 	if err != nil {
// 		return chats, err
// 	}

// 	for cur.Next(ctx) {
// 		var chat model.Chat
// 		err := cur.Decode(&chat)
// 		if err != nil {
// 			return chats, err
// 		}
// 		chats = append(chats, &chat)
// 	}

// 	if err := cur.Err(); err != nil {
// 		return chats, nil
// 	}

// 	cur.Close(ctx)
// 	if len(chats) == 0 {
// 		return chats, mongo.ErrNoDocuments
// 	}

// 	return chats, nil
// }
