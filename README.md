# TerReader

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
			log.Fatal(res.Error)
		}

		fmt.Printf("%+v\n", *res.Row)
	}
}
```
