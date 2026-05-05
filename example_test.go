package fred_test

import (
	"context"
	"fmt"
	"time"

	"github.com/shanehull/go-fred"
)

func ExampleClient_GetSeriesObservations() {
	client, _ := fred.New(fred.WithAPIKey("your-key"))
	obs, _ := client.GetSeriesObservations(context.Background(), "DGS20",
		fred.WithObservationStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
	)
	for _, o := range obs {
		if !o.IsNA {
			fmt.Printf("%s: %.2f\n", o.Date.Format("2006-01-02"), o.Value)
		}
	}
}

func ExampleClient_SearchSeries() {
	client, _ := fred.New(fred.WithAPIKey("your-key"))
	results, _ := client.SearchSeries(context.Background(), "GDP",
		fred.WithLimit(5),
		fred.WithSortOrder(fred.SortDesc),
	)
	for _, r := range results {
		fmt.Printf("%s: %s\n", r.ID, r.Title)
	}
}
