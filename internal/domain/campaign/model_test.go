package campaign

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	name     = "Nova campanha X fds"
	content  = "Corpo da Campanha para teste"
	contacts = []string{"email1@test.com.br", "email2@test.com"}
)

func Test_NewCampaign_CreateCampanign(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(name, content, contacts)
	assert.NotEmpty(t, campaign.Id)
	assert.Equal(campaign.Name, name)
	assert.Equal(len(campaign.Contacts), len(contacts))
}

func Test_NewCampaign_CreatedOnMustBeNow(t *testing.T) {
	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)
	campaign, _ := NewCampaign(name, content, contacts)

	assert.Greater(campaign.CreatedAt, now)
}

func Test_NewCampaign_MustValidateName(t *testing.T) {
	assert := assert.New(t)
	_, err := NewCampaign("", content, contacts)

	assert.Equal("name is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateContent(t *testing.T) {
	assert := assert.New(t)
	_, err := NewCampaign(name, "", contacts)
	assert.Equal("content is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateContact(t *testing.T) {
	assert := assert.New(t)
	_, err := NewCampaign(name, content, []string{})
	assert.Equal("contacts is required with min 1", err.Error())
}
