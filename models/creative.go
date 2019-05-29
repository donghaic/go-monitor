package models

/**
 * 广告创意，GT前端变更后同步到Redis中对应的JSON数据
 */
type CreativeInfo struct {
	AdvId      string `json:"adv_id"`
	ProductId  string `json:"product_id"`
	CampaignId string `json:"campaign_id"`
	GroupId    string `json:"group_id"`
	MediaId    string `json:"media_id"`
	CreativeId string `json:"creative_id"`

	LandingPage  string `json:"landingpage"`
	UserId       string `json:"user_id"`
	UserBudgetId string `json:"user_budget_id"`
	AdvImpLink   string `json:"adv_imp_link"`
	AdvClkLink   string `json:"adv_clk_link"`
	AdvLinkId    string `json:"adv_link_id"`
	AdvLinkTag   string `json:"adv_link_tag"`

	WaxExclusionCodes []string `json:"action_codes"`
}
