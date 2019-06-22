package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"os"
	"shopping-cart-kata/appservice"
	"shopping-cart-kata/cache"
	"shopping-cart-kata/cart"
	"shopping-cart-kata/catalog"
	"shopping-cart-kata/promotion"
)

func main() {
	var cfgFilePath = flag.String("config", "config.json", "Location of the config file")
	flag.Parse()
	cfg, err := loadConfig(*cfgFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	a := createApp(cfg)
	a.ConfigRoutes()
	a.ConfigURLBuilders()
	a.Run(cfg.ListenAddress)
}

func loadConfig(cfgFilePath string) (Config, error) {
	var cfg Config
	file, err := os.Open(cfgFilePath)
	if err != nil {
		return cfg, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	return cfg, err
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
