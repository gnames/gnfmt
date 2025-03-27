package gnlang

import "strings"

func Lang(code string) string {
	code = strings.ToLower(code)
	iso := iso639_3(code)
	return langCodeMap[iso]
}

func LangCode(lang string) string {
	lang = strings.ToLower(lang)
	if code, ok := langMap[lang]; ok {
		return string(code)
	}
	return ""
}

func Country(code string) string {
	code = strings.ToUpper(code)
	iso := iso3166(code)
	return countryCodeMap[iso]
}

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
