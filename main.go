package goxhparser

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type Service struct {
	XMLName    xml.Name `xml:"xml"`
	Title      string   `xml:"title"`
	Name       string   `xml:"name"`
	EntityID   string   `xml:"entity_id"`
	EntityType string   `xml:"entity_type"`
	Sources    []Source `xml:"source"`
}

type Source struct {
	URL   string `xml:"url"`
	Block string `xml:"block"`
	Title string `xml:"title"`
	Link  struct {
		Href   string `xml:",chardata"`
		Prefix string `xml:"prefix,attr"`
	} `xml:"link"`
	ShortContent  string `xml:"short_content"`
	FullContent   string `xml:"full_content"`
	Author        string `xml:"author"`
	Rating        string `xml:"rating"`
	SourceContent []SourceContent
}

type SourceContent struct {
	Title        string
	Link         string
	ShortContent string
	FullContent  string
	Author       string
	Rating       string
}

func XMLToStruct(filename string) (Service, error) {
	
	xmlFile, err := os.Open(filename)
	if err != nil {
		return Service{}, err
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var service Service

	err = xml.Unmarshal(byteValue, &service)
	if err != nil {
		return Service{}, err
	}
	return service, nil
}

func ParseSource(source Source) ([]SourceContent, error) {
	content, err := parse(source)
	if err != nil {
		return content, err
	}
	return content, nil
}

func ParseSources(sources []Source) ([]Source, error) {
	for key, source := range sources {
		content, err := ParseSource(source)
		if err != nil {
			return sources, err
		}
		sources[key].SourceContent = content
	}
	return sources, nil
}

func ParseByXMLFile(filename string) ([]Source, error) {
	service, err := XMLToStruct(filename)
	if err != nil {
		return []Source{}, err
	}
	sources, err := ParseSources(service.Sources)
	if err != nil {
		return []Source{}, err
	}
	return sources, nil
}
