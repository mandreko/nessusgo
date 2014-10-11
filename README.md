nessusgo
========

Nessus API Library written in Go

Installation
============

Run "go get" pointing to the package:

```
$ go get github.com/mandreko/nessusgo
```

Or clone the source code directly from GitHub:

```
$ git clone https://github.com/mandreko/nessusgo.git
```

Basic Usage
===========

An example of using the API library can be created in a standalone Go file:

```
package main

import (
	"fmt"
	"github.com/mandreko/nessusgo"
)

func main() {
	client := nessusgo.NewClient("https://nessushost:8834")

	client.Authenticate("user", "password")

	// List Scans
	scans := client.ListScans()

	client.LogOut()

	fmt.Printf("Scans: %v\n", scans)
}

```

This can be ran from Go by running:

```
$ go run test.go
```

Alternatively, it can be compiled to a binary for your OS:

```
$ go build test.go
```