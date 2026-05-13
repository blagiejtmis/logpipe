// Package window implements a sliding-window event counter.
//
// A Counter divides a configurable time window into a fixed number of
// buckets. Each call to Add records events in the current bucket; as
// time advances, expired buckets are zeroed so that Count always
// reflects only events that occurred within the most recent window.
//
// Typical usage:
//
//	c, err := window.New(time.Minute, 60) // 60 one-second buckets
//	if err != nil {
//		log.Fatal(err)
//	}
//	c.Add(1)              // record one event
//	fmt.Println(c.Count()) // events in the last minute
//
// Counter is safe for concurrent use.
package window
