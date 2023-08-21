package main

import (
	"github.com/NikhilSharmaWe/marketplace/proto"
	"google.golang.org/grpc"
)

func main() {
	// ctx := context.Background()

	// c, err := NewGRPCClient(":4000")
	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// For testing

	// ****** First create 3 shops
	// shop, err := c.CreateShop(ctx, &proto.CreateShopRequest{
	// 	Name:           "VS Computers",
	// 	Location:       "Meerut",
	// 	Operationhours: "6",
	// 	Coordinates: &proto.Coordinates{
	// 		Latitude:  34.0522,
	// 		Longitude: -118.2437,
	// 	},
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", shop)

	// time.Sleep(time.Second * 3)

	// shop, err = c.CreateShop(ctx, &proto.CreateShopRequest{
	// 	Name:           "Kamal Computers",
	// 	Location:       "Delhi",
	// 	Operationhours: "6",
	// 	Coordinates: &proto.Coordinates{
	// 		Latitude:  40.7128,
	// 		Longitude: -74.0060,
	// 	},
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", shop)

	// time.Sleep(time.Second * 3)

	// shop, err = c.CreateShop(ctx, &proto.CreateShopRequest{
	// 	Name:           "Nikhil Computers",
	// 	Location:       "Bangalore",
	// 	Operationhours: "6",
	// 	Coordinates: &proto.Coordinates{
	// 		Latitude:  51.5074,
	// 		Longitude: -0.1278,
	// 	},
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", shop)

	// time.Sleep(time.Second * 3)

	// ****** Then create 3 users
	// user, err := c.CreateUser(ctx, &proto.CreateUserRequest{
	// 	Name:     "Nikhil Sharma",
	// 	Location: "Gujarat",
	// 	Coordinates: &proto.Coordinates{
	// 		Latitude:  40.7125,
	// 		Longitude: -73.9975,
	// 	},
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", user)

	// user, err = c.CreateUser(ctx, &proto.CreateUserRequest{
	// 	Name:     "Rewak Tyagi",
	// 	Location: "Mumbai",
	// 	Coordinates: &proto.Coordinates{
	// 		Latitude:  40.705,
	// 		Longitude: -75.9975,
	// 	},
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", user)

	// user, err = c.CreateUser(ctx, &proto.CreateUserRequest{
	// 	Name:     "Prince Adhana",
	// 	Location: "Gujarat",
	// 	Coordinates: &proto.Coordinates{
	// 		Latitude:  38.7125,
	// 		Longitude: -70.9975,
	// 	},
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", user)

	// ****** Then search for the nearest neighbour for Nikhil Sharma
	// user, err := c.GetNearestNeighbour(ctx, &proto.GetNearestNeighbourRequest{
	// 	UserId: "64e2d62863453733885e3797",
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", user)

	// ****** Create a new product
	// product, err := c.CreateProduct(ctx, &proto.CreateProductRequest{
	// 	Name:        "Keyboard",
	// 	Description: "Used for typing",
	// 	Price:       1000,
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", product)

	// ****** Add this new product to one of the shop serviceable products
	// shop, err := c.AddServiceableProduct(ctx, &proto.AddServiceableProductRequest{
	// 	ShopId:    "64e2d5e963453733885e378d",
	// 	ProductId: "64e2d6bc63453733885e3799",
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", shop)

	// time.Sleep(time.Second * 3)

	// ****** Get the serviceable products for the shop
	// products, err := c.GetServiceableProducts(ctx, &proto.GetServiceableProductsRequest{
	// 	ShopId: "64e2d5e963453733885e378d",
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", products)

	// time.Sleep(time.Second * 3)

	// ****** Get the shops for users based on a max distance
	// shop, err := c.GetShopForUser(ctx, &proto.GetShopForUserRequest{
	// 	UserId:          "64e2d62863453733885e3793",
	// 	MaxDistanceInKM: 1000,
	// })

	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", shop)

	// ****** UpdateInventory by changing the quantity of serviceable products of a shop
	// inventory, err := c.UpdateInventory(ctx, &proto.UpdateInventoryRequest{
	// 	ShopId:    "64e2d5e963453733885e378d",
	// 	ProductId: "64e2d6bc63453733885e3799",
	// 	Change:    2,
	// 	Add:       true,
	// })
	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", inventory)

	// ****** Get the changes inventory
	// inventory, err := c.GetInventory(ctx, &proto.GetInventoryRequest{
	// 	ShopId:    "64e0aafaf0d439c6dabf6b46",
	// 	ProductId: "64e0ab9e84d19bab82fa6fcf",
	// })
	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", inventory)

	// ****** Get shops which provide a particular product
	// shops, err := c.GetShopsByServiceableProducts(ctx, &proto.GetShopsByServiceableProductsRequest{
	// 	ProductId: "64e2d6bc63453733885e3799",
	// })
	// if err != nil {
	// 	log.Println("Error:", err)
	// }

	// fmt.Printf("%+v\n", shops)
}

func NewGRPCClient(remoteAddr string) (proto.MarketplaceServiceClient, error) {

	conn, err := grpc.Dial(remoteAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := proto.NewMarketplaceServiceClient(conn)

	return c, nil
}
