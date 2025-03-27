package gnlang_test

import (
	"testing"

	"github.com/gnames/gnfmt/gnlang"
	"github.com/stretchr/testify/assert"
)

func TestCountryCode(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		cnt, code string
	}{
		{"Russia", "RUS"},
		{"GB", "GBR"},
		{"usa", "USA"},
		{"Vietnam", "VNM"},
		{"United States", "USA"},
		{"something", ""},
	}

	for _, v := range tests {
		code := gnlang.CountryCode(v.cnt)
		assert.Equal(v.code, code)
	}

}

func TestCountry(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		code, cnt string
	}{
		{"RUS", "Russian Federation"},
		{"GbR", "United Kingdom of Great Britain and Northern Ireland"},
		{"USA", "United States of America"},
		{"vnm", "Viet Nam"},
		{"TTT", ""},
	}

	for _, v := range tests {
		cnt := gnlang.Country(v.code)
		assert.Equal(v.cnt, cnt)
	}
}
func TestLang(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		code, lang string
	}{
		{"Eng", "English"},
		{"FRA", "French"},
		{"spa", "Spanish"},
		{"deu", "German"},
		{"zzz", ""},
	}

	for _, v := range tests {
		lang := gnlang.Lang(v.code)
		assert.Equal(v.lang, lang)
	}
}

func TestLangCode(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		lang, code string
	}{
		{"English", "eng"},
		{"french", "fra"},
		{"SPANISH", "spa"},
		{"German", "deu"},
		{"unknown", ""},
	}

	for _, v := range tests {
		code := gnlang.LangCode(v.lang)
		assert.Equal(v.code, code)
	}
}
