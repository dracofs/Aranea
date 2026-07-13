package crawler

import (
	"ARANEA/internal/crawler"
)

func main() {
	seed := "https/example.com"
	workers := 5

	c := crawler.NewCrawler(seed)
	c.start(workers)
}
