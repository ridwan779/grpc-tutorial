package main

import (
	"context"
	"time"
	"log"
	"net/http"
	"text/template"
	"google.golang.org/grpc"
	pb "github.com/ridwan779/grpc-tutorial/lib"
)

const (
	address     = "localhost:55551"
)

var tmpl = template.Must(template.ParseGlob("form/*"))

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")

		// Set up a connection to the server.
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		c := pb.NewCRUDClient(conn)
		r, err := c.Insert(ctx, &pb.InsertRequest{Name: name, City: city})

		if err != nil {
			log.Fatalf("could not insert data: %v", err)
		}
		
		log.Printf("Status: %s", r.Message)
	}

    http.Redirect(w, r, "/new", 301)
}

func main() {

    log.Println("Server started on: http://localhost:8080")
    // http.HandleFunc("/", Index)
    // http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    // http.HandleFunc("/edit", Edit)
    http.HandleFunc("/insert", Insert)
    // http.HandleFunc("/update", Update)
    // http.HandleFunc("/delete", Delete)
    http.ListenAndServe(":8080", nil)
}