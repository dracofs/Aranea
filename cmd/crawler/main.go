package main

import (
	"ARANEA/internal/crawler"
	"ARANEA/internal/utils"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	seed := flag.String("seed", "", "starter URL to crawl (http or https)")
	workers := flag.Int("workers", 5, "number of crawl workers")
	flag.Parse()

	seedURL := *seed
	if seedURL == "" {
		fmt.Fprintf(os.Stderr, "usage: %s -seed <url>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	normalized, err := utils.Normalize(seedURL, seedURL)
	if err != nil {
		log.Fatalf("seed must be an absolute http(s) URL, got %q: %v", seedURL, err)
	}

	c, err := crawler.NewCrawler(normalized)
	if err != nil {
		log.Fatalf("crawler init: %v", err)
	}
	c.Start(*workers)
}
