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
	UpdateContact(int, *domain.Contact) error
	Count() (int, error)
}

type service struct {
	db *sql.DB
}

func New() Service {
	db, err := sql.Open("sqlite3", "contact.db")
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
	defer rows.Close()

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
	defer q.Close()

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
	_, err := s.db.Exec("DELETE FROM contact WHERE id = $1", id)

	return err
}

func (s *service) UpdateContact(id int, c *domain.Contact) error {
	log.Println(c)
	query := `UPDATE contact SET name = $1, tel = $2 WHERE id = $3`
	updatedContact, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer updatedContact.Close()

	_, err = updatedContact.Exec(c.Name, c.Tel, c.ID)
	if err != nil {
		return err
	}

	return nil
}

// sql easy peasy lemon squeezy for miches
func (s *service) Count() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(tel) FROM contact").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
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
