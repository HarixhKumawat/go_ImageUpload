package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"
)

const MAX_IMAGE_SIZE = 1024 * 1024 * 5 //5MB

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	//i decide what the hell happends on the server
	// if post request then i show the html file if not then riun the upload thingy
	if r.Method != "POST" {
		outputHTML(w, "static/uploadImage.htm", "")
		return
	}

	// hackable os not use
	// if r.ContentLength > MAX_IMAGE_SIZE {
	// 	http.Error(w, "Upload file size EXCEEDED! expected <5MB", http.StatusBadRequest)
	// }

	//the good method
	r.Body = http.MaxBytesReader(w, r.Body, MAX_IMAGE_SIZE)
	if err := r.ParseMultipartForm(MAX_IMAGE_SIZE); err != nil {
		http.Error(w, "Upload file size EXCEEDED! expected <5MB", http.StatusBadRequest)
		return
	}

	//take the file with header from the file feild of html fiel
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//this is really smart as this says close the file once everuything unmderneath it has finsihded executing
	defer file.Close()

	// //this code only checks the file type is correct only
	// buff := make([]byte, 512)
	// _, err = file.Read(buff)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// filetype := http.DetectContentType(buff)
	// if filetype != "image/jpeg" && filetype != "image/png" { {
	// 	http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
	// 	return
	// }

	// _, err := file.Seek(0, io.SeekStart)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// //done checking the filetype

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//we create the file
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(fileHeader.Filename)
	dst, err := os.Create(fmt.Sprintf("./uploads/%s", fileName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	//copy the  file into the destination object we created
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful\ngo to [http://localhost:7000/viewImage?imageId=%v%s", fileName, "]")
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.URL.Query().Get("imageId")
	if imageId == "" {
		http.Error(w, "Image id not found", http.StatusNotFound)
		return
	}

	passImageId := map[string]interface{}{"imageId": imageId}
	outputHTML(w, "static/callImage.htm", passImageId)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userAgent := map[string]interface{}{"userAgent": r.UserAgent()}
		outputHTML(w, "static/index.htm", userAgent)
	})

	fs := http.FileServer(http.Dir("uploads/"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	http.HandleFunc("/name", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		age := r.URL.Query().Get("age")

		fmt.Fprint(w, "hellow, this is your first request mr ", name, " with age: ", age, "\nand the request object is->? ", r)
	})
	http.HandleFunc("/uploadImage", uploadHandler)
	http.HandleFunc("/viewImage", ImageHandler)

	log.Println("listening server on port http://localhost:7000")
	http.ListenAndServe(":7000", nil)

}
