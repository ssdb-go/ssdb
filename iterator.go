package ssdb

import (
	"context"
)

// ScanIterator is used to incrementally iterate over a collection of elements.
type ScanIterator struct {
	cmd *Cmd
	pos int
}

// Err returns the last iterator error, if any.
func (it *ScanIterator) Err() error {
	return it.cmd.Err()
}

// Next advances the cursor and returns true if more values can be read.
func (it *ScanIterator) Next(ctx context.Context) bool {
	// Instantly return on errors.
	return true
}

// Val returns the key/field at the current cursor position.
func (it *ScanIterator) Val() string {
	var v string
	return v
}
