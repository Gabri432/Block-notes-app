package main

import (
	"net/http"
	"testing"
)

func TestGetFormData(t *testing.T) {
	return
}
func TestGetPostById(t *testing.T) {
	var w http.ResponseWriter
	var posts Posts
	posts = append(posts, Post{PostId: "", Title: ""})
	post := GetPostById(w, Posts{}, "nonExistingId")
	if post.Title != "" {
		t.Fatalf("Got %s, expected ''.", post.Title)
	}
}
func TestSavePost(t *testing.T) {
	return
}
func TestRemovePost(t *testing.T) {
	return
}
func TestUnmarshalPost(t *testing.T) {
	return
}
func TestReversePosts(t *testing.T) {
	return
}
