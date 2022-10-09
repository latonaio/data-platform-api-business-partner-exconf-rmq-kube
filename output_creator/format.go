package output_creator

func ConvertToOutput(sessionId, partnerId string, exist bool) *Output {
	return &Output{
		BusinessPartner:  partnerId,
		ExistenceConf:    exist,
		RuntimeSessionID: sessionId,
	}
}
