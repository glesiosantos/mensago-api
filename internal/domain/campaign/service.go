package campaign

import (
	"mensago-api/internal/contract"
)

type Service struct {
	Repository Repository
}

func (s *Service) CreateCampaign(c contract.NewCampaignDto) (string, error) {
	campaign, _ := NewCampaign(c.Name, c.Content, c.Emails)
	return campaign.Id, nil
}
