package servicecall

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"vat/logger"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	setup()

	var cases = []struct {
		Name        string
		CountryName string
		VatNum      string
		Expected    bool
	}{
		{
			Name:        "home24 valid vat",
			CountryName: "DE",
			VatNum:      "266182271",
			Expected:    true,
		},
		{
			Name:        "In-valid vat",
			CountryName: "DE",
			VatNum:      "0266182271",
			Expected:    false,
		},
	}

	for _, input := range cases {
		actual, err := Validate(input.CountryName, input.VatNum)
		assert.NoError(t, err)
		if input.Expected != actual {
			t.Errorf("Test failed:  %s \n Want to be: %v but we got %v",
				input.Name,
				input.Expected,
				actual,
			)
		}
	}
	assert.Equal(t, 5, 6)

}

func TestUnmarshalVatReponse(t *testing.T) {
	setup()
	var cases = []struct {
		Name         string
		ResponseBody string
		Expected     bool
	}{
		{
			Name:         "home24 valid vat",
			ResponseBody: `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><checkVatResponse xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types"><countryCode>DE</countryCode><vatNumber>266182271</vatNumber><requestDate>2021-07-20+02:00</requestDate><valid>true</valid><name>---</name><address>---</address></checkVatResponse></soap:Body></soap:Envelope>`,
			Expected:     true,
		},
		{
			Name:         "In-valid vat",
			ResponseBody: `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><checkVatResponse xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types"><countryCode>DE</countryCode><vatNumber>0266182271</vatNumber><requestDate>2021-07-20+02:00</requestDate><valid>false</valid><name>---</name><address>---</address></checkVatResponse></soap:Body></soap:Envelope>`,
			Expected:     false,
		},
	}

	for _, input := range cases {
		res := http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(input.ResponseBody))}
		actual, err := unmarshalVatReponse(&res)
		assert.NoError(t, err)
		if input.Expected != actual {
			t.Errorf("Test failed:  %s \n Want to be: %v but we got %v",
				input.Name,
				input.Expected,
				actual,
			)
		}
	}
}

func setup() {
	log := logger.NewLogger()
	logger.InitLogger(log)
	defer logger.Sync()

	// hystrix config
	hystrix.ConfigureCommand(VatBreaker, hystrix.CommandConfig{
		Timeout: 10000,
	})
}
