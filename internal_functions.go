package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetFormData(w http.ResponseWriter, r *http.Request) Post {
	err := r.ParseForm()
	if err != nil {
		RespondError(w, http.StatusNoContent, err.Error())
	}
	data := r.PostForm
	title := data.Get("title")
	time := int(time.Now().Unix())
	postId := title + "-" + strconv.Itoa(time)
	isDraft, err := strconv.ParseBool(data.Get("isDraft"))
	if err != nil {
		isDraft = true
	}
	return Post{
		PostId: postId, Time: time, Title: title, Content: data.Get("content"), IsDraft: isDraft,
	}
}

func GetPostById(w http.ResponseWriter, posts Posts, id string) Post {
	if len(posts) < 1 {
		return Post{}
	}
	for _, post := range posts {
		if post.PostId == id {
			return post
		}
	}
	RespondError(w, http.StatusNotFound, "Post not found.")
	return Post{}
}

func SavePost(w http.ResponseWriter, post Post, fileName string) error {
	content, err := os.ReadFile(fileName)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return err
	}
	var posts Posts
	json.Unmarshal(content, &posts)
	posts = append(posts, post)
	copy(posts[1:], posts)
	posts[0] = post
	postByte, _ := json.MarshalIndent(posts, "", " ")
	if os.WriteFile(fileName, postByte, 0644) != nil {
		RespondError(w, http.StatusInternalServerError, "Error while saving post.")
	}
	return err
}

func RemovePost(posts Posts, post Post) Posts {
	for i, p := range posts {
		if p == post {
			posts = append(posts[:i], posts[i+1:]...)
		}
	}
	return posts
}

func SearchPostInJSON(w http.ResponseWriter, jsonContent []byte, postId string) (Post, Posts) {
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
	if err := os.Truncate(fileName, 0); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, p := range postList {
		SavePost(w, p, fileName)
	}
}

func FormattingPost(newPost Post) Post {
	newPostTime := int(time.Now().Unix())
	newPostTitle := strings.TrimSpace(newPost.Title)
	newPostId := newPostTitle + strconv.Itoa(newPostTime)
	return Post{PostId: newPostId, Time: newPostTime, Content: newPost.Content, Title: newPostTitle, IsDraft: newPost.IsDraft}
}

func ReadFrom(w http.ResponseWriter, fileName string, endPosition, startPosition int) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return []byte{}
	}
	defer file.Close()
	offset, errorMessage := file.Seek(int64(endPosition), startPosition)
	if errorMessage != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return []byte{}
	}
	newContent := make([]byte, offset)
	file.Read(newContent)
	return newContent
}

func RenderHTML(w http.ResponseWriter, htmlTemplate string, route string, post Post) {
	htmlPage, err := template.ParseFiles(htmlTemplate, "templates/header.html", "templates/footer.html")
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	data := Data{Content: post, Route: route}
	htmlPage.Execute(w, data)
}

func RespondError(w http.ResponseWriter, code int, errorMessage string) {
	errorPost := Post{Title: "Error " + strconv.Itoa(code), Content: errorMessage}
	RenderHTML(w, "templates/error.html", "/error", errorPost)
}

func RespondJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
