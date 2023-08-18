package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/NikhilSharmaWe/marketplace/proto"
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

var (
	ctx                          = context.Background()
	userCollection               = CreateMongoCollection(ctx, "user")
	shopCollection               = CreateMongoCollection(ctx, "shop")
	productCollection            = CreateMongoCollection(ctx, "product")
	inventoryCollection          = CreateMongoCollection(ctx, "inventory")
	serviceableProductCollection = CreateMongoCollection(ctx, "serviceableProduct")
	neighbourCollection          = CreateMongoCollection(ctx, "neighbour")
)

func main() {

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

	c := NewGRPCPriceFetcherServer(app)

	go func() {
		time.Sleep(time.Second * 3)

		// shop, err := c.CreateShop(ctx, &proto.CreateShopRequest{
		// 	Name:           "VS Computers",
		// 	Location:       "Meerut",
		// 	Operationhours: "6",
		// })

		// if err != nil {
		// 	log.Println("Error:", err)
		// }

		// fmt.Printf("%+v\n", shop)

		// time.Sleep(time.Second * 3)

		// product, err := c.CreateProduct(ctx, &proto.CreateProductRequest{
		// 	Name:        "Keyboard",
		// 	Description: "Used for typing",
		// 	Price:       1000,
		// })

		// if err != nil {
		// 	log.Println("Error:", err)
		// }

		// fmt.Printf("%+v\n", product)

		shop, err := c.GetShopByID(ctx, &proto.GetRequest{Id: "64dff6cff058592f09994b4b"})
		if err != nil {
			log.Println("Error:", err)
		}

		fmt.Printf("%+v\n", shop)

		time.Sleep(time.Second * 3)

		product, err := c.GetProductByID(ctx, &proto.GetRequest{Id: "64dff6d2f058592f09994b4d"})
		if err != nil {
			log.Println("Error:", err)
		}

		fmt.Printf("%+v\n", product)

		time.Sleep(time.Second * 3)

		resp, err := c.AddServiceableProduct(ctx, &proto.AddServiceableProductRequest{
			ShopId:    "64dff6cff058592f09994b4b",
			ProductId: "64dff6d2f058592f09994b4d",
		})

		if err != nil {
			log.Println("Error:", err)
		}

		fmt.Printf("%+v\n", resp)

		// time.Sleep(time.Second * 3)

		// shop, err = c.GetShopByID(ctx, &proto.GetRequest{Id: "64dff3125ae6ced6a729a67b"})
		// if err != nil {
		// 	log.Println("Error:", err)
		// }

		// fmt.Printf("%+v\n", shop)
	}()

	makeGRPCServerAndRun(*grpcAddr, app)
}
