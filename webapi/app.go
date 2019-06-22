package webapi

import (
	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"net/http"
	"shopping-cart-kata/appservice"
	"shopping-cart-kata/cache"
)

// App is the web api application
type App struct {
	AppSvc    appservice.AppService
	HashGen   *hashids.HashID
	Router    *mux.Router
	CartCache cache.Cache
}

// Run runs the application
func (a *App) Run(listenAddr string) {
	http.ListenAndServe(listenAddr, a.Router)
}
