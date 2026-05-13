// Package flatten implements a log-record processor that expands nested
// map-typed fields into dot-notation (or custom-separator) top-level keys.
//
// # Overview
//
// Many log producers emit structured payloads where metadata is grouped under
// a single key (e.g. "meta", "labels", "kubernetes").  Downstream sinks and
// routing rules often need to address individual sub-fields directly.  The
// flatten processor solves this by walking each configured field and writing
// its leaf values back into the top-level record.
//
// # Usage
//
//	f, err := flatten.New([]flatten.Rule{
//		{Field: "meta", Separator: ".", DropSource: true},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	out := f.Apply(record)
//
// Deeply nested maps are handled recursively; non-map values at any level are
// written as-is.  When DropSource is true the original nested key is removed
// from the record after expansion.
package flatten
