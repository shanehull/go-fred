// Example: get GeoFRED map metadata and regional data.
//
//	go run examples/geofred/main.go
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

	group, err := client.GetSeriesGroup(ctx, "SMU56000000500000001A")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Map: %s\n", group.Title)
	fmt.Printf("Group: %s  Region: %s  Frequency: %s\n\n",
		group.SeriesGroup, group.RegionType, group.Frequency)

	data, err := client.GetRegionalData(ctx,
		fred.WithSeriesGroup(group.SeriesGroup),
		fred.WithRegionType(group.RegionType),
		fred.WithRegionalDate(group.MaxDate),
		fred.WithRegionalFrequency("a"),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range data.Data {
		fmt.Printf("%-30s  value: %s\n", d.Region, d.Value)
	}
}
