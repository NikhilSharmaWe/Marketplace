package main

import (
	"context"
	"net"

	"github.com/NikhilSharmaWe/marketplace/proto"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	}
	err := s.svc.Create(ctx, s.svc.shopCollection, shop)
	if err != nil {
		return nil, err
	}

	return &proto.Shop{
		Id:             shop.ID,
		Name:           shop.Name,
		Location:       shop.Location,
		Operationhours: shop.OperationHours,
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
	}
	err := s.svc.Create(ctx, s.svc.userCollection, user)
	if err != nil {
		return nil, err
	}

	return &proto.User{
		Id:       user.ID,
		Name:     user.Name,
		Location: user.Location,
	}, nil
}

func (s *GRPCMarketPlaceServer) GetShopByID(ctx context.Context, req *proto.GetRequest) (*proto.Shop, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	id := req.Id
	shop := Shop{}

	err := s.svc.GetByID(id, s.svc.shopCollection, &shop)
	if err != nil {
		return nil, err
	}

	return &proto.Shop{
		Id:             shop.ID,
		Name:           shop.Name,
		Location:       shop.Location,
		Operationhours: shop.OperationHours,
	}, nil
}

func (s *GRPCMarketPlaceServer) GetProductByID(ctx context.Context, req *proto.GetRequest) (*proto.Product, error) {
	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	id := req.Id
	product := Product{}

	err := s.svc.GetByID(id, s.svc.productCollection, &product)
	if err != nil {
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
		return nil, err
	}

	return &proto.User{
		Id:       user.ID,
		Name:     user.Name,
		Location: user.Location,
	}, nil
}

func (s *GRPCMarketPlaceServer) AddServiceableProduct(ctx context.Context, req *proto.AddServiceableProductRequest) (*proto.ServiceableProduct, error) {
	// ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	// shopId := req.ShopId
	// productId := req.ProductId
	// shop := Shop{}

	// err := s.svc.GetByID(shopId, s.svc.shopCollection, &shop)
	// if err != nil {
	// 	return nil, err
	// }

	// product := Product{}
	// err = s.svc.GetByID(productId, s.svc.productCollection, &product)
	// if err != nil {
	// 	return nil, err
	// }

	// products := append(shop.ServiceableProducts, product)

	// err = s.svc.UpdateServiceableProductsForShop(ctx, shopId, products)
	// if err != nil {
	// 	return nil, err
	// }

	// return &proto.Shop{
	// 	Id:             shop.ID,
	// 	Name:           shop.Name,
	// 	Location:       shop.Location,
	// 	Operationhours: shop.OperationHours,
	// }, nil

	ctx = context.WithValue(ctx, "requestID", uuid.NewV4().String())
	shopId := req.ShopId
	productId := req.ProductId
	shop := Shop{}

	err := s.svc.GetByID(shopId, s.svc.shopCollection, &shop)
	if err != nil {
		return nil, err
	}

	product := Product{}
	err = s.svc.GetByID(productId, s.svc.productCollection, &product)
	if err != nil {
		return nil, err
	}

	serviceableProduct := ServiceableProduct{
		ID:        primitive.NewObjectID().Hex(),
		ProductID: productId,
		ShopID:    shopId,
	}
	err = s.svc.Create(ctx, s.svc.serviceableProductCollection, serviceableProduct)
	if err != nil {
		return nil, err
	}

	return &proto.ServiceableProduct{
		Id:        serviceableProduct.ID,
		ProductId: serviceableProduct.ProductID,
		ShopId:    serviceableProduct.ShopID,
	}, nil
}
