package cserver

import (
	"container/list"
	"fmt"
	"html/template"
	"log"
	"net/http"
	manager "simpleSurance/FileManager"
	window "simpleSurance/MovingWindow"
	"time"
)

//Page Content, it only contains int, umber of requests
type Page struct {
	Body  int
}

//
//Server which holds window and manager objects
//it is responsible for handling requests
//
type Server struct {
	window *window.MovingWindow
	fileManager *manager.FileManager

	Addr string
	Pattern string
}

var templates = template.Must(template.ParseFiles("Pages/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//requestHandler handles received requests
func (s *Server )requestHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().UnixNano()
	n := s.window.Request(t)
	p := &Page{Body: n}
	renderTemplate(w, "view", p)

	err := s.fileManager.UpdateFiles(t)
	if err!= nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Initialize Server, creates necessary objects
func (s *Server)initialization() error{
	s.window = &window.MovingWindow{List: list.List{}}
	s.fileManager = &manager.FileManager{LastChange1: 0, LastChange2:0,
		FileName1:"DB/db1.txt", FileName2:"DB/db2.txt"}
	err := s.fileManager.Init()
	return err
}

func (s * Server) backup() error{
	err := s.fileManager.Backup(s.window)
	return err
}

func (server * Server) Start() error{
	err := server.initialization()
	if err != nil{
		fmt.Println("Server init: %s", err)
		return err
	}
	err = server.backup()
	if err != nil{
		fmt.Println("Server backup: %s", err)
		return err
	}
	http.HandleFunc(server.Pattern, server.requestHandler)
	log.Fatal(http.ListenAndServe(server.Addr, nil))
	return nil
}