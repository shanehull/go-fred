package fred_test

import (
	"context"
	"os"
	"testing"

	"github.com/shanehull/go-fred"
)

func TestGetTags(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tags, err := client.GetTags(ctx, fred.WithTagLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		t.Fatal("expected at least one tag")
	}
}

func TestGetRelatedTags(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tags, err := client.GetRelatedTags(ctx, []string{"gdp"}, fred.WithTagLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		t.Fatal("expected at least one related tag for gdp")
	}
}

func TestGetTagsSeries(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	results, err := client.GetTagsSeries(ctx, []string{"gdp"}, fred.WithLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one series for gdp tag")
	}
}
