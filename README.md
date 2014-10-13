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
  n := nessusgo.New("https://nessushost:8834")

  if err := n.Login("user", "password"); err == nil {
    fmt.Printf("Login successful\n")
  }

  if record, err := n.Feed(); err == nil {
    fmt.Printf("Nessus Server Version: %v\n", record.Reply.Contents.ServerVersion)
  } else {
    fmt.Printf("Error: %v\n", err)
  }

  if record, err := n.Uuid(); err == nil {
    fmt.Printf("Server UUID: %s\n", *record)
  } else {
    fmt.Printf("Error: %v\n", err)
  }

  if record, err := n.Server.GetLoad(); err == nil {
    fmt.Printf("Server Platform: %v\n", record.Reply.Contents.Platform)
    fmt.Printf("Load Average: %v\n", record.Reply.Contents.Load.LoadAverage)
  } else {
    fmt.Printf("Error: %v\n", err)
  }

  if records, err := n.Users.List(); err == nil {
    fmt.Printf("Users: %v\n", len(*records))
    for _,e := range *records {
      fmt.Printf("\tName: %s\tAdmin: %v\tLast Login: %v\n", e.Name, e.IsAdmin, e.LastLogin)
    }
  } else {
    fmt.Printf("Error: %v\n", err)
  }

  if err := n.Logout(); err == nil {
    fmt.Printf("Logout successful\n")
  } else {
    fmt.Printf("Error: %v\n", err)
  }
}

```

This can be ran from Go by running:

```
$ go run test.go
Login successful
Nessus Server Version: 5.0.1
Server UUID: 151fa290-3618-f71d-146b-b145dc22f0e545f2a6dc426b05c7
Server Platform: LINUX
Load Average: 0
Users: 1
	Name: admin	Admin: true	Last Login: 2014-10-12 21:48:40 -0400 EDT
Logout successful
```

Alternatively, it can be compiled to a binary for your OS:

```
$ go build test.go
```