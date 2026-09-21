package campaign

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

var (
	fake = gofakeit.New(12345)

	name     = fake.Sentence(3)
	content  = fake.Paragraph(1, 2, 5, " ")
	contacts = []string{
		fake.Email(),
		fake.Email(),
	}
)

func TestNewCampaign_CreateCampaign(t *testing.T) {
	campaign, err := NewCampaign(name, content, contacts)

	require.NoError(t, err)
	require.NotNil(t, campaign)
	require.NotEmpty(t, campaign.Id)
	require.Equal(t, name, campaign.Name)
	require.Len(t, campaign.Contacts, len(contacts))
}

func TestNewCampaign_CreatedOnMustBeNow(t *testing.T) {
	now := time.Now().Add(-time.Minute)

	campaign, err := NewCampaign(name, content, contacts)

	require.NoError(t, err)
	require.NotNil(t, campaign)
	require.Greater(t, campaign.CreatedAt, now)
}

func TestNewCampaign_MustValidateName(t *testing.T) {
	campaign, err := NewCampaign("", content, contacts)

	require.Nil(t, campaign)
	require.EqualError(t, err, "name is required with min 5")
}

func TestNewCampaign_MustValidateContent(t *testing.T) {
	campaign, err := NewCampaign(name, "", contacts)

	require.Nil(t, campaign)
	require.EqualError(t, err, "content is required with min 5")
}

func TestNewCampaign_MustValidateContacts(t *testing.T) {
	campaign, err := NewCampaign(name, content, []string{})

	require.Nil(t, campaign)
	require.EqualError(t, err, "contacts is required with min 1")
}
