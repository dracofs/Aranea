package crawler

import (
	"ARANEA/internal/dedupe"
	"ARANEA/internal/frontier"
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
