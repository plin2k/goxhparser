package goxhparser

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type Parser struct {
	File    string
	Service Service
	Content []Content
}

type Service struct {
	XMLName    xml.Name `xml:"xml"`
	Title      string   `xml:"title"`
	Name       string   `xml:"name"`
	EntityID   string   `xml:"entity_id"`
	EntityType string   `xml:"entity_type"`
	Sources    []Source `xml:"source"`
	Rules      []Rule   `xml:"rule"`
}

type Source struct {
	Link     string `xml:",chardata"`
	RuleName string `xml:"rule,attr"`
	Rule     Rule
}

type Rule struct {
	Name  string `xml:"name,attr"`
	Block string `xml:"block"`
	Title string `xml:"title"`
	Link  struct {
		Href   string `xml:",chardata"`
		Prefix string `xml:"prefix,attr"`
	} `xml:"link"`
	ShortContent string `xml:"short_content"`
	FullContent  string `xml:"full_content"`
	Author       string `xml:"author"`
	Rating       string `xml:"rating"`
}

type Content struct {
	Title        string
	Link         string
	ShortContent string
	FullContent  string
	Author       string
	Rating       string
}

func NewParser(filename string) *Parser {
	return &Parser{
		File: filename,
	}
}

func (parser *Parser) XMLToStruct() error {
	xmlFile, err := os.Open(parser.File)
	if err != nil {
		return err
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	err = xml.Unmarshal(byteValue, &parser.Service)
	if err != nil {
		return err
	}
	parser.ruleToSource()
	return nil
}

func (parser *Parser) Exec() error {
	err := parser.XMLToStruct()
	if err != nil {
		return err
	}
	for _, source := range parser.Service.Sources {
		content, err := parser.Parse(source)
		if err != nil {
			return err
		}
		parser.Content = append(parser.Content, content...)
	}
	return nil
}

func (parser *Parser) ruleToSource() {
	for _, rule := range parser.Service.Rules {
		for skey, source := range parser.Service.Sources {
			if source.RuleName == rule.Name {
				parser.Service.Sources[skey].Rule = rule
			}
		}
	}
	parser.Service.Rules = nil
}
