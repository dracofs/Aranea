package crawler

import (
	"ARANEA/internal/dedupe"
	"ARANEA/internal/fetcher"
	"ARANEA/internal/metrics"
	"ARANEA/internal/parser"
	"ARANEA/internal/queue"
	"ARANEA/internal/redisclient"
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

func NewCrawler(seed string) (*Crawler, error) {
	rdb, err := redisclient.Connect()

	if err != nil {
		return nil, err
	}

	q := queue.NewRedisQueue(rdb)
	s := dedupe.NewRedisSet(rdb)

	if err := q.Clear(); err != nil {
		return nil, err
	}
	if err := s.Clear(); err != nil {
		return nil, err
	}

	c := &Crawler{queue: q, set: s}

	if _, err := s.Add(seed); err != nil {
		return nil, err
	}
	if err := q.Push(seed); err != nil {
		return nil, err
	}

	c.updateQueueDepth()
	return c, nil
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

		fmt.Printf("[Worker %d] Crawling: %s\n", index, curr)

		start := time.Now()
		content, err := fetcher.Fetch(curr)
		metrics.FetchDuration.Observe(time.Since(start).Seconds())

		if err != nil {
			metrics.FetchErrors.Inc()
			log.Printf("[Worker %d] fetch %s: %v", index, curr, err)
			continue
		}

		metrics.PagesCrawled.Inc()

		links, err := parser.GetLinks(content)

		if err != nil {
			log.Printf("[Worker %d] parse %s: %v", index, curr, err)
			continue
		}

		for _, link := range links {
			normalized, err := utils.Normalize(curr, link)
			if err != nil {
				continue
			}

			added, err := c.set.Add(normalized)
			if err != nil {
				log.Println("Error adding to set")
				continue
			}

			if added {
				if err := c.queue.Push(normalized); err != nil {
					log.Println("Error pushing to queue")
					continue
				}
				c.updateQueueDepth()
			}
		}

	}
}
