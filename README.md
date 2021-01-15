# TerReader
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>
[![Go Report Card](https://goreportcard.com/badge/github.com/will-evil/terreader)](https://goreportcard.com/report/github.com/will-evil/terreader)
[![Go Reference](https://pkg.go.dev/badge/github.com/will-evil/terreader.svg)](https://pkg.go.dev/github.com/will-evil/terreader)

Package for reading terrorists database by the records, not by lines.
Records in the database consist of multiple lines and if you want to get the correct record you need to join these lines correctly. But this package can do this work instead of you.

## Installing
Use `go get` for install package.

```bash
go get github.com/will-evil/terreader
```
After than include package to your project.

```
import "github.com/will-evil/terreader"
```

## Getting Started

```
package main

import (
	"fmt"
	"log"

	"github.com/will-evil/terreader"
)

func main() {
	tr, err := terreader.NewTerReader("/home/user/path_to_yor_file/file.dbf", "866")
	if err != nil {
		log.Fatal(err)
	}

	results, err := tr.Read(5)
	if err != nil {
		log.Fatal(err)
	}

	for res := range results {
		if res.Error != nil {
			errText := fmt.Sprintf("error for record with number '%d'. Error: '%s'", res.Number, res.Error)
			log.Fatal(errText)
		}

		fmt.Printf("%+v\n", *res.Row)
	}
}
```

## Article about the package

[Article on Medium.com](https://medium.com/rnds/114e9f6fadbb)

[Article on Teletype.in](https://blog.rnds.pro/W-GiIYj95)
