package server

import (
	"fmt"
	"log"
	"net/http"
	"pb/internal/domain"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.homePage)
	r.HandleFunc("/add", s.addContactPage)
	r.HandleFunc("/contacts/{id:[0-9]+}/edit", s.editPage)

	r.HandleFunc("/contacts", s.handleContacts)
	r.HandleFunc("/contacts/{id:[0-9]+}", s.handleContactsByID)

	return r
}

// pages
func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	s.templ.ExecuteTemplate(w, "base", nil)
}

func (s *Server) addContactPage(w http.ResponseWriter, r *http.Request) {
	s.templ.ExecuteTemplate(w, "contact-page", nil)
}

func (s *Server) editPage(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	contact, err := s.db.GetContact(id)

	if err != nil {
		log.Printf("No contact with id %v found", id)
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	s.templ.ExecuteTemplate(w, "edit", contact)
}

// handlers
func (s *Server) handleContacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.getContacts(w, r)
	case "POST":
		s.postContact(w, r)
	}
}

func (s *Server) handleContactsByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.getContact(w, r)
	case "DELETE":
		s.deleteContact(w, r)
	case "PUT":
		s.handleEdit(w, r)
	}
}

// getContacts serves the page
func (s *Server) getContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := s.db.GetContacts()
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	totalCount, err := s.db.Count()
	if err != nil {
		log.Println(err)
	}

	s.templ.ExecuteTemplate(w, "count", map[string]any{"Count": totalCount, "SwapOOB": true})
	s.templ.ExecuteTemplate(w, "list", contacts)
}

func (s *Server) getContact(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()


	contact, err := s.db.GetContact(id)
	if err != nil {
		log.Printf("No contact with id %v found", id)
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	s.templ.ExecuteTemplate(w, "person", contact)
}

func (s *Server) postContact(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	name := r.FormValue("name")
	tel := r.FormValue("tel")
	contact := domain.NewContact(name, tel)

	if err := s.db.Store(contact); err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) deleteContact(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		log.Fatal("err: unable to get id", id)
	}
	defer r.Body.Close()

	// TODO: work on validating ID provided
	// Example: id 30 does not exists
	// What do??????
	if err := s.db.Delete(id); err != nil {
		log.Printf("error deleteing contact with ID %d: %v", id, err)
	}

	totalCount, err := s.db.Count()
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusOK)
	s.templ.ExecuteTemplate(w, "count", map[string]any{"Count": totalCount, "SwapOOB": true})
}

// update
func (s *Server) handleEdit(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		log.Fatal("err: unable to get id", id)
	}


	r.ParseForm()
	q := r.Form["name"]
	log.Println(q)


	URL := fmt.Sprintf("/contact/%d", id)

	http.Redirect(w, r, URL, http.StatusOK)
}

func getID(r *http.Request) (int, error) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, fmt.Errorf("error: id %v unable to convert to int", idParam)
	}

	return id, nil
}
