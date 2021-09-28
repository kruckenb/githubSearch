package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/kruckenb/githubSearch/proto"
	"google.golang.org/grpc"
)

const bindHost = "localhost"

func main() {
	timeout := flag.Int("timeout", 1, "Connect timeout seconds")
	port := flag.Int("port", 8080, "Connect to localhost TCP port")
	searchUser := flag.String("user", "", "Search GitHub user (optional)")
	help := flag.Bool("help", false, "Show this help")
	flag.Parse()

	if flag.NArg() == 0 || *help {
		fmt.Println("Usage: [-flag1] [-flagN] searchTerm1 searchTermN")
		flag.PrintDefaults()
		os.Exit(0)
	}
	searchTerms := flag.Args()
	fmt.Printf("Search terms: %v\n", searchTerms)

	// Set up a connection to the server.
	address := fmt.Sprintf("%s:%d", bindHost, *port)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", address, err)
	}
	defer conn.Close()

	c := pb.NewGithubSearchServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout) * time.Second)
	defer cancel()

	r, err := c.Search(ctx,
		&pb.SearchRequest{User: *searchUser, SearchTerm: strings.Join(searchTerms, " ")})
	if err != nil {
		log.Fatalf("SEARCH_FAILURE: %v", err)
	}
	log.Printf("Search results: %v", r.GetResults())
}
