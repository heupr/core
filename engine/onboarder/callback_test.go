package onboarder

/*
import (
	// "context"
    // "fmt"
	// "io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	// "golang.org/x/oauth2"
)
*/

/*
func Test_githubCallbackHandle(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.String() != "/token" {
            t.Errorf("Unexpected exchange URL %s\n", r.URL.String())
        }
        headerAuth := r.Header.Get("Authorization")
        if headerAuth != "Basic Q0xJRU5UX0lEOkNMSUVOVF9TRUNSRVQ=" {
			t.Errorf("Unexpected authorization header %s\n", headerAuth)
		}
        headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/x-www-form-urlencoded" {
			t.Errorf("Unexpected Content-Type header %s\n", headerContentType)
		}
        body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed reading request body - body: %s error: %s\n", body, err)
		}
        expected := "code=exchange-code&grant_type=authorization_code&redirect_uri=REDIRECT_URL"
        if string(body) != expected {
            t.Errorf("Unexpected exchange payload\nFound: %v\nExpected: %v", string(body), expected)
		}
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer"))
    }))
    defer ts.Close()

    oaConfig = &oauth2.Config{
        // ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "REDIRECT_URL",
		// Scopes:       []string{"scope1", "scope2"},
		Endpoint:     oauth2.Endpoint{
			AuthURL:  ts.URL + "/auth",
			TokenURL: ts.URL + "/token",
		},
    }

	oaConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     oauth2.Endpoint,
	}

    rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/handler-check", nil)
	if err != nil {
		t.Error(
			"Failure generating testing request",
			"\n", err,
		)
	}
    req.Form = url.Values{}
    req.Form.Set("state", oaState)
    req.Form.Set("code", "exchange-code")
	h := RepoServer{}
	handler := http.HandlerFunc(h.githubCallbackHandle)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusPermanentRedirect {
		t.Errorf("Handler returning incorrect status code; returning %v", status)
	}
}

func Test_githubCallbackHandle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/token" {
			t.Errorf("Unexpected exchange request URL, %v is found.", r.URL)
		}
		headerAuth := r.Header.Get("Authorization")
		if headerAuth != "Basic Q0xJRU5UX0lEOkNMSUVOVF9TRUNSRVQ=" {
			t.Errorf("Unexpected authorization header, %v is found.", headerAuth)
		}
		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/x-www-form-urlencoded" {
			t.Errorf("Unexpected Content-Type header, %v is found.", headerContentType)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed reading request body: %s.", err)
		}
		if string(body) != "code=exchange-code&grant_type=authorization_code&redirect_uri=REDIRECT_URL" {
			t.Errorf("Unexpected exchange payload, %v is found.", string(body))
		}
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer"))
	}))
	defer ts.Close()
	conf := newConf(ts.URL)
	tok, err := conf.Exchange(context.Background(), "exchange-code")
	if err != nil {
		t.Error(err)
	}
	if !tok.Valid() {
		t.Fatalf("Token invalid. Got: %#v", tok)
	}
	if tok.AccessToken != "90d64460d14870c08c81352a05dedd3465940a7c" {
		t.Errorf("Unexpected access token, %#v.", tok.AccessToken)
	}
	if tok.TokenType != "bearer" {
		t.Errorf("Unexpected token type, %#v.", tok.TokenType)
	}
	scope := tok.Extra("scope")
	if scope != "user" {
		t.Errorf("Unexpected value for scope: %v", scope)
	}
}
*/

/*
func Test_githubCallbackHandle(t *testing.T) {
	h := RepoServer{}

    oaConfig = &oauth2.Config{
    	ClientID:     "",
        ClientSecret: "",
        Endpoint:     oauth2.Endpoint{},
    	// Endpoint:     oauth2.Endpoint{
		// 	AuthURL:  "http://test.com/auth",
		// 	TokenURL: "http://test.com/token",
		// },
    	Scopes:       []string{"test-scope"},
    }

	outputURL := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOffline)

	req, err := http.NewRequest("POST", "/handler-check", nil)
	if err != nil {
		t.Errorf("Failure generating test request %s", err)
    }
    // ctx := req.Context()
    // ctx = context.WithValue(ctx, "code", "exchange-code")
    // req = req.WithContext(ctx)
    req.Form = url.Values{}
    req.Form.Set("state", oaState)
    rec := httptest.NewRecorder()

    handler := http.HandlerFunc(h.githubCallbackHandle)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusPermanentRedirect {
        t.Errorf("HANDLER RETURNING WRONG CODE")
	}
}
*/

/*
func Test_githubCallbackHandle(t *testing.T) {
	h := RepoServer{}

	oaConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Endpoint:     oauth2.Endpoint{},
		Scopes:       []string{},
	}

	outputURL := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOffline)

	req, err := http.NewRequest("POST", "/handler-check", nil)
	if err != nil {
		t.Errorf("Failure generating test request: %s", err)
	}
    req.Form = url.Values{}
    req.Form.Set("code", outputURL)
    req.Form.Set("state", oaState)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(h.githubCallbackHandle)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusPermanentRedirect {
		t.Errorf("Handler returning incorrect status")
	}
}

// "code" url.Values needs to be populated
// request URL set via AuthCodeURL method
// pass in secret to populate "state"
// context set to Background
*/

