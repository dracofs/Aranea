package crawler

import (
	"ARANEA/internal/dedupe"
	"ARANEA/internal/fetcher"
	"ARANEA/internal/frontier"
	"ARANEA/internal/parser"
	"ARANEA/internal/utils"
	"fmt"
)

type Crawler struct {
	queue   *frontier.Queue
	visited *dedupe.Visited
}

func newCrawler(seed string) *Crawler {
	q := frontier.NewQueue()
	v := dedupe.NewVisited()

	q.Push(seed)

	return &Crawler{queue: q, visited: v}
}

func (c *Crawler) Crawl() {
	// main crawl loop, structure follows the outline below:
	/*
	 - get url from queue, mark as visited
	 - fetch page
	 - parse page for links, add to queue if not visited
	 - stop after a certain depth is reached or queue is empty
	*/

	for {
		curr, flag := c.queue.Pop()

		if !flag {
			break
		}

		c.visited.Insert(curr)
		fmt.Println("Crawling:", curr)

		content, err := fetcher.Fetch(curr)

		if err != nil {
			continue
		}

		links, err := parser.GetLinks(content)

		if err != nil {
			continue
		}

		for _, link := range links {
			normalized, err := utils.Normalize(curr, link)

			if err != nil {
				continue
			}

			if !c.visited.Seen(normalized) {
				c.queue.Push(normalized)
			}
		}

	}
}
