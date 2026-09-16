package campaign

import "mensago-api/internal/contract"

type Service struct {
	Repository Repository
}

func (s *Service) CreateCampaign(newCampaign contract.NewCampaignDto) error {
	return nil
}
