package config_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gnames/gnfmt/gncsv/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path string
		delim     rune
		fieldsNum int
		nilErr    bool
	}{
		{"comma", "comma-norm.csv", ',', 9, true},
		{"tab", "tab-norm.csv", '\t', 9, true},
		{"pipe", "pipe-norm.csv", '|', 9, true},
		// not existing file gets defaults
		{"no file", "nofile", 0, 0, false},
	}

	for _, v := range tests {
		path := filepath.Join("..", "testdata", v.path)
		c, err := config.New(config.OptPath(path))
		assert.Equal(v.nilErr, err == nil, v.msg)
		assert.Equal(v.delim, c.ColSep, v.msg)
		assert.Equal(v.fieldsNum, c.FieldsNum, v.msg)
	}
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	// Test case: Valid CSV file with comma delimiter
	t.Run("Valid CSV with comma delimiter", func(t *testing.T) {
		path := filepath.Join("..", "testdata", "comma-norm.csv")
		c, err := config.New(config.OptPath(path))
		assert.Nil(err)
		assert.Equal(',', c.ColSep)
		assert.Equal(9, c.FieldsNum)
		assert.Equal(
			"taxonID",
			c.Headers[0],
		)
	})

	// Test case: Valid CSV file with tab delimiter
	t.Run("Valid CSV with tab delimiter", func(t *testing.T) {
		path := filepath.Join("..", "testdata", "tab-norm.csv")
		c, err := config.New(config.OptPath(path))
		assert.Nil(err)
		assert.Equal('\t', c.ColSep)
		assert.Equal(9, c.FieldsNum)
		assert.Equal(
			"taxonID",
			c.Headers[0],
		)
	})

	// Test case: Valid CSV file with pipe delimiter
	t.Run("Valid CSV with pipe delimiter", func(t *testing.T) {
		path := filepath.Join("..", "testdata", "pipe-norm.csv")
		c, err := config.New(config.OptPath(path))
		assert.Nil(err)
		assert.Equal('|', c.ColSep)
		assert.Equal(9, c.FieldsNum)
		assert.Equal(
			"taxonID",
			c.Headers[0],
		)
	})

	// Test case: No input or output provided
	t.Run("No input or output", func(t *testing.T) {
		_, err := config.New()
		assert.NotNil(err)
		assert.Equal(config.ErrNoInputOrOutput, err)
	})

	// Test case: Empty path
	t.Run("Empty path", func(t *testing.T) {
		_, err := config.New(config.OptPath(""))
		assert.NotNil(err)
		assert.Equal(config.ErrNoInputOrOutput, err)
	})

	// Test case: Writer provided but no headers or delimiter
	t.Run("Writer with no headers or delimiter", func(t *testing.T) {
		writer := &strings.Builder{}
		_, err := config.New(config.OptWriter(writer))
		assert.NotNil(err)
		assert.Equal(config.ErrNoHeaders, err)
	})

	// Test case: Valid Writer with headers and delimiter
	t.Run("Writer with headers and delimiter", func(t *testing.T) {
		writer := &strings.Builder{}
		headers := []string{"header1", "header2", "header3"}
		c, err := config.New(
			config.OptWriter(writer),
			config.OptHeaders(headers),
			config.OptColSep(','),
		)
		assert.Nil(err)
		assert.Equal(headers, c.Headers)
		assert.Equal(',', c.ColSep)
		assert.Equal(len(headers), c.FieldsNum)
	})
}
