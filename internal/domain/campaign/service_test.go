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
	args := r.Called(campaign)
	return args.Error(0)
}

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)

	service := Service{}

	newCampaign := contract.NewCampaignDto{
		Name:    "Nova Campanha",
		Content: "Corpo da nova campanha",
		Emails:  []string{"teste@teste.com"},
	}

	id, err := service.CreateCampaign(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)

}

func Test_Create_Save_Campaign(t *testing.T) {
	assert := assert.New(t)

	repository := new(repositoryMock)
	service := Service{Repository: repository}

	repository.On("Save").Return(nil)
	newCampaign := contract.NewCampaignDto{
		Name:    "Nova Campanha",
		Content: "Corpo da nova campanha",
		Emails:  []string{"teste@teste.com"},
	}

	id, err := service.CreateCampaign(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)

}
