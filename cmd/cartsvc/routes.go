package main

import (
	"net/url"
)

var buildCartURL func(wid string) (*url.URL, error)

// ConfigRoutes configures the API routes
func (a *App) ConfigRoutes(authority string) {
	a.Router.HandleFunc("/carts", a.createCart).Host(authority).Methods("POST")
	a.Router.HandleFunc("/carts/{id}", a.getCart).Host(authority).Methods("GET").Name("cart")
	a.Router.HandleFunc("/carts/{id}", a.deleteCart).Host(authority).Methods("DELETE")
	a.Router.HandleFunc("/carts/{id}/items", a.addArticleToCart).Host(authority).Methods("POST")
}

// ConfigURLBuilders setup URL builders
func (a *App) ConfigURLBuilders() {
	buildCartURL = func(wid string) (*url.URL, error) {
		return a.Router.Get("cart").URL("id", wid)
	}
}
