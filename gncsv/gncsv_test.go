package gncsv_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnfmt/gncsv"
	"github.com/gnames/gnfmt/gncsv/config"
	"github.com/stretchr/testify/assert"
)

func TestReadCSV(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path, dataSlice, dataFirst string
		offset, limit, dataLen          int
	}{
		{"csv0", "comma-norm.csv", "2|Nothocercus bonapartei", "2", 0, 3, 10},
		{"csv1", "comma-norm.csv", "1|Tinamus major", "2", 1, 2, 10},
		{"csv2", "comma-norm.csv", "3|Crypturellus soui", "2", 2, 1, 10},
		{"tsv0", "tab-norm.csv", "2|Nothocercus bonapartei", "2", 0, 3, 10},
		{"tsv1", "tab-norm.csv", "1|Tinamus major", "2", 1, 2, 10},
		{"tsv2", "tab-norm.csv", "3|Crypturellus soui", "2", 2, 1, 10},
		{"psv0", "pipe-norm.csv", "2|Nothocercus bonapartei", "2", 0, 3, 10},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		opt := config.OptPath(path)
		cfg, err := config.New(opt)
		assert.Nil(err)
		c := gncsv.New(cfg)
		sl, err := c.ReadSlice(v.offset, v.limit)
		assert.Nil(err)
		assert.Equal(v.limit, len(sl))
		assert.Equal(v.dataSlice, sl[0][0]+"|"+sl[0][1])

		var wg sync.WaitGroup
		wg.Add(1)
		ch := make(chan []string)
		var res [][]string
		go func() {
			defer wg.Done()
			for row := range ch {
				res = append(res, row)
			}
		}()
		count, err := c.Read(context.Background(), ch)
		close(ch)
		wg.Wait()
		assert.Equal(count, len(res))
		assert.Equal(v.dataLen, len(res))
		assert.Equal(v.dataFirst, res[0][0])
	}
}

func TestWriteCSV(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path string
		dataLen   int
	}{
		{"csv", "comma-norm.csv", 10},
		{"pipe", "pipe-norm.csv", 10},
		{"tab", "tab-norm.csv", 10},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		opt := config.OptPath(path)
		cfg, err := config.New(opt)
		assert.Nil(err)
		headers := cfg.Headers
		r := gncsv.New(cfg)

		tmpDir := os.TempDir()
		pathWrite := filepath.Join(tmpDir, v.path)
		opts := []config.Option{
			config.OptPath(pathWrite),
			config.OptHeaders(headers),
		}
		cfgWrite, err := config.New(opts...)
		assert.Nil(err)
		w := gncsv.New(cfgWrite)
		var wg sync.WaitGroup
		wg.Add(1)
		ch := make(chan []string)
		go func() {
			defer wg.Done()
			err = w.WriteStream(context.Background(), ch)
			assert.Nil(err)
		}()
		count, err := r.Read(context.Background(), ch)
		assert.Nil(err)
		close(ch)
		wg.Wait()
		assert.Equal(v.dataLen, count)

		opt = config.OptPath(pathWrite)
		cfg, err = config.New(opt)
		assert.Nil(err)
		assert.Equal(cfg.FieldsNum, len(headers))
		assert.Equal(',', cfg.ColSep)
	}
}

func TestIOWriter(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path string
		dataLen   int
	}{
		{"csv", "comma-norm.csv", 10},
		{"pipe", "pipe-norm.csv", 10},
		{"tab", "tab-norm.csv", 10},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		opt := config.OptPath(path)
		cfg, err := config.New(opt)
		assert.Nil(err)
		headers := cfg.Headers
		r := gncsv.New(cfg)

		var b bytes.Buffer
		opts := []config.Option{
			config.OptWriter(&b),
			config.OptHeaders(headers),
		}
		cfgWrite, err := config.New(opts...)
		assert.Nil(err, v.msg)
		w := gncsv.New(cfgWrite)
		var wg sync.WaitGroup
		wg.Add(1)
		ch := make(chan []string)
		go func() {
			defer wg.Done()
			err = w.WriteStream(context.Background(), ch)
			assert.Nil(err)
		}()
		count, err := r.Read(context.Background(), ch)
		close(ch)
		wg.Wait()
		assert.Equal(v.dataLen, count)

		assert.Greater(len(b.String()), 500)
	}
}

