package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	isDraft, err := strconv.ParseBool(data.Get("isDraft"))
	if err != nil {
		isDraft = true
	}
	return Post{
		PostId: postId, Time: int(time.Unix()), Title: title, Content: data.Get("content"), IsDraft: isDraft,
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

func RemovePost(posts Posts, post Post) Posts {
	for i, p := range posts {
		if p == post {
			posts = append(posts[:i], posts[i+1:]...)
		}
	}
	return posts
}

func UnmarshalPost(w http.ResponseWriter, jsonContent []byte, postId string) (Post, Posts) {
	var posts Posts
	json.Unmarshal(jsonContent, &posts)
	return GetPostById(w, posts, postId), posts
}

func ReversePosts(list Posts) Posts {
	var newList Posts
	for i := len(list) - 1; i >= 0; i-- {
		newList = append(newList, list[i])
	}
	return newList
}

func UpdatePostList(w http.ResponseWriter, fileName string, postList Posts) {
	if err := os.Truncate("database/posts.json", 0); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, p := range postList {
		SavePost(w, p, "database/posts.json")
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
	//RespondJSON(w, code, map[string]string{"error": errorMessage})
	errorPost := Post{Title: "Error " + strconv.Itoa(code), Content: errorMessage}
	RenderHTML(w, "templates/error.html", "/error", errorPost)
}

func RespondJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
