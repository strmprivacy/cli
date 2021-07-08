package entity

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

func SetupGrpc(host string, token string) (*grpc.ClientConn, context.Context) {

	var err error
	var creds grpc.DialOption
	if strings.Contains(host, ":50051") {
		creds = grpc.WithInsecure()
	} else {
		creds = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}
	clientConnection, err := grpc.Dial(host, creds)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return clientConnection, ctx
}
