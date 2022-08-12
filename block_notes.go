package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Post struct {
	PostId  string `json:"postId"`
	Time    int    `json:"time"`
	Title   string `json:"title"`
	Content string `json:"content,omitempty"`
	IsDraft bool   `json:"isDraft"`
}

type Posts []Post

type Data struct {
	Content  Post
	ReadMode bool
}

func main() {
	http.HandleFunc("/", readPost)
	http.HandleFunc("/drafts", readPost)
	http.HandleFunc("/new", createPost)
	http.HandleFunc("/modify/", updatePost)
	http.HandleFunc("/delete/", deletePost)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func readPost(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("database/posts.json")
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var posts Posts
	var drafts Posts
	if err := json.Unmarshal(content, &posts); err != nil {
		log.Println(err.Error())
	}
	if route := strings.TrimPrefix(r.URL.Path, "/"); route == "saved" {
		drafts = getDrafts(posts)
	}
	htmlPage, err := template.ParseFiles("templates/main.html", "templates/header.html", "templates/footer.html", "templates/post.html")
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
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
	RenderHTML(w, "templates/form.html", false, Post{})
	if http.MethodPost != r.Method {
		return
	}
	post := GetFormData(w, r)
	if post.Title == "" {
		RespondError(w, http.StatusNoContent, "Not title provided.")
		return
	}
	SavePost(w, post, "database/posts.json")
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("./database/posts.json")
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/modify/")
	var posts Posts
	json.Unmarshal(content, &posts)
	log.Print(id, r.URL.Path)
	post := GetPostById(w, posts, id)
	RenderHTML(w, "templates/form.html", false, post)
	if http.MethodPut != r.Method {
		return
	}
	post = GetFormData(w, r)
	if post.Title == "" {
		RespondError(w, http.StatusNoContent, "Not title provided.")
		return
	}
	SavePost(w, post, "database/posts.json")
}
func deletePost(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("./database/posts.json")
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, err.Error())
	}
	id := strings.TrimPrefix(r.URL.Path, "/delete/")
	var posts Posts
	json.Unmarshal(content, &posts)
	post := GetPostById(w, posts, id)
	RenderHTML(w, "templates/form.html", true, post)
	for i, p := range posts {
		if p == post {
			posts = append(posts[:i], posts[i+1:]...)
		}
	}
	if err := os.Truncate("database/posts.json", 0); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, p := range posts {
		SavePost(w, p, "database/posts.json")
	}
}
