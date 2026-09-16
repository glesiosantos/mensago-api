package campaign

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	Email string
}

type Campaign struct {
	Id        string
	Name      string
	Content   string
	Contacts  []Contact
	CreatedAt time.Time
}

func NewCampaign(name string, content string, emails []string) (*Campaign, error) {

	contacts := make([]Contact, len(emails))

	for index, email := range emails {
		contacts[index].Email = email
	}

	if name == "" {
		return nil, errors.New("name is required")
	} else if content == "" {
		return nil, errors.New("content is required")
	} else if len(emails) == 0 {
		return nil, errors.New("contacts is required")
	}

	return &Campaign{
		Id:        uuid.NewString(),
		Name:      name,
		Content:   content,
		Contacts:  contacts,
		CreatedAt: time.Now(),
	}, nil
}
