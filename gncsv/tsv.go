package gncsv

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnfmt/gncsv/config"
)

// gntsv implements GnCSV interface.
type gntsv struct {
	cfg config.Config
}

// NewTSV creates a new GnCSV instance.
func NewTSV(cfg config.Config) GnCSV {
	res := gntsv{cfg: cfg}
	return &res
}

// Headers returns headers detected in the file.
func (g *gntsv) Headers() []string {
	return g.cfg.Headers
}

// ReadSlice reads a portion of the CSV data, starting at the given
// offset and reading up to the specified limit. It returns a slice
// of string slices, where each inner slice represents a row in the CSV.
func (g *gntsv) ReadSlice(offset, limit int) ([][]string, error) {
	f, err := os.Open(g.cfg.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewScanner(f)
	fieldsNum, lineNum := g.skipHeader(r)

	var res [][]string
	var count int
	for r.Scan() {
		lineNum++
		count++

		if limit > 0 && len(res) == limit {
			break
		}

		if offset > 0 && count <= offset {
			continue
		}

		line := r.Text()
		sep := string(g.cfg.ColSep)
		row := strings.Split(line, sep)
		rowFieldsNum := len(row)
		if fieldsNum == 0 {
			fieldsNum = rowFieldsNum
		}

		if fieldsNum != rowFieldsNum {
			skip, err := g.badRow(lineNum, fieldsNum, rowFieldsNum)
			if skip {
				continue
			}
			if err != nil {
				return nil, err
			}
			row = gnfmt.NormRowSize(row, fieldsNum)
		}

		res = append(res, row)
	}

	if err := r.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Read reads all CSV data and sends each row as a string slice to
// the provided channel. It returns the total number of rows read and
// any error encountered. It uses a context for cancellation.
func (g *gntsv) Read(ctx context.Context, ch chan<- []string) (int, error) {
	f, err := os.Open(g.cfg.Path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r := bufio.NewScanner(f)
	// some lines are huge, and generate an error
	// "bufio.Scanner: token too long". We are increasing the
	// maximum buffer size here.
	buf := make([]byte, 0, 64*1024)         // Set 64K buffer (same as default)
	r.Buffer(buf, bufio.MaxScanTokenSize*4) // Set 256k for maximum token size

	fieldsNum, lineNum := g.skipHeader(r)

	var count int64
	for r.Scan() {
		lineNum++

		if count%100_000 == 0 {
			fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 50))
			fmt.Fprintf(os.Stderr, "\rProcessed %s lines", humanize.Comma(count))
		}

		line := r.Text()
		sep := string(g.cfg.ColSep)
		row := strings.Split(line, sep)
		rowFieldsNum := len(row)
		if fieldsNum == 0 {
			fieldsNum = rowFieldsNum
		}

		if fieldsNum != rowFieldsNum {
			skip, err := g.badRow(lineNum, fieldsNum, rowFieldsNum)
			if skip {
				continue
			}
			if err != nil {
				return 0, err
			}
			row = gnfmt.NormRowSize(row, fieldsNum)
		}

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			count++
			ch <- row
		}
	}

	if err := r.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 50))
		slog.Error("Scanner error", "error", err)
		return int(count), err
	}

	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 50))
	return int(count), nil
}

// WriteStream writes CSV data received from the provided channel. Each
// string slice received from the channel represents a row in the CSV.
// It uses a context for cancellation.
func (g *gntsv) WriteStream(ctx context.Context, ch <-chan []string) error {
	f, err := os.Create(g.cfg.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	g.cfg.Headers = g.escapeFields(g.cfg.Headers)
	headers := strings.Join(g.cfg.Headers, string(g.cfg.ColSep)) + "\n"
	_, err = w.Write([]byte(headers))
	if err != nil {
		return err
	}
	for row := range ch {
		row = g.escapeFields(row)
		line := strings.Join(row, string(g.cfg.ColSep)) + "\n"
		_, err = w.Write([]byte(line))
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

func (g *gntsv) badRow(
	lineNum, fieldsNum, rowFieldsNum int,
) (bool, error) {
	switch g.cfg.BadRowMode {
	case gnfmt.ErrorBadRow:
		err := fmt.Errorf("wrong number of fieds: '%d'", lineNum)
		slog.Error("Bad row",
			"line", lineNum,
			"fieldsNum", fieldsNum,
			"rowFieldsNum", rowFieldsNum,
			"error", err,
		)
		return false, err
	case gnfmt.SkipBadRow:
		slog.Warn(
			"Wrong number of fields, SKIPPING row",
			"line", lineNum,
			"fieldsNum", fieldsNum,
			"rowFieldsNum", rowFieldsNum,
		)
		return true, nil
	case gnfmt.ProcessBadRow:
		slog.Warn(
			"Wrong number of fields, PROCESSING the row anyway",
			"line", lineNum,
			"fieldsNum", fieldsNum,
			"rowFieldsNum", rowFieldsNum,
		)
	}
	return false, nil
}

func (g *gntsv) skipHeader(r *bufio.Scanner) (int, int) {
	var fieldsNum, lineNum int
	// ignore headers gif they are given
	if len(g.cfg.Headers) > 0 {
		lineNum++
		r.Scan()
		line := r.Text()
		sep := string(g.cfg.ColSep)
		row := strings.Split(line, sep)
		fieldsNum = len(row)
	}
	return fieldsNum, lineNum
}

func (g *gntsv) escapeFields(ss []string) []string {
	res := make([]string, len(ss))
	for i := range ss {
		rs := []rune(ss[i])
		for ii := range rs {
			if rs[ii] == g.cfg.ColSep || rs[ii] == '\n' {
				// TODO find less destructive way to escape ColSep
				rs[ii] = '�'
			}
		}
		res[i] = string(rs)
	}
	return res
}
