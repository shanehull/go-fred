// Example: show recent release dates for GDP.
//
//	go run examples/releases/main.go
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

	rel, err := client.GetRelease(ctx, 53)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Release: %s\n\n", rel.Name)

	dates, err := client.GetReleaseDates(ctx, 53,
		fred.WithReleaseDateLimit(5),
		fred.WithReleaseDateSortOrder(fred.SortDesc),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Recent release dates:")
	for _, d := range dates {
		fmt.Printf("  %s\n", d.Date)
	}
}
