package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/gorilla/mux"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Homepage!") // write data to response
}
//user
func user(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("user.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of user login
		fmt.Println("id:", r.Form["id"])
		fmt.Println("name:", r.Form["name"])
		fmt.Println("email:", r.Form["email"])
		fmt.Println("password:", r.Form["password"])
	}
	
}
//userid
func Userid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Userid := vars["userId"]
	fmt.Fprintln(w, "User id:", Userid)
}
//posts
// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("posts.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}

}
//postid
func Postid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Postid := vars["postId"]
	fmt.Fprintln(w, "Post id:", Postid)
}

func main() {
	
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", sayhelloName) // setting router rule
	router.HandleFunc("/users", user)
	router.HandleFunc("/users/{userId}", Userid)
	router.HandleFunc("/posts", upload)
	router.HandleFunc("/posts/{postId}", Postid)
	
	log.Fatal(http.ListenAndServe(":8080", router))
}


