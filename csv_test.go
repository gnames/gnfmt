package gnfmt_test

import (
	"path/filepath"
	"testing"

	"github.com/gnames/gnfmt"
)

func TestReadHeaderCSV(t *testing.T) {
	path := filepath.Join("testdata", "test.tsv")
	header, err := gnfmt.ReadHeaderCSV(path, '\t')
	if err != nil {
		t.Errorf("Cannot read csv file '%s': %s", path, err)
	}
	headerTest := map[string]int{
		"Id":             0,
		"NameCanonical":  1,
		"NameAuthorship": 2,
		"NameYear":       3,
		"RefString":      4,
		"RefYear":        5,
	}
	for k, v := range headerTest {
		if header[k] != v {
			t.Errorf("Wrong header values: '%s', %d", k, v)
		}
	}
}

func TestToCSV(t *testing.T) {
	ss := []string{"one\"two", "three,four", "five"}
	res := gnfmt.ToCSV(ss, ',')
	testRes := `"one""two","three,four",five`
	if res != testRes {
		t.Errorf("ToCSV failed, got '%s' instad of '%s'", res, testRes)
	}
}

func TestToTSV(t *testing.T) {
	ss := []string{"one\"two", "three\tfour", "five"}
	res := gnfmt.ToCSV(ss, '\t')
	testRes := "\"one\"\"two\"\t\"three\tfour\"\tfive"
	if res != testRes {
		t.Errorf("ToTSV failed, got '%s' instad of '%s'", res, testRes)
	}
}
