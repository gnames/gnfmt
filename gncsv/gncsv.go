package gncsv

import "github.com/gnames/gnfmt/gncsv/config"

// New creates a new CSV or TSV/PSV reader/writer based on the provided
// configuration. If the ColSep in the config is a comma, it creates
// a CSV reader/writer. Otherwise, it creates a TSV reader/writer.
func New(cfg config.Config) GnCSV {
	if cfg.ColSep == ',' || cfg.WithQuotes {
		return NewCSV(cfg)
	}
	return NewTSV(cfg)
}
