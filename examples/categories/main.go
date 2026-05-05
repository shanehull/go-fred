// Example: explore the category hierarchy starting from root.
//
//	go run examples/categories/main.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shanehull/go-fred"
)

func main() {
	client, err := fred.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	root, err := client.GetCategory(ctx, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Root: %s\n\n", root.Name)

	children, err := client.GetCategoryChildren(ctx, 0)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range children {
		fmt.Printf("%4d  %s\n", c.ID, c.Name)
	}
}
