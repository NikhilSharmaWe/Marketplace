package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/NikhilSharmaWe/marketplace/proto"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"github.com/SaiNageswarS/go-api-boot/server"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

var (
	grpcAddr *string
	webAddr  *string
)

type application struct {
	ctx             context.Context
	userRepo        UserRepository
	shopRepo        ShopRepository
	productRepo     ProductRepository
	inventoryRepo   InventoryRepository
	serviceableRepo ServiceableProductRepository
	goApiBoot       *server.GoApiBoot
	grpcClient      proto.MarketplaceServiceClient
	logger          *log.Logger
	mongoClient     *mongo.Client
}

type UserRepository struct {
	odm.AbstractRepository[User]
}

type ProductRepository struct {
	odm.AbstractRepository[Product]
}

type ShopRepository struct {
	odm.AbstractRepository[Shop]
}

type InventoryRepository struct {
	odm.AbstractRepository[Inventory]
}

type ServiceableProductRepository struct {
	odm.AbstractRepository[ServiceableProduct]
}

func main() {
	grpcAddr = flag.String("grpc", ":4000", "listen address of the grpc transport")
	webAddr = flag.String("web", ":3000", "listen address of the web transport")
	flag.Parse()

	err := godotenv.Load("vars.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := newApplication()
	app.start()
}

func newApplication() *application {
	ctx := context.Background()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	userRepo := &UserRepository{
		AbstractRepository: odm.AbstractRepository[User]{
			Database:       "market",
			CollectionName: "user",
		},
	}

	productRepo := &ProductRepository{
		AbstractRepository: odm.AbstractRepository[Product]{
			Database:       "market",
			CollectionName: "product",
		},
	}

	shopRepo := &ShopRepository{
		AbstractRepository: odm.AbstractRepository[Shop]{
			Database:       "market",
			CollectionName: "shop",
		},
	}

	inventoryRepo := &InventoryRepository{
		AbstractRepository: odm.AbstractRepository[Inventory]{
			Database:       "market",
			CollectionName: "inventory",
		},
	}

	serviceableProductRepo := &ServiceableProductRepository{
		AbstractRepository: odm.AbstractRepository[ServiceableProduct]{
			Database:       "market",
			CollectionName: "serviceableProduct",
		},
	}

	grpcClient, err := newGRPCClient(*grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient := odm.GetClient()

	goApiBoot := server.NewGoApiBoot()

	return &application{
		ctx:             ctx,
		userRepo:        *userRepo,
		shopRepo:        *shopRepo,
		productRepo:     *productRepo,
		inventoryRepo:   *inventoryRepo,
		serviceableRepo: *serviceableProductRepo,
		goApiBoot:       goApiBoot,
		logger:          logger,
		grpcClient:      grpcClient,
		mongoClient:     mongoClient,
	}
}

func (app *application) start() {
	app.setupGoApiBoot()
	app.goApiBoot.Start(*grpcAddr, *webAddr)
}

func (app *application) setupGoApiBoot() {
	app.setupGRPCServer()
	app.setupWebserver()
}

func newGRPCClient(remoteAddr string) (proto.MarketplaceServiceClient, error) {
	conn, err := grpc.Dial(remoteAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := proto.NewMarketplaceServiceClient(conn)
	return c, nil
}
