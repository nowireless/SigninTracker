package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"signin3/api"
	"signin3/database"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port    int
	Network string

	ClientPath string

	API api.Config
}

func main() {
	// Logging setupt
	if runtime.GOOS == "windows" {
		// Setup logging for windows
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	}

	// Configuration
	config := Config{}
	config.Port = 8080
	config.Network = ""

	config.API = api.Config{}
	config.API.App.Database = database.Config{
		User:     "signin",
		Password: "foobar",
		Host:     "localhost",
		Port:     5432,
		Database: "signin",
	}

	// Start server with signal handing
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
		serve(ctx, config)
	}()

	wg.Wait()
	log.Info("Exiting")
}

func serve(ctx context.Context, config Config) {
	r := mux.NewRouter()

	// API setup
	log.Info("Initializing API")
	err := api.NewAPI(config.API).Initialize(r.PathPrefix("/api/v1").Subrouter())
	if err != nil {
		log.Error(err)
		log.Fatal("Failed to initializecpo API")
	}
	log.Info("API initialized")

	// Client setup
	log.Info("Client Path: ", config.ClientPath)
	r.PathPrefix("/client/").Handler(http.StripPrefix("/client/", http.FileServer(http.Dir(config.ClientPath))))

	// Middleware setup
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

	r.Use(handlers.RecoveryHandler(
		handlers.RecoveryLogger(&loggerImpl{}),
		handlers.PrintRecoveryStack(true),
	))

	// Server setup
	address := fmt.Sprintf("%s:%d", config.Network, config.Port)
	s := &http.Server{
		Addr:        address,
		Handler:     r,
		ReadTimeout: 2 * time.Minute,
	}

	// Start server
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

// Recovery
type loggerImpl struct{}

func (l *loggerImpl) Println(args ...interface{}) {
	log.Error(append([]interface{}{"Panic in handler: "}, args...))
}
