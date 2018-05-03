# SurfDAP - Go LDAP client for human beings

[![GoDoc](https://godoc.org/github.com/octogo/surfdap?status.svg)](https://godoc.org/github.com/octogo/surfdap)

## Installation

```bash
$ go get -v -u github.com/octogo/surfdap
$ go install github.com/octogo/surfdap/surfdap
```

## Usage

### Shell
```bash
$ surfdap search --help
```

### Embed
```go
package main

import (
    "fmt"
    "os"

    "github.com/octogo/surfdap"
)

func main() {
    // Obtain surfdap configuration
    config := surfdap.Config{
        Host: "localhost",
        Port: 389,
        BaseDN: dc=example,dc=org,
    }

    // Obtain root node.
    root, err := surfdap.New(config)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Navigate tree using root node.
    fmt.Println(root.DN())          // DN of the root node
    fmt.Println(root.Attributes())  // map[string][]string if all attributes of root node
    fmt.Println(root.Children())    // []Node with all children of root node
    fmt.Println(
        root.Children()[0].Parent() // Parent node of first child equals root node
    )
}
```
