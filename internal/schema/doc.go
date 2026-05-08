// Package schema validates structured log records against a declared
// field schema.
//
// # Overview
//
// A schema is a list of FieldDef entries, each specifying a field name,
// expected type, and whether the field is required. The Validator checks
// every incoming record and returns an error on the first violation.
//
// # Supported Types
//
//   - string  – Go string
//   - number  – int, int64, float32, float64
//   - bool    – Go bool
//
// # Usage
//
//	v, err := schema.New(schema.Config{
//		Fields: []schema.FieldDef{
//			{Name: "level",   Type: schema.FieldTypeString, Required: true},
//			{Name: "latency", Type: schema.FieldTypeNumber, Required: false},
//		},
//	})
//	if err != nil { /* bad config */ }
//	if err := v.Validate(record); err != nil { /* drop or flag record */ }
package schema
