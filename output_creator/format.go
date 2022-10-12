package output_creator

func ConvertToOutput(partnerId string, exist bool) *Output {
	return &Output{
		BusinessPartner: partnerId,
		ExistenceConf:   exist,
	}
}
