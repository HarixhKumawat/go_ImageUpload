# go_ImageUpload

This is a Go server that allows users to upload images. The server is built using best practices to ensure file size is limited to under 5MB and only images are allowed to be uploaded.


# Usage

To use the server simlu open the port 7000 in browser, the index route has all the example usages mentioned

/uploadImgage route will return an HTML page with a form that allows the user to upload an image. When the user submits the form (via a POST request), the server will process the uploaded file and return a success message with a unique image ID.


## ROUTES
- [/] index.htm 
- [/name] return name with age from url query
- [/uploadImage] the image upload form
- [/viewImage] returns the requested image(url query passed imageId)

# Installation

* To install the server, follow these steps:


1. Clone the repository to your local machine.
2. Install Go if you haven't already.
3. Run go run main.go to start the server. or serverImg if you are on unix operating system

# Access the server via http://localhost:7000/.

