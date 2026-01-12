package gnfmt

import (
	"bytes"
	"encoding/gob"
)

// GNgob is for serializing data into gob format.
type GNgob struct{}

// Encode takes an object and serializes it into gob blob. It returns an
// error in case of encoding problems.
func (e GNgob) Encode(input any) ([]byte, error) {
	var respBytes bytes.Buffer
	enc := gob.NewEncoder(&respBytes)
	if err := enc.Encode(input); err != nil {
		return nil, err
	}
	return respBytes.Bytes(), nil
}

// Decode deserializes gob bytes into a go object. It returns an error if
// decoding fails.
func (e GNgob) Decode(input []byte, output any) error {
	b := bytes.NewBuffer(input)
	dec := gob.NewDecoder(b)
	return dec.Decode(output)
}
