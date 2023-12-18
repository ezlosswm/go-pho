package domain 

type Contact struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Tel string `json:"tel"`
}

type Contacts struct {
	Contacts []Contact
}

func NewContact(name, tel string) *Contact {
	return &Contact{
		Name: name,
		Tel: tel,
	}
}
