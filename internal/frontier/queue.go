package frontier

type Queue struct {
	urls []string
}

func NewQueue() *Queue {
	return &Queue{urls: []string{}}
}

func (q *Queue) Push(url string) {
	q.urls = append(q.urls, url)
}

func (q *Queue) Pop() (string, bool) {
	if len(q.urls) == 0 {
		return "", false
	}

	url := q.urls[0]
	q.urls = q.urls[1:]

	return url, true
}
