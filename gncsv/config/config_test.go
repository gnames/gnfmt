package config_test

import (
	"path/filepath"
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
	}{
		{"comma", "comma-norm.csv", ',', 9},
		{"tab", "tab-norm.csv", '\t', 9},
		{"pipe", "pipe-norm.csv", '|', 9},
		// not existing file gets defaults
		{"no file", "nofile", ',', 0},
	}

	for _, v := range tests {
		path := filepath.Join("..", "testdata", v.path)
		c, err := config.New(config.OptPath(path))
		assert.Nil(err)
		assert.Equal(v.delim, c.ColSep, v.msg)
		assert.Equal(v.fieldsNum, c.FieldsNum, v.msg)
	}
}

func TestBadConfig(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path string
	}{
		{"? delim", "unknown.csv"},
	}

	for _, v := range tests {
		path := filepath.Join("..", "testdata", v.path)
		_, err := config.New(config.OptPath(path))
		assert.NotNil(err)
	}
}
