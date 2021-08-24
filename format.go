package gnfmt

import "fmt"

// Format sets available output formats
type Format int

const (
	// FormatNone is for cases when format is not set yet.
	FormatNone Format = iota

	// CSV sets output to comma-separated values.
	CSV

	// CompactJSON sets output into one-liner JSON.
	CompactJSON

	// PrettyJSON sets output into easier to read JSON with new lines and
	// indentations.
	PrettyJSON

	// TSV sets output to tab-separated values.
	TSV
)

var formatStringMap = map[string]Format{
	"csv": CSV, "compact": CompactJSON, "pretty": PrettyJSON, "tsv": TSV,
}

var formatMap = map[Format]string{
	FormatNone:  "",
	CSV:         "CSV",
	CompactJSON: "compact JSON",
	PrettyJSON:  "pretty JSON",
	TSV:         "TSV",
}

// String representation of a format.
func (f Format) String() string {
	return formatMap[f]
}

// NewFormat is a constructor that converts a string into a corresponding format.
// If string cannot be converted, the constructor returns an error and
// and FormatNone format.
func NewFormat(s string) (Format, error) {
	if f, ok := formatStringMap[s]; ok {
		return f, nil
	}

	err := fmt.Errorf(
		"cannot convert '%s' to format, use 'csv', 'tsv', 'compact' or 'pretty' as input",
		s,
	)
	return FormatNone, err
}
