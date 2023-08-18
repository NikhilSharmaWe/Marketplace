package main

import (
	"github.com/NikhilSharmaWe/marketplace/proto"
	"google.golang.org/grpc"
)

func NewGRPCClient(remoteAddr string) (proto.MarketplaceServiceClient, error) {

	conn, err := grpc.Dial(remoteAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := proto.NewMarketplaceServiceClient(conn)

	return c, nil
}
