package main

import (
	"context"
	"time"
	"io"
	"log"
	"net/http"
	"text/template"
	"google.golang.org/grpc"
	pb "github.com/ridwan779/grpc-tutorial/lib"
)

const (
	address     = "crud-server:55551"
)

type Employee struct {
	No int
    Id    string
    Name  string
    City string
}

type SubscriberRPC struct {
	conn      grpc.ClientConn
	err 	  error
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := pb.NewCRUDClient(conn)
	stream, err := c.List(ctx, &pb.Empty{})

	result := []Employee{}

	number := 1

	for {
		employee := Employee{}
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListData(_) = _, %v", c, err)
		}
		employee.No = number
		employee.Id = data.Id
		employee.Name = data.Name
		employee.City = data.City

		result = append(result, employee)

		number += 1
	}

    tmpl.ExecuteTemplate(w, "Index", result)
}

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
		stream, err := c.Insert(ctx, &pb.InsertRequest{Name: name, City: city})

		if err != nil {
			log.Fatalf("could not insert data: %v", err)
		}
		
		log.Printf("Status: %s", stream.Message)
	}

    http.Redirect(w, r, "/", 301)
}

func Show(w http.ResponseWriter, r *http.Request) {
	// Set up a connection to the server.
	nId := r.URL.Query().Get("id")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := pb.NewCRUDClient(conn)
	data, err := c.Show(ctx, &pb.DataId{Id: nId})

	employee := Employee{}

	employee.No = 1
	employee.Id = data.Id
	employee.Name = data.Name
	employee.City = data.City

	tmpl.ExecuteTemplate(w, "Show", employee)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	// Set up a connection to the server.
	nId := r.URL.Query().Get("id")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := pb.NewCRUDClient(conn)
	data, err := c.Show(ctx, &pb.DataId{Id: nId})

	employee := Employee{}

	employee.No = 1
	employee.Id = data.Id
	employee.Name = data.Name
	employee.City = data.City
	
	tmpl.ExecuteTemplate(w, "Edit", employee)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")
		// Set up a connection to the server.
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		c := pb.NewCRUDClient(conn)
		r, err := c.Update(ctx, &pb.InsertRequest{Id: id, Name: name, City: city})

		if err != nil {
			log.Fatalf("could not update data: %v", err)
		}
		
		log.Printf("Status: %s", r.Message)
	}

    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	// Set up a connection to the server.
	nId := r.URL.Query().Get("id")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := pb.NewCRUDClient(conn)
	data, err := c.Delete(ctx, &pb.DataId{Id: nId})

	if err != nil {
		log.Fatalf("could not delete data: %v", err)
	}
	
	log.Println(data.Message)

	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	
    http.HandleFunc("/", Index)
    http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    http.HandleFunc("/edit", Edit)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    http.ListenAndServe(":8080", nil)
}