package gnfmt

// BadRow type describes different scenarios of processing rows with wrong
// number of fields.
type BadRow int

const (
	// ErrorBadRow means that an error will be returned when a row with wrong
	// number of fields is encountered.
	ErrorBadRow BadRow = iota

	// SkipBadRow means that rows with wrong number of fields will not be
	// processed.
	SkipBadRow

	// ProcessBadRow means processing bad row hoping for the best.
	ProcessBadRow
)
