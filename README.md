# copilot-go

copilot-go is a Go library for GitHub Copilot extensions.

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

