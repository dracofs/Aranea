package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetLinks(html string) ([]string, error) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, err
	}

	links := []string{}

	doc.Find("a").Each(func(index int, s *goquery.Selection) {
		link, exists := s.Attr("href")

		if exists {
			links = append(links, link)
		}

	})

	return links, nil
}
