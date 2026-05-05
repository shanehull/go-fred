// Example: list FRED data sources and their releases.
//
//	go run examples/sources/main.go
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

	sources, err := client.GetSources(ctx,
		fred.WithSourceLimit(5),
		fred.WithSourceSortOrder(fred.SortAsc),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range sources {
		fmt.Printf("%4d  %s\n", s.ID, s.Name)
	}

	fmt.Println()

	src, err := client.GetSource(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Source 1: %s\n\n", src.Name)

	releases, err := client.GetSourceReleases(ctx, 1, fred.WithReleaseLimit(5))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Recent releases:")
	for _, r := range releases {
		fmt.Printf("  %4d  %s\n", r.ID, r.Name)
	}
}
