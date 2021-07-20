package main

import (
	"net/http"
	"vat/logger"
	"vat/service"
	"vat/servicecall"

	"github.com/afex/hystrix-go/hystrix"
)

const appName = "VAT-service"

func main() {

	// init the logger
	log := logger.NewLogger()
	logger.InitLogger(log)
	defer logger.Sync()

	// hystrix config
	hystrix.ConfigureCommand(servicecall.VatBreaker, hystrix.CommandConfig{
		Timeout: 10000,
	})

	//  create a new *router instance
	router := service.NewRouter()
	logger.Fatal(http.ListenAndServe(":7888", router))
}
