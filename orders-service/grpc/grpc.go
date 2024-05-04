package grpc

import (
	"context"
	"fmt"
	"os"

	proto "github.com/NikolaB131-org/logistics-backend/proto/warehouse"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcConn *grpc.ClientConn

func ConnectGrpc() {
	conn, err := grpc.Dial(os.Getenv("GRPC_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Unable to connect grpc server: %v\n", err)
		os.Exit(1)
	}
	grpcConn = conn

	whc := proto.NewWarehouseServiceClient(conn)
	_, err = whc.DecreaseProductQuantity(context.Background(), &proto.DecreaseProductQuantityRequest{Id: "1", Quantity: 2})
	fmt.Println(err)
}
