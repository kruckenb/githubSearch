# GitHub Search

Service that searches GitHub to find repos containing specific phrases

### API Spec

```
service GithubSearchService {
   rpc Search(SearchRequest) returns (SearchResponse);
}

message SearchRequest {
   string search_term = 1;
   string user = 2;
}

message SearchResponse {
   repeated Result results = 1;
}

message Result {
   string file_url = 1;
   string repo = 2;
}
```

## Build Instructions
```
brew install protobuf protoc-gen-go protoc-gen-go-grpc
git clone https://github.com/kruckenb/githubSearch
cd githubSearch/
cd server
go build
cd ../client
go build
```
### Update GRPC / Protobuf definition
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/githubSearch.proto
```

## Usage
### GRPC Server
```
Usage: ./server/server [-flag1] [-flagN]
  -help
    	Show this help
  -port int
    	Server TCP port (default 8080)
  -token string
    	GitHub auth token

Environment variables
	GITHUB_AUTH_TOKEN
	GITHUB_SEARCH_PORT
```	

### CLI GRPC client
```
Usage: ./client/client: [-flag1] [-flagN] searchTerm1 searchTermN
  -help
    	Show this help
  -port int
    	Connect to localhost TCP port (default 8080)
  -timeout int
    	Connect timeout seconds (default 1)
  -user string
    	Search GitHub user (optional)
```

## To-do
- [X] Template to start with (thanks grpc.io!)
- [X] Get something working
- [X] CLI client
- [X] GitHub API interface (thanks go-github!)
- [X] Some code organization
- [X] Configure through CLI
- [X] Configure through ENV
- [X] Github oauth
- [X] Some tests
  - [ ] Mock for rate-limit test
  - [ ] More tests?
- [ ] REST API / Web server
- [ ] Multi-instance rate limits via Redis
- [ ] CI stuff
- [ ] Deployment stuff
