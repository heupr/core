package frontend

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/schema"
	"go.uber.org/zap"

	"core/utils"
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

func staticHandler(filepath string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			if PROD {
				utils.SlackLog.Error(
					"Error generating landing page",
					zap.Error(err),
				)
			}
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(data)))
	})
}

func setupCompleteHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("website2/setup-complete.html")
	if err != nil {
		if PROD {
			utils.SlackLog.Error(
				"Error generating setup complete page",
				zap.Error(err),
			)
		}
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	utils.AppLog.Info("Completed user signed up")
	if PROD {
		utils.SlackLog.Info("Completed user signed up")
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(data)))
}

// NOTE: Depreciate this code.
var decoder = schema.NewDecoder()

// NOTE: Depreciate this code.
var mainHandler = http.StripPrefix(
	"/",
	http.FileServer(http.Dir("../website/")),
)
