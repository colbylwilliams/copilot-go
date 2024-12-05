package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/colbylwilliams/copilot-go"
	"github.com/colbylwilliams/copilot-go/_examples/copilot_api/agent"
	"github.com/colbylwilliams/copilot-go/_examples/copilot_api/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envFile      = ".env"
	defaultPort  = "3333"
	readTimeout  = 5 * time.Second   // 5 seconds
	writeTimeout = 300 * time.Second // 5 minutes
)

func main() {
	if err := realMain(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("failed to run service:", err)
		os.Exit(1)
	}
}

//nolint:maintidx // main can have a lot of code.
func realMain() error {
	fmt.Println("Starting api")

	fmt.Println("loading config from", envFile)

	// load the config
	cfg, err := copilot.LoadConfig(envFile)
	if err != nil {
		return err
	}
	if cfg.HTTPPort == "" {
		fmt.Println("no PORT environment variable specified, defaulting to", defaultPort)
		cfg.HTTPPort = defaultPort
	}
	fmt.Println("using port:", cfg.HTTPPort)

	// create the payload verifier
	verifier, err := copilot.NewPayloadVerifier()
	if err != nil {
		return fmt.Errorf("failed to create payload authenticator: %w", err)
	}

	myagent := agent.NewAgent(cfg)

	// create the router
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Get("/_ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	router.Post("/events", copilot.WebhookHandler)

	router.Post("/agent", copilot.AgentHandler(verifier, myagent))

	authHandlers := &auth.AuthHandlers{
		ClientID: cfg.GitHubAppClientID,
		Callback: cfg.GitHubAppFQDN + "/auth/callback",
	}

	router.Route("/auth", func(r chi.Router) {
		r.Get("/authorization", authHandlers.PreAuth)
		r.Get("/callback", authHandlers.PostAuth)
	})

	addr := ":" + cfg.HTTPPort
	if cfg.IsDevelopment() {
		addr = "127.0.0.1" + addr // Prevents MacOS from prompting you about accepting network connections.
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	fmt.Println("Starting server on port " + cfg.HTTPPort)

	return server.ListenAndServe()
}