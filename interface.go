package gnfmt

// Outputter interface us a uniform way to create an output of a datum
type Outputter interface {
	// FormattedOutput takes a record and returns a string representation of
	// the record accourding to supplied format.
	Output(record any, f Format) string
}

// Encoder interface allows to switch between different encoding types.
type Encoder interface {
	//Encode takes a Go object and converts it into bytes
	Encode(input any) ([]byte, error)
	// Decode takes an input of bytes and decodes it into Go object.
	Decode(input []byte, output any) error
}
