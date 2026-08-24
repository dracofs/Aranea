package main

import (
	"ARANEA/internal/crawler"
	"log"
)

func main() {
	seed := "https://redis.io/"
	workers := 5

	c, err := crawler.NewCrawler(seed)
	if err != nil {
		log.Fatalf("crawler init: %v", err)
	}
	c.Start(workers)
}
