package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/NikhilSharmaWe/marketplace/proto"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func makeGRPCServerAndRun(listenAddr string, svc application) error {
	grpcPriceFetcher := NewGRPCPriceFetcherServer(svc)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	proto.RegisterMarketplaceServiceServer(server, grpcPriceFetcher)

	return server.Serve(ln)
}

type GRPCMarketPlaceServer struct {
	svc application
	proto.UnimplementedMarketplaceServiceServer
}

func NewGRPCPriceFetcherServer(svc application) *GRPCMarketPlaceServer {
	return &GRPCMarketPlaceServer{
		svc: svc,
	}
}

func (s *GRPCMarketPlaceServer) CreateShop(ctx context.Context, req *proto.CreateShopRequest) (*proto.Shop, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())

	shop := Shop{
		ID:             primitive.NewObjectID().Hex(),
		Name:           req.Name,
		Location:       req.Location,
		OperationHours: req.Operationhours,
		Coordinates: [2]float64{
			req.Coordinates.Latitude,
			req.Coordinates.Longitude,
		},
	}

	err := s.svc.Create(ctx, s.svc.shopCollection, shop)
	if err != nil {
		log.Println("fail to creating new shop:", err)
		return nil, err
	}

	return &proto.Shop{
		Id:             shop.ID,
		Name:           shop.Name,
		Location:       shop.Location,
		Operationhours: shop.OperationHours,
		Coordinates: &proto.Coordinates{
			Latitude:  shop.Coordinates[0],
			Longitude: shop.Coordinates[1],
		},
	}, nil
}

func (s *GRPCMarketPlaceServer) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.Product, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())

	product := Product{
		ID:          primitive.NewObjectID().Hex(),
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
	}
	err := s.svc.Create(ctx, s.svc.productCollection, product)
	if err != nil {
		log.Println("fail to creating new product:", err)
		return nil, err
	}

	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       float32(product.Price),
	}, nil
}

func (s *GRPCMarketPlaceServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())

	user := User{
		ID:       primitive.NewObjectID().Hex(),
		Name:     req.Name,
		Location: req.Location,
		Coordinates: [2]float64{
			req.Coordinates.Latitude,
			req.Coordinates.Longitude,
		},
	}

	err := s.svc.Create(ctx, s.svc.userCollection, user)
	if err != nil {
		log.Println("fail to creating new user:", err)
		return nil, err
	}

	return &proto.User{
		Id:       user.ID,
		Name:     user.Name,
		Location: user.Location,
		Coordinates: &proto.Coordinates{
			Latitude:  user.Coordinates[0],
			Longitude: user.Coordinates[1],
		},
	}, nil
}

func (s *GRPCMarketPlaceServer) GetShopByID(ctx context.Context, req *proto.GetRequest) (*proto.Shop, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	id := req.Id
	shop := Shop{}

	err := s.svc.GetByID(id, s.svc.shopCollection, &shop)
	if err != nil {
		log.Printf("fail to get shop with id [%s]: %s", req.Id, err)
		return nil, err
	}

	products, err := s.GetServiceableProducts(ctx, &proto.GetServiceableProductsRequest{ShopId: id})
	if err != nil {
		return nil, err
	}

	return &proto.Shop{
		Id:                  shop.ID,
		Name:                shop.Name,
		Location:            shop.Location,
		Operationhours:      shop.OperationHours,
		ServiceableProducts: products.Products,
		Coordinates: &proto.Coordinates{
			Latitude:  shop.Coordinates[0],
			Longitude: shop.Coordinates[1],
		},
	}, nil
}

func (s *GRPCMarketPlaceServer) GetProductByID(ctx context.Context, req *proto.GetRequest) (*proto.Product, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	id := req.Id
	product := Product{}

	err := s.svc.GetByID(id, s.svc.productCollection, &product)
	if err != nil {
		log.Printf("fail to get product with id [%s]: %s", req.Id, err)
		return nil, err
	}

	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       float32(product.Price),
	}, nil
}

