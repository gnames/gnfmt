package gncsv_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnfmt/gncsv"
	"github.com/gnames/gnfmt/gncsv/config"
	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {
	assert := assert.New(t)
	headers := []string{
		"taxonID", "scientificName", "kingdom", "phylum", "class", "order",
		"family", "genus", "nomenclaturalCode",
	}
	tests := []struct {
		msg, path string
		headers   []string
	}{
		{"csv", "comma-norm.csv", headers},
		{"pipe", "pipe-norm.csv", headers},
		{"tab", "tab-norm.csv", headers},
	}
	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		opt := config.OptPath(path)
		cfg, err := config.New(opt)
		assert.Nil(err)
		c := gncsv.New(cfg)
		assert.Equal(v.headers, c.Headers())
	}
}

func TestNoHeadersCSV(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path, dataFirst string
		dataLen              int
	}{
		{"csv", "comma-no-headers.csv", "2|Nothocercus bonapartei", 10},
		{"tab", "tab-no-headers.csv", "2|Nothocercus bonapartei", 10},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		headers := []string{
			"taxonID", "scientificName", "kingdom", "phylum",
			"class", "order", "family", "genus", "nomenclaturalCode",
		}

		opts := []config.Option{
			config.OptPath(path),
			config.OptHeaders(headers),
		}
		cfg, err := config.New(opts...)
		assert.Nil(err)

		c := gncsv.New(cfg)

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
		assert.Nil(err)
		close(ch)
		wg.Wait()
		first := strings.Join(res[0][0:2], "|")
		assert.Equal(v.dataFirst, first)
		assert.Equal(count, len(res))
		assert.Equal(v.dataLen, len(res))
	}
}

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
		assert.Nil(err)
		close(ch)
		wg.Wait()
		assert.Equal(count, len(res))
		assert.Equal(v.dataLen, len(res))
		assert.Equal(v.dataFirst, res[0][0])
	}
}

func TestReadF(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path, id, name string
	}{
		{"csv0", "comma-norm.csv", "2", "Nothocercus bonapartei"},
		{"quoted", "comma-quoted.csv", "2", "Nothocercus bonapartei"},
		{"tsv0", "tab-norm.csv", "2", "Nothocercus bonapartei"},
		{"psv0", "pipe-norm.csv", "2", "Nothocercus bonapartei"},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		opt := config.OptPath(path)
		cfg, err := config.New(opt)
		assert.Nil(err)
		c := gncsv.New(cfg)
		sl, err := c.ReadSlice(0, 1)
		assert.Nil(err)
		id := c.F(sl[0], "taxoNiD")
		assert.Equal(v.id, id)
		name := c.F(sl[0], "ScientificName")
		assert.Equal(v.name, name)
		unknown := c.F(sl[0], "smth")
		assert.Equal("", unknown)
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
		if _, err := os.Stat(pathWrite); err == nil {
			os.Remove(pathWrite)
		}
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
}

func TestReadChunks(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path                     string
		chunkSize, chunksNum, dataLen int
	}{
		{"csv", "comma-norm.csv", 3, 4, 10},
		{"pipe", "pipe-norm.csv", 4, 3, 10},
		{"tab", "tab-norm.csv", 5, 2, 10},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.path)
		opt := config.OptPath(path)
		cfg, err := config.New(opt)
		assert.Nil(err)
		c := gncsv.New(cfg)

		chOut := make(chan [][]string)
		var wg sync.WaitGroup
		var res [][]string

		var chunksCount int
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chOut {
				chunksCount++
				res = append(res, chunk...)
			}
		}()

		ctx := context.Background()
		count, err := c.ReadChunks(ctx, chOut, v.chunkSize)
		close(chOut)
		wg.Wait()

		assert.Nil(err, v.msg)
		assert.Equal(v.dataLen, count, v.msg)
		assert.Equal(v.dataLen, len(res), v.msg)
		assert.Equal(v.chunksNum, chunksCount, v.msg)
	}
}

func ExampleReader_ReadSlice() {
	path := filepath.Join("testdata", "comma-norm.csv")
	opts := []config.Option{
		config.OptPath(path),
	}

	cfg, err := config.New(opts...)
	if err != nil {
		panic(err)
	}
	c := gncsv.New(cfg)

	rows, err := c.ReadSlice(0, 3)
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		fmt.Println(strings.Join(row, ","))
	}

	// Output:
	// 2,Nothocercus bonapartei,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Nothocercus,ICZN
	// 1,Tinamus major,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Tinamus,ICZN
	// 3,Crypturellus soui,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Crypturellus,ICZN
}

