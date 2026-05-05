// Example: search FRED for GDP-related series.
//
//	go run examples/search/main.go
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
	results, err := client.SearchSeries(ctx, "GDP",
		fred.WithLimit(10),
		fred.WithSortOrder(fred.SortDesc),
		fred.WithOrderBy(fred.OrderByPopularity),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		fmt.Printf("%-12s  %s\n", r.ID, r.Title)
	}
	fmt.Printf("\nFound %d results\n", len(results))
}
