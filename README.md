# go-update: Build self-updating Go programs

[![Go Reference](https://pkg.go.dev/badge/github.com/contriboss/go-update.svg)](https://pkg.go.dev/github.com/contriboss/go-update)

> Fork of [inconshreveable/go-update](https://github.com/inconshreveable/go-update) (unmaintained since 2016).

Package update provides functionality to implement secure, self-updating Go programs (or other single-file targets)
A program can update itself by replacing its executable file with a new version.

It provides the flexibility to implement different updating user experiences
like auto-updating, or manual user-initiated updates. It also boasts
advanced features like binary patching and code signing verification.

Example of updating from a URL:

```go
import (
    "net/http"

    "github.com/contriboss/go-update"
)

func doUpdate(url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return update.Apply(resp.Body, update.Options{})
}
```

## Features

- Cross platform support (Windows too!)
- Binary patch application
- Checksum verification
- Code signing verification
- Support for updating arbitrary files

## Installation

```bash
go get github.com/contriboss/go-update
```

Requires Go 1.25+.

## License
Apache
