package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestScanPosts(t *testing.T) {
	file, _ := os.Open("database/posts.json")
	bytes := ScanPosts(file, 15, 21)
	var posts []Post
	json.Unmarshal(bytes, &posts)
	if len(posts) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(posts))
	}
}
