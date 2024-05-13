package main

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func startClient() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := NewEntityServiceClient(conn)

	log.Print("--- Pagination Example")

	const pageSize = 2
	log.Println("Making the first call to get the first page of entities")
	entityResponse, err := c.GetEntityWithPagination(context.Background(), &EntityRequest{Page: 0, PageSize: pageSize})
	if err != nil {
		log.Fatalf("could not get entity: %v", err)
	}

	// Print each entity to the console
	for _, entity := range entityResponse.Entities {
		log.Printf("Entity ID: %s, Entity Name: %s\n", entity.Id, entity.Name)
	}
	log.Print("Ok, our page size is 2 and it is constant size. We have ", entityResponse.Total, " entities, that means we need to make 2 more calls")

	numberOfCallsToMake := 2
	for i := 1; i <= numberOfCallsToMake; i++ {
		log.Printf("Making a call#: %d\n", i+1)

		entityResponse, err := c.GetEntityWithPagination(context.Background(), &EntityRequest{Page: int32(i), PageSize: pageSize})
		if err != nil {
			log.Fatalf("could not get entity: %v", err)
		}

		// Print each entity to the console
		for _, entity := range entityResponse.Entities {
			log.Printf("Entity ID: %s, Entity Name: %s\n", entity.Id, entity.Name)
		}
	}

	log.Print("--- Stream Example")
	stream, err := c.GetEntityWithStream(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not get entity: %v", err)
	}

	for {
		entity, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not receive entity: %v", err)
		}

		log.Printf("From Stream: Entity ID: %s, Entity Name: %s\n", entity.Id, entity.Name)
	}
}
