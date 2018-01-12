package frontend

import (
	"html/template"
	"net/http"
)

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	if PROD {
		http.Redirect(w, r, "https://heupr.io", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "https://127.0.0.1:8081", http.StatusMovedPermanently)
	}
}

// baseHTML is necessary for unit testing purposes.
var baseHTML = "../templates/base.html"

func render(filepath string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles([]string{
			baseHTML,
			filepath,
		}...)
		if err != nil {
			slackErr("Error parsing "+filepath, err)
			http.Error(w, "error parsing static page", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, ""); err != nil {
			slackErr("Error rendering "+filepath, err)
			http.Error(w, "error rendering static page", http.StatusInternalServerError)
			return
		}
	})
}
