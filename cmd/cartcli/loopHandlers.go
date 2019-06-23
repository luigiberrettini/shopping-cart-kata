package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func handleCreateCart(a *App) {
	fmt.Println("\rAttempting cart creation...")
	c, err := a.createCart()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Created cart %s\n", c)
}

func handleAddArticleToCart(input *bufio.Scanner, a *App) {

}

func handleGetCartSubtotal(input *bufio.Scanner, a *App) {

}

func handleDeleteCart(input *bufio.Scanner, a *App) {
	msg := "Please input cart ID"
	fmt.Println(msg)
	var id string
	for input.Scan() {
		id = input.Text()
		if len(id) != 0 {
			break
		}
		fmt.Println(msg)
	}
	fmt.Printf("\rAttempting to delete cart %q...\n", id)
	code, msg, err := a.deleteCart(id)
	if err != nil {
		if code != 0 {
			fmt.Printf("Status code %d\n", code)
		}
		fmt.Println(err)
		return
	}
	if code != http.StatusNoContent {
		fmt.Printf("Status code %d: %s\n", code, msg)
		return
	}
	fmt.Println("Cart deleted")
}
