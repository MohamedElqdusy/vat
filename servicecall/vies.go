package servicecall

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"vat/logger"
	"vat/models"

	"vat/reselience"
)

const (
	// VatBreaker is the name for the VIES call breaker
	VatBreaker  = "vat_validation"
	viesBaseURL = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"
)

// Validate used to validate the vat number by calling the vies soap-service
func Validate(countryCode, vatNum string) (bool, error) {
	envelope := fmt.Sprintf(`
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v1="http://schemas.conversesolutions.com/xsd/dmticta/v1">
	<soapenv:Header/>
	<soapenv:Body>
	  <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
		<countryCode>%s</countryCode>
		<vatNumber>%s</vatNumber>
	  </checkVat>
	</soapenv:Body>
	</soapenv:Envelope>`, countryCode, vatNum)

	logger.Infof("Request to URL: %s with vat number %s%s", viesBaseURL, countryCode, vatNum)
	// construct the soap request
	req, err := http.NewRequest("POST", viesBaseURL, bytes.NewBufferString(envelope))
	if err != nil {
		logger.Error(err)
		return false, err
	}

	// do the soap request
	resp, err := reselience.CircuitBreaker(VatBreaker, req)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	return unmarshalVatReponse(resp)
}

func unmarshalVatReponse(res *http.Response) (bool, error) {
	// read the response
	xmlRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// checks for the invalid input error
	if bytes.Contains(xmlRes, []byte("INVALID_INPUT")) {
		return false, errors.New("INVALID VAT number")
	}

	// unmarshal the response
	var v models.VatSoapResponse
	if err := xml.Unmarshal(xmlRes, &v); err != nil {
		return false, err
	}

	// handle soap fault message if any
	if len(v.Soap.SoapFault.Message) != 0 {
		return false, errors.New(v.Soap.SoapFault.Message)
	}

	// return the validation if no soap fault
	return v.Soap.Soap.Valid, nil
}
