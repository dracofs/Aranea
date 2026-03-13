package dedupe

type Visited struct {
	urls map[string]bool
}

func NewVisited() *Visited {
	return &Visited{urls: make(map[string]bool)}
}

func (v *Visited) Seen(url string) bool {
	return v.urls[url]
}

func (v *Visited) Insert(url string) {
	v.urls[url] = true
}