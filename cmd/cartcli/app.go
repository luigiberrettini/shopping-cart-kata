package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrReqPreparation error on request preparation
var ErrReqPreparation = errors.New("Error on request preparation")

// ErrReqExecution error on request execution
var ErrReqExecution = errors.New("Error on request execution")

// ErrRespDecode error decoding the response body
var ErrRespDecode = errors.New("Error decoding the response body")

// App is the client application
type App struct {
	BaseURL    string
	HTTPClient http.Client
}

func (a *App) createCart() (cart, error) {
	var c cart
	url := fmt.Sprintf("%s/carts", a.BaseURL)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return c, ErrReqPreparation
	}
	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return c, ErrReqExecution
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, ErrRespDecode
	}
	return c, nil
}

func (a *App) addArticleToCart() {
	fmt.Println("Adding article to a cart...")
}

func (a *App) getCartSubtotal() {
	fmt.Println("Retrieving the cart...")
}

func (a *App) deleteCart(id string) (int, string, error) {
	url := fmt.Sprintf("%s/carts/%s", a.BaseURL, id)
	fmt.Println(url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, "", ErrReqPreparation
	}
	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return 0, "", ErrReqExecution
	}
	if resp.StatusCode != http.StatusNoContent {
		var e apiError
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return resp.StatusCode, "", ErrRespDecode
		}
		return resp.StatusCode, e.Msg, nil
	}
	return resp.StatusCode, "", nil
}
