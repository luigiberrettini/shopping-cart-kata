package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"shopping-cart-kata/appservice"
)

func (a *App) createCart(w http.ResponseWriter, r *http.Request) {
	id, err := a.AppSvc.CreateCart()
	if err == appservice.ErrNotInitialized {
		respondWithError(w, http.StatusInternalServerError, "The system is not configured properly")
		return
	}
	if err == appservice.ErrCartCreation {
		respondWithError(w, http.StatusInternalServerError, "The system is not operating properly")
		return
	}
	wid, err := a.encode(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "The system encountered an unxepected condition")
		return
	}
	url, err := buildCartURL(wid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "The system encountered an unxepected condition")
		return
	}
	c := &cartVM{ID: wid, URL: url.String()}
	c.ComputeEtag()
	a.CartCache.AddOrReplace(wid, c)
	w.Header().Set("Location", c.URL)
	respondWithPayload(w, http.StatusCreated, *c, c.GetEtag())
}

func (a *App) addArticleToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wid := vars["id"]
	if im := r.Header.Get("If-Match"); len(im) != 0 {
		if _, ok := a.CartCache.GetByEtagWithID(im, wid); !ok {
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}
	}
	id, err := a.decode(wid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var article itemCreateVM
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&article); err != nil {
		respondWithError(w, http.StatusBadRequest, "Request payload cannot be decoded")
		return
	}
	err = a.AppSvc.AddArticleToCart(id, article.ID, article.Quantity)
	if err == appservice.ErrNotInitialized {
		respondWithError(w, http.StatusInternalServerError, "The system is not configured properly")
		return
	}
	if err == appservice.ErrCartNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err == appservice.ErrArtNotFound {
		respondWithError(w, http.StatusUnprocessableEntity, "The article does not exist")
		return
	}
	if err == appservice.ErrArtAlreadyAdded {
		respondWithError(w, http.StatusBadRequest, "Article already added")
		return
	}
	if err == appservice.ErrNonPositiveArtQty {
		respondWithError(w, http.StatusUnprocessableEntity, "Article quantity must be positive")
		return
	}
	url, err := buildCartURL(wid)
	article.CartURL = url.String()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "The system encountered an unxepected condition")
		return
	}
	a.CartCache.Remove(wid)
	respondWithPayload(w, http.StatusOK, article, "")
}

func (a *App) getCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wid := vars["id"]
	if _, ok := a.CartCache.GetByEtagWithID(r.Header.Get("If-None-Match"), wid); ok {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	id, err := a.decode(wid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	pc, err := a.AppSvc.GetCart(id)
	if err == appservice.ErrNotInitialized {
		respondWithError(w, http.StatusInternalServerError, "The system is not configured properly")
		return
	}
	if err == appservice.ErrCartNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err == appservice.ErrPromoRulesApplication {
		respondWithError(w, http.StatusInternalServerError, "The system encountered an unxepected condition")
		return
	}
	c := fromPricedCart(pc, wid, r.URL.String())
	cp := &c
	cp.ComputeEtag()
	a.CartCache.AddOrReplace(wid, cp)
	respondWithPayload(w, http.StatusOK, *cp, cp.GetEtag())
}

func (a *App) deleteCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wid := vars["id"]
	if im := r.Header.Get("If-Match"); len(im) != 0 {
		if _, ok := a.CartCache.GetByEtagWithID(im, wid); !ok {
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}
	}
	id, err := a.decode(wid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = a.AppSvc.DeleteCart(id)
	if err == appservice.ErrNotInitialized {
		respondWithError(w, http.StatusInternalServerError, "The system is not configured properly")
		return
	}
	a.CartCache.Remove(wid)
	w.WriteHeader(http.StatusNoContent)
}

func (a *App) getArticles(w http.ResponseWriter, r *http.Request) {
	arts := a.AppSvc.Catalog.GetArticles()
	respondWithPayload(w, http.StatusOK, arts, "")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithPayload(w, code, map[string]string{"error": message}, "")
}

func respondWithPayload(w http.ResponseWriter, statusCode int, viewmodel interface{}, etag string) {
	response, _ := json.Marshal(viewmodel)
	if etag != "" {
		w.Header().Set("ETag", etag)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
