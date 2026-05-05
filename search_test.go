package fred_test

import (
	"context"
	"os"
	"testing"

	"github.com/shanehull/go-fred"
)

func TestSearchSeries(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	results, err := client.SearchSeries(ctx, "GDP", fred.WithLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(results) < 5 {
		t.Fatalf("expected at least 5 results, got %d", len(results))
	}
	for _, r := range results {
		if r.ID == "" {
			t.Error("expected non-empty ID in search result")
		}
		if r.Title == "" {
			t.Error("expected non-empty Title in search result")
		}
	}
}

func TestSearchSeriesPagination(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	results, err := client.SearchSeries(ctx, "GDP", fred.WithLimit(1500))
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1500 {
		t.Fatalf("expected 1500 results, got %d", len(results))
	}
}

func TestGetReleaseSeries(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	results, err := client.GetReleaseSeries(ctx, 53, fred.WithLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result for GDP release")
	}
	for _, r := range results {
		if r.ID == "" {
			t.Error("expected non-empty ID")
		}
	}
}

func TestGetCategorySeries(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	results, err := client.GetCategorySeries(ctx, 1, fred.WithLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result for category 1")
	}
	for _, r := range results {
		if r.ID == "" {
			t.Error("expected non-empty ID")
		}
	}
}
