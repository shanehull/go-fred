// Example: get the first-published values for a series at each date.
//
//	go run examples/vintage/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shanehull/go-fred"
)

func main() {
	client, err := fred.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)

	firstRelease, err := client.GetSeriesFirstRelease(ctx, "DGS20",
		fred.WithObservationStart(start),
		fred.WithObservationEnd(end),
	)
	if err != nil {
		log.Fatal(err)
	}

	current, err := client.GetSeriesObservations(ctx, "DGS20",
		fred.WithObservationStart(start),
		fred.WithObservationEnd(end),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i, fr := range firstRelease {
		if fr.IsNA || current[i].IsNA {
			continue
		}
		rev := current[i].Value - fr.Value
		if rev != 0 {
			fmt.Printf("%s  first: %.2f  latest: %.2f  revision: %+.2f\n",
				fr.Date.Format("2006-01-02"), fr.Value, current[i].Value, rev)
		}
	}
}
