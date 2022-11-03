package input_reader

type Input struct {
	BusinessPartner     BusinessPartner `json:"BusinessPartner"`
	APISchema           string          `json:"api_schema"`
	Accepter            []string        `json:"accepter"`
	BusinessPartnerCode string          `json:"business_partner_code"`
	Deleted             bool            `json:"deleted"`
}
type BusinessPartner struct {
	BusinessPartner int `json:"BusinessPartner"`
}
