# Echo

<div align="center">
	<div align="right">
		<strong><code>v0.1.2</code></strong>
	</div>
</div>

[![Go Reference](https://pkg.go.dev/badge/github.com/boratanrikulu/echo.svg)](https://pkg.go.dev/github.com/boratanrikulu/echo)

## What

TCP Echo server. RFC862 implementation.

## Usage of package

```go
package main

import (
	"context"

	"github.com/boratanrikulu/echo"
)

func main() {
	s := echo.NewServer().
		Address(":1337").
		Banner(true).
		Verbose(false)

	ctx := context.Background()
	s.Run(ctx)
}
```

## Usage of CLI

```
               _
              | |
  ___    ___  | |__     ___
 / _ \  / __| | '_ \   / _ \
|  __/ | (__  | | | | | (_) |
 \___|  \___| |_| |_|  \___/

Usage: echo [options]
Options:
	-a,  --address,   which interface and port will be used  *[example: ":1337"]
	-n,  --no-banner
	-v,  --verbose
	-h,  --help

* means "must be set".
```

## Live Demo

It lives at **echo.bora.sh:1337**.  
You can try it using telnet.
```
telnet echo.bora.sh 1337
```

## Development

To run tests, 
```sh
go test ./... -v -cover -count=1 -race
```

To run server,
```sh
go run ./cmd --address :1337
```

To run server with Docker,
```sh
docker build -t echo .
```
```sh
docker run --rm -i -t -p 1337:1337 echo --address :1337
```
