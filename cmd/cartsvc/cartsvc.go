package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"shopping-cart-kata/appservice"
	"shopping-cart-kata/cache"
	"shopping-cart-kata/cart"
	"shopping-cart-kata/catalog"
	"shopping-cart-kata/promotion"
)

func main() {
	cfg := loadConfig()
	a := createApp(cfg)
	a.ConfigRoutes(cfg.Authority)
	a.ConfigURLBuilders()
	a.Run(cfg.ListenAddress)
}

func loadConfig() Config {
	var companyName = flag.String("company", "AcME", "Company name for catalog articles")
	var hashSalt = flag.String("salt", "a9a21fd753f9431381c3980c7664aab6", "Hash salt for REST IDs")
	var listenAddress = flag.String("listen", "127.0.0.1:8000", "Address:port on which to listen")
	var authority = flag.String("authority", "127.0.0.1:8000", "Authority part of REST URLs")
	flag.Parse()
	return Config{
		CompanyName:   *companyName,
		HashSalt:      *hashSalt,
		ListenAddress: *listenAddress,
		Authority:     *authority,
	}
}

func createApp(cfg Config) *App {
	return &App{
		AppSvc: appservice.AppService{
			CartIDG: new(generator),
			CartDB:  cart.NewStore(),
			Catalog: createCatalog(cfg.CompanyName),
			PromEng: createPromoEngine(),
		},
		HashGen:   createHashGenerator(cfg.HashSalt),
		Router:    mux.NewRouter().StrictSlash(true),
		CartCache: cache.NewCache(),
	}
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
