package frontend

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	ghoa "golang.org/x/oauth2/github"

	"core/utils"
)

func slackErr(msg string, err error) {
	if PROD {
		utils.SlackLog.Error(msg, zap.Error(err))
	}
}

func slackMsg(msg string) {
	if PROD {
		utils.SlackLog.Info(msg)
	}
}

var (
	oauthConfig = &oauth2.Config{
		// NOTE: These will need to be added for production.
		// TODO: Try to configure RedirectURL without using Ngrok
		RedirectURL:  "http://127.0.0.1:8080/repos",              //This needs to match the "User authorization callback URL" in "Mike/JohnHeuprTest"
		ClientID:     "Iv1.83cc17f7f984aeec",                     //This needs to match the "ClientID" in "Mike/JohnHeuprTest"
		ClientSecret: "c9c5f71edcf1a85121ae86bae5295413dff46fad", //This needs to match the "ClientSecret" in "Mike/JohnHeuprTest"
		// Scopes:       []string{""},
		Endpoint: ghoa.Endpoint,
	}
	appID			 = 6807 //This needs to match the "ID" in "Mike/JohnHeuprTest"
	store                = sessions.NewCookieStore([]byte("yoda-dooku-jinn-kenobi-skywalker-tano"))
	oauthTokenSessionKey = "oauth_token"
	// templatePath is for testing purposes only; a better solution is needed.
	templatePath = "../"
)

const sessionName = "heupr-session"

func login(w http.ResponseWriter, r *http.Request) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "authorization error", http.StatusUnauthorized)
		return
	}
	sessionID := newUUID.String()
	oauthFlowSession, err := store.New(r, sessionID)
	if err != nil {
		http.Error(w, "authorization error", http.StatusUnauthorized)
		return
	}

	//TODO: Experiment with this setting.
	oauthFlowSession.Options.MaxAge = 1 * 60 // 1 minute

	// Use session ID for state params which protects against CSRF.
	redirectURL := oauthConfig.AuthCodeURL(sessionID)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

var newUserToServerClient = func(code string) (*github.Client, error) {
	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	client := github.NewClient(oauthConfig.Client(oauth2.NoContext, token))
	return client, nil
}

var newServerToServerClient = func(appId, installationId int) (*github.Client, error) {
	var key string
	if PROD {
		key = "heupr.2017-10-04.private-key.pem"
	} else {
		key = "mikeheuprtest.2017-11-16.private-key.pem"
	}
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appId, installationId, key)
	if err != nil {
		return nil, err
	}
	client := github.NewClient(&http.Client{Transport: itr})
	return client, nil
}

type label struct {
	Name     string
	Selected bool
}

type storage struct {
	FullName string   `schema:"FullName"`
	Labels   []string `schema:"Labels"`
	Buckets  map[string][]label
}


//TODO: Confirm this is needed for adding labels.
//TODO: False makes it into the Post Request.
func updateStorage(s *storage, labels []string) {
	for bcktName, bcktLabels := range s.Buckets {
		updated := []label{}
		for i := range labels {
			label := label{Name: labels[i]}
			for j := range bcktLabels {
				if labels[i] == bcktLabels[j].Name {
					label.Selected = bcktLabels[j].Selected
				}
			}
			updated = append(updated, label)
		}
		s.Buckets[bcktName] = updated
	}
}

