package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"encoding/json"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website\n");
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request \n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func createTable(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /createTabele request \n")
	io.WriteString(w, "Creating Table!\n")
	make_table_with_schema()
}

func createShortenedUrl(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request_body shorten_request
	err := decoder.Decode(&request_body)
	if err != nil {
		panic(err)
	}

	fmt.Println(request_body.Original_url)

	shortened_url := shorten_url(request_body.Original_url)

	insert_row_into_table(request_body.Original_url, shortened_url)

	fmt.Fprint(w, shortened_url)
}

type shorten_request struct {
	Original_url string
}

func main() {

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/create-table", createTable)
	http.HandleFunc("/createShortenedUrl", createShortenedUrl)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed \n")
	} else if err != nil {
		fmt.Printf("error starting server %s\n", err)
		os.Exit(1)
	}

}