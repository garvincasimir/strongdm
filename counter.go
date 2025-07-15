package main

import "time"

// Counter implements a leaky bucket algorithm to limit total calls per minute
// (CPM).
type Counter struct {
	// buckets is the current state of the rate limit buckets.
	buckets map[string]Bucket
}

// NewCounter creates a new rate limiting counter.
func NewCounter() *Counter {
	return &Counter{
		buckets: map[string]Bucket{},
	}
}

// Add checks the current value and size of the rate limit bucket specified by
// "key", based on the given limit per window. It returns Info about the bucket
// state, and true/false to indicate whether the value was successfully added to
// the bucket. If the limit is zero, it always returns success.
func (p *Counter) Add(key string, limitPerWindow int64, add int64) Info {
	if limitPerWindow == 0 {
		return Info{
			Bucket:  key,
			ResetAt: time.Now(),
			Allowed: true,
		}
	}

	now := time.Now()

	existingBucket := p.buckets[key]

	newBucket := existingBucket.Plus(now, limitPerWindow, add)

	newCount := newBucket.CountAt(now)
	bucketSize := BucketSize(limitPerWindow)
	if newCount > bucketSize {
		return Info{
			Bucket:     key,
			ResetAt:    existingBucket.WillReach(bucketSize-add, now),
			BucketSize: bucketSize,
			Remaining:  max(0, bucketSize-existingBucket.CountAt(now)),
			Allowed:    false,
		}
	}

	p.buckets[key] = newBucket

	remaining := bucketSize - newCount
	return Info{
		Bucket:     key,
		ResetAt:    newBucket.WillReach(bucketSize-1, now),
		BucketSize: bucketSize,
		Remaining:  remaining,
		Allowed:    true,
	}
}


// Info contains rate limit information produced by a rate limit check.
type Info struct {
	Bucket     string    `json:"bucket"`
	ResetAt    time.Time `json:"resetAt"`
	BucketSize int64     `json:"bucketSize"`
	Remaining  int64     `json:"remaining"`
	Allowed    bool      `json:"allowed"`
}