func (s *GRPCMarketPlaceServer) GetUserByID(ctx context.Context, req *proto.GetRequest) (*proto.User, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	id := req.Id
	user := User{}

	err := s.svc.GetByID(id, s.svc.userCollection, &user)
	if err != nil {
		log.Printf("fail to get user with id [%s]: %s", req.Id, err)
		return nil, err
	}

	return &proto.User{
		Id:       user.ID,
		Name:     user.Name,
		Location: user.Location,
		Coordinates: &proto.Coordinates{
			Latitude:  user.Coordinates[0],
			Longitude: user.Coordinates[1],
		},
	}, nil
}

func (s *GRPCMarketPlaceServer) AddServiceableProduct(ctx context.Context, req *proto.AddServiceableProductRequest) (*proto.Shop, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	shopId := req.ShopId
	productId := req.ProductId
	shop := Shop{}

	err := s.svc.GetByID(shopId, s.svc.shopCollection, &shop)
	if err != nil {
		return nil, err
	}

	var alreadyExists bool
	for _, id := range shop.ServiceableProductsId {
		if id == productId {
			alreadyExists = true
		}
	}

	if !alreadyExists {
		products := append(shop.ServiceableProductsId, productId)

		err = s.svc.UpdateServiceableProductsForShop(ctx, shopId, products)
		if err != nil {
			log.Printf("fail to update serviceable products for shop with id [%s] by adding product with id[%s]: %s", req.ShopId, req.ProductId, err)
			return nil, err
		}
	}

	_, err = s.GetInventory(ctx, &proto.GetInventoryRequest{
		ShopId:    shopId,
		ProductId: productId,
	})

	if err == mongo.ErrNoDocuments {
		inventory := Inventory{
			ID:        primitive.NewObjectID().Hex(),
			ShopID:    shopId,
			ProductID: productId,
			Quantity:  0,
		}

		err = s.svc.Create(ctx, s.svc.inventoryCollection, inventory)
		if err != nil {
			return nil, err
		}
	}

	return s.GetShopByID(ctx, &proto.GetRequest{
		Id: shopId,
	})
}

func (s *GRPCMarketPlaceServer) GetServiceableProducts(ctx context.Context, req *proto.GetServiceableProductsRequest) (*proto.Products, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	id := req.ShopId
	shop := Shop{}

	err := s.svc.GetByID(id, s.svc.shopCollection, &shop)
	if err != nil {
		log.Printf("fail to get shop with id [%s]: %s", req.ShopId, err)
		return nil, err
	}

	var serviceableProducts []*proto.Product

	for _, productId := range shop.ServiceableProductsId {
		product, err := s.GetProductByID(ctx, &proto.GetRequest{
			Id: productId,
		})

		if err != nil {
			log.Printf("fail to get product with id [%s]: %s", productId, err)
			return nil, err
		}
		serviceableProducts = append(serviceableProducts, product)
	}

	return &proto.Products{
		Products: serviceableProducts,
	}, nil
}

func (s *GRPCMarketPlaceServer) GetInventory(ctx context.Context, req *proto.GetInventoryRequest) (*proto.Inventory, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	shopId := req.ShopId
	productId := req.ProductId

	inventory := Inventory{}
	err := s.svc.GetByShopAndProductID(shopId, productId, s.svc.inventoryCollection, &inventory)
	if err != nil {
		log.Printf("fail to get inventory of shop with id [%s] for product with id [%s]: %s", shopId, productId, err)
		return nil, err
	}

	return &proto.Inventory{
		Id:        inventory.ID,
		ShopId:    inventory.ShopID,
		ProductId: inventory.ProductID,
		Quantity:  int32(inventory.Quantity),
	}, nil
}

func (s *GRPCMarketPlaceServer) UpdateInventory(ctx context.Context, req *proto.UpdateInventoryRequest) (*proto.Inventory, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	shopId := req.ShopId
	productId := req.ProductId

	inventory := Inventory{}
	err := s.svc.GetByShopAndProductID(shopId, productId, s.svc.inventoryCollection, &inventory)
	if err != nil {
		log.Printf("fail to get inventory of shop with id [%s] for product with id [%s]: %s", shopId, productId, err)
		return nil, err
	}

	updatedQuantity := inventory.Quantity
	if req.Add {
		updatedQuantity += int(req.Change)

	} else {
		updatedQuantity -= int(req.Change)
	}

	err = s.svc.UpdateInventory(ctx, inventory.ShopID, inventory.ProductID, updatedQuantity)
	if err != nil {
		log.Printf("fail to update inventory of shop with id [%s] for product with id [%s]: %s", shopId, productId, err)
		return nil, err
	}

	return &proto.Inventory{
		Id:        inventory.ID,
		ShopId:    inventory.ShopID,
		ProductId: inventory.ProductID,
		Quantity:  int32(updatedQuantity),
	}, nil
}

