package main

import (
	"testing"
	"time"
)

func TestNewCounter(t *testing.T) {
	counter := NewCounter()
	if counter == nil {
		t.Fatal("NewCounter() returned nil")
	}
	if counter.buckets == nil {
		t.Fatal("NewCounter() buckets map is nil")
	}
	if len(counter.buckets) != 0 {
		t.Error("NewCounter() buckets map should be empty")
	}
}

func TestCounter_Add_ZeroLimit(t *testing.T) {
	counter := NewCounter()

	info := counter.Add("test-key", 0, 10)

	if info.Bucket != "test-key" {
		t.Errorf("Expected bucket 'test-key', got '%s'", info.Bucket)
	}
	if !info.Allowed {
		t.Error("Expected allowed=true for zero limit")
	}
	if info.BucketSize != 0 {
		t.Errorf("Expected BucketSize=0, got %d", info.BucketSize)
	}
	if info.Remaining != 0 {
		t.Errorf("Expected Remaining=0, got %d", info.Remaining)
	}
}

func TestCounter_Add_Success(t *testing.T) {
	counter := NewCounter()
	limitPerWindow := int64(60) // 1 per second

	info := counter.Add("test-key", limitPerWindow, 1)

	if info.Bucket != "test-key" {
		t.Errorf("Expected bucket 'test-key', got '%s'", info.Bucket)
	}
	if !info.Allowed {
		t.Error("Expected allowed=true for first request")
	}
	if info.BucketSize != 1 {
		t.Errorf("Expected BucketSize=1, got %d", info.BucketSize)
	}
	if info.Remaining != 0 {
		t.Errorf("Expected Remaining=0, got %d", info.Remaining)
	}
}

func TestCounter_Add_Rejection(t *testing.T) {
	counter := NewCounter()
	limitPerWindow := int64(60) // 1 per second

	// Fill the bucket
	counter.Add("test-key", limitPerWindow, 1)

	// Try to add another token - should be rejected
	info := counter.Add("test-key", limitPerWindow, 1)

	if info.Allowed {
		t.Error("Expected allowed=false for exceeded limit")
	}
	if info.BucketSize != 1 {
		t.Errorf("Expected BucketSize=1, got %d", info.BucketSize)
	}
	if info.Remaining != 0 {
		t.Errorf("Expected Remaining=0, got %d", info.Remaining)
	}
}

func TestCounter_Add_LargeLimit(t *testing.T) {
	counter := NewCounter()
	limitPerWindow := int64(120) // 2 per second

	// First request should succeed
	info1 := counter.Add("test-key", limitPerWindow, 1)
	if !info1.Allowed {
		t.Error("First request should be allowed")
	}
	if info1.BucketSize != 2 {
		t.Errorf("Expected BucketSize=2, got %d", info1.BucketSize)
	}
	if info1.Remaining != 1 {
		t.Errorf("Expected Remaining=1, got %d", info1.Remaining)
	}

	// Second request should succeed
	info2 := counter.Add("test-key", limitPerWindow, 1)
	if !info2.Allowed {
		t.Error("Second request should be allowed")
	}
	if info2.Remaining != 0 {
		t.Errorf("Expected Remaining=0, got %d", info2.Remaining)
	}

	// Third request should fail
	info3 := counter.Add("test-key", limitPerWindow, 1)
	if info3.Allowed {
		t.Error("Third request should be rejected")
	}
}

func TestCounter_Add_MultipleKeys(t *testing.T) {
	counter := NewCounter()
	limitPerWindow := int64(60) // 1 per second

	// Add to first key
	info1 := counter.Add("key1", limitPerWindow, 1)
	if !info1.Allowed {
		t.Error("First key should be allowed")
	}

	// Add to second key - should be independent
	info2 := counter.Add("key2", limitPerWindow, 1)
	if !info2.Allowed {
		t.Error("Second key should be allowed")
	}

	// Try to add to first key again - should be rejected
	info3 := counter.Add("key1", limitPerWindow, 1)
	if info3.Allowed {
		t.Error("First key second request should be rejected")
	}
}

