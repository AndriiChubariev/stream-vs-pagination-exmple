package main

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	EntityServiceServer
}

var database = []Entity{
	{
		Id:   "1",
		Name: "Entity 1",
	},
	{
		Id:   "2",
		Name: "Entity 2",
	},
	{
		Id:   "3",
		Name: "Entity 3",
	},
	{
		Id:   "4",
		Name: "Entity 4",
	},
	{
		Id:   "5",
		Name: "Entity 5",
	},
	{
		Id:   "6",
		Name: "Entity 6",
	},
}

func (s *server) GetEntityWithPagination(ctx context.Context, req *EntityRequest) (*EntityResponse, error) {

	page := req.GetPage()
	pageSize := req.GetPageSize()

	from := page * pageSize
	to := page*pageSize + pageSize

	entities := make([]*Entity, pageSize)
	for i := from; i < to; i++ {
		entities[i-from] = &database[i]
	}

	result := &EntityResponse{
		Entities: entities,
		Total:    uint32(len(database)),
	}
	return result, nil
}

func (s *server) GetEntityWithStream(_ *emptypb.Empty, stream EntityService_GetEntityWithStreamServer) error {
	for i := range database {
		if err := stream.Send(&database[i]); err != nil {
			return err
		}
	}
	return nil
}

func StartServer() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterEntityServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
