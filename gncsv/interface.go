package gncsv

import "context"

// GnCSV combines Reader and Writer interfaces, providing a unified
// interface for reading and writing CSV data.
type GnCSV interface {
	Reader
	Writer
}

// Reader defines an interface for reading CSV data.
type Reader interface {
	// ReadSlice reads a portion of the CSV data, starting at the given
	// offset and reading up to the specified limit. It returns a slice
	// of string slices, where each inner slice represents a row in the CSV.
	ReadSlice(offset, limit int) ([][]string, error)

	// Read reads all CSV data and sends each row as a string slice to
	// the provided channel. It returns the total number of rows read and
	// any error encountered. It uses a context for cancellation.
	Read(context.Context, chan<- []string) (int, error)
}

// Writer defines an interface for writing CSV data.
type Writer interface {
	// Write writes CSV data received from the provided channel. Each
	// string slice received from the channel represents a row in the CSV.
	// It uses a context for cancellation.
	Write(context.Context, <-chan []string) error
}
