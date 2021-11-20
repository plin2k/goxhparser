package goxhparser

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	ErrorContentNotFound      = "content not found"
	ErrorServerResponseNot200 = "server response not 200"
)

// Parse ...
func (parser *Parser) Parse(source Source) ([]Content, error) {
	var contentOutput []Content
	res, err := http.Get(source.Link)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New(ErrorServerResponseNot200)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	selection := doc.Find(source.Rule.Block)
	if len(selection.Nodes) == 0 {
		err = errors.New(ErrorContentNotFound)
		return nil, err
	}

	selection.Each(func(i int, s *goquery.Selection) {

		var content = Content{
			Title:        source.getTitle(s),
			Link:         source.getLink(s),
			ShortContent: source.getShortContent(s),
			FullContent:  source.getFullContent(s),
			Author:       source.getAuthor(s),
			Rating:       source.getRating(s),
			Date:         source.getDate(s),

			SourceTagName: source.TagName,
		}

		if content.Title != "" && content.Link != "" {
			contentOutput = append(contentOutput, content)
		}
	})
	return contentOutput, nil
}

func (source Source) getTitle(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(source.Rule.Title).Text())
}

func (source Source) getLink(s *goquery.Selection) string {
	if link, exists := s.Find(source.Rule.Link.Href).Attr("href"); exists {
		if source.Rule.Link.Prefix != "" {
			return source.Rule.Link.Prefix + link
		}
		return link
	}
	return ""
}

func (source Source) getShortContent(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(source.Rule.ShortContent).Text())
}

func (source Source) getFullContent(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(source.Rule.FullContent).Text())
}

func (source Source) getAuthor(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(source.Rule.Author).Text())
}

func (source Source) getRating(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find(source.Rule.Rating).Text())
}

func (source Source) getDate(s *goquery.Selection) time.Time {
	var dateStr string
	if dt, ok := s.Find(source.Rule.Date.Time).Attr(source.Rule.Date.Attribute); ok {
		dateStr = dt
	} else {
		dateStr = s.Find(source.Rule.Date.Time).Text()
	}

	date, err := time.Parse(source.Rule.Date.Layout, strings.TrimSpace(dateStr))
	if err != nil {
		return time.Time{}
	}
	return date
}
