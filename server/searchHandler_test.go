package main

import (
	"context"
	"github.com/google/go-github/github"
	pb "github.com/kruckenb/githubSearch/proto"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	repoUrl1 = "repoUrl1"
	repoUrl2  = "repoUrl2"
	repoName1 = "repo1"
	repoName2 = "repo2"

	repositories = []github.Repository{
		{URL:&repoUrl1, FullName: &repoName1},
		{URL:&repoUrl2, FullName: &repoName2}}
)

func TestMapResults(t *testing.T) {
	searchResults := &github.RepositoriesSearchResult{Repositories: repositories}
	pbResults := mapSearchResults(searchResults)

	assert.Equal(t, repoName1, pbResults[0].Repo, "map repo1 name")
	assert.Equal(t, repoUrl1, pbResults[0].FileUrl, "map file1 name")
	assert.Equal(t, repoName2, pbResults[1].Repo, "map repo2 name")
	assert.Equal(t, repoUrl2, pbResults[1].FileUrl, "map file2 name")
}

func TestSearchTerms(t *testing.T) {
	assert.Equal(t, "", constructSearchTerms(&pb.SearchRequest{SearchTerm: ""}), "No search terms")
	assert.Equal(t, "foo", constructSearchTerms(&pb.SearchRequest{SearchTerm: "foo"}), "Search foo")
	assert.Equal(t, " user:bar", constructSearchTerms(&pb.SearchRequest{SearchTerm: "", User: "bar"}), "Search user bar")
	assert.Equal(t, "foo boo user:bar", constructSearchTerms(&pb.SearchRequest{SearchTerm: "foo boo", User: "bar"}), "Search foo boo for user bar")
}

func TestSearch(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.GetSearchRepositories,
			github.RepositoriesSearchResult{
				Repositories: repositories,
			}))
	c := github.NewClient(mockedHTTPClient)

	repos, err := searchGithub(context.Background(), &pb.SearchRequest{}, c)
	assert.Nil(t, err, "No search error")
	assert.Equal(t, 2, len(repos), "Find 2 repos")
	assert.Equal(t, repoName1, repos[0].Repo,"Search match first repo name")
	assert.Equal(t, repoName2, repos[1].Repo,"Search match second repo name")
	assert.Equal(t, repoUrl1, repos[0].FileUrl,"Search match first url")
	assert.Equal(t, repoUrl2, repos[1].FileUrl,"Search match second url")
}

func TestSearchRateLimitNoAuth(t *testing.T) {
	// Use Github to test this, no immediately-obvious way to test this with mocks
	client := githubClient(context.Background(), nil)

	// Search enough times to trigger a rate limit error
	for count := 0; count < 10; count++ {
		searchGithub(context.Background(), &pb.SearchRequest{SearchTerm: ""}, client)
	}

	_, err := searchGithub(context.Background(), &pb.SearchRequest{SearchTerm: ""}, client)
	assert.Condition(t, func() bool { return strings.HasPrefix(err.Error(), "GITHUB_RATE_LIMIT_EXCEEDED: ") }, "Rate limit error w/o auth")
}
