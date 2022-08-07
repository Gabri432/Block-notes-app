package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Post struct {
	PostId  string `json:"postId"`
	Time    int    `json:"time"`
	Title   string `json:"title"`
	Content string `json:"content,omitempty"`
	IsDraft bool   `json:"isDraft"`
}

type Posts []Post

func main() {
	http.HandleFunc("/", readPost)
	http.HandleFunc("/drafts", readPost)
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
	var drafts Posts
	json.Unmarshal(content, &posts)
	if route := strings.TrimPrefix(r.URL.Path, "/"); route == "saved" {
		drafts = getDrafts(posts)
	}
	htmlPage, err := template.ParseFiles("templates/main.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	if len(drafts) > 0 {
		htmlPage.Execute(w, drafts)
	} else {
		htmlPage.Execute(w, posts)
	}
}

func getDrafts(posts Posts) (drafts Posts) {
	for _, post := range posts {
		if !post.IsDraft {
			drafts = append(drafts, post)
		}
	}
	return
}

func createPost(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "templates/form.html", false)
	post := getFormData(w, r)
	if post.Title == "" {
		respondError(w, http.StatusNoContent, "Not title provided.")
	}
	savePost(w, post, "database/posts.json")
}

func getFormData(w http.ResponseWriter, r *http.Request) Post {
	data := r.PostForm
	title := data.Get("title")
	time := time.Now().Second()
	postId := title + "#" + string(rune(time))
	return Post{
		PostId: postId, Time: time, Title: title, Content: data.Get("content"), IsDraft: false,
	}
}
func updatePost(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "templates/form.html", false)
	content, err := ioutil.ReadFile("./database/posts.json")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	id := strings.TrimPrefix(r.URL.Path, "/modify/")
	var posts Posts
	json.Unmarshal(content, &posts)
	post := getPostById(w, posts, id)
	post = getFormData(w, r)
	if post.Title == "" {
		respondError(w, http.StatusNoContent, "Not title provided.")
	}
	savePost(w, post, "database/posts.json")
}
func deletePost(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "templates/form.html", true)
	content, err := ioutil.ReadFile("./database/posts.json")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
	}
	id := strings.TrimPrefix(r.URL.Path, "/delete/")
	var posts Posts
	json.Unmarshal(content, &posts)
	post := getPostById(w, posts, id)
	for i, p := range posts {
		if p == post {
			posts = append(posts[:i], posts[i+1:]...)
		}
	}
	if err := os.Truncate("database/posts.json", 0); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	for _, p := range posts {
		savePost(w, p, "database/posts.json")
	}
}

func getPostById(w http.ResponseWriter, posts Posts, id string) Post {
	for _, post := range posts {
		if post.PostId == id {
			respondJSON(w, http.StatusFound, post)
			return post
		}
	}
	respondError(w, http.StatusNotFound, "Post not found.")
	return Post{}
}

func savePost(w http.ResponseWriter, post Post, fileName string) {
	postByte, _ := json.Marshal(post)
	if ioutil.WriteFile(fileName, postByte, 0644) != nil {
		respondError(w, http.StatusInternalServerError, "Error while saving post.")
	}
}

func renderHTML(w http.ResponseWriter, htmlTemplate string, readMode bool) {
	htmlPage, err := template.ParseFiles(htmlTemplate, "templates/header.html", "templates/footer.html")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	htmlPage.Execute(w, readMode)
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
