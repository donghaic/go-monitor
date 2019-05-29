package models

//
// 广告检测地址
type LinkSwitch struct {
	PreLinkId  string `json:"preLinkId"`
	PreLinkImp string `json:"preLinkImp"`
	PreLinkClk string `json:"preLinkClk"`

	AfLinkId  string `json:"afLinkId"`
	AfLinkImp string `json:"afLinkImp"`
	AfLinkClk string `json:"afLinkClk"`

	Scale int `json:"scale"`
}
