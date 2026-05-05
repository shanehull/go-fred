package fred_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/shanehull/go-fred"
)

func TestGetSeriesInfo(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	s, err := client.GetSeriesInfo(ctx, "DGS20")
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != "DGS20" {
		t.Errorf("expected DGS20, got %s", s.ID)
	}
	if s.Title == "" {
		t.Error("expected non-empty title")
	}
}

func TestGetSeriesObservations(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	obs, err := client.GetSeriesObservations(ctx, "DGS20",
		fred.WithObservationStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		fred.WithObservationEnd(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)),
		fred.WithObservationSortOrder(fred.SortDesc),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(obs) == 0 {
		t.Fatal("expected at least one observation")
	}
	if obs[0].Date.Before(obs[len(obs)-1].Date) {
		t.Error("expected descending date order")
	}
	for _, o := range obs {
		if !o.IsNA && o.Value <= 0 {
			t.Errorf("unexpected non-positive value %f at %s", o.Value, o.Date)
		}
	}
}

func TestGetSeriesAllReleases(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	obs, err := client.GetSeriesAllReleases(ctx, "DGS20",
		fred.WithRealtimeStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		fred.WithObservationLimit(50),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(obs) == 0 {
		t.Fatal("expected at least one observation")
	}
	for _, o := range obs {
		if o.RealtimeStart.IsZero() {
			t.Error("expected non-zero realtime_start")
		}
	}
}

func TestGetSeriesFirstRelease(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	obs, err := client.GetSeriesFirstRelease(ctx, "DGS20",
		fred.WithRealtimeStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		fred.WithObservationLimit(50),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(obs) == 0 {
		t.Fatal("expected at least one observation")
	}
	seen := make(map[string]bool)
	for _, o := range obs {
		d := o.Date.Format("2006-01-02")
		if seen[d] {
			t.Errorf("duplicate date %s in first release", d)
		}
		seen[d] = true
	}
}

func TestGetSeriesAsOf(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	asOf := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	obs, err := client.GetSeriesAsOf(ctx, "DGS20", asOf,
		fred.WithRealtimeStart(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
		fred.WithObservationStart(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range obs {
		if o.RealtimeEnd.After(asOf) {
			t.Errorf("realtime_end %s is after asOf %s", o.RealtimeEnd, asOf)
		}
	}
}

func TestGetSeriesVintageDates(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	dates, err := client.GetSeriesVintageDates(ctx, "DGS20",
		fred.WithObservationLimit(20),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(dates) == 0 {
		t.Fatal("expected at least one vintage date")
	}
	for _, d := range dates {
		if d.IsZero() {
			t.Error("expected non-zero vintage date")
		}
	}
}

func TestGetSeriesCategories(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	cats, err := client.GetSeriesCategories(ctx, "DGS20")
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) == 0 {
		t.Fatal("expected at least one category")
	}
}

func TestGetSeriesRelease(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	rel, err := client.GetSeriesRelease(ctx, "DGS20")
	if err != nil {
		t.Fatal(err)
	}
	if rel.ID == 0 {
		t.Error("expected non-zero release ID")
	}
}

func TestGetSeriesTags(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tags, err := client.GetSeriesTags(ctx, "DGS20", fred.WithTagLimit(5))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		t.Fatal("expected at least one tag")
	}
}

func TestSearchSeriesTags(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tags, err := client.SearchSeriesTags(ctx, "GDP", fred.WithTagLimit(5))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		t.Fatal("expected at least one tag for GDP search")
	}
}

func TestSearchSeriesRelatedTags(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tags, err := client.SearchSeriesRelatedTags(ctx, "GDP", []string{"business"}, fred.WithTagLimit(5))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		t.Fatal("expected at least one related tag")
	}
}

func TestGetSeriesUpdates(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	results, err := client.GetSeriesUpdates(ctx, fred.WithUpdateLimit(5))
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one updated series")
	}
}
