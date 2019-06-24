package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"shopping-cart-kata/appservice"
	"shopping-cart-kata/cache"
	"shopping-cart-kata/cart"
	"strings"
	"testing"
)

type uncache struct{}

func (c *uncache) GetByEtagWithID(etag string, wid string) (cache.Etagger, bool) {
	var e cache.Etagger
	return e, false
}

func (c *uncache) AddOrReplace(wid string, e cache.Etagger) {
	return
}

func (c *uncache) Remove(wid string) {
	return
}

func TestHash(t *testing.T) {
	const id = 1
	a := testApp(new(uncache))
	wid, err := a.encode(id)
	if err != nil {
		t.Fatalf("Encoding error %v", err)
	}
	did, err := a.decode(wid)
	if err != nil {
		t.Fatalf("Decoding error %v", err)
	}
	if did != id {
		t.Errorf("Encoded id %d to value %s decoded to value %d", id, wid, did)
	}
}

func TestCreateWihtNonInitializedAppSvc(t *testing.T) {
	cfg := Config{HashSalt: "", ListenAddress: "127.0.0.1"}
	a := &App{
		AppSvc:    appservice.AppService{CartIDG: nil, CartDB: nil, Catalog: nil, PromEng: nil},
		HashGen:   createHashGenerator(cfg.HashSalt),
		Router:    mux.NewRouter().StrictSlash(true),
		CartCache: new(uncache),
	}
	a.ConfigRoutes(cfg.ListenAddress)
	a.ConfigURLBuilders()
	req, _ := http.NewRequest("POST", "http://127.0.0.1/carts", nil)
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusInternalServerError, response)
}

func TestUnsuccessfulAddArticle(t *testing.T) {
	a := testApp(new(uncache))
	req, _ := http.NewRequest("POST", "http://127.0.0.1/carts", nil)
	response := executeRequest(a, req)

	var c cartVM
	json.NewDecoder(response.Body).Decode(&c)
	url := fmt.Sprintf("%s/items", c.URL)
	art := a.AppSvc.Catalog.GetArticles()[0]
	item := itemCreateVM{ID: art.Code, Quantity: 5}

	item.ID = "notValid"
	j, _ := json.Marshal(item)
	b := bytes.NewBuffer(j)
	req, _ = http.NewRequest("POST", url, b)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusUnprocessableEntity, response)

	item.ID = art.Code
	item.Quantity = 0
	j, _ = json.Marshal(item)
	b = bytes.NewBuffer(j)
	req, _ = http.NewRequest("POST", url, b)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusUnprocessableEntity, response)

	item.Quantity = -2
	j, _ = json.Marshal(item)
	b = bytes.NewBuffer(j)
	req, _ = http.NewRequest("POST", url, b)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusUnprocessableEntity, response)
}

func TestUnsuccessfulSetArticleQuantity(t *testing.T) {
	a := testApp(new(uncache))
	req, _ := http.NewRequest("POST", "http://127.0.0.1/carts", nil)
	response := executeRequest(a, req)

	var c cartVM
	json.NewDecoder(response.Body).Decode(&c)
	url := fmt.Sprintf("%s/items", c.URL)
	art := a.AppSvc.Catalog.GetArticles()[0]
	item := itemCreateVM{ID: art.Code, Quantity: 5}

	j, _ := json.Marshal(item)
	b := bytes.NewBuffer(j)
	req, _ = http.NewRequest("POST", url, b)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusCreated, response)

	item.ID = "notValid"
	j, _ = json.Marshal(item)
	b = bytes.NewBuffer(j)
	req, _ = http.NewRequest("PUT", url, b)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusUnprocessableEntity, response)

	item.ID = art.Code
	item.Quantity = 0
	j, _ = json.Marshal(item)
	b = bytes.NewBuffer(j)
	req, _ = http.NewRequest("PUT", url, b)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusUnprocessableEntity, response)

	item.Quantity = -2
	j, _ = json.Marshal(item)
	b = bytes.NewBuffer(j)
	req, _ = http.NewRequest("PUT", url, b)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusUnprocessableEntity, response)
}

func TestGetCartWithArticles(t *testing.T) {
	a := testApp(new(uncache))
	req, _ := http.NewRequest("POST", "http://127.0.0.1/carts", nil)
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusCreated, response)

	var c cartVM
	if err := json.NewDecoder(response.Body).Decode(&c); err != nil {
		t.Errorf("Error decoding response.Body: %s", err)
	}
	url := fmt.Sprintf("%s/items", c.URL)
	art := a.AppSvc.Catalog.GetArticles()[0]
	item := itemCreateVM{ID: art.Code, Quantity: 5}
	j, err := json.Marshal(item)
	if err != nil {
		t.Errorf("Error marshalling the article; %s", err)
	}
	b := bytes.NewBuffer(j)
	req, _ = http.NewRequest("POST", url, b)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusCreated, response)

	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusBadRequest, response)

	b = bytes.NewBuffer(j)
	req, _ = http.NewRequest("PUT", url, b)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusOK, response)

	req, _ = http.NewRequest("GET", c.URL, nil)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusOK, response)

	rb, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error reading response.Body: %s", err)
	}
	respBody := string(rb)

	var dc cartVM
	if err := json.NewDecoder(bytes.NewReader(rb)).Decode(&dc); err != nil {
		t.Errorf("Error decoding response.Body: %s", err)
	}
	if dc.ID != c.ID {
		t.Errorf("Cart ID %s instead of %s\n%s", dc.ID, c.ID, respBody)
	}
	if dc.Subtotal == 0 {
		t.Errorf("Subtotal 0\n%s", respBody)
	}
	if !strings.Contains(dc.URL, c.ID) {
		t.Errorf("URL %s does not contain ID %s\n%s", dc.URL, c.ID, respBody)
	}
	if etag := response.Header().Get("ETag"); len(etag) == 0 {
		t.Errorf("No etag\n%s", respBody)
	}
	if len(dc.Items) == 0 {
		t.Fatalf("No article added\n%v\n%s", dc, respBody)
	}
	ci := dc.Items[0]
	if ci.ID != item.ID || ci.Quantity != item.Quantity {
		t.Errorf("Article to add does not match item:\n%v\n%v", item, ci)
	}
}

