package campaign

import (
	"mensago-api/internal/contract"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Save(campaign *Campaign) error {
	args := r.Mock.Called(campaign)
	return args.Error(0)
}

var (
	newCampaign = contract.NewCampaignDto{
		Name:    "Nova Campanha",
		Content: "Corpo da nova campanha",
		Emails:  []string{"teste@teste.com"},
	}
	service = Service{}
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(repositoryMock)
	service.Repository = repositoryMock
	id, err := service.CreateCampaign(newCampaign)
	assert.NotNil(id)
	assert.Nil(err)
}

func Test_Create_Save_Campaign(t *testing.T) {
	repositoryMock := new(repositoryMock)

	repositoryMock.On(
		"Save",
		mock.MatchedBy(func(campaign *Campaign) bool {
			if campaign.Name != newCampaign.Name {
				return false
			}

			if campaign.Content != newCampaign.Content {
				return false
			}

			if len(campaign.Contacts) != len(newCampaign.Emails) {
				return false
			}

			return true
		}),
	).Return(nil)
	service.Repository = repositoryMock
	service.CreateCampaign(newCampaign)
	repositoryMock.AssertExpectations(t)
}