func (s *GRPCMarketPlaceServer) GetShopsByServiceableProducts(ctx context.Context, req *proto.GetShopsByServiceableProductsRequest) (*proto.Shops, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	productId := req.ProductId

	shops, err := s.svc.GetShopsWithProduct(ctx, productId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("No shops with this product")
		} else {
			return nil, err
		}
	}

	resultShops := proto.Shops{}

	for _, shop := range shops {
		pshop, err := s.ParseShop(ctx, shop)
		if err != nil {
			log.Printf("fail to parse shop with id [%s]: %s", shop.ID, err)
			return nil, err
		}

		resultShops.Shops = append(resultShops.Shops, pshop)
	}

	return &proto.Shops{
		Shops: resultShops.Shops,
	}, nil
}

func (s *GRPCMarketPlaceServer) GetShopForUser(ctx context.Context, req *proto.GetShopForUserRequest) (*proto.Shops, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	userId := req.UserId
	user := User{}
	resultShops := proto.Shops{}

	err := s.svc.GetByID(userId, s.svc.userCollection, &user)
	if err != nil {
		log.Printf("fail to get user with id [%s]: %s", userId, err)
		return nil, err
	}

	shops, err := s.svc.GetAllShops(ctx)
	if err != nil {
		return nil, err
	}

	for _, shop := range shops {
		if calculateDistance(user.Coordinates, shop.Coordinates) < req.MaxDistanceInKM {
			pshop, err := s.ParseShop(ctx, shop)
			if err != nil {
				log.Printf("fail to parse shop with id [%s]: %s", shop.ID, err)
				return nil, err
			}

			resultShops.Shops = append(resultShops.Shops, pshop)
		}
	}

	return &proto.Shops{
		Shops: resultShops.Shops,
	}, nil
}

func (s *GRPCMarketPlaceServer) GetNearestNeighbour(ctx context.Context, req *proto.GetNearestNeighbourRequest) (*proto.User, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	userId := req.UserId
	user := User{}

	err := s.svc.GetByID(userId, s.svc.userCollection, &user)
	if err != nil {
		log.Printf("fail to get user with id [%s]: %s", userId, err)
		return nil, err
	}

	users, err := s.svc.GetAllUser(ctx)
	if err != nil {
		log.Println("fail to get all users", err)
		return nil, err
	}

	nearest := users[0]
	nearestDistance := calculateDistance(user.Coordinates, nearest.Coordinates)
	for _, kuser := range users {
		distance := calculateDistance(user.Coordinates, kuser.Coordinates)
		if distance < nearestDistance && kuser.ID != user.ID {
			nearest = kuser
			nearestDistance = distance
		}
	}

	return &proto.User{
		Id:       nearest.ID,
		Name:     nearest.Name,
		Location: nearest.Location,
		Coordinates: &proto.Coordinates{
			Latitude:  nearest.Coordinates[0],
			Longitude: nearest.Coordinates[1],
		},
	}, nil
}

func (s *GRPCMarketPlaceServer) ParseShop(ctx context.Context, shop *Shop) (*proto.Shop, error) {
	products, err := s.GetServiceableProducts(ctx, &proto.GetServiceableProductsRequest{
		ShopId: shop.ID,
	})
	if err != nil {
		return nil, err
	}

	return &proto.Shop{
		Id:                  shop.ID,
		Name:                shop.Name,
		Location:            shop.Location,
		Operationhours:      shop.OperationHours,
		ServiceableProducts: products.Products,
		Coordinates: &proto.Coordinates{
			Latitude:  shop.Coordinates[0],
			Longitude: shop.Coordinates[1],
		},
	}, nil
}
