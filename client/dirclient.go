package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/media-informatics/aufgabe04a/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	dirname = flag.String("dir", ".", "directory to list")
	server  = flag.String("server", service.Addr, "server address with port")
)

func main() {
	flag.Parse()
	// replace deprecated grpc.WithInsecure():
	cred := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(*server, cred)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := service.NewDirectoryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dir := &service.DirName{Name: *dirname}
	dirlist, err := c.GetDir(ctx, dir)
	if err != nil {
		log.Fatalf("did not receive dir: %v", err)
	}
	fmt.Printf("content of %s:\n", *dirname)
	for _, f := range dirlist.GetEntry() {
		fmt.Println(f)
	}
}
