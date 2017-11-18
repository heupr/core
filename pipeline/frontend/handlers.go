package frontend

import (
	"net/http"

	"github.com/gorilla/schema"

	"core/utils"
)

const setup = `
<html>
    <title>
        Heupr
    </title>
    <body>
        <p>Awesome! Setup is complete!</p>
        <p>Issue assignments will go out in a few minutes through GitHub</p>
        <p>Return to the <a href="/">main page</a></p>
    </body>
</html>
`

var decoder = schema.NewDecoder()

var mainHandler = http.StripPrefix(
	"/",
	http.FileServer(http.Dir("../website/")),
)

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	if PROD {
		http.Redirect(w, r, "https://heupr.io", http.StatusMovedPermanently)
	} else {
		http.Redirect(
			w,
			r,
			"https://127.0.0.1:8081",
			http.StatusMovedPermanently,
		)
	}
}

func setupCompleteHandler(w http.ResponseWriter, r *http.Request) {
	utils.AppLog.Info("Completed user signed up")
	if PROD {
		utils.SlackLog.Info("Completed user signed up")
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(setup))
}
