package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrReqPreparation error on request preparation
var ErrReqPreparation = errors.New("Error on request preparation")

// ErrReqExecution error on request execution
var ErrReqExecution = errors.New("Error on request execution")

// ErrRespRead error reading the response body
var ErrRespRead = errors.New("Error reading response body")

// ErrRespDecode error decoding the response body
var ErrRespDecode = errors.New("Error decoding the response body")

// App is the client application
type App struct {
	BaseURL    string
	HTTPClient http.Client
}

func (a *App) createCart() (cart, int, string, error) {
	var c cart
	url := fmt.Sprintf("%s/carts", a.BaseURL)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return c, 0, "", ErrReqPreparation
	}
	i, code, msg, err := performReq(a, req, &c)
	return *i.(*cart), code, msg, err
}

func (a *App) addArticleToCart(id string, etag string, aCod string, aQty int) (article, int, string, error) {
	var art article
	url := fmt.Sprintf("%s/carts/%s/items", a.BaseURL, id)
	j, err := json.Marshal(article{ID: aCod, Quantity: aQty})
	if err != nil {
		return art, 0, "", ErrReqPreparation
	}
	b := bytes.NewBuffer(j)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return art, 0, "", ErrReqPreparation
	}
	if etag != "" {
		req.Header.Set("If-Match", etag)
	}
	i, code, msg, err := performReq(a, req, &art)
	return *i.(*article), code, msg, err
}

func (a *App) getCart(id string, etag string) (cart, int, string, error) {
	var c cart
	url := fmt.Sprintf("%s/carts/%s", a.BaseURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c, 0, "", ErrReqPreparation
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	i, code, msg, err := performReq(a, req, &c)
	return *i.(*cart), code, msg, err
}

func (a *App) deleteCart(id string, etag string) (int, string, error) {
	url := fmt.Sprintf("%s/carts/%s", a.BaseURL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, "", ErrReqPreparation
	}
	if etag != "" {
		req.Header.Set("If-Match", etag)
	}

	_, code, msg, err := performReq(a, req, nil)
	return code, msg, err
}

func performReq(a *App, req *http.Request, i interface{}) (interface{}, int, string, error) {
	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return nil, resp.StatusCode, "", ErrReqExecution
	}

	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, "", ErrRespRead
	}
	if len(rb) == 0 {
		return nil, resp.StatusCode, "", nil
	}
	var e apiError
	err1 := json.NewDecoder(bytes.NewReader(rb)).Decode(&e)
	var err2 error
	if i != nil {
		err2 = json.NewDecoder(bytes.NewReader(rb)).Decode(&i)
	}
	if err1 != nil && err2 != nil {
		return nil, resp.StatusCode, "", ErrRespDecode
	}
	if etag := resp.Header.Get("ETag"); err2 == nil && len(etag) != 0 {
		c := i.(*cart)
		c.ETag = etag
	}
	return i, resp.StatusCode, e.Msg, nil
}