func TestBadRowsError(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert := assert.New(t)
	tests := []struct {
		msg, path string
		err       bool
	}{
		{"csv", "comma-norm.csv", false},
		{"csv less", "comma-less.csv", true},
		{"csv more", "comma-more.csv", true},
		{"tsv", "tab-norm.csv", false},
		{"tsv less", "tab-less.csv", true},
		{"tsv more", "tab-more.csv", true},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		opt := config.OptPath(path)
		cfg, err := config.New(opt)
		assert.Nil(err)
		c := gncsv.New(cfg)
		_, err = c.ReadSlice(0, 0)
		assert.Equal(v.err, err != nil, v.msg)
	}
}

func TestBadRowsSkip(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert := assert.New(t)
	tests := []struct {
		msg, path string
		err       bool
		rowsNum   int
	}{
		{"csv", "comma-norm.csv", false, 10},
		{"csv less", "comma-less.csv", false, 9},
		{"csv more", "comma-more.csv", false, 9},
		{"tsv", "tab-norm.csv", false, 10},
		{"tsv less", "tab-less.csv", false, 9},
		{"tsv more", "tab-more.csv", false, 9},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		opts := []config.Option{
			config.OptPath(path),
			config.OptBadRowMode(gnfmt.SkipBadRow),
		}

		cfg, err := config.New(opts...)
		assert.Nil(err)
		c := gncsv.New(cfg)
		rows, err := c.ReadSlice(0, 0)
		assert.Equal(v.err, err != nil, v.msg)
		assert.Equal(v.rowsNum, len(rows))
	}
}

func TestBadRowsProcess(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert := assert.New(t)
	tests := []struct {
		msg, path          string
		err                bool
		fieldsNum, rowsNum int
	}{
		{"csv", "comma-norm.csv", false, 9, 10},
		{"csv less", "comma-less.csv", false, 9, 10},
		{"csv more", "comma-more.csv", false, 9, 10},
		{"tsv", "tab-norm.csv", false, 9, 10},
		{"tsv less", "tab-less.csv", false, 9, 10},
		{"tsv more", "tab-more.csv", false, 9, 10},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		opts := []config.Option{
			config.OptPath(path),
			config.OptBadRowMode(gnfmt.ProcessBadRow),
		}

		cfg, err := config.New(opts...)
		assert.Nil(err)
		c := gncsv.New(cfg)
		rows, err := c.ReadSlice(0, 0)
		assert.Equal(v.err, err != nil, v.msg)
		for _, row := range rows {
			assert.Equal(v.fieldsNum, len(row))
		}
		assert.Equal(v.rowsNum, len(rows))
	}
}

func TestTabWithQuotes(t *testing.T) {
	assert := assert.New(t)
	path := filepath.Join("testdata", "tab-w-quotes.csv")
	opts := []config.Option{
		config.OptPath(path),
		config.OptWithQuotes(true),
	}

	cfg, err := config.New(opts...)
	assert.Nil(err)
	c := gncsv.New(cfg)
	rows, err := c.ReadSlice(0, 0)
	assert.Nil(err)
	l := len(cfg.Headers)
	for _, row := range rows {
		assert.Equal(l, len(row))
	}
	assert.Equal(10, len(rows))

	opts = []config.Option{
		config.OptPath(path),
	}

	cfg, err = config.New(opts...)
	assert.Nil(err)
	c = gncsv.New(cfg)
	rows, err = c.ReadSlice(0, 0)
	assert.Nil(rows)
	assert.NotNil(err)
}
