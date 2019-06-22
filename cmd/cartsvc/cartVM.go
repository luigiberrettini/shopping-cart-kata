package main

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"shopping-cart-kata/pricedcart"
)

type cartVM struct {
	ID       string      `json:"id"`
	Subtotal float64     `json:"subTotal"`
	Items    []itemGetVM `json:"items"`
	URL      string      `json:"url"`
	etag     string
}

func fromPricedCart(pc pricedcart.PricedCart, wid string, url string) cartVM {
	var c cartVM
	c.ID = wid
	c.Subtotal = pc.GetSubtotal()
	pcItems := pc.GetItems()
	c.Items = make([]itemGetVM, len(pcItems))
	for i, pci := range pcItems {
		c.Items[i] = fromPricedItem(pci)
	}
	c.URL = url
	return c
}

func (c *cartVM) GetEtag() string {
	return c.etag
}

func (c *cartVM) ComputeEtag() {
	data, err := json.Marshal(c)
	if err != nil {
		c.etag = ""
		return
	}
	crc := crc32.ChecksumIEEE(data)
	c.etag = fmt.Sprintf(`W/"%s-%d-%08X"`, "cart", len(data), crc)
}