func repos(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "bad request method", http.StatusBadRequest)
		return
	}

	oauthFlowSession, err := store.Get(r, r.FormValue("state"))
	fmt.Println("oauthFlow session : ")

	if err != nil {
		fmt.Println("invalid state: ", oauthFlowSession)
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	code := r.FormValue("code")
	client, err := newUserToServerClient(code)

	if err != nil {
		utils.AppLog.Error("failure creating frontend client", zap.Error(err))
		http.Error(w, "client failure", http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, sessionName)
	session.Options.MaxAge = 0
	if err != nil {
		http.Error(w, "error establishing session", http.StatusInternalServerError)
		return
	}
	session.Save(r, w)

	ctx := context.Background()
	opts := &github.ListOptions{Page: 1, PerPage: 50}
	installationID := 0
	installations, _, err := listUserInstallations(ctx, client, opts)

	for i := range installations {
		if *installations[i].AppID == appID {
			installationID = *installations[i].ID
			break
		}
	}
	if installationID == 0 {
		utils.AppLog.Warn("heupr installation not found")
		http.Error(w, "error detecting heupr installation", http.StatusInternalServerError)
		return
	}

	repos := make(map[int]string)
	for {
		repo, resp, err := client.Apps.ListUserRepos(ctx, installationID, opts)
		if err != nil {
			utils.AppLog.Error("error collecting user repos", zap.Error(err))
			http.Error(w, "error collecting user repos", http.StatusInternalServerError)
			return
		}
		for i := range repo {
			repos[*repo[i].ID] = *repo[i].FullName
		}

		if resp.NextPage == 0 {
			break
		} else {
			opts.Page = resp.NextPage
		}
	}

	client, err = newServerToServerClient(appID, installationID)

	if err != nil {
		utils.AppLog.Error("could not obtain github installation key", zap.Error(err))
		http.Error(w, "client failure", http.StatusInternalServerError)
		return
	}

	opts = &github.ListOptions{PerPage: 100}
	labels := make(map[int][]string)
	for key, value := range repos {
		name := strings.Split(value, "/")
		for {
			l, resp, err := client.Issues.ListLabels(ctx, name[0], name[1], opts)
			if err != nil {
				utils.AppLog.Error("error collecting repo labels", zap.Error(err))
				http.Error(w, "error collecting repo labels", http.StatusInternalServerError)
				return
			}
			for i := range l {
				labels[key] = append(labels[key], *l[i].Name)
			}

			if resp.NextPage == 0 {
				break
			} else {
				opts.Page = resp.NextPage
			}
		}
	}

	for id, name := range repos {
		filename := "gob/" + strconv.Itoa(id) + ".gob"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			file, err := os.Create(filename)
			defer file.Close()
			if err != nil {
				utils.AppLog.Error("error creating storage file", zap.Error(err))
				http.Error(w, "error creating storage file", http.StatusInternalServerError)
				return
			}

			s := storage{
				FullName: name,
				Buckets:  make(map[string][]label),
			}

			for _, l := range labels[id] {
				s.Buckets["typedefault"] = append(s.Buckets["typedefault"], label{Name: l})
			}

			encoder := gob.NewEncoder(file)
			if err := encoder.Encode(s); err != nil {
				utils.AppLog.Error("error encoding info to new file", zap.Error(err))
				http.Error(w, "error encoding info to new file", http.StatusInternalServerError)
				return
			}
		} else {
			file, err := os.Open(filename)//os.OpenFile(filename, os.O_WRONLY, 0644)
			defer file.Close()
			if err != nil {
				http.Error(w, "error opening storage file", http.StatusInternalServerError)
				return
			}
			decoder := gob.NewDecoder(file)
			s := storage{}
			err = decoder.Decode(&s)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Storage2", s)

			updateStorage(&s, labels[id])
			file2, err := os.OpenFile(filename, os.O_WRONLY, 0644)
			defer file2.Close()
			encoder := gob.NewEncoder(file2)
			if err := encoder.Encode(s); err != nil {
				utils.AppLog.Error("error re-encoding info to file", zap.Error(err))
				http.Error(w, "error storing user info", http.StatusInternalServerError)
				return
			}
		}
	}

	input := struct {
		Repos map[int]string
	}{
		Repos: repos,
	}

	t, err := template.ParseFiles(
		templatePath+"templates/base.html",
		templatePath+"templates/repos.html",
	)

	if err != nil {
		slackErr("Repos selection page", err)
		http.Error(w, "error loading repo selections", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
			"storage":				input,
			"csrf":           csrf.Token(r),
			csrf.TemplateTag: csrf.TemplateField(r),
	}
	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		slackErr("Repos selection page", err)
		http.Error(w, "error loading repo selections", http.StatusInternalServerError)
		return
	}
}

