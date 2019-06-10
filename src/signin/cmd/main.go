package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"signin/api"
	"signin/app"
	"sync"
	"time"

	"github.com/shiena/ansicolor"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	if runtime.GOOS == "windows" {
		// Setup logging for windows
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	}

	// Create application
	a, err := app.New()
	if err != nil {
		log.Panic(err)
	}
	defer a.Close()

	// Create RESTful API
	api, err := api.New(a)
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)

		// Wait for interrupt
		<-ch
		log.Info("Singal caught. Stopping...")
		cancel()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		// Start Serving API and Client
		serve(ctx, api)
	}()

	wg.Wait()
	log.Info("Exiting")
}

func serve(ctx context.Context, api *api.API) {
	r := mux.NewRouter()

	// Serve Client files
	// TODO Make client dir configurable
	clientDir := "./client"
	r.PathPrefix("/client/").Handler(http.StripPrefix("/client/", http.FileServer(http.Dir(clientDir))))

	// Serve API
	api.Init(r.PathPrefix("/api").Subrouter())

	/*
	 * Middleware Setup
	 */
	// Cross-Origin Resource Sharing middleware Setup
	r.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	))

	// Logging
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Request ", r.Method, r.RequestURI, " from ", r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	})

	/*
	 * Server Setup
	 */
	s := &http.Server{
		Addr:        fmt.Sprintf(":%d", api.Config.Port),
		Handler:     r,
		ReadTimeout: 2 * time.Minute,
	}

	// Server Start
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			log.Error(err)
		}
		close(done)
	}()

	log.Info("Listening on ", s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
	<-done
}
