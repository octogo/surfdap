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
    surfer, err := surfdap.New("localhost", 389, false, "dc=example,dc=com", "", "")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println(surfer)
    // dn: dc=example,dc=org
    // objectClass: dcObject
    // objectClass: organization
    // dc: example
    // o: Example Org

    surfer.Entry()
    // Return *ldap.Entry of the underlaying LDAP object.

    surfer.Parent()
    // Returns a surfdap.Surfer for the parent object of this object.

    surfer.Lookup(
        surfdap.One,
        surfdap.Filter("(objectClass=*)"),
        nil,
    )
}
```
