package main

import (
	"flag"
	pb "github.com/kruckenb/githubSearch/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
)

var port *int
var token *string
const (
	defaultPort = 8080
	defaultToken = ""
	envPort = "GITHUB_SEARCH_PORT"
	envToken = "GITHUB_AUTH_TOKEN"
)

func main() {
	readSettings()

	lis, err := net.Listen("tcp", "localhost:" + strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("LISTEN_FAIL: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGithubSearchServiceServer(s, &server{oauthToken: token})

	log.Printf("Server started on %s", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("SERVER_START_FAIL: %v", err)
	}
}

// Parse configuration settings from command-line flags and environment variables
func readSettings() {
	port = flag.Int("port", defaultPort, "Server TCP port")
	token = flag.String("token", defaultToken, "GitHub auth token")
	help := flag.Bool("help", false, "Show this help")
	flag.Parse()

	if *help {
		flag.Usage()

		os.Exit(0)
	}

	if *port == defaultPort && os.Getenv(envPort) != "" {
		if num, err := strconv.ParseInt(os.Getenv(envPort), 10, 8); err != nil {
			*port = int(num)
		}
	}

	if *token == defaultToken && os.Getenv(envToken) != "" {
		*token = os.Getenv(envToken)
	}
}