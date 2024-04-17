package metrics

import (
	"time"
)

// TimeVar ...
type TimeVar struct{ v time.Time }

// Set ...
func (o *TimeVar) Set(date time.Time) { o.v = date }

// Add ...
func (o *TimeVar) Add(duration time.Duration) { o.v = o.v.Add(duration) }

// String ...
func (o *TimeVar) String() string { return o.v.Format(time.RFC3339) }
