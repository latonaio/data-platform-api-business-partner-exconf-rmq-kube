package dpfm_api_input_reader

import (
	"data-platform-api-business-partner-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *GeneralSDC) ConvertToBusinessPartnerGeneral() *requests.BusinessPartnerGeneral {
	data := sdc.BusinessPartnerGeneral
	return &requests.BusinessPartnerGeneral{
		BusinessPartner: data.BusinessPartner,
	}
}
