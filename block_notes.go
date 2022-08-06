package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Post struct {
	PostId     string `json:"postId"`
	Time       int    `json:"time"`
	Title      string `json:"title"`
	Content    string `json:"content,omitempty"`
	IsFinished bool   `json:"isFinished"`
}

type Posts []Post

func main() {
	http.HandleFunc("/", readPost)
	http.HandleFunc("/saved", readPost)
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
	}
	var posts Posts
	var savedPosts Posts
	json.Unmarshal(content, &posts)
	if route := strings.TrimPrefix(r.URL.Path, "/"); route == "saved" {
		savedPosts = getSavedPosts(posts)
	}
	htmlPage, err := template.ParseFiles("main.html", "templates/header.html", "templates/input.html")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	if len(savedPosts) > 0 {
		htmlPage.Execute(w, savedPosts)
	} else {
		htmlPage.Execute(w, posts)
	}
}

func getSavedPosts(posts Posts) (savedPosts Posts) {
	for _, post := range posts {
		if !post.IsFinished {
			savedPosts = append(savedPosts, post)
		}
	}
	return
}

func createPost(w http.ResponseWriter, r *http.Request) {
	post := getFormData(w, r)
	if post.Title == "" {
		respondError(w, http.StatusNoContent, "Not title provided.")
	}
	respondJSON(w, http.StatusOK, "createPost is working!")
}

func getFormData(w http.ResponseWriter, r *http.Request) Post {
	data := r.PostForm
	title := data.Get("title")
	time := time.Now().Second()
	postId := title + "#" + string(rune(time))
	return Post{
		PostId: postId, Time: time, Title: title, Content: data.Get("content"), IsFinished: false,
	}
}
func updatePost(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "updatePost is working!")
}
func deletePost(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "deletePost is working!")
}

func respondError(w http.ResponseWriter, code int, errorMessage string) {
	respondJSON(w, code, map[string]string{"error": errorMessage})
	log.Fatal(errorMessage)
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
