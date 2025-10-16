# gnfmt

[![Go Report Card](https://goreportcard.com/badge/github.com/gnames/gnfmt)](https://goreportcard.com/report/github.com/gnames/gnfmt)
[![GoDoc](https://godoc.org/github.com/gnames/gnfmt?status.svg)](https://godoc.org/github.com/gnames/gnfmt)

`gnfmt` is a Go helper library designed to simplify data serialization and enhance command-line interface (CLI) output. It provides tools for converting Go data structures into various formats.

## Features

- **Data Serialization:** Convert Go objects to JSON (compact/pretty), CSV, TSV, and Gob formats
- **Pretty Printing:** Display Go objects in a human-readable JSON format in the terminal
- **CSV/TSV Utilities:** Read headers, convert records, and normalize row sizes
- **Time Formatting:** Convert seconds into human-readable duration strings
- **Flexible Encoders:** Pluggable encoder interface for easy format switching

## Installation

```shell
go get github.com/gnames/gnfmt
```

## Usage

### Data Serialization and Formatting

`gnfmt` supports multiple output formats for your data structures.

#### Format Types

```go
package main

import (
	"fmt"
	"github.com/gnames/gnfmt"
)

func main() {
	// Create format from string
	format, err := gnfmt.NewFormat("pretty")
	if err != nil {
		panic(err)
	}

	fmt.Println(format) // Output: "pretty JSON"

	// Available formats:
	// - "csv"     -> CSV format
	// - "tsv"     -> TSV (tab-separated) format
	// - "compact" -> Compact JSON (single line)
	// - "pretty"  -> Pretty JSON (indented)
}
```

#### JSON Encoding

```go
package main

import (
	"fmt"
	"github.com/gnames/gnfmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	person := Person{Name: "Alice", Age: 30}

	// Compact JSON
	encoder := gnfmt.GNjson{Pretty: false}
	data, _ := encoder.Encode(person)
	fmt.Println(string(data)) // {"name":"Alice","age":30}

	// Pretty JSON
	encoder.Pretty = true
	data, _ = encoder.Encode(person)
	fmt.Println(string(data))
	// {
	//   "name": "Alice",
	//   "age": 30
	// }

	// Quick pretty print
	fmt.Println(gnfmt.Ppr(person))
}
```

#### Gob Encoding

```go
package main

import (
	"fmt"
	"github.com/gnames/gnfmt"
)

func main() {
	data := map[string]int{"apples": 5, "oranges": 3}

	encoder := gnfmt.GNgob{}

	// Encode to gob
	gobData, err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}

	// Decode from gob
	var result map[string]int
	err = encoder.Decode(gobData, &result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result) // map[apples:5 oranges:3]
}
```

### CSV/TSV Utilities

#### Reading CSV Headers

```go
package main

import (
	"fmt"
	"github.com/gnames/gnfmt"
)

func main() {
	// Read CSV header and get field positions
	header, err := gnfmt.ReadHeaderCSV("data.csv", ',')
	if err != nil {
		panic(err)
	}

	// header is a map[string]int with field names and their indices
	fmt.Println(header["name"])  // 0
	fmt.Println(header["email"]) // 1

	// For TSV files, use '\t' as separator
	tsvHeader, _ := gnfmt.ReadHeaderCSV("data.tsv", '\t')
	fmt.Println(tsvHeader)
}
```

#### Converting Records to CSV/TSV

```go
package main

import (
	"fmt"
	"github.com/gnames/gnfmt"
)

func main() {
	record := []string{"John Doe", "john@example.com", "Has \"quotes\""}

	// Convert to CSV (comma-separated)
	csvLine := gnfmt.ToCSV(record, ',')
	fmt.Println(csvLine)
	// John Doe,john@example.com,"Has ""quotes"""

	// Convert to TSV (tab-separated)
	tsvLine := gnfmt.ToCSV(record, '\t')
	fmt.Println(tsvLine)
	// John Doe	john@example.com	"Has ""quotes"""
}
```

#### Normalizing Row Sizes

```go
package main

import (
	"fmt"
	"github.com/gnames/gnfmt"
)

func main() {
	// Truncate or expand rows to match expected size
	row := []string{"a", "b", "c", "d", "e"}

	// Truncate to 3 fields
	normalized := gnfmt.NormRowSize(row, 3)
	fmt.Println(normalized) // [a b c]

	// Expand to 7 fields (adds empty strings)
	expanded := gnfmt.NormRowSize(row, 7)
	fmt.Println(expanded) // [a b c d e  ]
}
```

### Time Formatting

Convert seconds into human-readable duration strings:

```go
package main

import (
	"fmt"
	"github.com/gnames/gnfmt"
)

func main() {
	fmt.Println(gnfmt.TimeString(45))        // 00:00:45
	fmt.Println(gnfmt.TimeString(3661))      // 01:01:01
	fmt.Println(gnfmt.TimeString(86400))     // 1 day 00:00:00
	fmt.Println(gnfmt.TimeString(172800))    // 2 days 00:00:00
	fmt.Println(gnfmt.TimeString(90061))     // 1 day 01:01:01
}
```

### Implementing Custom Encoders

The `Encoder` interface allows you to create custom serialization formats:

```go
package main

import (
	"github.com/gnames/gnfmt"
)

type MyCustomEncoder struct{}

func (e MyCustomEncoder) Encode(input interface{}) ([]byte, error) {
	// Your encoding logic here
	return []byte{}, nil
}

func (e MyCustomEncoder) Decode(input []byte, output interface{}) error {
	// Your decoding logic here
	return nil
}

func main() {
	var encoder gnfmt.Encoder = MyCustomEncoder{}
	// Use your custom encoder
	_, _ = encoder.Encode("data")
}
```

### Implementing the Outputter Interface

Create custom output formatters by implementing the `Outputter` interface:

```go
package main

import (
	"github.com/gnames/gnfmt"
)

type MyOutputter struct{}

func (o MyOutputter) Output(record interface{}, f gnfmt.Format) string {
	// Your custom output logic based on format
	return ""
}

func main() {
	var outputter gnfmt.Outputter = MyOutputter{}
	format, _ := gnfmt.NewFormat("csv")
	result := outputter.Output("data", format)
	_ = result
}
```

## API Documentation

For complete API documentation, visit [GoDoc](https://godoc.org/github.com/gnames/gnfmt).

## License

See the LICENSE file for details.
