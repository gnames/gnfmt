package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gnames/gnfmt"
)

// Config provides settings for processing CSV data, including
// file path, headers, delimiter, field count, and bad row handling mode.
type Config struct {
	// Path to the CSV file.
	Path string

	// Writer can be used for writing output instead of file.
	Writer io.Writer

	// Headers are the names of fields in the CSV file.
	Headers []string

	// ColSep is the delimiter character used in the CSV file.
	ColSep rune

	// FieldsNum is the expected number of fields in each row.
	FieldsNum int

	// BadRowMode specifies how to handle rows with an incorrect
	// number of fields. Options include processing, ignoring, or
	// raising an error (default).
	BadRowMode gnfmt.BadRow

	// WithQuotes is true if `"` is used when need arises to
	// protect field separators, new lines inside a field.
	WithQuotes bool
}

// Update creates a copy of the Config and applies the provided
// options to it, returning the updated copy.
func (c Config) Update(opts ...Option) Config {
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// Option is a function that modifies a Config.
type Option func(*Config)

// OptPath sets the Path field of the Config.
func OptPath(s string) Option {
	return func(cfg *Config) {
		cfg.Path = s
	}
}

// OptWriter sets Writer field of the Config.
func OptWriter(w io.Writer) Option {
	return func(cfg *Config) {
		cfg.Writer = w
	}
}

// OptHeaders sets the Headers field of the Config.
func OptHeaders(ss []string) Option {
	return func(cfg *Config) {
		cfg.Headers = ss
	}
}

// OptColSep sets the ColSep field of the Config.
func OptColSep(r rune) Option {
	return func(cfg *Config) {
		cfg.ColSep = r
	}
}

// OptBadRowMode sets the BadRowMode field of the Config.
func OptBadRowMode(br gnfmt.BadRow) Option {
	return func(cfg *Config) {
		cfg.BadRowMode = br
	}
}

// OptWithQuotes sets WithQuotes field of the Config.
func OptWithQuotes(b bool) Option {
	return func(cfg *Config) {
		cfg.WithQuotes = b
	}
}

// OptFieldsNum sets the FieldsNum field of the Config.
func OptFieldsNum(i int) Option {
	return func(cfg *Config) {
		cfg.FieldsNum = i
	}
}

// detectDelimiter analyzes a line of text and determines the most
// likely delimiter character (comma, tab, or pipe).
func detectDelimiter(line string) rune {
	count := make(map[rune]int)
	rs := []rune(line)
	for _, r := range rs {
		switch r {
		case ',':
			count[',']++
		case '\t':
			count['\t']++
		case '|':
			count['|']++
		}
	}
	res := '?'
	var maxCount int
	for k, v := range count {
		if v > maxCount {
			maxCount = v
			res = k
		}
	}
	return res
}

// readLine reads the first line from a file.
func readLine(path string) string {
	if path == "" {
		return ""
	}

	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		return line
	}
	return ""
}

// New creates a new Config by analyzing the first line of a CSV file
// to determine the delimiter and headers. Options can be provided to
// override the detected settings.
func New(opts ...Option) (Config, error) {
	res := Config{
		ColSep: ',',
	}

	for _, opt := range opts {
		opt(&res)
	}
	if res.Path == "" && res.Writer == nil {
		return res, errors.New("no input or output provided")
	}

	// try to open file
	firstLine := readLine(res.Path)
	if firstLine != "" {
		delimiter := detectDelimiter(firstLine)
		if delimiter == '?' {
			return res, fmt.Errorf("cannot determine delimiter: '%s'", firstLine)
		}

		headers := strings.Split(firstLine, string(delimiter))
		res = Config{
			Headers:   headers,
			ColSep:    delimiter,
			FieldsNum: len(headers),
		}
	}

	// we have to run opts again, because  Config is updated
	for _, opt := range opts {
		opt(&res)
	}
	return res, nil
}
