package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
)

func handleCreateCart(a *App) {
	fmt.Println("\rAttempting cart creation...")
	c, code, msg, err := a.createCart()
	if code == http.StatusCreated {
		fmt.Printf("Created cart %s\n", c)
		return
	}
	printOtherInfo(code, msg, err)
}

func handleAddArticleToCart(input *bufio.Scanner, a *App) {
	id := inputString(input, "cart ID", false)
	etag := inputString(input, "cart ETag (press return to skip)", true)
	artCode := inputString(input, "article code", false)
	artQty := inputInt(input, "article quantity (must be positive)", false)
	format := "\rAttempting to add article %q with quantity %d to cart %q ETag %s...\n"
	fmt.Printf(format, artCode, artQty, id, etag)
	art, code, msg, err := a.addArticleToCart(id, etag, artCode, artQty)
	if code == http.StatusCreated {
		fmt.Printf("Added article %s\n", art)
		return
	}
	if code == http.StatusNotFound {
		fmt.Println("No cart found with that ID")
		return
	}
	if code == http.StatusPreconditionFailed {
		fmt.Println("No cart found with that ETag")
		return
	}
	if code == http.StatusConflict {
		fmt.Println("Article already added to the cart")
		return
	}
	if code == http.StatusUnprocessableEntity {
		fmt.Println("The article does not exist")
		return
	}
	printOtherInfo(code, msg, err)
}

func handleGetCartSubtotal(input *bufio.Scanner, a *App) {
	id := inputString(input, "cart ID", false)
	etag := inputString(input, "cart ETag (press return to skip)", true)
	fmt.Printf("\rAttempting to get subtotal for cart %q with ETag %s...\n", id, etag)
	c, code, msg, err := a.getCart(id, etag)
	if code == http.StatusOK {
		fmt.Printf("Cart subtotal %g\nCart %s", c.Subtotal, c)
		return
	}
	if code == http.StatusNotFound {
		fmt.Println("No cart found with that ID")
		return
	}
	if code == http.StatusNotModified {
		fmt.Println("Cart with that ETag was not modified: omit ETag to get the cart")
		return
	}
	printOtherInfo(code, msg, err)
}

func handleDeleteCart(input *bufio.Scanner, a *App) {
	id := inputString(input, "cart ID", false)
	etag := inputString(input, "cart ETag (press return to skip)", true)
	fmt.Printf("\rAttempting to delete cart %q with ETag %s...\n", id, etag)
	code, msg, err := a.deleteCart(id, etag)
	if code == http.StatusNoContent {
		fmt.Println("Cart deleted")
		return
	}
	if code == http.StatusNotFound {
		fmt.Println("No cart found with that ID")
		return
	}
	if code == http.StatusPreconditionFailed {
		fmt.Println("No cart found with that ETag")
		return
	}
	printOtherInfo(code, msg, err)
}

func inputInt(input *bufio.Scanner, suffix string, optional bool) int {
	for {
		s := inputString(input, suffix, optional)
		i, err := strconv.Atoi(s)
		if err == nil && i > 0 {
			return i
		}
	}
}

func inputString(input *bufio.Scanner, suffix string, optional bool) string {
	msg := fmt.Sprintf("Please input %s", suffix)
	fmt.Println(msg)
	var text string
	for input.Scan() {
		text = input.Text()
		if len(text) != 0 || optional {
			break
		}
		fmt.Println(msg)
	}
	return text
}

func printOtherInfo(code int, msg string, err error) {
	if code != 0 {
		fmt.Printf("Status code: %d\n", code)
	}
	if len(msg) != 0 {
		fmt.Printf("Message: %s\n", msg)
	}
	if err != nil {
		fmt.Println(err)
	}
}
