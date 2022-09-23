package gnfmt

import (
	"bytes"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// GNjson allows to decode and encode JSON format.
type GNjson struct {
	Pretty bool
}

// Encode takes an object and coverts it into JSON. It returns an error
// if the encoding fails.
func (e GNjson) Encode(input interface{}) ([]byte, error) {
	if e.Pretty {
		return jsoniter.MarshalIndent(input, "", "  ")
	}
	return jsoniter.Marshal(input)
}

// Decode converts JSON into a go object. If decoding breaks, it
// returns an error.
func (e GNjson) Decode(input []byte, output interface{}) error {
	r := bytes.NewReader(input)
	err := jsoniter.NewDecoder(r).Decode(output)
	return err
}

// Output converts an object into a JSON string. It takes an object and
// a format and returns the corresponding JSON string. In case of a problem
// it returns an empty string.
func (e GNjson) Output(input interface{}, f Format) string {
	switch f {
	case CompactJSON:
		e.Pretty = false
	case PrettyJSON:
		e.Pretty = true
	default:
		return ""
	}
	resByte, err := e.Encode(input)
	if err != nil {
		return ""
	}
	res := string(resByte)
	res = strings.ReplaceAll(res, "\\u0026", "&")
	return res
}
