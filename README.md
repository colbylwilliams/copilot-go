# copilot-go

copilot-go is a Go library for GitHub Copilot extensions.

[![Go Reference](https://pkg.go.dev/badge/github.com/colbylwilliams/copilot-go.svg)](https://pkg.go.dev/github.com/colbylwilliams/copilot-go)

## Installing

copilot-go is compatible with modern Go releases in module mode, with Go installed:

```sh
go get github.com/colbylwilliams/copilot-go
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/colbylwilliams/copilot-go"
```

and run `go get` without parameters.

Finally, to use the top-of-trunk version of this repo, use the following command:

```sh
go get github.com/colbylwilliams/copilot-go@main
```

## Usage

See the [_examples](./_examples/) directory for complete and runnable examples.

```go
package agent

import (
    "github.com/colbylwilliams/copilot-go"
    "github.com/colbylwilliams/copilot-go/sse"
)

type MyAgent struct {
    cfg *copilot.Config
}

func NewAgent(cfg *copilot.Config) *MyAgent {
    return &MyAgent{
        cfg: cfg,
    }
}

func (a *MyAgent) Execute(ctx context.Context, token string, req *copilot.Request, w http.ResponseWriter) error {

    session := copilot.GetSession(ctx)
    if session == nil {
        return fmt.Errorf("session not found in context")
    }

    // write the sse headers
    sse.WriteStreamingHeaders(w)

    for _, m := range req.Messages {
        // ...
    }

    // respond to user messages

    return nil
}
```

```go
package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "github.com/colbylwilliams/copilot-go"
)

func main() {
    // load the config from .env file
    cfg, err := copilot.LoadConfig()
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

    router.Post("/events", copilot.WebhookHandler)
    router.Post("/agent", copilot.AgentHandler(verifier, myagent))

    fmt.Println("Starting server on port " + cfg.HTTPPort)

    return http.ListenAndServe(":" + cfg.HTTPPort, router)
}
```