package dpfm_api_input_reader

import (
	"data-platform-api-business-partner-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *SDC) ConvertToBusinessPartner() *requests.BusinessPartner {
	data := sdc.BusinessPartnerID
	return &requests.BusinessPartner{
		BusinessPartnerID: data.BusinessPartnerID,
	}
}
