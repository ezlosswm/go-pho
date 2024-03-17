package mock

import (
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
	db *domain.Contacts
}

func New() Service {
	db := &domain.Contacts{
		Contacts: []domain.Contact{
			{ID: 0, Name: "Carlos", Tel: "243-298-0084"},
			{ID: 1, Name: "Foo", Tel: "231-090-4075"},
			{ID: 2, Name: "Bar", Tel: "271-091-2827"},
		},
	}

	return &service{db: db}
}

func (s *service) GetContacts() (*domain.Contacts, error) {
	return s.db, nil
}

func (s *service) GetContact(id int) (*domain.Contact, error) {
	return &s.db.Contacts[id], nil
}

func (s *service) Store(c *domain.Contact) error {
	contact := &domain.Contact{
		ID:   len(s.db.Contacts),
		Name: c.Name,
		Tel:  c.Tel,
	}
	s.db.Contacts = append(s.db.Contacts, *contact)

	return nil
}

func (s *service) Delete(id int) error {
	s.db.Contacts = append(s.db.Contacts[:id], s.db.Contacts[id+1:]...)

	return nil
}

func (s *service) UpdateContact(id int, c *domain.Contact) error {
	return nil
}

// sql easy peasy lemon squeezy for miches
func (s *service) Count() (int, error) {
	return len(s.db.Contacts), nil
}
