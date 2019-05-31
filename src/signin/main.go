package main

import (
	"net/http"
	"os"
	"runtime"
	"signin/api"

	"github.com/shiena/ansicolor"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func main() {
	if runtime.GOOS == "windows" {
		// Setup logging for windows
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	}

	clientDir := "./client"

	r := mux.NewRouter()
	r.PathPrefix("/client/").Handler(http.StripPrefix("/client/", http.FileServer(http.Dir(clientDir))))

	apiRouter := r.PathPrefix("/api").Subrouter()
	api.InitSubRouter(apiRouter)

	log.Info("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
