package output_creator

func ConvertToOutput(businessPartner string, exist bool) *Output {
	return &Output{
		BusinessPartner: businessPartner,
		ExistenceConf:   exist,
	}
}
