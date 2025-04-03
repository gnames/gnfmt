package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib"
	"github.com/gnames/gnsys"
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

	// SkipHeaders
	SkipHeaders bool

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

// OptSkipHeaders sets OptSkipHeaders field of the config.
func OptSkipHeaders(b bool) Option {
	return func(cfg *Config) {
		cfg.SkipHeaders = b
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
	res := rune(0)
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

var ErrNoInputOrOutput = errors.New("no input or output provided")
var ErrNoHeadersOrColSep = errors.New("provide headers and/or delimiter manually")
var ErrNoHeaders = errors.New("provide headers manually")
var ErrFileMissing = errors.New("provide valid input file path")
var ErrEmptyFirstLine = fmt.Errorf("empty first line")

// New creates a new Config by analyzing the first line of a CSV file
// to determine the delimiter and headers. Options can be provided to
// override the detected settings.
func New(opts ...Option) (Config, error) {
	res := Config{}

	for _, opt := range opts {
		opt(&res)
	}
	res.FieldsNum = len(res.Headers)

	// we need either file to read from of output for new CSV.
	if res.Path == "" && res.Writer == nil {
		return res, ErrNoInputOrOutput
	}

	if res.ColSep != 0 && len(res.Headers) > 0 {
		return res, nil
	}

	if res.Writer != nil {
		if res.ColSep == 0 {
			res.ColSep = ','
		}
		if len(res.Headers) == 0 {
			return res, ErrNoHeaders
		}

		return res, nil
	}

	if exists, _ := gnsys.FileExists(res.Path); !exists {
		return res, ErrFileMissing
	}

	// try to open file
	firstLine := readLine(res.Path)
	delimiter := res.ColSep
	headers := res.Headers
	skipHeaders := res.SkipHeaders

	if firstLine == "" {
		return res, ErrEmptyFirstLine
	}

	if delimiter == 0 {
		delimiter = detectDelimiter(firstLine)
	}

	if len(headers) == 0 {
		headers = strings.Split(firstLine, string(delimiter))
		headers = gnlib.Map(headers, func(s string) string {
			return strings.Trim(s, `"`)
		})
		skipHeaders = true
	}

	res = Config{
		Headers:     headers,
		SkipHeaders: skipHeaders,
		ColSep:      delimiter,
		FieldsNum:   len(headers),
	}

	// we have to run opts again, because  Config is updated
	for _, opt := range opts {
		opt(&res)
	}

	if res.ColSep == 0 || len(res.Headers) == 0 {
		return res, ErrNoHeadersOrColSep
	}
	res.FieldsNum = len(res.Headers)

	return res, nil
}
