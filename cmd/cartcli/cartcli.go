package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var welcomeMessage = `
OUR CATALOG CONTAINS THE FOLLOWING ARTICLES
   Code         | Name                     |  Price
   --------------------------------------------------
   VOUCHER      | CompanyName Voucher      |   5.00 €
   TSHIRT       | CompanyName T-Shirt      |  20.00 €
   MUG          | CompanyName Coffee Mug   |   7.50 €

PLEASE CHOOSE AN OPERATION
 1) Create a cart
 2) Add an article to a cart
 3) Get the cart subtotal
 4) Delete a cart
 5) Quit
`

func main() {
	var authority = flag.String("address", "127.0.0.1:8000", "Address:port of the server")
	flag.Parse()
	fmt.Println(welcomeMessage)
	mainLoop(*authority)
}

func mainLoop(authority string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "1":
			createCart(scanner, authority)
		case "2":
			addArticleToCart(scanner, authority)
		case "3":
			getCartSubtotal(scanner, authority)
		case "4":
			deleteCart(scanner, authority)
		case "5":
			os.Exit(0)
		default:
			fmt.Printf("Invalid choice %s", input)
		}
		fmt.Println(welcomeMessage)
	}
}

func createCart(b *bufio.Scanner, authority string) {
	fmt.Println("Creating a cart...")
}

func addArticleToCart(b *bufio.Scanner, authority string) {
	fmt.Println("Adding article to a cart...")
}

func getCartSubtotal(b *bufio.Scanner, authority string) {
	fmt.Println("Retrieving the cart...")
}

func deleteCart(b *bufio.Scanner, authority string) {
	fmt.Println("Deleting the cart...")
}
