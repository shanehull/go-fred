// Example: fetch and display recent 20-year Treasury yield observations.
//
//	go run examples/series/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shanehull/go-fred"
)

func main() {
	client, err := fred.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	obs, err := client.GetSeriesObservations(ctx, "DGS20",
		fred.WithObservationStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		fred.WithObservationSortOrder(fred.SortDesc),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, o := range obs {
		if o.IsNA {
			fmt.Printf("%s  missing\n", o.Date.Format("2006-01-02"))
		} else {
			fmt.Printf("%s  %.2f%%\n", o.Date.Format("2006-01-02"), o.Value)
		}
	}

	fmt.Printf("\n%d observations retrieved\n", len(obs))
	if os.Getenv("FRED_API_KEY") == "" {
		fmt.Println("Set $FRED_API_KEY to use your own key.")
	}
}
