package campaign

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	Email string `validate:"email"`
}

type Campaign struct {
	Id        string    `validate:"required"`
	Name      string    `validate:"min=5,max=100"`
	Content   string    `validate:"min=5,max=1024"`
	Contacts  []Contact `validate:"min=1,dive"`
	CreatedAt time.Time `validate:"required"`
}

func NewCampaign(name string, content string, emails []string) (*Campaign, error) {

	contacts := make([]Contact, len(emails))

	for index, email := range emails {
		contacts[index].Email = email
	}

	return &Campaign{
		Id:        uuid.NewString(),
		Name:      name,
		Content:   content,
		Contacts:  contacts,
		CreatedAt: time.Now(),
	}, nil
}
