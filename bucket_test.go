package main

import (
	"testing"
	"time"
)

func TestBucket_CountAt(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		bucket   Bucket
		at       time.Time
		expected int64
	}{
		{
			name: "no leakage",
			bucket: Bucket{
				UpdatedAt:      now,
				LimitPerWindow: 100,
				Count:          50.0,
			},
			at:       now,
			expected: 50,
		},
		{
			name: "partial leakage",
			bucket: Bucket{
				UpdatedAt:      now,
				LimitPerWindow: 60, // 1 per second
				Count:          30.0,
			},
			at:       now.Add(10 * time.Second),
			expected: 20,
		},
		{
			name: "complete leakage",
			bucket: Bucket{
				UpdatedAt:      now,
				LimitPerWindow: 60,
				Count:          30.0,
			},
			at:       now.Add(time.Minute),
			expected: 0,
		},
		{
			name: "ceiling behavior",
			bucket: Bucket{
				UpdatedAt:      now,
				LimitPerWindow: 60,
				Count:          10.1,
			},
			at:       now,
			expected: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bucket.CountAt(tt.at)
			if result != tt.expected {
				t.Errorf("CountAt() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestBucket_Plus(t *testing.T) {
	now := time.Now()

	bucket := Bucket{
		UpdatedAt:      now.Add(-10 * time.Second),
		LimitPerWindow: 60,
		Count:          20.0,
	}

	result := bucket.Plus(now, 60, 5)

	if result.UpdatedAt != now {
		t.Errorf("Plus() UpdatedAt = %v, expected %v", result.UpdatedAt, now)
	}
	if result.LimitPerWindow != 60 {
		t.Errorf("Plus() LimitPerWindow = %d, expected 60", result.LimitPerWindow)
	}
	// After 10 seconds, 10 tokens should have leaked, so 20 - 10 + 5 = 15
	expected := 15.0
	if result.Count != expected {
		t.Errorf("Plus() Count = %f, expected %f", result.Count, expected)
	}
}

func TestBucket_WillReach(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		bucket   Bucket
		count    int64
		now      time.Time
		expected time.Time
	}{
		{
			name: "negative count returns now",
			bucket: Bucket{
				UpdatedAt:      now,
				LimitPerWindow: 60,
				Count:          30.0,
			},
			count:    -1,
			now:      now,
			expected: now,
		},
		{
			name: "already at or below count",
			bucket: Bucket{
				UpdatedAt:      now,
				LimitPerWindow: 60,
				Count:          20.0,
			},
			count:    30,
			now:      now,
			expected: now,
		},
		{
			name: "exact count",
			bucket: Bucket{
				UpdatedAt:      now,
				LimitPerWindow: 60,
				Count:          30.0,
			},
			count:    30,
			now:      now,
			expected: now,
		},
		{
			name: "future reset time",
			bucket: Bucket{
				UpdatedAt:      now,
				LimitPerWindow: 60, // 1 per second
				Count:          30.0,
			},
			count:    20,
			now:      now,
			expected: now.Add(10 * time.Second),
		},
		{
			name: "reset time in past returns now",
			bucket: Bucket{
				UpdatedAt:      now.Add(-2 * time.Minute),
				LimitPerWindow: 60,
				Count:          30.0,
			},
			count:    20,
			now:      now,
			expected: now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bucket.WillReach(tt.count, tt.now)
			if !result.Equal(tt.expected) {
				t.Errorf("WillReach() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestBucketSize(t *testing.T) {
	tests := []struct {
		name           string
		limitPerWindow int64
		expected       int64
	}{
		{
			name:           "small limit",
			limitPerWindow: 1,
			expected:       1, // (1 * 1 + 60 - 1) / 60 = 1
		},
		{
			name:           "medium limit",
			limitPerWindow: 60,
			expected:       1, // (60 * 1 + 60 - 1) / 60 = 1
		},
		{
			name:           "large limit",
			limitPerWindow: 120,
			expected:       2, // (120 * 1 + 60 - 1) / 60 = 2
		},
		{
			name:           "very large limit",
			limitPerWindow: 3600,
			expected:       60, // (3600 * 1 + 60 - 1) / 60 = 60
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BucketSize(tt.limitPerWindow)
			if result != tt.expected {
				t.Errorf("BucketSize() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestBucket_Integration(t *testing.T) {
	now := time.Now()
	limitPerWindow := int64(60) // 1 per second

	// Start with empty bucket
	bucket := Bucket{
		UpdatedAt:      now,
		LimitPerWindow: limitPerWindow,
		Count:          0.0,
	}

	// Add some tokens
	bucket = bucket.Plus(now, limitPerWindow, 10)
	if bucket.CountAt(now) != 10 {
		t.Errorf("After adding 10 tokens, count should be 10, got %d", bucket.CountAt(now))
	}

	// Check count after 5 seconds (5 tokens should have leaked)
	later := now.Add(5 * time.Second)
	if bucket.CountAt(later) != 5 {
		t.Errorf("After 5 seconds, count should be 5, got %d", bucket.CountAt(later))
	}

	// Check when bucket will reach 2 tokens
	willReach := bucket.WillReach(2, now)
	expected := now.Add(8 * time.Second)
	if !willReach.Equal(expected) {
		t.Errorf("WillReach(2) should be %v, got %v", expected, willReach)
	}
}