func TestGetNonExistentCart(t *testing.T) {
	a := testApp(new(uncache))
	req, _ := http.NewRequest("GET", "http://127.0.0.1/carts/myHash", nil)
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusNotFound, response)
}

func TestGetCreatedCartWithNormalReq(t *testing.T) {
	a := testApp(new(uncache))
	req, _ := http.NewRequest("POST", "http://127.0.0.1/carts", nil)
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusCreated, response)

	var c cartVM
	if err := json.NewDecoder(response.Body).Decode(&c); err != nil {
		t.Errorf("Error decoding response.Body: %s", err)
	}
	req, _ = http.NewRequest("GET", c.URL, nil)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusOK, response)

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Error reading response.Body: %s", err)
	}
	respBody := string(b)

	var dc cartVM
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&dc); err != nil {
		t.Errorf("Error decoding response.Body: %s", err)
	}
	if dc.ID != c.ID {
		t.Errorf("Cart ID %s instead of %s\n%s", dc.ID, c.ID, respBody)
	}
	if dc.Subtotal != 0 {
		t.Errorf("Subtotal %g instead of 0\n%s", dc.Subtotal, respBody)
	}
	if !strings.Contains(dc.URL, c.ID) {
		t.Errorf("URL %s does not contain ID %s\n%s", dc.URL, c.ID, respBody)
	}
	if etag := response.Header().Get("ETag"); len(etag) == 0 {
		t.Errorf("No etag\n%s", respBody)
	}
	if len(dc.Items) != 0 {
		t.Fatalf("Articles found in created cart\n%v\n%s", dc, respBody)
	}
}

func TestGetCreatedCartWithCondReq(t *testing.T) {
	a := testApp(cache.NewCache())
	req, _ := http.NewRequest("POST", "http://127.0.0.1/carts", nil)
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusCreated, response)
	var c cartVM
	if err := json.NewDecoder(response.Body).Decode(&c); err != nil {
		t.Errorf("Error decoding response.Body: %s", err)
	}
	req, _ = http.NewRequest("GET", c.URL, nil)
	req.Header.Add("If-None-Match", response.Header().Get("ETag"))
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusNotModified, response)
}

func TestDeletedCachedCart(t *testing.T) {
	a := testApp(cache.NewCache())
	req, _ := http.NewRequest("POST", "http://127.0.0.1/carts", nil)
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusCreated, response)

	var c cartVM
	if err := json.NewDecoder(response.Body).Decode(&c); err != nil {
		t.Errorf("Error decoding response.Body: %s", err)
	}

	req, _ = http.NewRequest("DELETE", c.URL, nil)
	req.Header.Add("If-Match", "fakeEtag")
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusPreconditionFailed, response)

	req, _ = http.NewRequest("DELETE", c.URL, nil)
	req.Header.Add("If-Match", response.Header().Get("ETag"))
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusNoContent, response)
}

func TestDeleteNonCachedCart(t *testing.T) {
	a := testApp(new(uncache))
	req, _ := http.NewRequest("DELETE", "http://127.0.0.1/carts/1", nil)
	req.Header.Add("If-Match", "fakeEtag")
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusPreconditionFailed, response)
}

func TestGetDeletedCart(t *testing.T) {
	a := testApp(new(uncache))
	req, _ := http.NewRequest("POST", "http://127.0.0.1/carts", nil)
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusCreated, response)

	var c cartVM
	if err := json.NewDecoder(response.Body).Decode(&c); err != nil {
		t.Errorf("Error decoding response.Body: %s", err)
	}
	req, _ = http.NewRequest("DELETE", c.URL, nil)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusNoContent, response)

	req, _ = http.NewRequest("GET", c.URL, nil)
	response = executeRequest(a, req)
	checkResponseCode(t, http.StatusNotFound, response)
}

func testApp(c cache.Cache) *App {
	cfg := Config{HashSalt: "a9a21fd753f94", ListenAddress: "127.0.0.1"}
	a := &App{
		AppSvc: appservice.AppService{
			CartIDG: new(generator),
			CartDB:  cart.NewStore(),
			Catalog: createCatalog(),
			PromEng: createPromoEngine(),
		},
		HashGen:   createHashGenerator(cfg.HashSalt),
		Router:    mux.NewRouter().StrictSlash(true),
		CartCache: c,
	}
	a.ConfigRoutes(cfg.ListenAddress)
	a.ConfigURLBuilders()
	return a
}

func executeRequest(a *App, req *http.Request) *httptest.ResponseRecorder {
	r := httptest.NewRecorder()
	a.Router.ServeHTTP(r, req)
	return r
}

func checkResponseCode(t *testing.T, expected int, response *httptest.ResponseRecorder) {
	if expected != response.Code {
		t.Fatalf("Expected status code %d instead of %d\n%s", expected, response.Code, response.Body)
	}
}
