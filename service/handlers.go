package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"vat/logger"
	"vat/models"
	"vat/servicecall"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, welcome to the VAT validation service")
}

func VatNumer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vatNum := ps.ByName("vat_num")

	// validate the vat number input and check for Germany only
	countryName, vatNum, err := checkVatNum(vatNum)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 400)
		return
	}

	// validate by VIES soap
	isValid, err := servicecall.Validate(countryName, vatNum)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// encode the end result
	if err = json.NewEncoder(w).Encode(models.VatValidation{isValid}); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 500)
	}
}

// validate the vat number initaliy before validating by VIES soap
func checkVatNum(vatNum string) (string, string, error) {
	if len(vatNum) < 3 {
		return "", "", errors.New("short VAT Num")
	}
	countryName := vatNum[0:2]
	vatNumber := vatNum[2:]

	// Valid for Germany only
	if countryName != "DE" {
		return "", "", errors.New("we only support DE")
	}

	return countryName, vatNumber, nil
}
