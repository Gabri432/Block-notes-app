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
	PostId  string `json:"postId"`
	Time    int    `json:"time"`
	Title   string `json:"title"`
	Content string `json:"content,omitempty"`
	IsDraft bool   `json:"isDraft"`
}

type Posts []Post

type Data struct {
	Content Post
	Route   string
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
	if route := strings.TrimPrefix(r.URL.Path, "/"); route == "drafts" {
		drafts = getDrafts(posts)
	}
	htmlPage, err := template.ParseFiles("templates/main.html", "templates/header.html", "templates/footer.html", "templates/post.html")
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(drafts) > 0 {
		htmlPage.Execute(w, ReversePosts(drafts))
	} else {
		htmlPage.Execute(w, ReversePosts(posts))
	}
}

func getDrafts(posts Posts) (drafts Posts) {
	for _, post := range posts {
		if post.IsDraft {
			drafts = append(drafts, post)
		}
	}
	return
}

func createPost(w http.ResponseWriter, r *http.Request) {
	post := GetFormData(w, r)
	post.Title = strings.TrimSpace(post.Title)
	RenderHTML(w, "templates/form.html", "/new", Post{})
	if http.MethodPost != r.Method {
		return
	}
	if post.Title == "" {
		RespondError(w, http.StatusNoContent, "No title provided.")
		return
	}
	SavePost(w, post, "database/posts.json")
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	contentALL, _ := ioutil.ReadAll(r.Body)
	var newPost Post
	json.Unmarshal(contentALL, &newPost)
	content, err := ioutil.ReadFile("./database/posts.json")
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/modify/")
	post, posts := UnmarshalPost(w, content, id)
	RenderHTML(w, "templates/form.html", "/modify/", post)
	if http.MethodPut != r.Method {
		return
	}
	newPost = Post{PostId: post.PostId, Time: int(time.Now().Unix()), Content: newPost.Content, Title: newPost.Title, IsDraft: newPost.IsDraft}
	newPostList := RemovePost(posts, post)
	newPostList = append(newPostList, newPost)
	if newPost.Title == "" {
		RespondError(w, http.StatusNoContent, "No title provided.")
		return
	}
	UpdatePostList(w, "database/posts.json", newPostList)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("./database/posts.json")
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, err.Error())
	}
	id := strings.TrimPrefix(r.URL.Path, "/delete/")
	post, posts := UnmarshalPost(w, content, id)
	RenderHTML(w, "templates/form.html", "/delete/", post)
	newPostList := RemovePost(posts, post)
	UpdatePostList(w, "database/posts.json", newPostList)
}
