package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"GetPolicies",
		"GET",
		"/auth/manager/queryAll",
		getPolicyHandler,
	},
	Route{
		"EnforceParameters",
		"POST",
		"/auth/hasAccess/{sub}/{obj}/{act}",
		paramEnforce,
	},
	Route{
		"EnforceJSON",
		"POST",
		"/auth/hasAccess",
		postEnforce,
	},
	Route{
		"NewPolicy",
		"POST",
		"/auth/manager/new",
		postPolicyHandler,
	},
	Route{
		"RemovePolicy",
		"DELETE",
		"/auth/manager/remove",
		removePolicyHandler,
	},
	Route{
		"QueryPolicy",
		"POST",
		"/auth/manager/queryMatching",
		queryPolicyHandler,
	},
	Route{
		"QueryPolicy",
		"POST",
		"/auth/manager/querySingle",
		querySinglePolicyHandler,
	},
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
