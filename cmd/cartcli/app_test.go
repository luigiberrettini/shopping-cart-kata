package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIErrors(t *testing.T) {
	const noMsg = ""
	someMsg := "Some message"
	apiErrors(t, http.StatusInternalServerError, someMsg)
	apiErrors(t, http.StatusBadRequest, someMsg)
	apiErrors(t, http.StatusUnprocessableEntity, someMsg)
	apiErrors(t, http.StatusConflict, someMsg)
	apiErrors(t, http.StatusNotFound, noMsg)
	apiErrors(t, http.StatusPreconditionFailed, noMsg)
}

func TestCreateCartSuccess(t *testing.T) {
	etag := `W/"123456789"`
	sCod := http.StatusCreated
	crt := cart{ID: "ABC", URL: "http://127.0.0.1/carts/ABC", ETag: etag}
	hf := func(w http.ResponseWriter, r *http.Request) {
		respondWithPayload(w, sCod, crt, etag)
	}
	ts := httptest.NewServer(http.HandlerFunc(hf))
	defer ts.Close()
	a := &App{BaseURL: ts.URL, HTTPClient: *ts.Client()}
	c, sc, m, err := a.createCart()
	if c.String() != crt.String() || sc != sCod || err != nil {
		t.Errorf("EXPECTED\nCart: %s\nStatus code: %d\nMessage: %s\n\n", crt, sCod, "")
		t.Errorf("GOT\nCart: %s\nStatus code: %d\nMessage: %s\nError: %s", c, sc, m, err)
	}
}

func TestAddArtToCartSuccessWithoutIfMatchEtag(t *testing.T) {
	addArticleToCartSuccess(t, "")
}

func TestAddArtToCartSuccessWithIfMatchEtag(t *testing.T) {
	addArticleToCartSuccess(t, `W/"123456789"`)
}

func TestGetCartSuccessWithoutIfNoneMatchEtag(t *testing.T) {
	getCartSuccess(t, "")
}

func TestGetCartSuccessWithIfNoneMatchEtagDiffFromRespEtag(t *testing.T) {
	getCartSuccess(t, `W/"XXXXXXXXX"`)
}

func TestGetCartSuccessWithIfNoneMatchEtagEqualToRespEtag(t *testing.T) {
	getCartSuccess(t, `W/"123456789"`)
}

func TestDeleteCartSuccessWithoutIfMatchEtag(t *testing.T) {
	deleteCartSuccess(t, "")
}

func TestDeleteCartSuccessWithIfMatchEtag(t *testing.T) {
	deleteCartSuccess(t, `W/"123456789"`)
}

func apiErrors(t *testing.T, sCod int, msg string) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		if msg != "" {
			respondWithError(w, sCod, msg)
		} else {
			w.WriteHeader(sCod)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(hf))
	defer ts.Close()
	a := &App{BaseURL: ts.URL, HTTPClient: *ts.Client()}

	c, sc, m, err := a.createCart()
	if sc != sCod || m != msg || err != nil {
		t.Errorf("EXPECTED\nStatus code: %d\nMessage: %s\n\n", sCod, msg)
		t.Errorf("GOT\nCart: %s\nStatus code: %d\nMessage: %s\nError: %s", c, sc, m, err)
	}
	art, sc, m, err := a.addArticleToCart("A", "B", "C", 1)
	if sc != sCod || m != msg || err != nil {
		t.Errorf("EXPECTED\nStatus code: %d\nMessage: %s\n\n", sCod, msg)
		t.Errorf("GOT\nArticle: %s\nStatus code: %d\nMessage: %s\nError: %s", art, sc, m, err)
	}
	crt, sc, m, err := a.getCart("A", "B")
	if sc != sCod || m != msg || err != nil {
		t.Errorf("EXPECTED\nStatus code: %d\nMessage: %s\n\n", sCod, msg)
		t.Errorf("GOT\nArticle: %s\nStatus code: %d\nMessage: %s\nError: %s", crt, sc, m, err)
	}
	sc, m, err = a.deleteCart("A", "B")
	if sc != sCod || m != msg || err != nil {
		t.Errorf("EXPECTED\nStatus code: %d\nMessage: %s\n\n", sCod, msg)
		t.Errorf("GOT\nArticle: %s\nStatus code: %d\nMessage: %s\nError: %s", art, sc, m, err)
	}
}

func addArticleToCartSuccess(t *testing.T, reqEtag string) {
	sCod := http.StatusCreated
	art := article{ID: "A", Quantity: 2}
	hf := func(w http.ResponseWriter, r *http.Request) {
		respondWithPayload(w, sCod, art, "")
	}
	ts := httptest.NewServer(http.HandlerFunc(hf))
	defer ts.Close()
	a := &App{BaseURL: ts.URL, HTTPClient: *ts.Client()}
	ar, sc, m, err := a.addArticleToCart("A", "B", art.ID, art.Quantity)
	if ar.String() != art.String() || sc != sCod || err != nil {
		t.Errorf("EXPECTED\nCart: %s\nStatus code: %d\nMessage: %s\n\n", art, sCod, "")
		t.Errorf("GOT\nCart: %s\nStatus code: %d\nMessage: %s\nError: %s", ar, sc, m, err)
	}
}

func getCartSuccess(t *testing.T, reqEtag string) {
	respEtag := `W/"123456789"`
	sCod := http.StatusOK
	crt := cart{
		ID:       "ABC",
		Subtotal: 30,
		Items: []item{
			item{ID: "A", Quantity: 2, UnitPrice: 10, TotalPrice: 20},
			item{ID: "B", Quantity: 5, UnitPrice: 2, TotalPrice: 10},
		},
		URL:  "http://127.0.0.1/carts/ABC",
		ETag: respEtag,
	}
	hf := func(w http.ResponseWriter, r *http.Request) {
		if reqEtag != respEtag {
			respondWithPayload(w, sCod, crt, respEtag)
		} else {
			sCod = http.StatusNotModified
			crt = *new(cart)
			w.WriteHeader(http.StatusNotModified)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(hf))
	defer ts.Close()
	a := &App{BaseURL: ts.URL, HTTPClient: *ts.Client()}
	c, sc, m, err := a.getCart(crt.ID, reqEtag)
	if c.String() != crt.String() || sc != sCod || err != nil {
		t.Errorf("EXPECTED\nCart: %s\nStatus code: %d\nMessage: %s\n\n", crt, sCod, "")
		t.Errorf("GOT\nCart: %s\nStatus code: %d\nMessage: %s\nError: %s", c, sc, m, err)
	}
}

func deleteCartSuccess(t *testing.T, reqEtag string) {
	wid := "ABC"
	sCod := http.StatusNoContent
	hf := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(sCod)
	}
	ts := httptest.NewServer(http.HandlerFunc(hf))
	defer ts.Close()
	a := &App{BaseURL: ts.URL, HTTPClient: *ts.Client()}
	sc, m, err := a.deleteCart(wid, reqEtag)
	if sc != sCod || m != "" || err != nil {
		t.Errorf("EXPECTED\nStatus code: %d\n\n", sCod)
		t.Errorf("GOT\nStatus code: %d\nMessage: %s\nError: %s", sc, m, err)
	}
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
