package campaign

import "time"

type Contacts struct {
	Email string
}

type Campaign struct {
	Id        string
	Name      string
	Content   string
	Contacts  []Contacts
	CreatedAt time.Time
}
