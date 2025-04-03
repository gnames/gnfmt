package gncsv

import (
	"fmt"
	"strings"

	"github.com/gnames/gnfmt/gncsv/config"
)

// New creates a new CSV or TSV/PSV reader/writer based on the provided
// configuration. If the ColSep in the config is a comma, it creates
// a CSV reader/writer. Otherwise, it creates a TSV reader/writer.
func New(cfg config.Config) GnCSV {
	if cfg.ColSep == ',' || cfg.WithQuotes {
		return NewCSV(cfg)
	}
	return NewTSV(cfg)
}

// getField is a field accessor. If the field with the given name exists, it returns
// the value of the field in the row. If not, returns empty string and an
// error.
func getField(headerMap map[string]int, row []string, field string) (string, error) {
	fieldLow := strings.ToLower(field)
	if fieldNum, ok := headerMap[fieldLow]; ok {
		return row[fieldNum], nil
	}
	return "", fmt.Errorf("unknown field: '%s'", field)
}
