package webapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"shopping-cart-kata/appservice"
	"shopping-cart-kata/cache"
	"shopping-cart-kata/cart"
	"shopping-cart-kata/catalog"
	"shopping-cart-kata/promotion"
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
	a := testApp()
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
	cfg := Config{CompanyName: "", HashSalt: ""}
	a := &App{
		AppSvc:    appservice.AppService{CartIDG: nil, CartDB: nil, Catalog: nil, PromEng: nil},
		HashGen:   createHashGenerator(cfg.HashSalt),
		Router:    mux.NewRouter().StrictSlash(true),
		CartCache: new(uncache),
	}
	a.ConfigRoutes()
	a.ConfigURLBuilders()
	req, _ := http.NewRequest("POST", "/carts", nil)
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusInternalServerError, response)
}

func TestGetCartWithArticles(t *testing.T) {
	a := testApp()
	req, _ := http.NewRequest("POST", "/carts", nil)
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
	a := testApp()
	req, _ := http.NewRequest("GET", "/carts/myHash", nil)
	response := executeRequest(a, req)
	checkResponseCode(t, http.StatusNotFound, response)
}

func TestGetCreatedCart(t *testing.T) {
	a := testApp()
	req, _ := http.NewRequest("POST", "/carts", nil)
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

func TestGetDeletedCart(t *testing.T) {
	a := testApp()
	req, _ := http.NewRequest("POST", "/carts", nil)
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

func testApp() *App {
	cfg := Config{CompanyName: "CmPnY", HashSalt: "a9a21fd753f94"}
	a := createApp(cfg)
	a.ConfigRoutes()
	a.ConfigURLBuilders()
	return a
}

func createApp(cfg Config) *App {
	return &App{
		AppSvc: appservice.AppService{
			CartIDG: new(dummyIDGenerator),
			CartDB:  cart.NewStore(),
			Catalog: createCatalog(cfg.CompanyName),
			PromEng: createPromoEngine(),
		},
		HashGen:   createHashGenerator(cfg.HashSalt),
		Router:    mux.NewRouter().StrictSlash(true),
		CartCache: new(uncache),
	}
}

type dummyIDGenerator struct {
	id int64
}

// NextID id generation
func (g *dummyIDGenerator) NextID() int64 {
	g.id++
	return g.id
}

func createHashGenerator(salt string) *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 32
	hg, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	return hg
}

func createCatalog(company string) catalog.Catalog {
	c := catalog.NewCatalog()
	c.AddArticle(catalog.Article{
		Code:  "VOUCHER",
		Name:  fmt.Sprintf("%s Voucher", company),
		Price: 5.0,
	})
	c.AddArticle(catalog.Article{
		Code:  "TSHIRT",
		Name:  fmt.Sprintf("%s T-Shirt", company),
		Price: 20.0,
	})
	c.AddArticle(catalog.Article{
		Code:  "MUG",
		Name:  fmt.Sprintf("%s Coffee Mug", company),
		Price: 7.5,
	})
	return c
}

func createPromoEngine() promotion.Engine {
	e := promotion.NewEngine()
	f1 := promotion.TwoForOne
	f2 := promotion.DiscountForThreeOrMore
	e.AddRule(&f1)
	e.AddRule(&f2)
	return e
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
