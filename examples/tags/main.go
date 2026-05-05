// Example: explore popular FRED tags.
//
//	go run examples/tags/main.go
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

	tags, err := client.GetTags(ctx,
		fred.WithTagLimit(10),
		fred.WithTagOrderBy(fred.OrderByPopularity),
		fred.WithTagSortOrder(fred.SortDesc),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range tags {
		fmt.Printf("%-30s  series: %d  popularity: %d\n", t.Name, t.SeriesCount, t.Popularity)
	}

	fmt.Println("\nRelated to 'gdp':")
	related, err := client.GetRelatedTags(ctx, []string{"gdp"}, fred.WithTagLimit(5))
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range related {
		fmt.Printf("  %s\n", t.Name)
	}
}