func TestCounter_Add_LeakageOverTime(t *testing.T) {
	counter := NewCounter()
	limitPerWindow := int64(60) // 1 per second

	// Fill the bucket
	counter.Add("test-key", limitPerWindow, 1)

	// Simulate time passage by directly manipulating the bucket
	now := time.Now()
	counter.buckets["test-key"] = Bucket{
		UpdatedAt:      now.Add(-2 * time.Second),
		LimitPerWindow: limitPerWindow,
		Count:          1.0,
	}

	// Should be allowed now due to leakage
	info := counter.Add("test-key", limitPerWindow, 1)
	if !info.Allowed {
		t.Error("Request should be allowed after leakage")
	}
}

func TestCounter_Add_ResetAt(t *testing.T) {
	counter := NewCounter()
	limitPerWindow := int64(60) // 1 per second

	// Fill the bucket
	counter.Add("test-key", limitPerWindow, 1)

	// Try to add another - should be rejected with reset time
	info := counter.Add("test-key", limitPerWindow, 1)
	if info.Allowed {
		t.Error("Request should be rejected")
	}

	now := time.Now()
	if info.ResetAt.Before(now) {
		t.Error("ResetAt should be in the future")
	}
}

func TestCounter_Add_AddMultipleTokens(t *testing.T) {
	counter := NewCounter()
	limitPerWindow := int64(180) // 3 per second

	// Add 2 tokens at once
	info := counter.Add("test-key", limitPerWindow, 2)
	if !info.Allowed {
		t.Error("Adding 2 tokens should be allowed")
	}
	if info.BucketSize != 3 {
		t.Errorf("Expected BucketSize=3, got %d", info.BucketSize)
	}
	if info.Remaining != 1 {
		t.Errorf("Expected Remaining=1, got %d", info.Remaining)
	}

	// Try to add 2 more tokens - should be rejected
	info2 := counter.Add("test-key", limitPerWindow, 2)
	if info2.Allowed {
		t.Error("Adding 2 more tokens should be rejected")
	}
}

func TestCounter_Add_NonExistentKey(t *testing.T) {
	counter := NewCounter()
	limitPerWindow := int64(60) // 1 per second

	// Verify the key doesn't exist initially
	if _, exists := counter.buckets["new-key"]; exists {
		t.Error("Key should not exist initially")
	}

	// Add to non-existent key - should create new bucket and allow
	info := counter.Add("new-key", limitPerWindow, 1)

	if !info.Allowed {
		t.Error("First request to non-existent key should be allowed")
	}
	if info.Bucket != "new-key" {
		t.Errorf("Expected bucket 'new-key', got '%s'", info.Bucket)
	}
	if info.BucketSize != 1 {
		t.Errorf("Expected BucketSize=1, got %d", info.BucketSize)
	}
	if info.Remaining != 0 {
		t.Errorf("Expected Remaining=0, got %d", info.Remaining)
	}

	// Verify the key now exists in the map
	if _, exists := counter.buckets["new-key"]; !exists {
		t.Error("Key should exist after first request")
	}
}

func TestCounter_Add_EdgeCases(t *testing.T) {
	counter := NewCounter()

	// Test with very small limit
	info1 := counter.Add("small", 1, 1)
	if !info1.Allowed {
		t.Error("Small limit should allow first request")
	}
	if info1.BucketSize != 1 {
		t.Errorf("Expected BucketSize=1, got %d", info1.BucketSize)
	}

	// Test with large limit
	info2 := counter.Add("large", 3600, 1)
	if !info2.Allowed {
		t.Error("Large limit should allow request")
	}
	if info2.BucketSize != 60 {
		t.Errorf("Expected BucketSize=60, got %d", info2.BucketSize)
	}
}

func TestInfo_Fields(t *testing.T) {
	counter := NewCounter()
	limitPerWindow := int64(120) // 2 per second

	info := counter.Add("test-key", limitPerWindow, 1)

	// Check all fields are set correctly
	if info.Bucket != "test-key" {
		t.Errorf("Expected Bucket='test-key', got '%s'", info.Bucket)
	}
	if info.ResetAt.IsZero() {
		t.Error("ResetAt should not be zero")
	}
	if info.BucketSize != 2 {
		t.Errorf("Expected BucketSize=2, got %d", info.BucketSize)
	}
	if info.Remaining != 1 {
		t.Errorf("Expected Remaining=1, got %d", info.Remaining)
	}
	if !info.Allowed {
		t.Error("Expected Allowed=true")
	}
}
