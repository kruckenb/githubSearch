package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	pb "github.com/kruckenb/githubSearch/proto"
	"golang.org/x/oauth2"
	"net/http"
)

type server struct {
	oauthToken *string

	pb.UnimplementedGithubSearchServiceServer
}

func (s *server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	results, err := searchGithub(ctx, req, githubClient(ctx,s.oauthToken))
	if err != nil {
		return nil, err
	}
	return &pb.SearchResponse{ Results: results }, nil
}

func githubClient(ctx context.Context, oauthToken *string) *github.Client {
	if oauthToken != nil && *oauthToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *oauthToken},
		)
		return github.NewClient(oauth2.NewClient(ctx, ts))
	}
	return github.NewClient(nil)
}

func searchGithub (ctx context.Context, req *pb.SearchRequest, client *github.Client) ([]*pb.Result, error) {
	searchResults, searchResponse, err :=
		client.Search.Repositories(ctx,
			constructSearchTerms(req),
			nil)

	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("GITHUB_RATE_LIMIT_EXCEEDED: %s", err.Error())
		}
		return nil, fmt.Errorf("GITHUB_API_ERROR: %s", err.Error())
	} else if searchResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GITHUB_API_HTTP_ERROR %d: %s", searchResponse.StatusCode, searchResponse.Status)
	}
	return mapSearchResults(searchResults), nil
}

func constructSearchTerms (req *pb.SearchRequest) string {
	terms := req.GetSearchTerm()
	if req.User != "" {
		terms += " user:" + req.GetUser()
	}
	return terms
}

// Map results from Github Search Result structure to response for this service
func mapSearchResults(searchResults *github.RepositoriesSearchResult) []*pb.Result {
	var results []*pb.Result
	for _, result := range searchResults.Repositories {
		results = append(results, &pb.Result{FileUrl: result.GetURL(), Repo: result.GetFullName()})
	}
	return results
}
