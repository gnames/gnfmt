package gnfmt

import (
	"bytes"
	"encoding/csv"
	"os"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ReadHeaderCSV only reads the first line of a CSV input. It takes
// a path to a CSV-encoded file and a separator character. It returns
// back a map with name of a field, and its index. It also returns an
// error, if the header could not be read.
func ReadHeaderCSV(path string, sep rune) (map[string]int, error) {
	res := make(map[string]int)
	f, err := os.Open(path)
	if err != nil {
		return res, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = sep

	// skip header
	header, err := r.Read()
	for i, v := range header {
		res[v] = i
	}
	return res, nil
}

// ToCSV takes a slice of strings and converts them to a comma-separated
// CSV row with '"' as a field quote. On Windows machine the row ends with
// '\r\n', in all other cases with '\n'
func ToCSV(record []string) string {
	var b bytes.Buffer
	useCRLF := runtime.GOOS == "windows"

	for i, field := range record {
		if i > 0 {
			b.WriteRune(',')
		}

		if !fieldNeedsQuotes(field) {
			b.WriteString(field)
			continue
		}

		b.WriteByte('"')
		for len(field) > 0 {
			// Search for special characters.
			ii := strings.IndexAny(field, "\"\r\n")
			if ii < 0 {
				ii = len(field)
			}

			// Copy verbatim everything before the special character.
			b.WriteString(field[:ii])
			field = field[ii:]

			// Encode the special character.
			if len(field) > 0 {
				switch field[0] {
				case '"':
					b.WriteString(`""`)
				case '\r':
					if !useCRLF {
						b.WriteByte('\r')
					}
				case '\n':
					if useCRLF {
						b.WriteString("\r\n")
					} else {
						b.WriteByte('\n')
					}
				}
				field = field[1:]
			}
		}
		b.WriteByte('"')
	}
	return b.String()
}

func fieldNeedsQuotes(field string) bool {
	if field == "" {
		return false
	}
	if field == `\.` || strings.ContainsRune(field, ',') || strings.ContainsAny(field, "\"\r\n") {
		return true
	}

	r1, _ := utf8.DecodeRuneInString(field)
	return unicode.IsSpace(r1)
}
