package main

import (
	pb "calculator/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"strconv"
)

// server is used to implement the server.
type server struct {
	pb.UnimplementedSumServer
}

func (s *server) Addition(ctx context.Context, in *pb.CalculationRequest) (*pb.CalculationResponse, error) {

	return &pb.CalculationResponse{Result: in.GetInt1() + in.GetInt2()}, nil
}

func (s *server) Subtraction(ctx context.Context, in *pb.CalculationRequest) (*pb.CalculationResponse, error) {

	return &pb.CalculationResponse{Result: in.GetInt1() - in.GetInt2()}, nil
}
func (s *server) Multiplication(ctx context.Context, in *pb.CalculationRequest) (*pb.CalculationResponse, error) {

	return &pb.CalculationResponse{Result: in.GetInt1() * in.GetInt2()}, nil
}
func (s *server) Division(ctx context.Context, in *pb.CalculationRequest) (*pb.CalculationResponse, error) {

	return &pb.CalculationResponse{Result: in.GetInt1() / in.GetInt2()}, nil
}

func main() {
	flag.Parse()

	content, _ := ioutil.ReadFile("port")
	text := string(content)
	log.Println(text)
	portNum, _ := strconv.Atoi(text)
	var (
		port = flag.Int("port", portNum, "The server port")
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSumServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
