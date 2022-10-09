package input_reader

type Input struct {
	ConnectionKey       string          `json:"connection_key"`
	Result              bool            `json:"result"`
	RedisKey            string          `json:"redis_key"`
	RuntimeSessionID    string          `json:"runtime_session_id"`
	Filepath            string          `json:"filepath"`
	BusinessPartner     BusinessPartner `json:"BusinessPartner"`
	APISchema           string          `json:"api_schema"`
	Accepter            []string        `json:"accepter"`
	BusinessPartnerCode string          `json:"business_partner_code"`
	Deleted             bool            `json:"deleted"`
}
type BusinessPartner struct {
	BusinessPartner  string `json:"BusinessPartner"`
	RuntimeSessionID string `json:"runtime_session_id"`
}
