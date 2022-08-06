package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Post struct {
	PostId     string `json:"postId"`
	Time       int    `json:"time"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	IsFinished string `json:"isFinished"`
}

type Posts []Post

func main() {
	http.HandleFunc("/", readPost)
	http.HandleFunc("/new", createPost)
	http.HandleFunc("/modify/:id", updatePost)
	http.HandleFunc("/delete/:id", deletePost)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func readPost(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "readPost is working!")
	content, err := ioutil.ReadFile("database/posts.json")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var posts Post
	json.Unmarshal(content, &posts)
	// getHtml(w, "templates/jobOffers.html", jobOffers)
}
func createPost(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "createPost is working!")
}
func updatePost(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "updatePost is working!")
}
func deletePost(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "deletePost is working!")
}

func respondError(w http.ResponseWriter, code int, errorMessage string) {
	respondJSON(w, code, map[string]string{"error": errorMessage})
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
