package bucket

import (
	"math"
	"time"
)

const (
	// WindowDuration is the duration which all limits are specified in terms of
	// (calls per minute or CPM).
	WindowDuration = 1 * time.Minute

	// BurstTolerance determines the size of each bucket. the bucket is sized to
	// allow a burst of 1 second's worth of tokens, on top of the steady rate
	// limit. If 1 second's worth of tokens is less than 1, then the bucket size
	// is 1 and there is no burst tolerance.
	BurstTolerance = 1 * time.Second
)

// Bucket represents a single rate limit bucket.
type Bucket struct {
	UpdatedAt      time.Time
	LimitPerWindow int64
	Count          float64
}

// CountAt returns the count of the bucket at the given time.
func (b Bucket) CountAt(now time.Time) int64 {
	return int64(math.Ceil(b.countAt(now)))
}

// Plus returns a new copy of the Bucket with the given amount of tokens added
// to it.
func (b Bucket) Plus(now time.Time, limitPerWindow int64, add int64) Bucket {
	return Bucket{
		UpdatedAt:      now,
		LimitPerWindow: limitPerWindow,
		Count:          b.countAt(now) + float64(add),
	}
}

func (b Bucket) countAt(now time.Time) float64 {
	leakage := (float64(b.LimitPerWindow) * float64(now.Sub(b.UpdatedAt))) / float64(WindowDuration)
	return max(0.0, b.Count-leakage)
}

// WillReach returns the time at which the bucket will leak enough to reach the
// given count, or now if it is already at or below the count. If you pass a
// negative count, it will also return now.
func (b Bucket) WillReach(count int64, now time.Time) time.Time {
	if count < 0 {
		return now
	}
	needToLeak := b.Count - float64(count)
	if needToLeak <= 0 {
		return now
	}
	resetAt := b.UpdatedAt.Add(time.Duration(needToLeak*float64(WindowDuration)) / time.Duration(b.LimitPerWindow))
	if resetAt.Before(now) {
		return now
	}
	return resetAt
}

// Size determines the size of each bucket automatically based on the
// given limit per window. The bucket is sized to allow a burst of 1 second's
// worth of tokens, on top of the steady rate limit. If 1 second's worth of
// tokens is less than 1, then the bucket size is 1 and there is no burst
// tolerance.
func Size(limitPerWindow int64) int64 {
	a := limitPerWindow * int64(BurstTolerance)
	b := int64(WindowDuration)
	// positive integer ceiling division
	return (a + b - 1) / b
}
