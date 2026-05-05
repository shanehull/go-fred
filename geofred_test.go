package fred_test

import (
	"context"
	"os"
	"testing"

	"github.com/shanehull/go-fred"
)

// SMU56000000500000001A is a known series with GeoFRED data
const geoSample = "SMU56000000500000001A"

func TestGetSeriesGroup(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	group, err := client.GetSeriesGroup(ctx, geoSample)
	if err != nil {
		t.Fatal(err)
	}
	if group.Title == "" {
		t.Error("expected non-empty title")
	}
	if group.SeriesGroup == "" {
		t.Error("expected non-empty series_group")
	}
}

func TestGetSeriesData(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	data, err := client.GetSeriesData(ctx, "WIPCPI",
		fred.WithMapDate("2012-01-01"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if data.Title == "" {
		t.Error("expected non-empty title")
	}
	if len(data.Data) == 0 {
		t.Fatal("expected at least one data entry")
	}
}

func TestGetRegionalData(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	data, err := client.GetRegionalData(ctx,
		fred.WithSeriesGroup("882"),
		fred.WithRegionType("state"),
		fred.WithRegionalDate("2013-01-01"),
		fred.WithMapUnits("Dollars"),
		fred.WithSeason("NSA"),
		fred.WithRegionalFrequency("a"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Data) == 0 {
		t.Fatal("expected at least one data entry")
	}
}