func generateWalkFunc(file *string, repoID string) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if info.Name() == repoID+".gob" {
			*file = info.Name()
		}
		return nil
	}
}

func console(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "bad request method", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	repoID := r.Form["repo-selection"][0]
	if repoID == "" {
		http.Error(w, "request error", http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, sessionName)
	session.Options.MaxAge = 0
	if err != nil {
		http.Error(w, "error establishing session", http.StatusInternalServerError)
		return
	}
	session.Values["repoID"] = repoID
	session.Save(r, w)

	file := repoID + ".gob"
	_, err = os.Stat(file)
	if err != nil {
		utils.AppLog.Error("error retrieving user settings", zap.Error(err))
		http.Error(w, "error retrieving user settings", http.StatusInternalServerError)
		return
	}

	f, err := os.Open(file)
	if err != nil {
		http.Error(w, "error opening user settings", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	s := storage{}
	err = decoder.Decode(&s)
	if err != nil {
		utils.AppLog.Error("error decoding user settings", zap.Error(err))
		http.Error(w, "error decoding user settings", http.StatusInternalServerError)
		return
	}

	fmt.Println("Storage", s)

	t, err := template.ParseFiles("../templates/base.html","../templates/console.html")
	if err != nil {
		slackErr("Settings console page", err)
		fmt.Println(err)
		http.Error(w, "error loading console", http.StatusInternalServerError)
		return
	}

	csrfToken := csrf.Token(r)
	fmt.Println(csrfToken)
	data := map[string]interface{}{
			"storage":				s,
			"csrf":           csrfToken,
			csrf.TemplateTag: csrf.TemplateField(r),
	}
	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		slackErr("Settings console page", err)
		http.Error(w, "error loading console", http.StatusInternalServerError)
		return
	}
}

func updateSettings(s *storage, form map[string][]string) {
	for name, bucket := range s.Buckets {
		for i := range bucket {
			for j := range form[name] {
				if bucket[i].Name == form[name][j] {
					bucket[i].Selected = true
					break
				} else {
					bucket[i].Selected = false
				}
			}
		}
	}
}

func complete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "bad request method", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	session, err := store.Get(r, sessionName)
	session.Options.MaxAge = 0
	if err != nil {
		http.Error(w, "error establishing session", http.StatusInternalServerError)
		return
	}
	repoID := session.Values["repoID"]
	delete(session.Values, "repoID")
	session.Save(r, w)

	file := repoID.(string) + ".gob"
	_, err = os.Stat(file)
	if err != nil {
		utils.AppLog.Error("error retrieving user settings", zap.Error(err))
		http.Error(w, "error retrieving user settings", http.StatusInternalServerError)
		return
	}

	f, err := os.Open(file)
	if err != nil {
		http.Error(w, "error opening user settings", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	s := storage{}
	err = decoder.Decode(&s)
	if err != nil {
		http.Error(w, "error decoding user settings", http.StatusInternalServerError)
		return
	}

	updateSettings(&s, r.Form)

	file2, err := os.OpenFile(file, os.O_WRONLY, 0644)
	defer file2.Close()
	encoder := gob.NewEncoder(file2)
	if err := encoder.Encode(s); err != nil {
		utils.AppLog.Error("error re-encoding info to file", zap.Error(err))
		http.Error(w, "error storing user info", http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles(
		templatePath+"templates/base.html",
		templatePath+"templates/complete.html",
	)
	if err != nil {
		slackErr("Error parsing signup complete page", err)
		http.Error(w, "error parsing signup complete page", http.StatusInternalServerError)
		return
	}
	buf := new(bytes.Buffer)
	err = t.Execute(w, "")
	if err != nil {
		slackErr("Error rendering template", err)
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
	utils.AppLog.Info("Completed user signed up")
	slackMsg("Completed user signed up")
	buf.WriteTo(w)
}
