package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
)

func GetFormData(w http.ResponseWriter, r *http.Request) Post {
	err := r.ParseForm()
	if err != nil {
		RespondError(w, http.StatusNoContent, err.Error())
	}
	data := r.PostForm
	title := data.Get("title")
	time := time.Now()
	postId := title + "-" + string(rune(time.Second()))
	return Post{
		PostId: postId, Time: int(time.Unix()), Title: title, Content: data.Get("content"), IsDraft: false,
	}
}

func GetPostById(w http.ResponseWriter, posts Posts, id string) Post {
	for _, post := range posts {
		if post.PostId == id {
			return post
		}
	}
	RespondError(w, http.StatusNotFound, "Post not found.")
	return Post{}
}

func SavePost(w http.ResponseWriter, post Post, fileName string) {
	content, err := ioutil.ReadFile("database/posts.json")
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var posts Posts
	json.Unmarshal(content, &posts)
	posts = append(posts, post)
	postByte, _ := json.MarshalIndent(posts, "", " ")
	if ioutil.WriteFile(fileName, postByte, 0644) != nil {
		RespondError(w, http.StatusInternalServerError, "Error while saving post.")
	}
}

func RenderHTML(w http.ResponseWriter, htmlTemplate string, route string, post Post) {
	htmlPage, err := template.ParseFiles(htmlTemplate, "templates/header.html", "templates/footer.html")
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
	}
	data := Data{Content: post, Route: route}
	htmlPage.Execute(w, data)
}

func RespondError(w http.ResponseWriter, code int, errorMessage string) {
	RespondJSON(w, code, map[string]string{"error": errorMessage})
	fmt.Println(errorMessage)
}

func RespondJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
