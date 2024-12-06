# copilot-go

copilot-go copilot provides the types and functions for authoring [GitHub Copilot Extensions].

Use this package to build GitHub Copilot [agents] and [skillsets].

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

This simple example agent uses [`chi`][chi] and [`go-github`][go-github]:

```go
package main

import (
    "context"
    "net/http"

    "github.com/colbylwilliams/copilot-go"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/google/go-github/v67/github"
)

func main() {
    cfg, _ := copilot.LoadConfig()
    v, _ := copilot.NewPayloadVerifier()

    a := &MyAgent{cfg: cfg}

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.RequestID)

    r.Post("/agent", copilot.AgentHandler(v, a))

    return http.ListenAndServe(":3333", r)
}

type MyAgent struct{ cfg *copilot.Config }

func (a *MyAgent) Execute(ctx context.Context, token string, req *copilot.Request, w http.ResponseWriter) error {
    rid := middleware.GetReqID(ctx)
    gh := github.NewClient(nil).WithAuthToken(token)

    me, _, _ := gh.Users.Get(ctx, "")

    sse.WriteStreamingHeaders(w)

    sse.WriteDelta(w, rid, fmt.Sprintf("Hello %s\n", me.GetName()))
    sse.WriteDelta(w, rid, "How are you today?\n")

    sse.WriteStop(w, rid)
}
```

### Configuration

[`LoadConfig`][LoadConfig] above loads .env file(s) with [godotenv], then loads the required configuration information from environment variables. Here's an example `.env` file with the required variables:

```sh
GITHUB_APP_CLIENT_ID=Iv23abcdef1234567890
GITHUB_APP_PRIVATE_KEY_PATH=my-agent.2024-08-20.private-key.pem
GITHUB_APP_FQDN=https://my_unique_devtunnelid-3333.use2.devtunnels.ms
```

These values will come from your Copilot extension's GitHub App. See _"[Creating a GitHub App for your Copilot Extension]"_ and _"[Configuring your GitHub App for your Copilot extension]"_ for details.

See the [`Config`][Config] and [`azure.Config`][azure.Config] types for a full set of available configuration variables.

### Copilot's LLM

In the first example, the agent replies with a static message. We can easily enhance it to use Copilot's LLM to generate the response.

```go
package main

import (
    "context"
    "net/http"

    "github.com/colbylwilliams/copilot-go"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/google/go-github/v67/github"
)

func main() {
    cfg, _ := copilot.LoadConfig()
    v, _ := copilot.NewPayloadVerifier()

    a := &MyAgent{cfg: cfg}

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.RequestID)

    r.Post("/agent", copilot.AgentHandler(v, a))

    return http.ListenAndServe(":3333", r)
}

type MyAgent struct{ cfg *copilot.Config }

func (a *MyAgent) Execute(ctx context.Context, token string, req *copilot.Request, w http.ResponseWriter) error {

    const prompt = "You are a helpful AI assistant."

    msg := &copilot.Message{Role: copilot.ChatRoleSystem, Content: prompt}

    comp := &copilot.CompletionsRequest{
        Messages: []*copilot.Message{msg},
        Model:    copilot.CopilotModelGPT4o,
    }

    comp.Messages = append(comp.Messages, req.Messages...)

    stream, _ := copilot.ChatCompletionsStream(ctx, token, *comp, w)

    _, err = io.Copy(w, stream)

    return nil
}
```


[GitHub Copilot Extensions]: https://github.com/features/copilot/extensions
[skillsets]: https://docs.github.com/copilot/building-copilot-extensions/building-a-copilot-agent-for-your-copilot-extension/about-copilot-agents
[agents]: https://docs.github.com/copilot/building-copilot-extensions/building-a-copilot-agent-for-your-copilot-extension/about-copilot-agents
[chi]: https://pkg.go.dev/github.com/go-chi/chi/v5
[godotenv]: https://pkg.go.dev/github.com/joho/godotenv
[go-github]: https://pkg.go.dev/github.com/google/go-github/v67
[Creating a GitHub App for your Copilot Extension]: https://docs.github.com/en/copilot/building-copilot-extensions/creating-a-copilot-extension/creating-a-github-app-for-your-copilot-extension
[Configuring your GitHub App for your Copilot extension]: https://docs.github.com/en/copilot/building-copilot-extensions/creating-a-copilot-extension/configuring-your-github-app-for-your-copilot-extension
[LoadConfig]: https://pkg.go.dev/github.com/colbylwilliams/copilot-go#LoadConfig
[Config]: https://pkg.go.dev/github.com/colbylwilliams/copilot-go#Config
[azure.Config]: https://pkg.go.dev/github.com/colbylwilliams/copilot-go/azure#Config