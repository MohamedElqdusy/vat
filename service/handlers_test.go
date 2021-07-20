package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vat/logger"
	"vat/models"
	"vat/servicecall"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestValidVatNumer(t *testing.T) {
	setup()
	req, err := http.NewRequest("GET", "/vat/DE266182271", nil)

	assert.NoError(t, err)
	recorder := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle("GET", "/vat/:vat_num", VatNumer)
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := models.VatValidation{true}
	actual := models.VatValidation{}

	err = json.NewDecoder(recorder.Body).Decode(&actual)
	assert.NoError(t, err)

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestVatNumerErrors(t *testing.T) {
	setup()

	var cases = []struct {
		Name          string
		URL           string
		ExpectedCode  int
		ExpectedError string
	}{
		{
			Name:          "Germany validation only",
			URL:           "/vat/AT266182271",
			ExpectedCode:  http.StatusBadRequest,
			ExpectedError: "we only support DE",
		},
		{
			Name:          "too short vat",
			URL:           "/vat/01",
			ExpectedCode:  http.StatusBadRequest,
			ExpectedError: "short VAT Num",
		},
	}

	for _, c := range cases {
		req, err := http.NewRequest("GET", c.URL, nil)

		assert.NoError(t, err)
		recorder := httptest.NewRecorder()

		router := httprouter.New()
		router.Handle("GET", "/vat/:vat_num", VatNumer)
		router.ServeHTTP(recorder, req)

		if status := recorder.Code; status != c.ExpectedCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, c.ExpectedCode)
		}

		res, err := ioutil.ReadAll(recorder.Body)
		assert.NoError(t, err)
		actual := string(res)
		if strings.TrimSpace(actual) != c.ExpectedError {
			t.Errorf("handler returned unexpected body: got %v want %v", actual, c.ExpectedError)
		}
	}

}

func setup() {
	log := logger.NewLogger()
	logger.InitLogger(log)
	defer logger.Sync()

	// hystrix config
	hystrix.ConfigureCommand(servicecall.VatBreaker, hystrix.CommandConfig{
		Timeout: 10000,
	})
}
