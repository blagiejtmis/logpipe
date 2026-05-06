// Package ratelimit implements a token-bucket-style rate limiter for
// controlling the maximum number of log lines ingested per source within
// a configurable time window.
//
// Usage:
//
//	// Allow at most 1000 lines per second per source.
//	limiter, err := ratelimit.New(1000, time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if limiter.Allow(record.Source) {
//		// forward record to pipeline
//	}
//
// Each source maintains its own independent counter that resets at the
// start of each window. The limiter is safe for concurrent use.
//
// A zero or negative limit is rejected by New and will return an error.
// A zero or negative window duration is likewise rejected.
package ratelimit
