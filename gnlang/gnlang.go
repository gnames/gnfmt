package gnlang

import (
	"errors"
	"strings"

	"golang.org/x/text/language"
)

var ErrInvalidCodeLength = errors.New("invalid language code length")
var ErrUnknownCode = errors.New("unknown language code")

// LangCode2To3Letters converts a two-letter ISO 639-1 language code to
// a three-letter ISO 639-3 language code.
func LangCode2To3Letters(code string) (string, error) {
	if len(code) != 2 {
		return "", ErrInvalidCodeLength
	}
	tag, err := language.Parse(code)
	if err != nil {
		return "", ErrUnknownCode
	}
	base, _ := tag.Base()
	threeLetterCode := base.ISO3()
	return threeLetterCode, nil
}

// LangCode3To2Letters converts a three-letter ISO 639-3 language code to
// a two-letter ISO 639-1 language code.
func LangCode3To2Letters(code string) (string, error) {
	if len(code) != 3 {
		return "", ErrInvalidCodeLength
	}
	tag, err := language.Parse(code)
	if err != nil {
		return "", ErrUnknownCode
	}
	base, _ := tag.Base()
	twoLetterCode := base.String()
	return twoLetterCode, nil
}

// Lang returns the language name for the given ISO 639-3 code.
func Lang(code string) string {
	code = strings.ToLower(code)
	iso := iso639_3(code)
	return langCodeMap[iso]
}

// LangCode returns the ISO 639-3 code for the given language name.
func LangCode(lang string) string {
	lang = strings.ToLower(lang)
	if code, ok := langMap[lang]; ok {
		return string(code)
	}
	return ""
}

// Country returns the country name for the given ISO 3166-1 alpha-2 code.
func Country(code string) string {
	code = strings.ToUpper(code)
	iso := iso3166(code)
	return countryCodeMap[iso]
}

// CountryCode returns the ISO 3166-1 alpha-2 code for the given country name.
func CountryCode(country string) string {
	country = strings.ToLower(country)
	if full, ok := countryAbbr[country]; ok {
		country = full
	}
	if code, ok := countryMap[country]; ok {
		return string(code)
	}
	return ""
}
