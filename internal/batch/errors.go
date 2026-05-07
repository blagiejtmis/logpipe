package batch

import "errors"

// Sentinel errors returned by New.
var (
	ErrInvalidMaxSize = errors.New("batch: maxSize must be >= 1")
	ErrInvalidInterval = errors.New("batch: interval must be > 0")
	ErrNilFlushFunc    = errors.New("batch: flush function must not be nil")
)
