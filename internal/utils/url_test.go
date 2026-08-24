package utils

import "testing"

func TestNormalizeResolvesRelativeAgainstBase(t *testing.T) {
	got, err := Normalize("https://en.wikipedia.org/wiki/WebCrawler", "/wiki/Search_engine")
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}
	want := "https://en.wikipedia.org/wiki/Search_engine"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNormalizeStripsFragment(t *testing.T) {
	got, err := Normalize("https://en.wikipedia.org/wiki/WebCrawler", "#History")
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}
	want := "https://en.wikipedia.org/wiki/WebCrawler"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNormalizeRejectsNonHTTP(t *testing.T) {
	if _, err := Normalize("https://en.wikipedia.org/wiki/WebCrawler", "mailto:me@example.com"); err == nil {
		t.Fatal("expected error for mailto link")
	}
}
