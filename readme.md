# Block-notes-app
A simple CRUD application in plain GO, without any external dependency such as 'gorilla/mux'.

## Functionalities
* The app will allow the user to create, read, update and delete posts.

* Posts can also be saved as drafts and can be viewed in the drafts section.

* Posts are paginated, they will be 5 per page.


## Future features
* It will show the last access time to a post.

## Why Golang?
Because it is easy to use but also a really powerful tool.

Also thanks to the 'html/template' Golang library rendering frontend content is much simpler and fast, allowing me to focus more on the backend side of the project.

## How to use it
### Clone the Project
```
git clone https://github.com/Gabri432/Block-notes-app.git
```
### Run the code
```
go run block_notes
```
### Run the tests
```
go test
```

## Routes
- `http://localhost:8081` where you will see the list of your posts.
- Go to `http://localhost:8081/new` or click "Create New Post" to create a new post.
- Go to `http://localhost:8081/drafts` or click "See Drafts" to see all of the posts you created but didn't finish.
- Go to `http://localhost:8081/modify/:id` or click on a post to modify it.
- Go to `http://localhost:8081/delete/:id` or click "Delete (this post)" to cancel a post.


## Project Structure
### Directories
- `block_notes` (main directory), which contains:
  - three '.go' files: 'block_notes.go', 'internal_functions.go' 'pagination.go';
  - their respective test files;
  - an executable 'block-notes-app.exe', the license and the readme files.

- `database`, which contains:
  - file 'posts.json' that will behave like a sort of database holding all of our posts.

- `templates`:
  - where all of the '.html' files are contained to render our frontend;
  - 'main.html' to serve all the existing posts and/or drafts;
  - 'error.html' where the user will be redirected in case of error;
  - 'form.html' where posts can be created, modified or deleted;
  - 'post.html' that describes how a single post should appear;
  - 'header.html' and 'footer.html' for rendering each page.

## Notes
- This project isn't perfect. This was just an opportunity for me to further explore the GO programming language.
- Because of that I didn't focus much of my attention to the appearence of the front-end, like responsitivity.
- There is of course room for improvements on both front-end and back-end sides.
- There are some tests but they do not cover most of the code.
- I hope this work could still be helpful or useful for somewhat porpuse.