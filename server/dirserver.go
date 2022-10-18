package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/media-informatics/aufgabe04a/service"
	"google.golang.org/grpc"
)

type DirServer struct {
	service.UnimplementedDirectoryServer
}

var (
	server = flag.String("server", service.Addr, "server address with port")
)

func (ds *DirServer) GetDir(ctx context.Context, in *service.DirName) (*service.FileList, error) {
	dir := in.GetName()
	log.Printf("received: %s", dir)
	f, err := os.Open(dir)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("could not open %s: %w", dir, err)
	}
	defer f.Close()

	list, err := f.ReadDir(0)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("could not read %s: %w", dir, err)
	}

	files := make([]string, len(list))
	for i, v := range list {
		files[i] = v.Name()
	}
	return &service.FileList{Entry: files}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *server)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s := grpc.NewServer()
	service.RegisterDirectoryServer(s, &DirServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
