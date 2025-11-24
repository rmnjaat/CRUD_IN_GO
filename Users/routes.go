package main

import (
	"github.com/go-chi/chi"
	"github.com/rmnjaat/users/common"
)

func AddOtherRoutes(router chi.Router) {
	common.RegisterRoutes(router)
}

func setUpAllRoutes(main_router chi.Router) {

	main_router.Route("/user", AddOtherRoutes)
}
