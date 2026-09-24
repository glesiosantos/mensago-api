package campaign

import (
	"errors"
	"testing"

	"mensago-api/internal/contract"
	"mensago-api/internal/utils"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Save(campaign *Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

var newCampaign = contract.NewCampaignDto{
	Name:    "Test Y",
	Content: "Body valid", // mínimo de 5 caracteres
	Emails:  []string{"teste1@test.com"},
}

func Test_Create_Campaign(t *testing.T) {
	repository := new(repositoryMock)
	service := Service{Repository: repository}

	repository.
		On("Save", mock.AnythingOfType("*campaign.Campaign")).
		Return(nil).
		Once()

	id, err := service.Create(newCampaign)

	require.NoError(t, err)
	require.NotEmpty(t, id)
	repository.AssertExpectations(t)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	repository := new(repositoryMock)
	service := Service{Repository: repository}

	invalidCampaign := newCampaign
	invalidCampaign.Name = ""

	id, err := service.Create(invalidCampaign)

	require.Error(t, err)
	require.Empty(t, id)
	require.EqualError(t, err, "name is required with min 5")

	// O repositório não deve ser chamado quando o domínio é inválido.
	repository.AssertNotCalled(t, "Save", mock.Anything)
}

func Test_Create_SaveCampaign(t *testing.T) {
	repository := new(repositoryMock)
	service := Service{Repository: repository}

	repository.
		On("Save", mock.MatchedBy(func(campaign *Campaign) bool {
			return campaign != nil &&
				campaign.Name == newCampaign.Name &&
				campaign.Content == newCampaign.Content &&
				len(campaign.Contacts) == len(newCampaign.Emails)
		})).
		Return(nil).
		Once()

	id, err := service.Create(newCampaign)

	require.NoError(t, err)
	require.NotEmpty(t, id)
	repository.AssertExpectations(t)
}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	repository := new(repositoryMock)
	service := Service{Repository: repository}

	repository.
		On("Save", mock.AnythingOfType("*campaign.Campaign")).
		Return(errors.New("error to save on database")).
		Once()

	id, err := service.Create(newCampaign)

	require.Error(t, err)
	require.Empty(t, id)
	require.ErrorIs(t, err, utils.ServerError)
	repository.AssertExpectations(t)
}
