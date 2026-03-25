package crawler

import (
	"ARANEA/internal/crawler"
)

func main() {
	seed := "https/example.com"
	c := crawler.NewCrawler(seed)
	c.Crawl()
}
