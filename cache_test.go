package main

import "testing"

func newTestCache() *Cache {
	freqNode := NewFrequnecyNode()
	list := NewFrequencyList(freqNode)
	return NewCache(list)
}

func TestNewCacheInitializes(t *testing.T) {
	cache := newTestCache()
	if cache == nil {
		t.Fatal("expected non-nil cache")
	}
	if cache.Items == nil {
		t.Error("expected Items map to be initialized")
	}
	if cache.FreqList.Head == nil {
		t.Error("expected FreqList.Head to be non-nil")
	}
}

func TestInsertFirstItem(t *testing.T) {
	cache := newTestCache()
	msg, err := cache.InsertToCache("David", "Engineer")
	if err != nil {
		t.Fatalf("expected no error on first insert, got: %v", err)
	}
	if msg == "" {
		t.Error("expected non-empty message on insert")
	}
	if _, ok := cache.Items["David"]; !ok {
		t.Error("expected key 'David' to be in cache.Items after insert")
	}
}

func TestInsertSecondItem(t *testing.T) {
	cache := newTestCache()
	if _, err := cache.InsertToCache("David", "Engineer"); err != nil {
		t.Fatalf("first insert failed: %v", err)
	}
	if _, err := cache.InsertToCache("musa", "Doctor"); err != nil {
		t.Fatalf("expected no error on second insert, got: %v", err)
	}
	if _, ok := cache.Items["David"]; !ok {
		t.Error("expected key 'David' to still be in cache.Items")
	}
	if _, ok := cache.Items["musa"]; !ok {
		t.Error("expected key 'musa' to be in cache.Items after second insert")
	}
}

func TestAccessExistingKey(t *testing.T) {
	cache := newTestCache()
	if _, err := cache.InsertToCache("David", "Engineer"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	msg, node := cache.AccessAndIncrement("David")
	if msg != "successful" {
		t.Errorf("expected message 'successful', got %q", msg)
	}
	if node == nil {
		t.Fatal("expected non-nil node on cache hit")
	}
	if node.Cachedkey != "David" {
		t.Errorf("expected node key 'David', got %q", node.Cachedkey)
	}
}

func TestAccessMissingKey(t *testing.T) {
	cache := newTestCache()
	if _, err := cache.InsertToCache("David", "Engineer"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	msg, node := cache.AccessAndIncrement("Ghost")
	if msg != "failed to find node" {
		t.Errorf("expected message 'failed to find node', got %q", msg)
	}
	if node != nil {
		t.Errorf("expected nil node on cache miss, got %v", node)
	}
}

func TestAccessIncrementsFrequency(t *testing.T) {
	cache := newTestCache()
	if _, err := cache.InsertToCache("David", "Engineer"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	_, node := cache.AccessAndIncrement("David")
	if node == nil {
		t.Fatal("expected non-nil node after access")
	}
	if node.Parent == nil {
		t.Fatal("expected node.Parent to be non-nil after access")
	}
	if node.Parent.Frequency != 2 {
		t.Errorf("expected frequency 2 after one access, got %d", node.Parent.Frequency)
	}
}
