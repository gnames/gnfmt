package gncsv

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnfmt/gncsv/config"
)

// gncsv implements GnCSV interface.
type gncsv struct {
	cfg config.Config
}

// NewCSV creates a new GnCSV instance.
func NewCSV(cfg config.Config) GnCSV {
	res := gncsv{cfg: cfg}
	return &res
}

// Headers returns headers detected in the file.
func (g *gncsv) Headers() []string {
	return g.cfg.Headers
}

// ReadSlice reads a portion of the CSV data, starting at the given
// offset and reading up to the specified limit. It returns a slice
// of string slices, where each inner slice represents a row in the CSV.
func (g *gncsv) ReadSlice(offset, limit int) ([][]string, error) {
	r, f, err := g.newReader()
	r.Comma = g.cfg.ColSep

	if err != nil {
		return nil, err
	}
	defer f.Close()

	fieldsNum, lineNum, err := g.skipHeader(r)
	if err != nil {
		return nil, err
	}

	var res [][]string
	var row []string

	var countLimit, countOffset int

	for {
		lineNum++
		countOffset++

		if limit > 0 && countLimit == limit {
			break
		}

		row, err = r.Read()
		if err == io.EOF {
			break
		}

		if fieldsNum == 0 {
			fieldsNum = len(row)
		}

		if err != nil {
			return nil, err
		}

		if offset > 0 && countOffset <= offset {
			continue
		}
		rowFieldsNum := len(row)
		if fieldsNum == 0 {
			fieldsNum = rowFieldsNum
		}

		if rowFieldsNum != fieldsNum {
			skip := g.badRow(lineNum, fieldsNum, rowFieldsNum)
			if skip {
				continue
			} else {
				// set row to the required size
				row = gnfmt.NormRowSize(row, fieldsNum)
			}
		}

		countLimit++
		res = append(res, row)
	}
	return res, nil
}

// Read reads all CSV data and sends each row as a string slice to
// the provided channel. It returns the total number of rows read and
// any error encountered. It uses a context for cancellation.
func (g *gncsv) Read(ctx context.Context, ch chan<- []string) (int, error) {
	r, f, err := g.newReader()
	r.Comma = g.cfg.ColSep
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// ignore headers if they are given
	fieldsNum, lineNum, err := g.skipHeader(r)
	if err != nil {
		return 0, err
	}

	var count int64
	for {
		lineNum++
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		rowFieldsNum := len(row)

		if fieldsNum == 0 {
			fieldsNum = rowFieldsNum
		}

		if fieldsNum != rowFieldsNum {
			skip := g.badRow(lineNum, fieldsNum, rowFieldsNum)
			if skip {
				continue
			} else {
				row = gnfmt.NormRowSize(row, fieldsNum)
			}
		}

		count++
		if count%100_000 == 0 {
			fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 50))
			fmt.Fprintf(os.Stderr, "\rProcessed %s lines", humanize.Comma(count))
		}

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			ch <- row
		}
	}

	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 50))
	return int(count), nil
}

// ToRow converts a slice of strings representing a CSV row
// into a single string with comma separation.
func ToRow(fields []string) string {
	// Create a new CSV writer with a comma separator
	var b bytes.Buffer
	w := csv.NewWriter(&b)

	// Write the row to the CSV writer
	err := w.Write(fields)
	if err != nil {
		// very unlikely
		return ""
	}
	w.Flush() // Ensure all data is written

	// Get the resulting string from the buffer
	return b.String()
}

// WriteStream writes CSV data received from the provided channel. Each
// string slice received from the channel represents a row in the CSV.
// It uses a context for cancellation.
func (g *gncsv) WriteStream(ctx context.Context, ch <-chan []string) error {
	var err error
	var w *csv.Writer
	if g.cfg.Path != "" {
		f, err := os.Create(g.cfg.Path)
		if err != nil {
			return err
		}
		defer f.Close()

		w = csv.NewWriter(f)
	} else {
		w = csv.NewWriter(g.cfg.Writer)
	}

	// Add headers, if they exist
	if len(g.cfg.Headers) > 0 {
		err = w.Write(g.cfg.Headers)
		if err != nil {
			return err
		}
	}

	for row := range ch {
		err = w.Write(row)
		if err != nil {
			for range ch {
			}
			return err
		}
		select {
		case <-ctx.Done():
			for range ch {
			}
			return ctx.Err()
		default:
		}
	}
	w.Flush()
	return nil
}

func (g *gncsv) newReader() (*csv.Reader, *os.File, error) {
	f, err := os.Open(g.cfg.Path)
	if err != nil {
		return nil, nil, err
	}
	r := csv.NewReader(f)
	if g.cfg.BadRowMode != gnfmt.ErrorBadRow {
		r.FieldsPerRecord = -1
	}
	return r, f, nil
}

func (g *gncsv) skipHeader(r *csv.Reader) (fieldsNum, lineNum int, err error) {
	if !g.cfg.SkipHeaders {
		return len(g.cfg.Headers), 0, nil
	}
	// ignore headers if they are given
	if len(g.cfg.Headers) > 0 {
		lineNum++
		row, err := r.Read()
		if err != nil {
			return 0, 0, err
		}
		fieldsNum = len(row)
	}
	return
}

func (g *gncsv) badRow(
	lineNum, fieldsNum, rowFieldsNum int,
) bool {
	msg := "SKIPPING row"
	skip := true
	if g.cfg.BadRowMode == gnfmt.ProcessBadRow {
		msg = "PROCESSING the row anyway"
		skip = false
	}

	slog.Warn(
		"Wrong number of fields, "+msg,
		"line", lineNum,
		"fieldsNum", fieldsNum,
		"rowFieldsNum", rowFieldsNum,
	)
	return skip
}
