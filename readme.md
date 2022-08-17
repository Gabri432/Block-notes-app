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
- Run "go run ." to open the application on http://localhost:8081, where you will see the list of your posts.
- Go to http://localhost:8081/new or click "Create New Post" to create a new post.
- Go to http://localhost:8081/drafts or click "See Drafts" to see all of the posts you created but didn't finish.
- Go to http://localhost:8081/modify/:id or click on a post to modify it.
- Go to http://localhost:8081/delete/:id or click "Delete (this post)" to cancel a post.


## Project Structure
- A folder named "database" where the file 'posts.json' will hold all of the posts.
- A folder named "templates" where all of the .html files are contained.
- The main folder will contain two .go files: 'block_notes.go' and 'internale_functions.go'.
- 'block_notes.go' contains the api-rest handlers.
- 'internal_functions.go' contains other functions used to simplifying the code writing.

## Notes
- This project isn't perfect. This was just an opportunity for me to further explore the GO programming language.
- Because of that I didn't focus much of my attention to the appearence of the front-end, like responsitivity.
- There is of course room for improvements on both front-end and back-end sides.
- I hope this work could still be helpful or useful for somewhat porpuse.