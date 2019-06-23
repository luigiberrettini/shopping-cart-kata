package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"shopping-cart-kata/catalog"
	"strings"
)

func main() {
	var baseURL = flag.String("baseUrl", "http://127.0.0.1:8000", "Address:port of the server")
	flag.Parse()
	welcomeMsg := buildWelcomeMsg(*baseURL)
	fmt.Println(welcomeMsg)
	input := bufio.NewScanner(os.Stdin)
	a := &App{BaseURL: *baseURL, HTTPClient: http.Client{}}
	for input.Scan() {
		choice := input.Text()
		switch choice {
		case "1":
			handleCreateCart(a)
		case "2":
			handleAddArticleToCart(input, a)
		case "3":
			handleGetCartSubtotal(input, a)
		case "4":
			handleDeleteCart(input, a)
		case "5":
			fmt.Println()
			os.Exit(0)
		default:
			fmt.Printf("Invalid choice %s", choice)
		}
		fmt.Printf("\n\n%s\n", welcomeMsg)
	}
}

func buildWelcomeMsg(baseURL string) string {
	arts, err := getArticles(baseURL)
	if err != nil {
		fmt.Printf("Error retrieving the catalog:\n%s", err)
		os.Exit(1)
	}
	return welcomeMsgForArticles(arts)
}

func getArticles(baseURL string) ([]catalog.Article, error) {
	var arts []catalog.Article
	url := fmt.Sprintf("%s/articles", baseURL)
	resp, err := http.Get(url)
	if err != nil {
		return arts, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&arts); err != nil {
		return arts, err
	}
	return arts, nil
}

func welcomeMsgForArticles(articles []catalog.Article) string {
	var sb strings.Builder
	sb.WriteString("\nOUR CATALOG CONTAINS THE FOLLOWING ARTICLES\n")
	sb.WriteString("      Code             |   Name             |   Price\n")
	sb.WriteString("   ---------------------------------------------------------------\n")
	indent := 3
	colWidth := 20
	for _, a := range articles {
		sCode := strings.Repeat(" ", colWidth-indent-len(a.Code))
		sName := strings.Repeat(" ", colWidth-indent-len(a.Name))
		sb.WriteString(fmt.Sprintf("      %s%s|   %s%s|  %6.2f €\n", a.Code, sCode, a.Name, sName, a.Price))
	}
	sb.WriteString("\nPLEASE SELECT AN OPERATION\n")
	sb.WriteString(" 1) Create a cart\n")
	sb.WriteString(" 2) Add an article to a cart\n")
	sb.WriteString(" 3) Get the cart subtotal\n")
	sb.WriteString(" 4) Delete a cart\n")
	sb.WriteString(" 5) Quit\n")
	return sb.String()
}
