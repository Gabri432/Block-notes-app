package main

import (
	"encoding/json"
	"net/http"
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

func TestReadFromTo(t *testing.T) {
	var w http.ResponseWriter
	data, _ := os.Stat("database/posts.json")
	bytes := ReadFromTo(w, "database/posts.json", int(data.Size()), 0)
	content, _ := os.ReadFile("database/posts.json")
	if len(bytes) != len(content) {
		t.Fatalf("Expected len(bytes) == %d, got len(bytes) == %d.", len(content), len(bytes))
	}
}

func TestPaginatePosts(t *testing.T) {
	file, _ := os.Open("database/posts.json")
	bytes := PaginatePosts(file, 1)
	var posts []Post
	json.Unmarshal(bytes, &posts)
	if len(posts) <= 5 {
		t.Fatalf("Expected at maximum 5 elements, got %d.", len(posts))
	}
}
