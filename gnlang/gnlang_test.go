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

func TestLangCode2To3Letters(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		iso1, iso3 string
	}{
		{"en", "eng"},
		{"En", "eng"},
		{"eN", "eng"},
		{"fr", "fra"},
		{"de", "deu"},
		{"es", "spa"},
		{"it", "ita"},
		{"pt", "por"},
		{"ru", "rus"},
		{"ja", "jpn"},
		{"ko", "kor"},
		{"zh", "zho"},
	}

	for _, v := range tests {
		code, err := gnlang.LangCode2To3Letters(v.iso1)
		assert.NoError(err)
		assert.Equal(v.iso3, code)
	}

	// Test invalid codes
	_, err := gnlang.LangCode2To3Letters("abc")
	assert.Error(err)
	assert.Equal(gnlang.ErrInvalidCodeLength, err)

	_, err = gnlang.LangCode2To3Letters("zz")
	assert.Error(err)
	assert.Equal(gnlang.ErrUnknownCode, err)
}

func TestLangCode3To2Letters(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		iso3, iso1 string
	}{
		{"eng", "en"},
		{"Eng", "en"},
		{"ENG", "en"},
		{"eNG", "en"},
		{"fra", "fr"},
		{"deu", "de"},
		{"spa", "es"},
		{"ita", "it"},
		{"por", "pt"},
		{"rus", "ru"},
		{"jpn", "ja"},
		{"kor", "ko"},
		{"zho", "zh"},
	}

	for _, v := range tests {
		code, err := gnlang.LangCode3To2Letters(v.iso3)
		assert.NoError(err)
		assert.Equal(v.iso1, code)
	}

	// Test invalid codes
	_, err := gnlang.LangCode3To2Letters("ab")
	assert.Error(err)
	assert.Equal(gnlang.ErrInvalidCodeLength, err)

	_, err = gnlang.LangCode3To2Letters("zzz")
	assert.Error(err)
	assert.Equal(gnlang.ErrUnknownCode, err)
}
