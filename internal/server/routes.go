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
	r.HandleFunc("/contacts/{id:[0-9]+}/display", s.handleDisplay)

	r.HandleFunc("/contacts", s.handleContacts)
	r.HandleFunc("/contacts/{id:[0-9]+}", s.handleContactsByID)

	return r
}

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

// pages
// homePage executes the home page template
func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	s.templ.ExecuteTemplate(w, "base", nil)
}

// addContactPage executes the contact page template
func (s *Server) addContactPage(w http.ResponseWriter, r *http.Request) {
	s.templ.ExecuteTemplate(w, "contact-page", nil)
}

// getContact retrieves the requested contact and the page
func (s *Server) getContact(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	contact, err := s.db.GetContact(id)
	if err != nil {
		log.Printf("No contact with id %v found", id)
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	s.templ.ExecuteTemplate(w, "person", contact)
}

// editPage returns a partial components to edit the form
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

// actions
// postContact makes a post request with the form data
func (s *Server) postContact(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	tel := r.FormValue("tel")
	contact := domain.NewContact(name, tel)

	defer r.Body.Close()

	if err := s.db.Store(contact); err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// deleteContact makes a delete request based on the specified ID
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

// getContacts returns the list of contacts
func (s *Server) getContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := s.db.GetContacts()
	if err != nil {
		log.Fatal(err)
	}

	totalCount, err := s.db.Count()
	if err != nil {
		log.Println(err)
	}

	s.templ.ExecuteTemplate(w, "count", map[string]any{"Count": totalCount, "SwapOOB": true})
	s.templ.ExecuteTemplate(w, "list", contacts)
}

// partials
// handleEdit make the put request to update the speciifed contact
func (s *Server) handleEdit(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		log.Fatal("err: unable to get id", id)
	}

	contact, err := s.db.GetContact(id)
	if err != nil {
		log.Printf("No contact with id %v found", id)
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	contact.Name = r.FormValue("name")
	contact.Tel = r.FormValue("tel")
	defer r.Body.Close()

	if err := s.db.UpdateContact(id, contact); err != nil {
		log.Println(err)
	}

	s.templ.ExecuteTemplate(w, "display", contact)
}

// handleDisplay manages the cancel functionility in the edit page
func (s *Server) handleDisplay(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	contact, err := s.db.GetContact(id)
	if err != nil {
		log.Printf("No contact with id %v found", id)
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	s.templ.ExecuteTemplate(w, "display", contact)
}

// getID retrieves the id from the URL
func getID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, fmt.Errorf("error: id %v unable to convert to int", idParam)
	}

	return id, nil
}