func ExampleReader_Read() {
	path := filepath.Join("testdata", "comma-norm.csv")
	opts := []config.Option{
		config.OptPath(path),
	}

	cfg, err := config.New(opts...)
	if err != nil {
		panic(err)
	}
	c := gncsv.New(cfg)

	ch := make(chan []string)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for row := range ch {
			fmt.Println(strings.Join(row, ","))
		}
	}()
	_, err = c.Read(context.Background(), ch)
	if err != nil {
		panic(err)
	}
	close(ch)
	wg.Wait()

	// Output:
	// 2,Nothocercus bonapartei,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Nothocercus,ICZN
	// 1,Tinamus major,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Tinamus,ICZN
	// 3,Crypturellus soui,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Crypturellus,ICZN
	// 4,Crypturellus cinnamomeus,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Crypturellus,ICZN
	// 5,Crypturellus boucardi,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Crypturellus,ICZN
	// 6,Crypturellus kerriae,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Crypturellus,ICZN
	// 7,Dendrocygna viduata,Animalia,Chordata,Aves,Anseriformes,Anatidae,Dendrocygna,ICZN
	// 8,Dendrocygna autumnalis,Animalia,Chordata,Aves,Anseriformes,Anatidae,Dendrocygna,ICZN
	// 9,Dendrocygna arborea,Animalia,Chordata,Aves,Anseriformes,Anatidae,Dendrocygna,ICZN
	// 10,Dendrocygna bicolor,Animalia,Chordata,Aves,Anseriformes,Anatidae,Dendrocygna,ICZN
}

func ExampleReader_ReadChunks() {
	path := filepath.Join("testdata", "comma-norm.csv")
	opts := []config.Option{
		config.OptPath(path),
	}

	cfg, err := config.New(opts...)
	if err != nil {
		panic(err)
	}
	c := gncsv.New(cfg)

	chOut := make(chan [][]string)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for chunk := range chOut {
			for _, row := range chunk {
				fmt.Println(strings.Join(row, ","))
			}
		}
	}()

	ctx := context.Background()
	if _, err := c.ReadChunks(ctx, chOut, 3); err != nil {
		panic(err)
	}
	close(chOut)
	wg.Wait()

	// Output:
	// 2,Nothocercus bonapartei,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Nothocercus,ICZN
	// 1,Tinamus major,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Tinamus,ICZN
	// 3,Crypturellus soui,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Crypturellus,ICZN
	// 4,Crypturellus cinnamomeus,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Crypturellus,ICZN
	// 5,Crypturellus boucardi,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Crypturellus,ICZN
	// 6,Crypturellus kerriae,Animalia,Chordata,Aves,Tinamiformes,Tinamidae,Crypturellus,ICZN
	// 7,Dendrocygna viduata,Animalia,Chordata,Aves,Anseriformes,Anatidae,Dendrocygna,ICZN
	// 8,Dendrocygna autumnalis,Animalia,Chordata,Aves,Anseriformes,Anatidae,Dendrocygna,ICZN
	// 9,Dendrocygna arborea,Animalia,Chordata,Aves,Anseriformes,Anatidae,Dendrocygna,ICZN
	// 10,Dendrocygna bicolor,Animalia,Chordata,Aves,Anseriformes,Anatidae,Dendrocygna,ICZN
}

func ExampleWriter() {
	var b bytes.Buffer
	writeOpts := []config.Option{
		config.OptWriter(&b),
		config.OptHeaders([]string{"header1", "header2", "header3"}),
	}
	cfg, err := config.New(writeOpts...)
	if err != nil {
		panic(err)
	}
	w := gncsv.New(cfg)

	ch := make(chan []string)
	go func() {
		defer close(ch)
		for range 3 {
			ch <- []string{"val1", "val2", "val3"}
		}
	}()

	if err := w.WriteStream(context.Background(), ch); err != nil {
		panic(err)
	}

	fmt.Println(b.String())
	// Output:
	// header1,header2,header3
	// val1,val2,val3
	// val1,val2,val3
	// val1,val2,val3
}
