# copilot-go/github

copilot-go/github is a companion Go library for use with copilot-go. It provides a set of utilities for working with the GitHub API.

[![Go Reference](https://pkg.go.dev/badge/github.com/colbylwilliams/copilot-go/github.svg)](https://pkg.go.dev/github.com/colbylwilliams/copilot-go/github)

## Installing

copilot-go is compatible with modern Go releases in module mode, with Go installed:

```sh
go get github.com/colbylwilliams/copilot-go/github
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/colbylwilliams/copilot-go/github"
```

and run `go get` without parameters.

Finally, to use the top-of-trunk version of this repo, use the following command:

```sh
go get github.com/colbylwilliams/copilot-go/github@main
```

## Usage

See the [/_examples/github](/_examples/gitbub) directory for a complete and runnable example.
