package gncsv_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

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
		{"csv0", "comma.csv", "1|Tinamus major", "1", 0, 3, 13},
		{"csv1", "comma.csv", "2|Nothocercus bonapartei", "1", 1, 2, 13},
		{"csv2", "comma.csv", "3|Crypturellus soui", "1", 2, 1, 13},
		{"tsv0", "tab.csv",
			"leptogastrinae:tid:42|http://leptogastrinae.lifedesks.org/pages/42",
			"leptogastrinae:tid:42", 0, 3, 12},
		{"tsv1", "tab.csv",
			"leptogastrinae:tid:2044|http://leptogastrinae.lifedesks.org/pages/2044",
			"leptogastrinae:tid:42", 2, 1, 12},
		{"pipe0", "pipe.csv",
			"leptogastrinae:tid:42|http://leptogastrinae.lifedesks.org/pages/42",
			"leptogastrinae:tid:42", 0, 3, 12},
		{"pipe1", "pipe.csv",
			"leptogastrinae:tid:2044|http://leptogastrinae.lifedesks.org/pages/2044",
			"leptogastrinae:tid:42", 2, 1, 12},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", "colsep", v.path)
		cfg, err := config.New(path)
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
		{"csv", "comma.csv", 13},
		{"pipe", "pipe.csv", 12},
		{"tab", "tab.csv", 12},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", "colsep", v.path)
		cfg, err := config.New(path)
		assert.Nil(err)
		headers := cfg.Headers
		r := gncsv.New(cfg)

		tmpDir := os.TempDir()
		pathWrite := filepath.Join(tmpDir, v.path)
		hdrs := config.OptHeaders(headers)
		cfgWrite, err := config.New(pathWrite, hdrs)
		assert.Nil(err)
		w := gncsv.New(cfgWrite)
		var wg sync.WaitGroup
		wg.Add(1)
		ch := make(chan []string)
		go func() {
			defer wg.Done()
			err = w.Write(context.Background(), ch)
			assert.Nil(err)
		}()
		count, err := r.Read(context.Background(), ch)
		close(ch)
		wg.Wait()
		assert.Equal(v.dataLen, count)

		cfg, err = config.New(pathWrite)
		assert.Nil(err)
		assert.Equal(cfg.FieldsNum, len(headers))
		assert.Equal(',', cfg.ColSep)
	}
}
