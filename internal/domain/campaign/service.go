package campaign

import (
	"mensago-api/internal/contract"
)

type Service struct {
	Repository Repository
}

func (s *Service) CreateCampaign(c contract.NewCampaignDto) (string, error) {
	campaign, _ := NewCampaign(c.Name, c.Content, c.Emails)
	s.Repository.Save(campaign)
	return campaign.Id, nil
}
