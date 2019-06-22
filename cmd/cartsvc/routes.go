package main

import (
	"net/url"
	"strings"
)

var buildCartURL func(wid string, reqURL *url.URL) (string, error)

// ConfigRoutes configures the API routes
func (a *App) ConfigRoutes() {
	a.Router.HandleFunc("/carts", a.createCart).Methods("POST")
	a.Router.HandleFunc("/carts/{id}", a.getCart).Methods("GET").Name("cart")
	a.Router.HandleFunc("/carts/{id}", a.deleteCart).Methods("DELETE")
	a.Router.HandleFunc("/carts/{id}/items", a.addArticleToCart).Methods("POST")
}

// ConfigURLBuilders setup URL builders
func (a *App) ConfigURLBuilders() {
	buildCartURL = func(wid string, reqURL *url.URL) (string, error) {
		u, err := a.Router.Get("cart").URL("id", wid)
		if err != nil {
			return "", err
		}
		return strings.Replace(reqURL.String(), reqURL.RequestURI(), u.String(), 1), nil
	}
}
