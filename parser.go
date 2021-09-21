package goxhparser

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

// Parse ...
func parse(source Source) ([]SourceContent, error) {
	var content []SourceContent
	res, err := http.Get(source.URL)
	if err != nil {
		return content, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("status code error: %d %s",res.StatusCode,res.Status))
		return content, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return content, err
	}

	selection := doc.Find(source.Block)
	if len(selection.Nodes) == 0 {
		err = errors.New("content not found")
		return content, err
	}
	selection.Each(func(i int, s *goquery.Selection) {
		title := s.Find(source.Title).Text()
		link, _ := s.Find(source.Link.Href).Attr("href")

		if source.Link.Prefix != "" {
			link = source.Link.Prefix + link
		}

		if title != "" && link != "" {
			content = append(content, SourceContent{
				Title:        title,
				Link:         link,
				ShortContent: s.Find(source.ShortContent).Text(),
				FullContent:  s.Find(source.FullContent).Text(),
				Author:       s.Find(source.Author).Text(),
				Rating:       strings.TrimSpace(s.Find(source.Rating).Text()),
			})
		}
	})
	return content, nil
}
