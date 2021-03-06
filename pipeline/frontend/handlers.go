package frontend

import (
	"bytes"
	"context"
	"encoding/gob"
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

func init() {
	if PROD {
		appID = 5535
		oauthConfig = &oauth2.Config{
			RedirectURL:  "https://heupr.io/repos",
			ClientID:     "Iv1.08a7e522bf043e73",
			ClientSecret: "",
			Endpoint:     ghoa.Endpoint,
		}
		domain = "https://heupr.io"
	} else {
		appID = 6807 //This needs to match the "ID" in "Mike/JohnHeuprTest"
		oauthConfig = &oauth2.Config{
			RedirectURL:  "https://127.0.0.1:8081/repos",             //This needs to match the "User authorization callback URL" in "Mike/JohnHeuprTest"
			ClientID:     "",                     //This needs to match the "ClientID" in "Mike/JohnHeuprTest"
			ClientSecret: "", //This needs to match the "ClientSecret" in "Mike/JohnHeuprTest"
			Endpoint:     ghoa.Endpoint,
		}
		domain = "https://127.0.0.1:8081"
	}
}

var (
	appID                int
	oauthConfig          *oauth2.Config
	store                = sessions.NewCookieStore([]byte("add key"))
	oauthTokenSessionKey = "oauth_token"
	// templatePath is for testing purposes only; a better solution is needed.
	templatePath = "../"
	gobPath      = "gob/"
	domain       string
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

	oauthFlowSession.Options.MaxAge = 0

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

var newServerToServerClient = func(appID int, installationID int64) (*github.Client, error) {
	var key string
	if PROD {
		key = "heupr.2017-10-04.private-key.pem"
	} else {
		key = "mikeheuprtest.2017-11-16.private-key.pem"
	}
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, int(installationID), key)
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
	RepoID   int64
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

	_, err := store.Get(r, r.FormValue("state"))
	if err != nil {
		utils.AppLog.Error("invalid state", zap.Error(err))
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
	installationFound := false
	var installationIDs []int64
	installations, _, err := listUserInstallations(ctx, client, opts)

	for i := range installations {
		if *installations[i].AppID == appID {
			installationFound = true
			installationIDs = append(installationIDs, *installations[i].ID)
		}
	}
	if !installationFound {
		utils.AppLog.Warn("heupr installation not found")
		http.Error(w, "error detecting heupr installation", http.StatusInternalServerError)
		return
	}

	repos := make(map[int64]string)


	for i := range installationIDs {
		for {
			repo, resp, err := client.Apps.ListUserRepos(ctx, installationIDs[i], opts)
			if err != nil {
				utils.AppLog.Error("error collecting user repos", zap.Error(err))
				http.Error(w, "error collecting user repos", http.StatusInternalServerError)
				return
			}
			for j := range repo {
				repos[*repo[j].ID] = *repo[j].FullName
			}

			if resp.NextPage == 0 {
				break
			} else {
				opts.Page = resp.NextPage
			}
		}
	}

	labels := make(map[int64][]string)
	for i := range installationIDs {
		client, err = newServerToServerClient(appID, installationIDs[i])

		if err != nil {
			utils.AppLog.Error("could not obtain github installation key", zap.Error(err))
			http.Error(w, "client failure", http.StatusInternalServerError)
			return
		}

		opts = &github.ListOptions{PerPage: 100}
		for key, value := range repos {
			name := strings.Split(value, "/")
			for {
				l, resp, err := client.Issues.ListLabels(ctx, name[0], name[1], opts)
				if err != nil {
					utils.AppLog.Error("error collecting repo labels", zap.Error(err))
					http.Error(w, "error collecting repo labels", http.StatusInternalServerError)
					return
				}
				for k := range l {
					labels[key] = append(labels[key], *l[k].Name)
				}

				if resp.NextPage == 0 {
					break
				} else {
					opts.Page = resp.NextPage
				}
			}
		}
	}

	for id, name := range repos {
		filename := gobPath + strconv.FormatInt(id, 10) + ".gob"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			file, err := os.Create(filename)
			defer file.Close()
			if err != nil {
				utils.AppLog.Error("error creating storage file", zap.Error(err))
				http.Error(w, "error creating storage file", http.StatusInternalServerError)
				return
			}

			s := storage{
				RepoID:   id,
				FullName: name,
				Buckets:  make(map[string][]label),
			}

			for _, l := range labels[id] {
				s.Buckets["typedefault"] = append(s.Buckets["typedefault"], label{Name: l})
				s.Buckets["typebug"] = append(s.Buckets["typebug"], label{Name: l})
				s.Buckets["typeimprovement"] = append(s.Buckets["typeimprovement"], label{Name: l})
				s.Buckets["typefeature"] = append(s.Buckets["typefeature"], label{Name: l})
			}

			encoder := gob.NewEncoder(file)
			if err := encoder.Encode(s); err != nil {
				utils.AppLog.Error("error encoding info to new file", zap.Error(err))
				http.Error(w, "error encoding info to new file", http.StatusInternalServerError)
				return
			}
		} else {
			file, err := os.Open(filename)
			defer file.Close()
			if err != nil {
				http.Error(w, "error opening storage file", http.StatusInternalServerError)
				return
			}
			decoder := gob.NewDecoder(file)
			s := storage{}
			err = decoder.Decode(&s)
			if err != nil {
				utils.AppLog.Error("eerror opening storage file", zap.Error(err))
				http.Error(w, "error opening storage file", http.StatusInternalServerError)
				return
			}
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
		Repos map[int64]string
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
		"storage":        input,
		"csrf":           csrf.Token(r),
		csrf.TemplateTag: csrf.TemplateField(r),
		"domain":         domain,
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		slackErr("Repos selection page", err)
		http.Error(w, "error loading repo selections", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
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

	file := gobPath + repoID + ".gob"
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

	t, err := template.ParseFiles(
		templatePath+"templates/base.html",
		templatePath+"templates/console.html",
	)
	if err != nil {
		slackErr("Settings console page", err)
		utils.AppLog.Error("settings console page", zap.Error(err))
		http.Error(w, "error loading console", http.StatusInternalServerError)
		return
	}

	csrfToken := csrf.Token(r)
	data := map[string]interface{}{
		"storage":        s,
		"csrf":           csrfToken,
		csrf.TemplateTag: csrf.TemplateField(r),
		"domain":         domain,
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		slackErr("Settings console page", err)
		http.Error(w, "error loading console", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func updateSettings(s *storage, form map[string][]string) {
	//Set everything to False
	for _, bucket := range s.Buckets {
		for i := range bucket {
			bucket[i].Selected = false
		}
	}
	//Mark selections to true
	for key, bucket := range s.Buckets {
		for i := range bucket {
			for _, selection := range form[key] {
				if bucket[i].Name == selection {
					bucket[i].Selected = true
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

	file := gobPath + repoID.(string) + ".gob"
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
	ingestorFile := utils.Config.IngestorGobs + "/" + repoID.(string) + ".gob"
	CopyFile(ingestorFile, file, 0664)
	t, err := template.ParseFiles(
		templatePath+"templates/base.html",
		templatePath+"templates/complete.html",
	)
	if err != nil {
		slackErr("Error parsing signup complete page", err)
		http.Error(w, "error parsing signup complete page", http.StatusInternalServerError)
		return
	}

	csrfToken := csrf.Token(r)
	data := map[string]interface{}{
		"csrf":           csrfToken,
		csrf.TemplateTag: csrf.TemplateField(r),
		"domain":         domain,
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		slackErr("Error rendering template", err)
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
	utils.AppLog.Info("Preparing to write out complete handler")
	buf.WriteTo(w)
	utils.AppLog.Info("Completed user signed up")
	slackMsg("Completed user signed up")
}
