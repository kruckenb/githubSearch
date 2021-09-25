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

