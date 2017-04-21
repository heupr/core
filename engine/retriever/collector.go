package retriever

import "github.com/google/go-github/github"

var Workload = make(chan github.Issue, 100)

func Collector(issues []github.Issue) {
	for _, i := range issues {
		Workload <- i
	}
}
