package gncsv

import "context"

// GnCSV combines Reader and Writer interfaces, providing a unified
// interface for reading and writing CSV data.
type GnCSV interface {
	Reader
	Writer
}

// Reader is an interface that defines methods for reading and processing CSV
// data. It provides functionality to read data in slices, stream rows through
// a channel, and retrieve headers.
type Reader interface {
	// ReadSlice reads a portion of the CSV data starting at the specified offset
	// and up to the given limit. Each row is represented as a slice of strings.
	//
	// Parameters:
	//   - offset: The starting position of the rows to read.
	//   - limit: The maximum number of rows to read.
	//
	// Returns:
	//   - [][]string: A slice of string slices, where each inner slice
	//     represents a row in the CSV.
	//   - error: An error if any issue occurs during reading; otherwise, nil.
	ReadSlice(offset, limit int) ([][]string, error)

	// Read reads all rows from the CSV data and sends each row as a slice of
	// strings to the provided channel. The function processes the data until the
	// end of the file or an error occurs.
	//
	// Parameters:
	//   - ctx: A context to control the cancellation or timeout of the read
	//     operation.
	//   - chOut: A channel to which the function sends each row of the CSV as a
	//     slice of strings.
	//
	// Returns:
	//   - int: The total number of rows read.
	//   - error: An error if any issue occurs during reading; otherwise, nil.
	Read(context.Context, chan<- []string) (int, error)

	// ReadChunks reads data in chunks of the specified size and sends each chunk
	// to the provided channel chOut. The function returns the total number of
	// records read and an error, if any occurs during the process.
	//
	// Parameters:
	//   - ctx: A context to control the cancellation or timeout of the read
	//     operation.
	//   - chOut: A channel to which the function sends slices of string slices,
	//     representing chunks of data.
	//   - chunkSize: The size of each chunk to be read.
	//
	// Returns:
	//   - int: The total number of records read.
	//   - error: An error if any issue occurs during reading; otherwise, nil.
	ReadChunks(
		ctx context.Context,
		chOut chan<- [][]string,
		chunkSize int,
	) (int, error)

	// Headers returns the headers of the CSV file as a slice of strings.
	Headers() []string

	// F is a field accessor. If the field with the given name exists, it returns
	// the value of the field in the row. If not, returns empty string and an
	// error.
	F(row []string, f string) (string, error)
}

// Writer defines an interface for writing CSV data.
type Writer interface {
	// WriteStream writes CSV data received from the provided channel. Each
	// string slice received from the channel represents a row in the CSV.
	// It uses a context for cancellation.
	WriteStream(context.Context, <-chan []string) error
}
