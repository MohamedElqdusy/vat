package models

import "encoding/xml"

// VatValidation represents the validation response
type VatValidation struct {
	IsValid bool `json:"valid"`
}

// VatSoapResponse represents the vies saop response
type VatSoapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Soap    struct {
		XMLName xml.Name `xml:"Body"`
		Soap    struct {
			XMLName xml.Name `xml:"checkVatResponse"`
			Valid   bool     `xml:"valid"`
		}
		SoapFault struct {
			XMLName string `xml:"Fault"`
			Code    string `xml:"faultcode"`
			Message string `xml:"faultstring"`
		}
	}
}
