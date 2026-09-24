package campaign

type CampaignRepository struct {
	campaigns []Campaign
}

func (c *CampaignRepository) Create(campaign *Campaign) error {
	c.campaigns = append(c.campaigns, *campaign)
	return nil
}
