package main

import (
	"context"
	"errors"

	"github.com/NikhilSharmaWe/marketplace/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (svc application) setupGRPCServer() {
	grpcPriceFetcher := NewGRPCPriceFetcherServer(svc)
	proto.RegisterMarketplaceServiceServer(svc.goApiBoot.GrpcServer, grpcPriceFetcher)
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

func (u *GRPCMarketPlaceServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (s *GRPCMarketPlaceServer) CreateShop(ctx context.Context, req *proto.CreateShopRequest) (*proto.Shop, error) {
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

	err := getErrorFromChan(s.svc.shopRepo.Save(shop))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to create shop")
	}

	return &proto.Shop{
		Id:             shop.ID,
		Name:           shop.Name,
		Location:       shop.Location,
		OperationHours: shop.OperationHours,
		Coordinates: &proto.Coordinates{
			Latitude:  shop.Coordinates[0],
			Longitude: shop.Coordinates[1],
		},
	}, nil
}

func (s *GRPCMarketPlaceServer) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.Product, error) {
	product := Product{
		ID:          primitive.NewObjectID().Hex(),
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
	}

	err := getErrorFromChan(s.svc.productRepo.Save(product))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to create product")
	}

	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       float32(product.Price),
	}, nil
}

func (s *GRPCMarketPlaceServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	user := User{
		ID:       primitive.NewObjectID().Hex(),
		Name:     req.Name,
		Location: req.Location,
		Coordinates: [2]float64{
			req.Coordinates.Latitude,
			req.Coordinates.Longitude,
		},
	}

	err := getErrorFromChan(s.svc.userRepo.Save(user))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to create user")
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
	id := req.Id
	shop := &Shop{}

	shop, err := getShopOrError(s.svc.shopRepo.FindOneById(id))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get shop")
	}

	products, err := s.GetServiceableProducts(ctx, &proto.GetServiceableProductsRequest{ShopId: id})
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get serviceable products for shop")
	}

	return &proto.Shop{
		Id:                  shop.ID,
		Name:                shop.Name,
		Location:            shop.Location,
		OperationHours:      shop.OperationHours,
		ServiceableProducts: products.Products,
		Coordinates: &proto.Coordinates{
			Latitude:  shop.Coordinates[0],
			Longitude: shop.Coordinates[1],
		},
	}, nil
}

func (s *GRPCMarketPlaceServer) GetProductByID(ctx context.Context, req *proto.GetRequest) (*proto.Product, error) {
	id := req.Id
	product := &Product{}

	product, err := getProductOrError(s.svc.productRepo.FindOneById(id))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get product")
	}

	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       float32(product.Price),
	}, nil
}

func (s *GRPCMarketPlaceServer) GetUserByID(ctx context.Context, req *proto.GetRequest) (*proto.User, error) {
	id := req.Id
	user := &User{}

	user, err := getUserOrError(s.svc.userRepo.FindOneById(id))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get user: %s")
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
	shopId := req.ShopId
	productId := req.ProductId
	shop := &Shop{}

	if !s.svc.productRepo.IsExistsById(productId) {
		s.svc.logger.Printf("Error: product[%s] does not exists", productId)
		return nil, errors.New("product does not exists")
	}

	shop, err := getShopOrError(s.svc.shopRepo.FindOneById(shopId))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get shop")
	}

	var alreadyExists bool
	for _, id := range shop.ServiceableProductsId {
		if id == productId {
			alreadyExists = true
		}
	}

	if !alreadyExists {
		products := append(shop.ServiceableProductsId, productId)
		shop.ServiceableProductsId = products

		errCh := s.svc.shopRepo.Save(shop)
		err := getErrorFromChan(errCh)
		if err != nil {
			s.svc.logger.Println("Error: ", err)
			return nil, errors.New("failed to add serviceable product")
		}
	}

	_, err = getInventoryOrError(s.svc.inventoryRepo.FindOne(primitive.M{"shop_id": shopId, "product_id": productId}))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			inventory := Inventory{
				ID:        primitive.NewObjectID().Hex(),
				ShopID:    shopId,
				ProductID: productId,
				Quantity:  0,
			}

			errCh := s.svc.inventoryRepo.Save(inventory)
			err := getErrorFromChan(errCh)
			if err != nil {
				s.svc.logger.Println("Error: ", err)
				return nil, errors.New("failed to create inventory")
			}
		} else {
			s.svc.logger.Println("Error: ", err)
			return nil, errors.New("failed while looking for inventory")
		}
	}

	return s.GetShopByID(ctx, &proto.GetRequest{
		Id: shopId,
	})
}

