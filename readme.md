# Block-notes-app
A simple CRUD application in GO.

## Functionalities
* The app will allow the user to create, read, update and delete posts.

* The user will also be capable of saving a post to finish it later.


## Future features
* It will show the last access time to a post.

* The app will also paginate the posts, showing just 20 post per time.

## Why Golang?
Because it is easy to use but also a really powerful tool.

Also thanks to the 'html/template' Golang library rendering frontend content is much simpler and fast, allowing me to focus more on the backend side of the project.

## How to use it
- Use "git clone https://github.com/Gabri432/Block-notes-app.git" to clone the project.
- Run "go run block_notes.go" to open the application on http://localhost:8081, where you will see the list of your posts.
- Go to http://localhost:8081/new or click "Create New Post" to create a new post.
- Go to http://localhost:8081/drafts or click "Unfinished Posts" to see all of the posts you created but didn't finish.
- Go to http://localhost:8081/modify/:id or click "Modify (this post)" to modify a post.
- Go to http://localhost:8081/delete/:id or click "Delete (this post)" to cancel a post.
