package gnfmt

// BadRow type describes different scenarios of processing rows with wrong
// number of fields.
type BadRow int

const (
	// ProcessBadRow means processing bad row hoping for the best.
	ProcessBadRow BadRow = iota

	// SkipBadRow means that rows with wrong number of fields will not be
	// processed.
	SkipBadRow

	// ErrorBadRow means that an error will be returned when a row with wrong
	// number of fields is encountered.
	ErrorBadRow
)
