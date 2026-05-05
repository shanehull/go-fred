package fred_test

import (
	"context"
	"os"
	"testing"

	"github.com/shanehull/go-fred"
)

func TestGetCategory(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	cat, err := client.GetCategory(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if cat.ID != 1 {
		t.Errorf("expected ID 1, got %d", cat.ID)
	}
	if cat.Name == "" {
		t.Error("expected non-empty name")
	}
}

func TestGetCategoryChildren(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	children, err := client.GetCategoryChildren(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(children) == 0 {
		t.Fatal("expected at least one child category for root")
	}
	for _, c := range children {
		if c.ID == 0 {
			t.Error("expected non-zero ID")
		}
	}
}

func TestGetCategoryRelated(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	related, err := client.GetCategoryRelated(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range related {
		if c.ID == 0 {
			t.Error("expected non-zero ID")
		}
	}
}

func TestGetCategoryTags(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tags, err := client.GetCategoryTags(ctx, 1, fred.WithTagLimit(5))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		t.Fatal("expected at least one tag")
	}
	for _, tag := range tags {
		if tag.Name == "" {
			t.Error("expected non-empty tag name")
		}
	}
}

func TestGetCategoryRelatedTags(t *testing.T) {
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tags, err := client.GetCategoryRelatedTags(ctx, 1, []string{"business"}, fred.WithTagLimit(5))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) == 0 {
		t.Fatal("expected at least one related tag")
	}
	for _, tag := range tags {
		if tag.Name == "" {
			t.Error("expected non-empty tag name")
		}
	}
}
