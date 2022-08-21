package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetPostById(t *testing.T) {
	var w http.ResponseWriter
	var posts Posts
	posts = append(posts, Post{PostId: "", Title: ""})
	post := GetPostById(w, Posts{}, "nonExistingId")
	if post.Title != "" {
		t.Fatalf("Got %s, expected ''.", post.Title)
	}
}
func TestRemovePost(t *testing.T) {
	var posts Posts
	post := Post{Title: "Lorem Ispum", Content: "Lorem Ispum"}
	posts = append(posts, post)
	newPostsList := RemovePost(posts, post)
	if len(newPostsList) != 0 {
		t.Fatal("Expected to have a list with zero posts remaining.")
	}
	return
}
func TestSearchPostInJSON(t *testing.T) {
	var w http.ResponseWriter
	var posts Posts
	postInput := Post{PostId: "id"}
	posts = append(posts, postInput)
	jsonContent, _ := json.Marshal(posts)
	postOutput, _ := SearchPostInJSON(w, jsonContent, postInput.PostId)
	if postOutput.PostId != postInput.PostId {
		t.Fatalf("Expected to find %s, got %s", postInput.PostId, postOutput.PostId)
	}
}
func TestReversePosts(t *testing.T) {
	var posts Posts
	post1 := Post{Title: "Lorem Ispum", Content: "Lorem Ispum"}
	post2 := Post{Title: "Ipsum Lorem", Content: "Ipsum Lorem"}
	posts = append(posts, post1, post2)
	newPostsList := ReversePosts(posts)
	if newPostsList[0] != posts[1] {
		t.Fatal("Expected new slice first element to be the old slice second element.")
	}
}
