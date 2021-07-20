package service

import (
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Method string            //HTTP method
	Path   string            //url endpoint
	Handle httprouter.Handle //Controller function which dispatches the right data for each route
}

type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/",
		Index,
	},
	Route{
		"GET",
		"/vat/:vat_num",
		VatNumer,
	},
}
