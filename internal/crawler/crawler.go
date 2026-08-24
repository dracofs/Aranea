package crawler

import (
	"ARANEA/internal/dedupe"
	"ARANEA/internal/fetcher"
	"ARANEA/internal/metrics"
	"ARANEA/internal/parser"
	"ARANEA/internal/queue"
	"ARANEA/internal/utils"
	"fmt"
	"log"
	"time"
)

type Crawler struct {
	queue *queue.RedisQueue
	set   *dedupe.RedisSet
}

func (c *Crawler) Start(workers int) {
	go func() {
		if err := metrics.Serve(":2112"); err != nil {
			log.Printf("metrics server: %v", err)
		}
	}()

	for i := 0; i < workers; i++ {
		go c.Crawl(i)
	}
	select {}
}

func NewCrawler(seed string) *Crawler {
	q := queue.NewRedisQueue()
	s := dedupe.NewRedisSet()

	q.Push(seed)

	c := &Crawler{queue: q, set: s}
	c.updateQueueDepth()
	return c
}

func (c *Crawler) updateQueueDepth() {
	n, err := c.queue.Len()
	if err != nil {
		return
	}
	metrics.QueueDepth.Set(float64(n))
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
		curr, err := c.queue.Pop()
		if err != nil {
			log.Println("Error popping from queue")
			continue
		}
		c.updateQueueDepth()

		c.set.Add(curr)
		fmt.Println("[Worker %d] Crawling: %s\n", index, curr)

		start := time.Now()
		content, err := fetcher.Fetch(curr)
		metrics.FetchDuration.Observe(time.Since(start).Seconds())

		if err != nil {
			metrics.FetchErrors.Inc()
			continue
		}

		metrics.PagesCrawled.Inc()

		links, err := parser.GetLinks(content)

		if err != nil {
			continue
		}

		for _, link := range links {
			normalized, err := utils.Normalize(curr, link)

			if err != nil {
				continue
			}

			exists, err := c.set.Seen(normalized)
			if err != nil {
				log.Println("Error checking set")
				continue
			}

			if !exists {
				c.queue.Push(normalized)
				c.updateQueueDepth()
			}
		}

	}
}
