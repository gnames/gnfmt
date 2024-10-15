package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gnames/gnfmt"
)

// Config provides settings for processing CSV data, including
// file path, headers, delimiter, field count, and bad row handling mode.
type Config struct {
	// Path to the CSV file.
	Path string

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
func readLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		return line, nil
	} else if scanner.Err() != nil {
		return "", scanner.Err()
	} else {
		return "", fmt.Errorf("empty file: %s", path)
	}
}

// New creates a new Config by analyzing the first line of a CSV file
// to determine the delimiter and headers. Options can be provided to
// override the detected settings.
func New(csvPath string, opts ...Option) (Config, error) {
	res := Config{
		Path:   csvPath,
		ColSep: ',',
	}

	// try to open file
	firstLine, err := readLine(csvPath)
	if err == nil {
		delimiter := detectDelimiter(firstLine)
		if delimiter == '?' {
			return res, fmt.Errorf("cannot determine delimiter: '%s'", firstLine)
		}

		headers := strings.Split(firstLine, string(delimiter))
		res = Config{
			Path:      csvPath,
			Headers:   headers,
			ColSep:    delimiter,
			FieldsNum: len(headers),
		}
	}

	for _, opt := range opts {
		opt(&res)
	}
	return res, nil
}