func (s *GRPCMarketPlaceServer) GetServiceableProducts(ctx context.Context, req *proto.GetServiceableProductsRequest) (*proto.Products, error) {
	var serviceableProducts []*proto.Product
	id := req.ShopId
	shop := &Shop{}

	shop, err := getShopOrError(s.svc.shopRepo.FindOneById(id))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get shop")
	}

	for _, productId := range shop.ServiceableProductsId {
		product, err := s.GetProductByID(ctx, &proto.GetRequest{
			Id: productId,
		})

		if err != nil {
			s.svc.logger.Println("Error: ", err)
			return nil, errors.New("failed to get product")
		}
		serviceableProducts = append(serviceableProducts, product)
	}

	return &proto.Products{
		Products: serviceableProducts,
	}, nil
}

func (s *GRPCMarketPlaceServer) GetInventory(ctx context.Context, req *proto.GetInventoryRequest) (*proto.Inventory, error) {
	shopId := req.ShopId
	productId := req.ProductId
	inventory := &Inventory{}

	inventory, err := getInventoryOrError(s.svc.inventoryRepo.FindOne(primitive.M{"shop_id": shopId, "product_id": productId}))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get inventory")
	}

	return &proto.Inventory{
		Id:        inventory.ID,
		ShopId:    inventory.ShopID,
		ProductId: inventory.ProductID,
		Quantity:  int32(inventory.Quantity),
	}, nil
}

func (s *GRPCMarketPlaceServer) UpdateInventory(ctx context.Context, req *proto.UpdateInventoryRequest) (*proto.Inventory, error) {
	shopId := req.ShopId
	productId := req.ProductId
	inventory := &Inventory{}

	inventory, err := getInventoryOrError(s.svc.inventoryRepo.FindOne(primitive.M{"shop_id": shopId, "product_id": productId}))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get inventory")
	}

	if req.Add {
		inventory.Quantity = inventory.Quantity + int(req.Change)
	} else {
		inventory.Quantity = inventory.Quantity - int(req.Change)
	}

	if inventory.Quantity < 0 {
		return nil, errors.New("inventory's quantity cannot be negative")
	}

	errCh := s.svc.inventoryRepo.Save(inventory)
	err = getErrorFromChan(errCh)
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to update inventory")
	}

	return &proto.Inventory{
		Id:        inventory.ID,
		ShopId:    inventory.ShopID,
		ProductId: inventory.ProductID,
		Quantity:  int32(inventory.Quantity),
	}, nil
}

func (s *GRPCMarketPlaceServer) GetShopsByServiceableProducts(ctx context.Context, req *proto.GetShopsByServiceableProductsRequest) (*proto.Shops, error) {
	var shops []Shop
	productId := req.ProductId
	filter := bson.M{"serviceableProductsId": productId}

	shops, err := getShopsOrError(s.svc.shopRepo.Find(filter, nil, 0, 0))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get shops")
	}

	resultShops := proto.Shops{}

	for _, shop := range shops {
		pshop, err := s.ParseShop(ctx, &shop)
		if err != nil {
			s.svc.logger.Println("Error: ", err)
			return nil, errors.New("failed to parse data")
		}

		resultShops.Shops = append(resultShops.Shops, pshop)
	}

	return &proto.Shops{
		Shops: resultShops.Shops,
	}, nil
}

func (s *GRPCMarketPlaceServer) GetShopForUser(ctx context.Context, req *proto.GetShopForUserRequest) (*proto.Shops, error) {
	var shops []Shop
	userId := req.UserId
	user := &User{}
	resultShops := proto.Shops{}

	user, err := getUserOrError(s.svc.userRepo.FindOneById(userId))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get user")
	}

	shops, err = getShopsOrError(s.svc.shopRepo.Find(nil, nil, 0, 0))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get all shops")
	}

	for _, shop := range shops {
		if calculateDistance(user.Coordinates, shop.Coordinates) < req.MaxDistanceInKM {
			pshop, err := s.ParseShop(ctx, &shop)
			if err != nil {
				s.svc.logger.Println("Error: ", err)
				return nil, errors.New("failed to parse data")
			}

			resultShops.Shops = append(resultShops.Shops, pshop)
		}
	}

	return &proto.Shops{
		Shops: resultShops.Shops,
	}, nil
}

func (s *GRPCMarketPlaceServer) GetNearestNeighbour(ctx context.Context, req *proto.GetNearestNeighbourRequest) (*proto.User, error) {
	var users []User
	userId := req.UserId
	user := &User{}

	user, err := getUserOrError(s.svc.userRepo.FindOneById(userId))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get user")
	}

	users, err = getUsersOrError(s.svc.userRepo.Find(nil, nil, 0, 0))
	if err != nil {
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get all users")
	}

	var nearest User
	if user.ID == users[0].ID {
		nearest = users[1]
	} else {
		nearest = users[0]
	}

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
		s.svc.logger.Println("Error: ", err)
		return nil, errors.New("failed to get serviceable products")
	}

	return &proto.Shop{
		Id:                  shop.ID,
		Name:                shop.Name,
		Location:            shop.Location,
		OperationHours:      shop.OperationHours,
		ServiceableProducts: products.Products,
		Coordinates: &proto.Coordinates{
			Latitude:  shop.Coordinates[0],
			Longitude: shop.Coordinates[1],
		},
	}, nil
}
