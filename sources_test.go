package fred_test

import (
	"context"
	"os"
	"testing"

	"github.com/shanehull/go-fred"
)

func TestGetSources(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	sources, err := client.GetSources(ctx, fred.WithSourceLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(sources) == 0 {
		t.Fatal("expected at least one source")
	}
}

func TestGetSource(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	src, err := client.GetSource(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if src.ID != 1 {
		t.Errorf("expected ID 1, got %d", src.ID)
	}
}

func TestGetSourceReleases(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	releases, err := client.GetSourceReleases(ctx, 1, fred.WithReleaseLimit(10))
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) == 0 {
		t.Fatal("expected at least one release for source 1")
	}
}
