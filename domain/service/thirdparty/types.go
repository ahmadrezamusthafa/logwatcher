package thirdparty

import (
	"github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"
)

type QueryInput struct {
	MessageQuery string
	ContextQuery string
	Limit        int
}

type QueryOutput struct {
	Timestamp *string `json:"timestamp"`
	Hostname  *string `json:"hostname"`
	FlowID    *string `json:"flowid"`
	Type      *string `json:"type"`
	Part      *string `json:"part"`
	Message   *string `json:"message"`
	Context   *string `json:"context"`
}

type ASIPCNTContext struct {
	CorrelationID   *string `json:"correlationid"`
	Event           *string `json:"event"`
	Uri             *string `json:"uri"`
	ProviderID      *string `json:"providerid"`
	ProviderHotelID *string `json:"providerhotelid"`
	Locale          *string `json:"locale"`
	ProviderBrandID *string `json:"providerbrandid"`
	ProviderChainID *string `json:"providerchainid"`
}

type ASIPSRCContext struct {
	CorrelationID *string `json:"correlationid"`
	ProviderID    *string `json:"providerid"`
	Event         *string `json:"event"`
	SourceMarket  *string `json:"sourcemarket"`
	CheckinDate   *string `json:"checkindate"`
	CheckoutDate  *string `json:"checkoutdate"`
	Locale        *string `json:"locale"`
	Currency      *string `json:"currency"`
	NoOfAdult     *int    `json:"noofadult"`
	NoOfChild     *int    `json:"noofchild"`
	NoOfRoom      *int    `json:"noofroom"`
}

type ASIPSRCDetailContext struct {
	CorrelationID *string `json:"correlationid"`
}

type LogAttribute struct {
	Attribute string              `json:"attribute"`
	Name      string              `json:"name"`
	Type      valuetype.ValueType `json:"type"`
}
