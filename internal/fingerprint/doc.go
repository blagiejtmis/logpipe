// Package fingerprint provides deterministic SHA-256 fingerprinting of log
// records.
//
// A [Fingerprinter] is constructed from a [Config] that specifies which record
// fields contribute to the hash and what separator is placed between their
// string representations.
//
// # Usage
//
//	cfg := fingerprint.Config{
//		Fields:    []string{"host", "level", "message"},
//		Separator: "|",
//	}
//	fp, err := fingerprint.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	hash := fp.Compute(record)
//
// If Fields is empty, all keys present in the record are used in sorted order,
// guaranteeing a stable result regardless of map iteration order.
//
// The returned fingerprint is a lowercase hexadecimal SHA-256 digest (64
// characters).
package fingerprint