/*
func Test_githubCallbackHandle(t *testing.T) {
    h := RepoServer{}
    oaConfig = &oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "",
		Endpoint:     oauth2.Endpoint{},
		Scopes:       []string{},
	}
    outputURL := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOffline)
    t.Error(outputURL) // TEMPORARY

    req, err := http.NewRequest("POST", "/handler-check", nil)
	if err != nil {
		t.Errorf("Failure generating test request: %v", err)
	}
    // req.Form = url.Values{}
    // req.Form.Set("state", oaState)

    rec := httptest.NewRecorder()

    handler := http.HandlerFunc(h.githubCallbackHandle)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusPermanentRedirect {
		t.Errorf("Handler returning incorrect status: %v", status)
	}
}
*/

/*
func Test_githubCallbackHandle(t *testing.T) {
    h := RepoServer{}

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/token" {
			t.Errorf("Unexpected exchange request URL, %v is found.", r.URL)
		}
		headerAuth := r.Header.Get("Authorization")
		if headerAuth != "Basic Q0xJRU5UX0lEOkNMSUVOVF9TRUNSRVQ=" {
			t.Errorf("Unexpected authorization header, %v is found.", headerAuth)
		}
		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/x-www-form-urlencoded" {
			t.Errorf("Unexpected Content-Type header, %v is found.", headerContentType)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed reading request body: %s.", err)
		}
		if string(body) != "code=exchange-code&grant_type=authorization_code&redirect_uri=REDIRECT_URL" {
			t.Errorf("Unexpected exchange payload, %v is found.", string(body))
		}
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer"))
	}))

    oaConfig = &oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		RedirectURL:  "REDIRECT_URL",
		Scopes:       []string{"scope1", "scope2"},
		Endpoint:     oauth2.Endpoint{
			AuthURL:  ts.URL + "/auth",
			TokenURL: ts.URL + "/token",
		},
    }

    // outputURL := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOffline)

    req, err := http.NewRequest("POST", "/handler-check", nil)
	if err != nil {
		t.Errorf("Failure generating test request: %v", err)
	}
    // // req.Form = url.Values{}
    // // req.Form.Set("state", oaState)

    rec := httptest.NewRecorder()

    handler := http.HandlerFunc(h.githubCallbackHandle)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusPermanentRedirect {
		t.Errorf("Handler returning incorrect status: %v", status)
	}
}
*/

/*
func Test_githubCallbackHandle(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Println("AUTH TEST")
    })
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("TOKEN TEST")
        w.WriteHeader(http.StatusOK)
    })
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/token" {
			t.Errorf("Unexpected exchange request URL, %v is found.", r.URL)
		}
		headerAuth := r.Header.Get("Authorization")
		if headerAuth != "Basic Q0xJRU5UX0lEOkNMSUVOVF9TRUNSRVQ=" {
			t.Errorf("Unexpected authorization header, %v is found.", headerAuth)
		}
		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/x-www-form-urlencoded" {
			t.Errorf("Unexpected Content-Type header, %v is found.", headerContentType)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed reading request body: %s.", err)
		}
		if string(body) != "code=exchange-code&grant_type=authorization_code&redirect_uri=REDIRECT_URL" {
			t.Errorf("Unexpected exchange payload, %v is found.", string(body))
		}
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer"))
	})
    ts := httptest.NewServer(mux)
    defer ts.Close()
    // oaConfig = &oauth2.Config{
	// 	ClientID:     "CLIENT_ID",
	// 	ClientSecret: "CLIENT_SECRET",
	// 	RedirectURL:  "REDIRECT_URL",
	// 	Scopes:       []string{"scope1", "scope2"},
	// 	Endpoint:     oauth2.Endpoint{
	// 		AuthURL:  ts.URL + "/auth",
	// 		TokenURL: ts.URL + "/token",
	// 	},
    // }

    h := RepoServer{}

    // outputURL := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOffline)

    req, err := http.NewRequest("POST", "/handler-check", nil)
	if err != nil {
		t.Errorf("Failure generating test request: %v", err)
	}
    req.Form = url.Values{}
    req.Form.Set("state", oaState)

    rec := httptest.NewRecorder()

    handler := http.HandlerFunc(h.githubCallbackHandle)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusPermanentRedirect {
		t.Errorf("Handler returning incorrect status: %v", status)
	}



}
*/

/*
func Test_githubCallbackHandle(t *testing.T) {
    h := RepoServer{}
    req, err := http.NewRequest("POST", "/handler-check", nil)
	if err != nil {
		t.Errorf("Failure generating test request: %v", err)
	}
    req.Form = url.Values{}
    req.Form.Set("state", oaState)
    rec := httptest.NewRecorder()
    handler := http.HandlerFunc(h.githubCallbackHandle)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusPermanentRedirect {
		t.Errorf("Handler returning incorrect status: %v", status)
	}
}
*/
