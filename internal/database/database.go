package database

import (
	"database/sql"
	"fmt"
	"log"
	"pb/internal/domain"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type Service interface {
	Store(*domain.Contact) error
	GetContact(int) (*domain.Contact, error)
	GetContacts() (*domain.Contacts, error)
	Delete(int) error
}

type service struct {
	db *sql.DB
}

func New() Service {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	return &service{db: db}
} 

func (s *service) GetContacts() (*domain.Contacts, error) {
	rows, err := s.db.Query("SELECT * FROM contact")
	if err != nil {
		return nil, err
	}

	contacts := new(domain.Contacts)
	for rows.Next() {
		contact, err := scan(rows)
		if err != nil {
			return nil, err
		}

		contacts.Contacts = append(contacts.Contacts, *contact)
	}

	return contacts, nil
}

func (s *service) GetContact(id int) (*domain.Contact, error) {
	q, err := s.db.Query("SELECT * FROM contact WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for q.Next() {
		return scan(q)
	}

	return nil, fmt.Errorf("id %d not found", id) 
}

func (s *service) Store(c *domain.Contact) error {
	q := `INSERT INTO contact 
	(name, tel)
	values ($1, $2)
	`	

	_, err := s.db.Exec(
		q, 
		c.Name,
		c.Tel)
	if err != nil {
		return err
	}
		
	return nil
}

func (s *service) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM contact where id = $1", id)

	return err
}

func scan(rows *sql.Rows) (*domain.Contact, error) {
	contact := new(domain.Contact)
	err := rows.Scan(
		&contact.ID,
		&contact.Name,
		&contact.Tel,
	)
	if err != nil {
		return nil, err
	}

	return contact, nil
}
