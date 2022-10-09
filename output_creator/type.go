package output_creator

type Output struct {
	BusinessPartner  string `json:"BusinessPartner"`
	ExistenceConf    bool   `json:"ExistenceConf"`
	RuntimeSessionID string `json:"runtime_session_id"`
}
