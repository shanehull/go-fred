package fred_test

import (
	"context"
	"os"
	"testing"

	"github.com/shanehull/go-fred"
)

func TestGetReleases(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	releases, err := client.GetReleases(ctx, fred.WithReleaseLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) == 0 {
		t.Fatal("expected at least one release")
	}
}

func TestGetReleasesDates(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	dates, err := client.GetReleasesDates(ctx, fred.WithReleaseDateLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(dates) == 0 {
		t.Fatal("expected at least one release date")
	}
}

func TestGetRelease(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	rel, err := client.GetRelease(ctx, 53)
	if err != nil {
		t.Fatal(err)
	}
	if rel.ID != 53 {
		t.Errorf("expected ID 53, got %d", rel.ID)
	}
}

func TestGetReleaseDates(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	dates, err := client.GetReleaseDates(ctx, 53, fred.WithReleaseDateLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(dates) == 0 {
		t.Fatal("expected at least one release date for GDP")
	}
}

func TestGetReleaseSources(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	sources, err := client.GetReleaseSources(ctx, 53)
	if err != nil {
		t.Fatal(err)
	}
	if len(sources) == 0 {
		t.Fatal("expected at least one source for GDP release")
	}
}

func TestGetReleaseTags(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tags, err := client.GetReleaseTags(ctx, 53, fred.WithTagLimit(5))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		t.Fatal("expected at least one tag for GDP release")
	}
}

func TestGetReleaseRelatedTags(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tags, err := client.GetReleaseRelatedTags(ctx, 53, []string{"gdp"}, fred.WithTagLimit(5))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		t.Fatal("expected at least one related tag for GDP")
	}
}

func TestGetReleaseTables(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	elements, err := client.GetReleaseTables(ctx, 53)
	if err != nil {
		t.Fatal(err)
	}
	if len(elements) == 0 {
		t.Fatal("expected at least one table element for GDP release")
	}
}
