package crawler

import (
	"ARANEA/internal/dedupe"
	"ARANEA/internal/fetcher"
	"ARANEA/internal/queue"
	"ARANEA/internal/parser"
	"ARANEA/internal/utils"
	"fmt"
)

type Crawler struct {
	queue   *queue.RedisQueue
	set *dedupe.RedisSet
}

func start(c *Crawler, workers int) {
	for i := 0; i < workers; i++ {
		go c.Crawl(i)
	}
	select {}
}
func newCrawler(seed string) *Crawler {
	q := queue.NewRedisQueue()
	s := dedupe.NewRedisSet()

	q.Push(seed)

	return &Crawler{queue: q, set: s}
}

func (c *Crawler) Crawl(index int) {
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

		c.set.Add(curr)
		fmt.Println("[Worker %d] Crawling: %s\n", index, curr)

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

			if !c.set.Seen(normalized) {
				c.queue.Push(normalized)
			}
		}

	}
}
