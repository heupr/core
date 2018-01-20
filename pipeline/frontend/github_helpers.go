package frontend

import (
	"context"
	"net/url"
	"reflect"

	"github.com/google/go-github/github"
	"github.com/google/go-querystring/query"
)

const mediaTypeIntegrationPreview = "application/vnd.github.machine-man-preview+json"

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// HeuprInstallation provides a workaround to GitHub API holes.
type HeuprInstallation struct {
	github.Installation
	AppID *int `json:"app_id,omitempty"`
}

func listUserInstallations(ctx context.Context, client *github.Client, opt *github.ListOptions) ([]*HeuprInstallation, *github.Response, error) {
	u, err := addOptions("user/installations", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeIntegrationPreview)

	var i struct {
		Installations []*HeuprInstallation `json:"installations"`
	}
	resp, err := client.Do(ctx, req, &i)
	if err != nil {
		return nil, resp, err
	}

	return i.Installations, resp, nil
}
