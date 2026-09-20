package campaign

import (
	"mensago-api/internal/contract"
	"mensago-api/internal/domain/utils"
)

type Service struct {
	Repository Repository
}

func (s *Service) Create(newCampaign contract.NewCampaignDto) (string, error) {

	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	if err != nil {
		return "", err
	}
	err = s.Repository.Save(campaign)
	if err != nil {
		return "", utils.ServerError
	}

	return campaign.Id, nil
}
