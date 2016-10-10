package grades

import (
	"encoding/csv"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	//"bytes"
	"strconv"
)

func generate() {

	// authentication
	//we may not need OAUTH for this method
	//see https://developer.github.com/changes/2014-02-28-issue-and-pull-query-enhancements/
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "56e5fcc7bec01a3f3f797d528f08a83a5e3fec74"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	opt := &github.IssueListByRepoOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	// Create a csv file
	f, err := os.Create("./trainingset_corefx.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	// Write Unmarshaled json data to CSV file
	w := csv.NewWriter(f)

	for i := 0; ; i++ {
		issues, resp, err := client.Issues.ListByRepo("dotnet", "corefx", opt)
		if err != nil {
			fmt.Printf("error fetching rate limit (%v)\n", err)
		}
		for _, issue := range issues {
			//used Buffer because it is faster than concating strings
			//var buffer bytes.Buffer
			//may want to replace record[] string with buffer
			//see http://golang-examples.tumblr.com/post/86169510884/fastest-string-contatenation
			var record []string
			//var buffer bytes.Buffer

			//url column (information column that will only need to be included for now to help develop labels)
			record = append(record, *issue.URL)

			//url id column (information column that will only need to be included for now to help develop labels)
			//Format: https://github.com/dotnet/coreclr/pull/{number}
			record = append(record, strconv.Itoa(*issue.Number))

			//title column (information column that will only need to be included for now to help develop labels)
			record = append(record, *issue.Title)

			if (issue.Body != nil)	{
				record = append(record, *issue.Body)
			} else {
				record = append(record, "No description")
			}
                        //username column (prediction value)
			record = append(record, *issue.User.Login)

			/*
                        //labels column
			if len(issue.Labels) <= 0 {
				record = append(record, "NOLABEL")
			} else {
                        	sep := []byte(",")
				buffer.WriteString(*issue.Labels[0].Name)
				for _, label := range issue.Labels[1:] {
					buffer.Write(sep)
					buffer.WriteString(*label.Name)
				}
				record = append(record, buffer.String())
				buffer.Reset()
			} */
                        w.Write(record)
		}
		if resp.NextPage == 0 {
			break
		} else {
			opt.ListOptions.Page = resp.NextPage
		}
	}
	w.Flush()
}
